package render

import (
    "buildacoin/data"
    "buildacoin/source"
    "bytes"
    "encoding/hex"
    "encoding/gob"
    "errors"
    "net/http"
    "strconv"
    "strings"
    cointemplate "buildacoin/template"
    webutil "buildacoin/web/util"
)

// web page where base coin template inputs are presented to the user on GET
// and a coin is built from input values on POST
type CoinPage struct {
    base *data.Meta
    conf *data.Conf
    markupTemplate *basePage
    inputs [][]data.Input
    db data.DB
}

func NewCoinPage(conf *data.Conf, base *data.Meta) (*CoinPage, error) {
    db, err := data.DBConnect(conf)
    if err != nil {
        return nil, err
    }

    markup, err := SimplePage(conf, "markup/coin.html", "style/coin.css")
    if err != nil {
        return nil, err
    }

    // build input descriptors
    groupCache := make(map[string][]data.Input)
    // populate map with slices
    for ii := 0; ii < base.InGroupCount(); ii++ {
        groupCache[base.InGroup(ii)] = make([]data.Input, 0, 8)
    }
    // organize inputs by group
    for ii := 0; ii < base.InputCount(); ii++ {
        input := base.Input(ii)
        groupCache[input.Group] = append(groupCache[input.Group], input)
    }
    inputs := make([][]data.Input, 0, len(groupCache))
    for ii := 0; ii < base.InGroupCount(); ii++ {
        inputs = append(inputs, groupCache[base.InGroup(ii)])
    }

    return &CoinPage { base, conf, markup, inputs, db }, nil
}

func (tt *CoinPage) ServeHTTP(out http.ResponseWriter, req *http.Request) {
    switch req.Method {
    case "GET":
        tt.serveForm(out, req, nil, nil)
    case "POST":
        tt.serveCoin(out, req)
    default:
        out.WriteHeader(http.StatusMethodNotAllowed)
    }
}

func (tt *CoinPage) serveForm(out http.ResponseWriter, req *http.Request,
        values map[string]string, errs []error) {
    inputs := tt.inputs
    // override default form values if supplied
    if values != nil {
        inputs = make([][]data.Input, len(tt.inputs))
        for groupIdx, group := range tt.inputs {
            inputs[groupIdx] = make([]data.Input, len(group))
            for inputIdx, input := range group {
                newInput := input
                override, ok := values[input.Id]
                if ok {
                    newInput.Default = override
                }
                inputs[groupIdx][inputIdx] = newInput
            }
        }
    }
    err := tt.markupTemplate.Execute(out, map[string]interface{} {
        "groups": inputs,
    }, errs)
    if err != nil {
        NewErrorPage(tt.conf,
            "error creating coin form: " + err.Error()).ServeHTTP(out, req)
        // TODO log
        return
    }
}

func (tt *CoinPage) serveCoin(out http.ResponseWriter, req *http.Request) {
    err := req.ParseForm()
    if err != nil {
        err = errors.New("error reading coinfiguration data: " + err.Error())
        tt.serveForm(out, req, nil, []error { err })
        // TODO log
        return
    }

    // unpack the www form data
    values := make(map[string]string)
    for key, value := range req.Form {
        values[key] = value[0]
    }

    // create a unique coin ID; this will be the only distinct handle on a
    // previously generated coin
    coinID := data.NewCoinID()

    coinName := values["name"]
    if len(coinName) < 1 {
        coinName = "build-a-coin"
    }
    coinName = strings.ToLower(coinName)

    filterMap, err := source.BuildFilterMap(tt.base, values)
    if err != nil {
        tt.serveForm(out, req, values, []error { err })
        // TODO log
        return
    }
    // If the coin template has a place for the coin ID, insert it
    // appropriately
    for _, sub := range tt.base.Subs() {
        if sub.Comment == "unique generated coin id" {
            filterMap[sub.Idx] = []byte(hex.EncodeToString(coinID.Bytes()))
            break
        }
    }

    template, streamType, err := tt.base.Template()
    if err != nil {
        const terse = "error getting coin template"
        if tt.conf.Debug() {
            err = errors.New(terse + ": " + err.Error())
        } else {
            err = errors.New(terse)
        }
        tt.serveForm(out, req, values, []error { err })
        // TODO log
        return
    }

    runner := cointemplate.GetRunner(streamType)

    outHeader := out.Header()
    outHeader["Content-Disposition"] = []string { "attachment; filename=" +
            coinName + "." + streamType }

    err = runner.Run(out, template, filterMap)
    if err != nil {
        // TODO log
        return
    }

    template.Close()

    // add the newborn coin to the db
    addrId, err := strconv.ParseUint(values["versionbyte"], 0, 8)
    if err != nil {
        // TODO log
    }
    protoPort, err := strconv.ParseUint(values["port"], 0, 16)
    if err != nil {
        // TODO log
    }
    initReward, err := strconv.ParseFloat(values["initial reward"], 64)
    if err != nil {
        // TODO log
    }
    halving, err := strconv.ParseInt(values["reward halving"], 0, 32)
    if err != nil {
        // TODO log
    }
    blockTime, err := strconv.ParseInt(values["block time"], 0, 64)
    if err != nil {
        // TODO log
    }
    diffTime, err := strconv.ParseInt(values["retarget window"], 0, 64)
    if err != nil {
        // TODO log
    }
    initDiff, err := strconv.ParseFloat(values["difficulty"], 64)
    if err != nil {
        // TODO log
    }
    serialBuf := new(bytes.Buffer)
    enc := gob.NewEncoder(serialBuf)
    err = enc.Encode(map[uint][]byte(filterMap))
    if err != nil {
        // TODO log
    }

    newCoin := &data.CoinSummary {
        coinID,
        tt.base.Id(),
        tt.base.Version(),
        values["name"],
        values["shortname"],
        byte(addrId),
        uint16(protoPort),
        initReward,
        int32(halving),
        blockTime,
        diffTime,
        initDiff,
        values["genesis message"],
        serialBuf.Bytes(),
    }

    err = tt.db.PutCoinSummary(newCoin, webutil.RequestOrigin(tt.conf, req),
            req.UserAgent())
    if err != nil {
        println(err.Error())
        // TODO log
    }
}

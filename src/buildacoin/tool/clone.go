package tool

import (
    "buildacoin/data"
    "buildacoin/template"
    "bytes"
    "encoding/gob"
    "fmt"
    "os"
    "strings"
)

func Clone(conf *data.Conf, id data.CoinID) {
    // Set up external dependencies.
    db, err := data.DBConnect(conf)
    if err != nil {
        fmt.Fprintln(os.Stderr, "failed to connect to db: ", err.Error())
        return
    }

    // Using the given coin id, find the rest of the coin's parameters,
    // including a gob-serialized copy of the substitution map used originally.
    coin, err := db.GetCoinSummary(id)
    if err != nil {
        fmt.Fprintln(os.Stderr, "failed to fetch coin info: ", err.Error())
        return
    }

    // Deserialize the substitution map.
    gob_buf := bytes.NewBuffer(coin.Serialized)
    dec := gob.NewDecoder(gob_buf)
    var filter_map template.FilterMap
    err = dec.Decode(&filter_map)
    if err != nil {
        fmt.Fprintln(os.Stderr, "failed to inflate filter map: ", err.Error())
        return
    }

    // Load the metadata of the base coin the target was created with.
    meta, err := data.LoadMeta(conf, coin.TemplateID)
    if err != nil {
        fmt.Fprintln(os.Stderr,
            "failed to load base coin info: ", err.Error())
        return
    }
    if meta.Version() != coin.TemplateVer {
        fmt.Fprintf(os.Stderr, "warning: available template version %s does " +
                               "not match version originally used to create " +
                               "coin (%s)\n", meta.Version(), coin.TemplateVer)
    }

    // Open a stream to read the template data.
    coin_template, inStreamType, err := meta.Template()
    if err != nil {
        fmt.Fprintln(os.Stderr,
            "failed to load coin template: ", err.Error())
        return
    }
    runner := template.GetRunner(inStreamType)

    // Open an output stream for the cloned coin tarball.
    out_file, err := os.Create(strings.ToLower(coin.Name) + "." + inStreamType)
    if err != nil {
        fmt.Fprintln(os.Stderr,
            "failed to open output file: ", err.Error())
        return
    }
    defer out_file.Close()

    err = runner.Run(out_file, coin_template, filter_map)
    if err != nil {
        fmt.Fprintln(os.Stderr,
            "failed to render coin: ", err.Error())
        return
    }
}

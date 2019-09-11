package data

import (
    "crypto/rand"
    _ "github.com/lib/pq"
    "database/sql"
    "encoding/hex"
    "strconv"
    "time"
)

//
// Database connection
//

type DB struct {
    *sql.DB
}
var nilDB DB

func DBConnect(conf *Conf) (DB, error) {
    confString := "dbname=" + conf.DBName() + " " +
            "host=" + conf.DBHost() + " " +
            "port=" + strconv.FormatUint(uint64(conf.DBPort()), 10) + " " +
            "user=" + conf.DBUser() + " " +
            "password=" + conf.DBPass() + " " +
            "sslmode=disable"
    rawDB, err := sql.Open("postgres", confString)
    if err != nil {
        return nilDB, err
    }
    return DB { rawDB }, nil
}

//
// coins table
//

type CoinSummary struct {
    ID CoinID
    TemplateID string
    TemplateVer string
    Name string
    Code string
    AddrId byte
    ProtoPort uint16
    InitReward float64
    Halving int32
    BlockTime int64
    DiffTime int64
    InitDiff float64
    EmbedMsg string
    Serialized []byte
}

func (tt DB) GetCoinSummary(id CoinID) (CoinSummary, error) {
    out := CoinSummary { ID: id }
    row := tt.QueryRow("SELECT template_id, template_ver, label, " +
            "currency_code, addr_id, proto_port, initial_reward, " +
            "halving_interval, block_time, diff_time, initial_diff, " +
            "embed_msg, subs_gob FROM coins WHERE id=$1",
            hex.EncodeToString(id.Bytes()))
    err := row.Scan(&out.TemplateID, &out.TemplateVer, &out.Name, &out.Code,
            &out.AddrId, &out.ProtoPort, &out.InitReward, &out.Halving,
            &out.BlockTime, &out.DiffTime, &out.InitDiff, &out.EmbedMsg,
            &out.Serialized)
    return out, err
}

func (tt DB) PutCoinSummary(input *CoinSummary, requestOrigin string,
        requestAgent string) error {
    requestOrigin_p := &requestOrigin
    if requestOrigin == "" {
        requestOrigin_p = nil
    }
    requestAgent_p := &requestAgent
    if requestAgent == "" {
        requestAgent_p = nil
    }

    _, err := tt.Exec("INSERT INTO coins VALUES ( $1, $2, $3, $4, $5, $6, " +
            "$7, $8, $9, $10, $11, $12, $13, $14, $15, $16 )",
            hex.EncodeToString(input.ID.Bytes()), input.TemplateID,
            input.TemplateVer, input.Name, input.Code, input.AddrId,
            input.ProtoPort, input.InitReward, input.Halving, input.BlockTime,
            input.DiffTime, input.InitDiff, input.EmbedMsg, input.Serialized,
            requestOrigin_p, requestAgent_p)
    return err
}

//
// compile_orders table
//

type CompileOrder struct {
    ID BuildID
    Coin CoinID
    DepositAddr *string
    Created time.Time
    Started *time.Time
    Ended *time.Time
    Deposited float64
    Credited float64
    Status string
}

//
// Identifiers
//

const CoinIDLen = 16
type CoinID [CoinIDLen]byte
var zeroCoinID CoinID

func NewCoinID() CoinID {
    var output CoinID
    n, err := rand.Read(output[:])
    if n < CoinIDLen || err != nil {
        panic("failed to get random bytes for new coin id")
    }
    if output == zeroCoinID {
        panic("the zero CoinID is illegal")
    }
    return output
}

func (tt CoinID) Bytes() []byte {
    return tt[:]
}

type BuildID [CoinIDLen]byte

func (tt BuildID) Bytes() []byte {
    return tt[:]
}

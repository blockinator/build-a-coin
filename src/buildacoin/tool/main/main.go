package main

import (
    "buildacoin/data"
    "buildacoin/tool"
    "encoding/hex"
    "flag"
    "fmt"
    "os"
)

func main() {
    //
    // Command line arguments and their destination variables:
    //
    var cloneId string
    flag.StringVar(&cloneId, "clone", "",
        "clone the coin with the given id into a tarball in the current " +
        "directory")

    var migrationName string
    flag.StringVar(&migrationName, "migrate", "",
        "perform a named migration, or 'list'")

    var confPath string
    flag.StringVar(&confPath, "conf", "", "load conf file instead of default")

    flag.Parse()

    conf, err := data.LoadConfFromArg(confPath)
    if err != nil {
        fmt.Fprintln(os.Stderr, "can't load config: " + err.Error())
        return
    }

    //
    // Dispatch flow based on arguments.  If multiple commands are given,
    // precedence is decided by this if-else chain.
    //
    if cloneId != "" {
        // Clone command: clone a coin in the database by its id.
        coin_id, err := coinIdFromHex(cloneId)
        if err != nil {
            fmt.Fprintln(os.Stderr, "bad coin id: " + err.Error())
            return
        }
        tool.Clone(conf, coin_id)
    } else if migrationName != "" {
        // Migrate command: perform some transformation on the persistent state
        // of a server instance.
        tool.Migrate(conf, migrationName)
    } else {
        // No command given.
        fmt.Fprintln(os.Stderr, "no command.  try -help")
        return
    }
}

func coinIdFromHex(id_string string) (data.CoinID, error) {
    out := data.CoinID {}

    bytes, err := hex.DecodeString(id_string)
    if err != nil {
        return out, err
    }

    if len(bytes) != len(out) {
        return out, fmt.Errorf("CoinID must be %d bytes but got %d",
            len(out), len(bytes))
    }

    copy(out[:], bytes)

    return out, nil
}

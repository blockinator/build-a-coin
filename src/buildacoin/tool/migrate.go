package tool

import (
    "buildacoin/data"
    "buildacoin/template"
    "bytes"
    "database/sql"
    "encoding/gob"
    "errors"
    "fmt"
    "os"
)

type migration struct {
    name string
    description string
    execute func(*data.Conf)
}

var all_migrations = []migration {
    {
        name: "fix_disjoint_coin_ids",
        description: "Older versions of buildacoin erroneously generated " +
                     "separate coin ids for the database and for the output " +
                     "coin.  Overwrite the database ids with those inserted " +
                     "into the generated coins",
        execute: fixDisjointCoinIds,
    },
}

func fixDisjointCoinIds(conf *data.Conf) {
    // Set up external dependencies.
    db, err := data.DBConnect(conf)
    if err != nil {
        fmt.Fprintln(os.Stderr, "failed to connect to db: ", err.Error())
        return
    }
    defer db.Close()

    edit_tx, err := db.Begin()
    if err != nil {
        fmt.Fprintln(os.Stderr, "failed to open a db transaction: ",
            err.Error())
        return
    }

    // Since the coin id actually used in the source code is only stored in the
    // serialized subs_gob, it is necessary to iterate over all coins,
    // deserialize their substitutions, and compare that id with the on in the
    // database.
    coin_row, err := db.Query("SELECT id,subs_gob FROM coins")
    if err != nil {
        fmt.Fprintln(os.Stderr, "failed to enumerate coins: ", err.Error())
        return
    }
    defer coin_row.Close()

    for coin_row.Next() {
        err = fixSingleDisjointCoinIdRow(conf, coin_row, edit_tx)
        if err != nil {
            edit_tx.Rollback()
            fmt.Fprintln(os.Stderr, err.Error())
            return
        }
    }
    edit_tx.Commit()
}

// Helper for fixDisjointCoinIds.  Extract both the database coin id and the id
// used in the output source code, and overwrite the database coin id if they
// differ.
func fixSingleDisjointCoinIdRow(conf *data.Conf,
                                row *sql.Rows,
                                edit_tx *sql.Tx) error {
    const CoinIdIndex = 51
    var err error

    // Grab the coin's current id and the gob of all of its source
    // substitutions.
    var db_id string
    var gob_bytes []byte

    err = row.Scan(&db_id, &gob_bytes)
    if err != nil {
        return errors.New("bad db coin data: " + err.Error())
    }

    // Deserialize the gob to get the substitution map back.
    gob_buffer := bytes.NewBuffer(gob_bytes)
    gob_decoder := gob.NewDecoder(gob_buffer)
    var filter_map template.FilterMap
    err = gob_decoder.Decode(&filter_map)
    if err != nil {
        return errors.New("failed to inflate filter map: " + err.Error())
    }

    // Grab the coin id that the user actually received out of the substitution
    // map.
    _, ok := filter_map[CoinIdIndex]
    if !ok {
        return errors.New("coin identifier value not found in substution map")
    }
    source_id := string(filter_map[CoinIdIndex])

    if (source_id == db_id) {
        fmt.Printf("%s already consistent\n", db_id)
    } else {
        fmt.Printf("%s to be updated to %s\n", db_id, source_id)
        result, err := edit_tx.Exec("UPDATE coins SET id = $1 WHERE id = $2",
                                    source_id, db_id)
        if err != nil {
            return errors.New("failed to update database id: " + err.Error())
        }
        rows_affected, err := result.RowsAffected()
        if err != nil { panic(err.Error()) }
        if rows_affected != 1 {
            return fmt.Errorf("changed %d rows with id %s, but it only " +
                              "makes sense for exactly one to change",
                              rows_affected)
        }
    }

    return nil
}

func Migrate(conf *data.Conf, migrationName string) {
    // 'list' is a special pseudo-migration
    if migrationName == "list" {
        for _, migration := range all_migrations {
            fmt.Printf("%s: %s\n", migration.name, migration.description)
        }
        return
    }

    // Look up the target migration and run it.
    target := findMigration(migrationName)
    if target == nil {
        fmt.Fprintln(os.Stderr, "unknown migration '" + migrationName +
                                "'.  try 'list'")
        return
    }
    target.execute(conf)
}

func findMigration(migrationName string) *migration {
    for _, candidate := range all_migrations {
        if candidate.name == migrationName {
            return &candidate
        }
    }
    return nil
}

package data

import (
    "testing"
)

func TestDBConnect(t *testing.T) {
    conf := DefaultConf.WithDBPass("testpw")

    db, err := DBConnect(conf)
    if err != nil {
        t.Fatal(err.Error())
    }

    db.Ping()
}

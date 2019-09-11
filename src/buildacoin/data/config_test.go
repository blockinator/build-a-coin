package data

import (
    "fmt"
    "testing"
)

const (
    TestDir = "../../../testing"
)

var goodConf *Conf = DefaultConf.WithBasesDir("/foo/bases")

func TestLoadConf(t *testing.T) {
    conf, err := LoadConfExplicit(TestDir + "/good_conf.json")
    if err != nil {
        t.Fatal(err.Error())
    }

    // jury-rigged equality in the face of limited external use and embedded
    // slices
    if fmt.Sprint(conf) != fmt.Sprint(goodConf) {
        t.Fatalf("good conf mismatch:\nexpected\n%v\nactual\n%v\n", *goodConf,
                *conf)
    }
}

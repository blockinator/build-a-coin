package source

import (
    "buildacoin/data"
    "buildacoin/template"
    "bytes"
    "testing"
)

func TestSimpleGenerate(t *testing.T) {
    expected := template.FilterMap {
        2: []byte("Hello"),
        9: []byte("World"),
        13: []byte("!"),
    }

    actual, err := BuildFilterMap(simpleMeta, simpleMap)
    if err != nil {
        t.Fatal(err.Error())
    }

    for key, val := range actual {
        if bytes.Compare(val, expected[key]) != 0 {
            t.Fatalf("filter map mismatch\nexpected\n%v\nactual\n%v\n", expected, actual)
        }
    }
}

var simpleMap = map[string]string {
    "first": "Hello",
    "second": "World",
}
var simpleMeta *data.Meta = data.NewMeta("", "", "", make([]string, 0),
        make([]data.Input, 0),
        []data.Sub { data.Sub { 2, "first", "", "-", "literal", nil },
            data.Sub { 13, "third", "", "!", "", nil },
            data.Sub { 9, "second", "", "-", "literal", nil } })

package data

import (
    "fmt"
    "testing"
)

var metadataConf = &Conf {
    conf_ {
        BasesDir: TestDir + "/bases",
    },
}

func TestDeserializeMeta(t *testing.T) {
    expected := &Meta {
        meta_: meta_ {
            Id: "simple",
            Label: "Simple Base Example",
            Version: "1.0",
            InGroups: []string {
                "basics",
            },
            Inputs: []Input {
                Input {
                    Group: "basics",
                    Id: "name",
                    Label: "the name of that thing you're making",
                    Default: "roflcopter.com",
                },
            },
            Subs: []Sub {
                Sub {
                    Idx: 1138,
                    Input: "name",
                    Comment: "this is a comment",
                    Default: "default.com",
                    Type: "literal",
                },
            },
        },
        conf: metadataConf,
    }

    actual, err := LoadMeta(metadataConf, "simple")
    if err != nil {
        t.Fatal(err.Error())
    }

    // jury-rigged equality in the face of limited external use and embedded
    // slices
    if fmt.Sprint(expected) != fmt.Sprint(actual) {
        t.Fatalf("metadata mismatch:\nexpected\n%v\nactual\n%v\n", *expected,
                *actual)
    }
}

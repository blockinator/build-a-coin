package data

import (
    "bytes"
    "errors"
    "io"
    "testing"
)

var testConf = &Conf {
    conf_ {
        ConfPath: DefaultConfPath,
        BasesDir: TestDir + "/bases",
    },
}

func TestPlaintextTemplate(t *testing.T) {
    testTemplateOk(t, "plaintext data", "plaintext")
}

func TestMissingTemplate(t *testing.T) {
    testTemplateErr(t, ErrNoSrc, "missing")
}

func TestAmbiguousTemplate(t *testing.T) {
    testTemplateErr(t, ErrAmbiguousSrc, "ambiguous")
}

//
// helper functions
//

func testTemplate(t *testing.T, expected string, baseName string) error {
    buf := new(bytes.Buffer)

    src, _, err := LoadTemplate(testConf, baseName)
    if err != nil {
        return err
    }

    _, err = io.Copy(buf, src)
    if err != nil {
        return err
    }

    actual := buf.String()
    if actual != expected {
        return errors.New("base mismatch: expected '" + expected +
                "' / actual '" + actual + "'")
    }
    return nil
}

func testTemplateOk(t *testing.T, expected string, baseName string) {
    err := testTemplate(t, expected, baseName)
    if err != nil {
        t.Fatal(err.Error())
    }
}

func testTemplateErr(t *testing.T, expected error, baseName string) {
    actual := testTemplate(t, "", baseName)
    if actual != expected {
        t.Fatalf("Wrong error: expected '%v' / actual '%v'\n", expected, actual)
    }
}

package template

import (
    "testing"
)

func TestTarCompress(t *testing.T) {
    unknown := TarRunner { "lolpression" }

    switch unknown.Run(nil, nil, nil).(type) {
    case ErrUnknCompress:
        // good!
    default:
        t.Fatal("tar runner failed to catch weird compression")
    }
}

package template

import (
    "io"
    "strings"
)

type Runner interface {
    Run(io.Writer, io.Reader, FilterMap) error
}

func GetRunner(archiveType string) Runner {
    parts := strings.Split(strings.ToLower(archiveType), ".")
    if len(parts) < 1 {
        return FilterRunner{}
    }
    switch parts[0] {
    case "tar":
        if len(parts) > 2 {
            return nil
        }
        compress := Uncompressed
        if len(parts) == 2 {
            compress = parts[1]
        }
        return TarRunner { compress }
    default:
        return FilterRunner{}
    }
}

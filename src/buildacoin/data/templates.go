package data

import (
    "compress/gzip"
    "errors"
    "io"
    "os"
    "path/filepath"
    "strings"
)

const (
    // Filename suffix indicating gzip compression
    GzipSuffix = ".gz"
)

var (
    // Error when more than one file appears to be the template for a base coin
    ErrAmbiguousSrc error = errors.New("template source name is ambiguous")
    // Error when no files appear to be the template for a base coin
    ErrNoSrc error = errors.New("no such base or template source")
    // Error when the format of a template file does not match the format
    // implied by its name
    ErrBadSrcFormat error = errors.New("bad format for template source file")
)

// Get the template stream for the named base coin.  Returns the stream and the
// file extension of the source file.
func LoadTemplate(conf *Conf, name string) (io.ReadCloser, string, error) {
    glob := conf.BasesDir() + "/" + name + "/template.*"
    matches, err := filepath.Glob(glob)

    if err != nil {
        return nil, "", err
    }

    if len(matches) < 1 {
        return nil, "", ErrNoSrc
    }
    if len(matches) > 1 {
        return nil, "", ErrAmbiguousSrc
    }
    path := matches[0]

    file, err := os.Open(path)
    if err != nil {
        return nil, "", err
    }

    _, filename := filepath.Split(path)
    ext := strings.Join(strings.Split(filename, ".")[1:], ".")

    return file, ext, nil
}

// baseReader is a chain of io.ReadClosers that are chained together and all
// closed with baseReader.Close()
type baseReader struct {
    closers []io.Closer
    reader io.Reader
}

func newBaseReader(origin io.ReadCloser) *baseReader {
    return &baseReader {
        closers: []io.Closer { origin },
        reader: origin,
    }
}

func (tt *baseReader) wrapGzip() error {
    reader, err := gzip.NewReader(tt.reader)
    if err != nil {
        return err
    }
    tt.reader = reader
    tt.closers = append(tt.closers, reader)
    return nil
}

func (tt baseReader) Read(input []byte) (int, error) {
    return tt.reader.Read(input)
}
func (tt baseReader) Close() error {
    closeErrs := []error(nil)
    // close readers in reverse order (starting with the outermost)
    for ii := len(tt.closers)-1; ii >= 0; ii-- {
        closer := tt.closers[ii]

        err := closer.Close()
        if err != nil {
            closeErrs = append(closeErrs, err)
        }
    }

    err := error(nil)

    if len(closeErrs) > 0 {
        errStr := ""
        for idx, err := range closeErrs {
            errStr += err.Error()
            if idx < len(closeErrs) - 1 {
                errStr += " + "
            }
        }
    }

    return err
}

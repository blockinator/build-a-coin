package template

import (
    "archive/tar"
    "compress/gzip"
    "io"
    "io/ioutil"
    "strings"
)

const Uncompressed = ""

type TarRunner struct {
    Compression string
}

func (tt TarRunner) Run(dst io.Writer, src io.Reader, values FilterMap) error {
    // wrap with appropriate compression
    var (
        compressIn io.ReadCloser
        compressOut io.WriteCloser
    )
    switch tt.Compression {
    case "gz":
        var err error
        compressIn, err = gzip.NewReader(src)
        if err != nil {
            return err
        }
        compressOut = gzip.NewWriter(dst)
    case Uncompressed:
        compressIn = ioutil.NopCloser(src)
        compressOut = WriterNopCloser(dst)
    default:
        return ErrUnknCompress { tt.Compression }
    }
    defer func() {
        compressIn.Close()
        compressOut.Close()
    }()

    // open tar streams
    tarIn := tar.NewReader(compressIn)
    tarOut := tar.NewWriter(compressOut)
    defer func() {
        tarOut.Close()
    }()

    // initialize filter with a null reader; it will be reset for each stream
    // in the tar
    filter := NewFilter(nil, values)

    // the comm channel is used to allow writing out to the tar concurrently
    // while templating the next file.  Since the first file has no previous
    // file to wait for, start a goroutine giving the first file the go-ahead
    // immediately
    comm := make(chan error)
    go func() {
        comm <- nil
    }()

    filterloop:
    for {
        // get next file from tar
        header, err := tarIn.Next()
        if err == io.EOF {
            // wait for the previous write-out to finish
            err = <-comm
            if err != nil {
                return err
            }
            break filterloop
        }
        if err != nil {
            return err
        }

        // filter header
        buf, _ := ioutil.ReadAll(filter.Reset(strings.NewReader(header.Name)))
        header.Name = string(buf)
        buf, _ =
            ioutil.ReadAll(filter.Reset(strings.NewReader(header.Linkname)))
        header.Linkname = string(buf)

        // non-streaming filtration of file to get size for tar header.
        fileData, err := ioutil.ReadAll(filter.Reset(tarIn))
        if err != nil {
            return ErrFilterFailure { header.Name, err }
        }
        header.Size = int64(len(fileData))

        // wait for the previous write-out to finish
        err = <-comm
        if err != nil {
            return err
        }

        // write out to tar
        go func(header *tar.Header, body []byte) {
            var err error
            defer func() {
                comm <- err
            }()
            err = tarOut.WriteHeader(header)
            if err != nil {
                return
            }
            _, err = tarOut.Write(body)
        }(header, fileData)
    }

    return nil
}

//
// Errors
//

type ErrUnknCompress struct { Compression string }
func (tt ErrUnknCompress) Error() string {
    return "unknown compression format '" + tt.Compression + "'"
}
type ErrFilterFailure struct { templName string; cause error }
func (tt ErrFilterFailure) Error() string {
    return "failed to filter '" + tt.templName + "': " + tt.cause.Error()
}

//
// nopWriteCloser: for some reason, ioutil only has the reader analog of this.
//

type nopWriteCloser struct {
    io.Writer
}
func (tt nopWriteCloser) Write(buf []byte) (int, error) {
    return tt.Write(buf)
}
func (tt nopWriteCloser) Close() error {
    return nil
}
func WriterNopCloser(writer io.Writer) io.WriteCloser {
    return nopWriteCloser { writer }
}

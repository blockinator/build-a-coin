package template

import (
    "io"
    "strings"
    "testing"
)

//
// testing
//

func simpleTemplateTest(t *testing.T, template string, sub string,
        expec string, bufSize uint) bool {
    subs := map[uint][]byte {
        1 : []byte(sub),
    }
    source := strings.NewReader(template)

    filter := newFilterDebug(source, subs, bufSize)

    buf := make([]byte, bufSize)
    res := ""

    n, err := 0, error(nil)
    for err != io.EOF {
        if err != nil {
            t.Error("error reading from filter: " + err.Error())
            return false
        }

        n, err = filter.Read(buf)

        res += string(buf[:n])
    }

    if res != expec {
        t.Error("bad template result: expect '" + expec + "' got '" +
                res + "'")
        return false
    }
    return true
}

func TestSingleSubstitution(t *testing.T) {
    simpleTemplateTest(t, "__._1-", "substitution!", "substitution!",
            InitBufSize)
}

func TestMultiSubstitution(t *testing.T) {
    str := "substitution!"
    for ii := 0; ii < 4; ii++ {
        template := ""
        result := ""
        for jj := 0; jj < ii; jj++ {
            template += " "
            result += " "
        }
        template += "__._1-"
        result += str
        for jj := 0; jj < ii; jj++ {
            template += " "
            result += " "
        }
        template += "__._1-"
        result += str
        for jj := 0; jj < ii; jj++ {
            template += " "
            result += " "
        }
        if !simpleTemplateTest(t, template, str, result,
                InitBufSize) {
            return
        }
    }
}

func TestSplitMarker(t *testing.T) {
    simpleTemplateTest(t, "__._1-", "a", "a", 2)
}

func TestOutputOverflow(t *testing.T) {
    for ii := uint(1); ii < 128; ii++ {
        if !simpleTemplateTest(t, "__._1-",
                "supercalifragilisticexpialidocious",
                "supercalifragilisticexpialidocious", ii) {
            return
        }
    }
}

func TestInputOverflow(t *testing.T) {
    for ii := uint(1); ii < 256; ii++ {
        if !simpleTemplateTest(t, "__._1- __._1-",
                "supercalifragilisticexpialidocious",
                "supercalifragilisticexpialidocious " +
                    "supercalifragilisticexpialidocious", ii) {
            return
        }
    }
}

func TestPartialTotalMatchOverlap(t *testing.T) {
    simpleTemplateTest(t, "___._1-", "substitution!", "_substitution!", 10)
}

func TestDistinctSubs(t *testing.T) {
    subA := "lollerskates"
    subB := "roflcopter"
    subC := "bbq"
    template := "  __._98765432- __._7-         __._1138- "
    expected := "  bbq lollerskates         roflcopter "

    subs := map[uint][]byte {
        7 : []byte(subA),
        1138 : []byte(subB),
        98765432 : []byte(subC),
    }
    source := strings.NewReader(template)

    filter := newFilterDebug(source, subs, InitBufSize)

    buf := make([]byte, InitBufSize)

    n, err := filter.Read(buf)
    if err != nil && err != io.EOF {
        t.Error("error reading from filter: " + err.Error())
        return
    }

    if string(buf[:n]) != expected {
        t.Error("bad template result: expect '" + expected + "' got '" +
                string(buf) + "'")
        return
    }
}

//
// benchmarking
//

type spaceReader struct {
    length uint
    count uint
}
func (s *spaceReader) Read(out []byte) (int, error) {
    n := uint(0)
    outByte := []byte(" ")[0]
    for ; int(n) < len(out) && n < s.length - s.count; n++ {
        out[n] = outByte
    }
    s.count += n
    err := error(nil)
    if s.count >= s.length {
        err = io.EOF
    }
    return int(n), err
}

func benchmarkSeek(b *testing.B, bufsize uint) {
    subs := make(map[uint][]byte)
    buf := make([]byte, bufsize)
    for ii := 0; ii < b.N; ii++ {
        input := &spaceReader{ length : 1 * (1024 * 1024) }
        filter := NewFilter(input, subs)
        err := error(nil)
        for err == nil {
            _, err = filter.Read(buf)
        }
        if err != io.EOF {
            b.Fatal("templating failed: " + err.Error())
        }
    }
}

func BenchmarkSeek8KBuf(b *testing.B) {
    benchmarkSeek(b, 8 * 1024)
}
func BenchmarkSeek4KBuf(b *testing.B) {
    benchmarkSeek(b, 4 * 1024)
}
func BenchmarkSeek16KBuf(b *testing.B) {
    benchmarkSeek(b, 16 * 1024)
}
func BenchmarkSeek1MBuf(b *testing.B) {
    benchmarkSeek(b, 1 * 1024 * 1024)
}

type subMarkerReader struct {
    target uint
    count uint
}
func (s *subMarkerReader) Read(out []byte) (int, error) {
    n := uint(0)
    outBytes := []byte("__._1-")
    thisCount := uint(len(out) / len(outBytes))
    for ; n < thisCount && n < s.target - s.count; n++ {
        copy(out[n*uint(len(outBytes)):], outBytes)
    }
    s.count += n
    err := error(nil)
    if s.count >= s.target {
        err = io.EOF
    }
    return int(n) * len(outBytes), err
}

func benchmarkSubstitution(b *testing.B, bufsize uint) {
    subs := map[uint][]byte {
        1 : []byte("123456"),
    }
    buf := make([]byte, bufsize)
    for ii := 0; ii < b.N; ii++ {
        input := &subMarkerReader{ target : (1 * (1024 * 1024)) / 6 }
        filter := NewFilter(input, subs)
        err := error(nil)
        for err == nil {
            _, err = filter.Read(buf)
        }
        if err != io.EOF {
            b.Fatal("templating failed: " + err.Error())
        }
    }
}

func BenchmarkSubstitution8KBuf(b *testing.B) {
    benchmarkSubstitution(b, 8 * 1024)
}
func BenchmarkSubstitution4KBuf(b *testing.B) {
    benchmarkSubstitution(b, 4 * 1024)
}
func BenchmarkSubstitution16KBuf(b *testing.B) {
    benchmarkSubstitution(b, 16 * 1024)
}
func BenchmarkSubstitution1MBuf(b *testing.B) {
    benchmarkSubstitution(b, 1 * 1024 * 1024)
}

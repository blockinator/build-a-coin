package types

import (
    "errors"
    "math/rand"
    "regexp"
    "strconv"
    "strings"
    "time"
)

type literalType struct{}
func (tt literalType) Produce(inputs ...string) (string, error) {
    if len(inputs) != 1 {
        return "", ErrWrongArity
    }
    return inputs[0], nil
}

func doUint(input string, size int) (string, error) {
    val, err := strconv.ParseUint(input, 0, size)
    if err != nil {
        return "", err
    }
    return strconv.FormatUint(val, 10), nil
}

func doInt(input string, size int) (string, error) {
    val, err := strconv.ParseInt(input, 0, size)
    if err != nil {
        return "", err
    }
    return strconv.FormatInt(val, 10), nil
}

var ErrWrongArity error = errors.New("wrong arity for type")

type byteType struct{}
func (tt byteType) Produce(inputs ...string) (string, error) {
    if len(inputs) != 1 {
        return "", ErrWrongArity
    }
    return doUint(inputs[0], 8)
}
type sevenBitType struct{}
func (tt sevenBitType) Produce(inputs ...string) (string, error) {
    if len(inputs) != 1 {
        return "", ErrWrongArity
    }
    return doUint(inputs[0], 7)
}
type uint16Type struct{}
func (tt uint16Type) Produce(inputs ...string) (string, error) {
    if len(inputs) != 1 {
        return "", ErrWrongArity
    }
    return doUint(inputs[0], 16)
}
type uint32Type struct{}
func (tt uint32Type) Produce(inputs ...string) (string, error) {
    if len(inputs) != 1 {
        return "", ErrWrongArity
    }
    return doUint(inputs[0], 32)
}
type uint64Type struct{}
func (tt uint64Type) Produce(inputs ...string) (string, error) {
    if len(inputs) != 1 {
        return "", ErrWrongArity
    }
    return doUint(inputs[0], 64)
}
type int32Type struct{}
func (tt int32Type) Produce(inputs ...string) (string, error) {
    if len(inputs) != 1 {
        return "", ErrWrongArity
    }
    return doInt(inputs[0], 32)
}
type int64Type struct{}
func (tt int64Type) Produce(inputs ...string) (string, error) {
    if len(inputs) != 1 {
        return "", ErrWrongArity
    }
    return doInt(inputs[0], 64)
}

type doubleType struct{}
func (tt doubleType) Produce(inputs ...string) (string, error) {
    if len(inputs) != 1 {
        return "", ErrWrongArity
    }
    val, err := strconv.ParseFloat(inputs[0], 64)
    if err != nil {
        return "", err
    }
    return strconv.FormatFloat(val, 'f', -1, 64), nil
}

var ErrIllegalChar error = errors.New("illegal characters in string")
var ErrStrTooLong error = errors.New("string too long")
const MaxStrLen = 256
type strType struct{}
func (tt strType) Produce(inputs ...string) (string, error) {
    if len(inputs) != 1 {
        return "", ErrWrongArity
    }
    if len(inputs[0]) > MaxStrLen {
        return "", ErrStrTooLong
    }
    for _, char := range []byte(inputs[0]) {
        if char < 0x20 || char > 0x7e || char == '"' || char == '\\' ||
                char == '\'' {
            return "", ErrIllegalChar
        }
    }
    return inputs[0], nil
}

var strAlphaRegex *regexp.Regexp = regexp.MustCompile("^[a-zA-Z]*$")
func checkStrAlpha(input string) error {
    if _, err := Str.Produce(input); err != nil {
        return err
    }
    if !strAlphaRegex.MatchString(input) {
        return ErrIllegalChar
    }
    return nil
}

type strAlphaType struct{}
func (tt strAlphaType) Produce(inputs ...string) (string, error) {
    if len(inputs) != 1 {
        return "", ErrWrongArity
    }
    if err := checkStrAlpha(inputs[0]); err != nil {
        return "", err
    }
    return inputs[0], nil
}
type strAlphaLowerType struct{}
func (tt strAlphaLowerType) Produce(inputs ...string) (string, error) {
    if len(inputs) != 1 {
        return "", ErrWrongArity
    }
    if err := checkStrAlpha(inputs[0]); err != nil {
        return "", err
    }
    return strings.ToLower(inputs[0]), nil
}
type strAlphaUpperType struct{}
func (tt strAlphaUpperType) Produce(inputs ...string) (string, error) {
    if len(inputs) != 1 {
        return "", ErrWrongArity
    }
    if err := checkStrAlpha(inputs[0]); err != nil {
        return "", err
    }
    return strings.ToUpper(inputs[0]), nil
}

type unixtimeCurrentType struct{}
func (tt unixtimeCurrentType) Produce(inputs ...string) (string, error) {
    now := time.Now().UTC()
    return strconv.FormatInt(now.Unix(), 10), nil
}

type randomUint32Type struct{}
func (tt randomUint32Type) Produce(inputs ...string) (string, error) {
    return strconv.FormatUint(uint64(rand.Uint32()), 10), nil
}

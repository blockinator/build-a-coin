package template

import (
    "bytes"
    "io"
    "strconv"
)

const (
    // Size of the filter buffer upon construction
    InitBufSize = 8 * 1024
    // Maximum number of digits representing index a substitution marker can
    // contain
    MaxMarkDigits = 8
    // Magic string indicating the beginning of a substitution marker
    LeadMarker = "__._"
    // Magic string indicating the end of a substitution marker
    TailMarker = "-"
)
var (
    leadMarker []byte = []byte(LeadMarker)
    tailMarker []byte = []byte(TailMarker)
    maxMarkLen uint = uint(len(leadMarker) + MaxMarkDigits + len(tailMarker))
)

// filter state machine info
const (
    // Filter state when skipping over non-mark data
    S_SEEK = iota
    // Filter state when attempting to match a lead marker
    S_MATCHING_LEAD
    // Filter state when reading in index number digits
    S_COLLECTING_DIGITS
    // Filter state when attempting to match a tail marker
    S_MATCHING_TAIL
)
type filterState struct {
    state uint
    inOverflow []byte
    outOverflow []byte
    matchOffset int
    digitBuf []byte
}

// Mapping of substitution marker indices to the bytes that will replace them
type FilterMap map[uint][]byte
// Streaming template processor for turning a base coin and values into a
// build-a-coin
type Filter struct {
    lastState filterState
    substitutions FilterMap
    input io.Reader
    buf []byte
    bytesProcessed uint64
}

// Construct a new filter that replaces markers in the input stream with values
// in the substitutions map.
func NewFilter(input io.Reader, substitutions FilterMap) *Filter {
    return &Filter{filterState{}, substitutions, input, make([]byte,
            InitBufSize), 0}
}
// Construct a new filter with an explicit buffer size for performance testing.
func newFilterDebug(input io.Reader, subs FilterMap,
        initBufSize uint) *Filter {
    return &Filter{filterState{}, subs, input, make([]byte, initBufSize), 0}
}

func (tt *Filter) Reset(input io.Reader) *Filter {
    tt.lastState = filterState{}
    tt.bytesProcessed = 0
    tt.input = input
    return tt
}

// Error when the filter is found to be in a state that it shouldn't
type BadStateError struct {
    stateName string
    reason string
    index uint64
}
func (e *BadStateError) Error() string {
    return "bad filter state @" + string(e.index) + " " + e.stateName + ": " +
            e.reason
}
// Error when an issue with the template stream prevents the filter from
// producing output
type TemplateError struct {
    reason string
    index uint64
}
func (e *TemplateError) Error() string {
    return "templating interrupted @" + strconv.FormatUint(e.index, 10) +
            ": " + e.reason
}

func (f *Filter) Read(out []byte) (int, error) {
    if len(out) < 1 {
        return 0, nil
    }
    // load filter state
    s := f.lastState

    // the input buffer needs to be at least as long as the requested output
    // buffer
    for len(f.buf) < len(out) {
        f.buf = make([]byte, (2 * len(f.buf)) + 1)
    }

    var in []byte
    var readErr error = nil
    inCount, outCount, inIdx, outIdx := 0, len(out), 0, 0

    // clean up on exit
    defer func() {
        f.bytesProcessed += uint64(inIdx)
        f.lastState = s
    }()

    // dump any output that overflowed from previous invocations
    if len(s.outOverflow) > 0 {
        outIdx = copy(out, s.outOverflow)
        if outIdx < len(s.outOverflow) {
            s.outOverflow = s.outOverflow[outIdx:]
        } else {
            s.outOverflow = make([]byte, 0)
        }
        if outIdx >= outCount {
            return outCount, readErr
        }
    }

    // filtering loop
    for {
        // locate next input
        if len(s.inOverflow) > 0 {
            in = s.inOverflow
            inCount = len(s.inOverflow)
            s.inOverflow = make([]byte, 0)
        } else {
            in = f.buf
            inCount, readErr = f.input.Read(in)
            if inCount < 1 {
                return outIdx, readErr
            }
        }

        // process each input byte
        for inIdx = 0; inIdx < inCount; inIdx++ {
            switch s.state {
            // the seek state copies bytes from in to out unless the beginning
            // of a lead marker is encountered
            case S_SEEK:
                if in[inIdx] == leadMarker[0] {
                    s.state = S_MATCHING_LEAD
                    s.matchOffset = 0
                    inIdx--
                    break
                }
                out[outIdx] = in[inIdx]
                outIdx++
                // if this seek write fills the output buffer, save the rest of
                // in as overflow for next time and return
                if outIdx >= outCount {
                    if inIdx < inCount {
                        s.inOverflow = in[inIdx+1:inCount]
                    }
                    return outCount, readErr
                }
            // the matching lead state attempts to match the entire lead marker
            // or dumps the partial match to out if not
            case S_MATCHING_LEAD:
                // check for overshoot.
                if s.matchOffset >= len(leadMarker) {
                    s.state = S_SEEK
                    return outIdx, &BadStateError{"MATCHING_LEAD",
                            "impossible match offset (too big)",
                            f.bytesProcessed + uint64(inIdx)}
                }
                // is there a break in the match?  Prepend the partial match to
                // input buffer.
                if in[inIdx] != leadMarker[s.matchOffset] {
                    s.state = S_SEEK
                    //inIdx -= s.matchOffset
                    out[outIdx] = leadMarker[0]
                    outIdx++
                    if len (leadMarker) > 1 {
                        in = bytes.Join([][]byte {
                            []byte(leadMarker[1:s.matchOffset]),
                            in[inIdx:inCount],
                            }, nil)
                        inIdx = 0
                        inCount = len(in)
                    }
                    inIdx--
                    s.matchOffset = -1
                    if outIdx >= outCount {
                        if inIdx < inCount {
                            s.inOverflow = in[inIdx+1:inCount]
                        }
                        return outCount, readErr
                    }
                } else {
                    s.matchOffset++
                    if s.matchOffset == len(leadMarker) {
                        s.state = S_COLLECTING_DIGITS
                    }
                }
            // the collecting digits state reads in UTF-8 decimal digits for
            // future conversion into an integer
            case S_COLLECTING_DIGITS:
                switch in[inIdx] {
                case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
                    if len(s.digitBuf) >= MaxMarkDigits {
                        s.state = S_SEEK
                        return outIdx, &TemplateError{
                                "too many digits in map specifier",
                                f.bytesProcessed + uint64(inIdx)}
                    }
                    s.digitBuf = append(s.digitBuf, in[inIdx])
                case tailMarker[0]:
                    s.state = S_MATCHING_TAIL
                    s.matchOffset = 0
                    inIdx--
                default:
                    s.state = S_SEEK
                    return outIdx, &TemplateError{"illegal character '" +
                                string(in[inIdx]) +
                                "' in template substitution specifier",
                            f.bytesProcessed + uint64(inIdx)}
                }
            case S_MATCHING_TAIL:
                // check for overshoot.
                if s.matchOffset >= len(tailMarker) {
                    s.state = S_SEEK
                    return outIdx, &BadStateError{"MATCHING_TAIL",
                            "impossible match offset (too big)",
                            f.bytesProcessed + uint64(inIdx)}
                }
                // is there a break in the match? this time it's fatal
                if in[inIdx] != tailMarker[s.matchOffset] {
                    s.state = S_SEEK
                    return outIdx, &TemplateError{"malformed tail marker",
                            f.bytesProcessed + uint64(inIdx)}
                }
                s.matchOffset++
                // once the whole tail marker is matched it's time to output
                // the substituted value
                if s.matchOffset >= len(tailMarker) {
                    s.state = S_SEEK
                    // condense the collected digit bytes into an unsigned
                    // integer
                    sub_int, err := strconv.Atoi(string(s.digitBuf))
                    sub_id := uint(sub_int)
                    s.digitBuf = s.digitBuf[:0]
                    if err != nil {
                        return outIdx, &TemplateError{
                                "failed to determine substitution specifier: " +
                                err.Error(),
                                f.bytesProcessed + uint64(inIdx)}
                    }
                    subValue, ok := f.substitutions[sub_id]
                    if !ok {
                        return outIdx, &TemplateError{
                                "unknown substitution specifier " +
                                    strconv.Itoa(int(sub_id)),
                                f.bytesProcessed + uint64(inIdx)}
                    }
                    outSlice := out[outIdx:outCount]
                    count := copy(outSlice, subValue)
                    if count < len(subValue) {
                        s.outOverflow = subValue[count:]
                    }
                    outIdx += count
                    if outIdx >= outCount {
                        if inIdx < inCount {
                            s.inOverflow = in[inIdx+1:inCount]
                        }
                        return outCount, readErr
                    }
                }
            }
        }
    }
    panic("template filter loop exited; this should never happen")
}

type FilterRunner struct {}
func (tt FilterRunner) Run(dst io.Writer, src io.Reader, values FilterMap) error {
    _, err := io.Copy(dst, NewFilter(src, values))
    return err
}

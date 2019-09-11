package bitcoin

import (
    "bytes"
    "encoding/binary"
    "math/big"
)

const (
    // Proof of work target compact representation of 1.0 difficulty
    Diff1Bits = 0x1d00ffff
)

// Compute a difficulty from a compact proof of work target.
func Difficulty(target uint32) float64 {
    numer := new(big.Rat).SetInt(TargetFull(Diff1Bits))
    denom := new(big.Rat).SetInt(TargetFull(target))

    diff, _ := new(big.Rat).Quo(numer, denom).Float64()
    return diff
}

// Compute a compact proof of work target from a difficulty.
func Target(difficulty float64) uint32 {
    baseTarget := new(big.Rat).SetInt(TargetFull(Diff1Bits))
    factor := new(big.Rat).SetFloat64(difficulty)

    newTarget := new(big.Rat).Quo(baseTarget, factor)

    intTarget := new(big.Int).Div(newTarget.Num(), newTarget.Denom())

    return TargetBits(intTarget)
}

// Convert a full-length proof of work target to compact form.
func TargetBits(full *big.Int) uint32 {
    fullBytes := full.Bytes()
    // If the leftmost bit in a target is 1, prepend a zero byte so it isn't
    if fullBytes[0] > 0x7F {
        fullBytes = bytes.Join([][]byte{ []byte{ 0x00 }, fullBytes }, nil)
    }
    // prepend the untruncated length of the target bytes
    fullBytes = bytes.Join([][]byte{ []byte { byte(len(fullBytes)) },
            fullBytes }, nil)

    // truncate the target bytes to a 32 bit unsigned int (binary.Read will
    // only consume 4 bytes when reading to a uint32)
    var bits uint32
    binary.Read(bytes.NewBuffer(fullBytes), binary.BigEndian, &bits)

    return bits
}

// Convert a compect proof of work target to full form.
func TargetFull(bits uint32) *big.Int {
    // break bits integer into its four constituent bytes
    byteBuf := new(bytes.Buffer)
    binary.Write(byteBuf, binary.BigEndian, bits)
    compactBytes := byteBuf.Bytes()

    // first byte is the length of the full target in bytes, so make a byte
    // slice of that length
    fullBytes := make([]byte, compactBytes[0])

    // copy the remaining compact bytes into the full slice and create a
    // big.Int from that
    fullBytes[0] = compactBytes[1]
    fullBytes[1] = compactBytes[2]
    fullBytes[2] = compactBytes[3]

    return new(big.Int).SetBytes(fullBytes)
}

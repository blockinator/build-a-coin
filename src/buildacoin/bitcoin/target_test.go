package bitcoin

import (
    "bytes"
    "encoding/hex"
    "math"
    "testing"
)

func TestCompactTarget(t *testing.T) {
    diff1FullBytes := []byte { 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
            0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, }

    actualFullTarget := TargetFull(Diff1Bits)

    if !bytes.Equal(diff1FullBytes, actualFullTarget.Bytes()) {
        t.Fatal("full target mismatch:\nexpected:\n" +
            hex.EncodeToString(diff1FullBytes) +
            "\nactual:\n" + hex.EncodeToString(actualFullTarget.Bytes()))
    }

    actualBits := TargetBits(actualFullTarget)

    if Diff1Bits != actualBits {
        t.Fatalf("target bits mismatch: expected: %08x / actual: %08x\n",
                Diff1Bits, actualBits)
    }
}

func TestDifficultyToBits(t *testing.T) {
    actualBits := Target(1.0)
    if actualBits != Diff1Bits {
        t.Fatalf("target bits mismatch: expected: %08x / actual: %08x\n",
                Diff1Bits, actualBits)
    }

    actualBits = Target(112628548.666347)
    if actualBits != 0x19262222 {
        t.Fatalf("target bits mismatch: expected: %08x / actual: %08x\n",
                0x19262222, actualBits)
    }
}

func TestBitsToDifficulty(t *testing.T) {
    actualDiff := Difficulty(Diff1Bits)
    if math.Abs(actualDiff - 1.0) > 0.001 {
        t.Fatalf("difficulty mismatch: expected: %f / actual: %f\n",
                1.0, actualDiff)
    }

    actualDiff = Difficulty(0x1972dbf2)
    if math.Abs(actualDiff - 37392766.136474) > 0.000001 {
        t.Fatalf("difficulty mismatch: expected: %f / actual: %f\n",
                37392766.136474, actualDiff)
    }
}

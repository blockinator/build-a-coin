package bitcoin

import (
    "bytes"
    "encoding/binary"
    "io"
)

const (
    // Number of value units (satoshis) in one coin
    Coin  = 100000000
    // Number of value units (satoshis) in one millicoin
    Milli = 100000
    // Number of value units (satoshis) in one microcoin
    Micro = 100
    // Geometric sum multiplier for total coins
    GeoFactor = 2.0
    // Precision
    Precision = 8
)

// Write out the bitcoin variable length int serialization of an integer.
func WriteVarint(out io.Writer, value uint64) (int, error) {
    // 0xfd, 0xfe, 0xff are magic values; all lower values are serialized
    // literally
    if value < 0xfd {
        return out.Write([]byte{ byte(value) })
    } else if value <= 0xffff {
        // values less than 2^16 are represented by the 0xfd magic number
        // followed by a 2 byte serialization of the value
        n, err := out.Write([]byte{ 0xfd })
        if err != nil {
            return n, err
        }
        err = binary.Write(out, binary.LittleEndian, uint16(value))
        if err != nil {
            return n, err
        }
        return n+2, nil
    } else if value <= 0xffffffff {
        // values less than 2^32 are represented by the 0xfe magic number
        // followed by a 4 byte serialization of the value
        n, err := out.Write([]byte{ 0xfe })
        if err != nil {
            return n, err
        }
        err = binary.Write(out, binary.LittleEndian, uint32(value))
        if err != nil {
            return n, err
        }
        return n+4, nil
    } else {
        // all other values (max 2^64-1 enforced by Go's type system) are
        // represented by the 0xff magic number followed by a 8 byte
        // serialization of the value
        n, err := out.Write([]byte{ 0xff })
        if err != nil {
            return n, err
        }
        err = binary.Write(out, binary.LittleEndian, value)
        if err != nil {
            return n, err
        }
        return n+8, nil
    }
    panic("impossible WriteVarint() flow")
}

// Get the bitcoin variable length int serialization of an integer.
func Varint(value uint64) []byte {
    buf := new(bytes.Buffer)
    _, err := WriteVarint(buf, value)
    if err != nil{
        return nil
    }
    return buf.Bytes()
}

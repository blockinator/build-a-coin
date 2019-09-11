package bitcoin

import (
    "bytes"
    "encoding/binary"
    "io"
)

const (
    // Bitcoin network transaction version
    TxVersion = 1
    // Transaction sequence number indicating a finalized transaction (other
    // values result in a nonstandard tx)
    FinalSequence = 0xffffffff
    // Transaction unlock time/block value for an unlocked transaction (other
    // values result in a nonstandard tx)
    UnlockedTime = 0x00000000
    // The length of an uncompressed public key in bytes
    PubkeyLen = 65
    // The length of a compressed public key in bytes
    CompPubkeyLen = 33
)

// A bitcoin transaction
type Tx struct {
    inputs [][]byte
    outputs [][]byte
}

// Write out the serialization of a transaction.
func (t *Tx) WriteBytes(out io.Writer) (int, error) {
    outCount := 0

    // version bytes
    err := binary.Write(out, binary.LittleEndian, uint32(TxVersion))
    if err != nil {
        return outCount, err
    }
    outCount += 4
    // inputs and outputs
    for _, category := range [][][]byte{ t.inputs, t.outputs } {
        // varint of the number of inputs or outputs
        n, err := WriteVarint(out, uint64(len(category)))
        outCount += n
        if err != nil {
            return outCount, err
        }
        // each input or output itself
        for _, inoutput := range category {
            n, err := out.Write(inoutput)
            outCount += n
            if err != nil {
                return outCount, err
            }
        }
    }
    // locktime
    err = binary.Write(out, binary.LittleEndian, uint32(UnlockedTime))
    if err != nil {
        return outCount, err
    }
    outCount += 4

    return outCount, nil
}

// Get the serialization of a transaction.
func (t *Tx) Bytes() []byte {
    buf := new(bytes.Buffer)
    _, err := t.WriteBytes(buf)
    if err != nil{
        return nil
    }
    return buf.Bytes()
}

// Construct a new input and append it to the transaction.
func (t *Tx) Input(srcTx Hash, outputIdx uint, scriptSig []byte) *Tx {
    buf := new(bytes.Buffer)
    buf.Write(srcTx.Bytes())
    binary.Write(buf, binary.LittleEndian, uint32(outputIdx))
    WriteVarint(buf, uint64(len(scriptSig)))
    buf.Write(scriptSig)
    binary.Write(buf, binary.LittleEndian, uint32(FinalSequence))
    t.inputs = append(t.inputs, buf.Bytes())
    return t
}

// Construct a new output and append it to the transaction.
func (t *Tx) Output(value uint64, scriptPubKey []byte) *Tx {
    buf := new(bytes.Buffer)
    binary.Write(buf, binary.LittleEndian, value)
    WriteVarint(buf, uint64(len(scriptPubKey)))
    buf.Write(scriptPubKey)
    t.outputs = append(t.outputs, buf.Bytes())
    return t
}

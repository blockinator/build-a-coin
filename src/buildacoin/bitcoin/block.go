package bitcoin

import (
    "bytes"
    "encoding/binary"
    "io"
    "time"
)

// A bitcoin block
type Block struct {
    version uint32
    prevBlock Hash
    timestamp time.Time
    targetBits uint32
    nonce uint32
    merkleRoot *Hash
    txs []*Tx
}

// Create a new block.
func NewBlock(version, targetBits, nonce uint32, prevBlock Hash,
        timestamp time.Time) *Block {
    return &Block{
        version: version,
        prevBlock: prevBlock,
        timestamp: timestamp,
        targetBits: targetBits,
        nonce: nonce,
        txs: make([]*Tx, 0, 1),
    }
}

// Add a transaction to the end of a block.
func (tt *Block) AddTx(tx *Tx) *Block {
    tt.txs = append(tt.txs, tx)
    return tt
}

func (tt *Block) SetMerkleRoot(root Hash) *Block {
    tt.merkleRoot = &root
    return tt
}

// Compute the merkle root from the transactions in a block.
func (tt *Block) MerkleRoot() Hash {
    if len(tt.txs) < 1 {
        return Hash{}
    }
    hashes := make([]Hash, 0, 8)
    for _, tx := range tt.txs {
        hashes = append(hashes, Sha256d(tx.Bytes()))
    }

    tree := MerkleTree(hashes, HashSha256d)

    return tree[len(tree)-1]
}

// Write out the serialization of a block's header.
func (tt *Block) WriteHeader(out io.Writer) (int, error) {
    outCount := 0

    // version bytes
    err := binary.Write(out, binary.LittleEndian, uint32(tt.version))
    if err != nil {
        return outCount, err
    }
    outCount += 4
    // previous block
    n, err := out.Write(tt.prevBlock.Bytes())
    outCount += n
    if err != nil {
        return outCount, err
    }
    // merkle root: from cache or computed
    var merkleRoot Hash
    if tt.merkleRoot != nil {
        merkleRoot = *tt.merkleRoot
    } else {
        merkleRoot = tt.MerkleRoot()
    }
    n, err = out.Write(merkleRoot.Bytes())
    outCount += n
    if err != nil {
        return outCount, err
    }
    // timestamp
    err = binary.Write(out, binary.LittleEndian, uint32(tt.timestamp.Unix()))
    if err != nil {
        return outCount, err
    }
    outCount += 4
    // target bits
    err = binary.Write(out, binary.LittleEndian, uint32(tt.targetBits))
    if err != nil {
        return outCount, err
    }
    outCount += 4
    // nonce
    err = binary.Write(out, binary.LittleEndian, uint32(tt.nonce))
    if err != nil {
        return outCount, err
    }
    outCount += 4

    return outCount, nil
}

// Get the serialization of a block's header.
func (tt *Block) Header() []byte {
    buf := new(bytes.Buffer)
    _, err := tt.WriteHeader(buf)
    if err != nil{
        return nil
    }
    return buf.Bytes()
}

// Write out the serialization of an entire block.
func (tt *Block) WriteBytes(out io.Writer) (int, error) {
    outCount := 0
    //header
    n, err := tt.WriteHeader(out)
    outCount += n
    if err != nil {
        return outCount, err
    }
    // tx count
    n, err = WriteVarint(out, uint64(len(tt.txs)))
    outCount += n
    if err != nil {
        return outCount, err
    }
    // txs
    for _, tx := range tt.txs {
        n, err = tx.WriteBytes(out)
        outCount += n
        if err != nil {
            return outCount, err
        }
    }

    return outCount, nil
}

// Get the serialization of an entire block.
func (tt *Block) Bytes() []byte {
    buf := new(bytes.Buffer)
    _, err := tt.WriteBytes(buf)
    if err != nil{
        return nil
    }
    return buf.Bytes()
}

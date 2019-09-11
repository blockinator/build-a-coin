package altcoins

import (
    "buildacoin/bitcoin"
    "bytes"
    "encoding/binary"
    "time"
)

const (
    // Current bitcoin protocol block version number
    BlockVersion = 2
    // Default value for block difficulty in compact format
    DefaultTargetBits = 0x1e0ffff0
    // Default block header nonce
    DefaultNonce = 0
    // Longest permissable length for genesis coinbase message strings and
    // output pubkeys
    MaxDataLen = 75
    // Script opcode for the checksig operation as used in genesis coinbase tx
    // outputs
    OpChecksig = 0xac
)

// Create a new genesis block.
func Genesis(genValue uint64, timestamp time.Time, coinbaseMsg string,
        pubKey []byte, difficulty float64) *bitcoin.Block {
    block := bitcoin.NewBlock(1, bitcoin.Target(difficulty),
            DefaultNonce, bitcoin.Hash{}, timestamp)
    block.AddTx(GenesisTx(genValue, coinbaseMsg, pubKey))
    return block
}

// Create a genesis coinbase transaction.
func GenesisTx(genValue uint64, coinbaseMsg string, pubKey []byte) *bitcoin.Tx {
    return new(bitcoin.Tx,
            ).Input(bitcoin.Hash{}, 4294967295, GenesisCoinbase(coinbaseMsg),
            ).Output(genValue, CoinbaseTxScriptPubKey(pubKey),
            )
}

// Create the genesis coinbase itself (the coinbase itself is only the input
// for the coinbase transaction).
func GenesisCoinbase(coinbaseMsg string) []byte {
    messageBytes := []byte(coinbaseMsg)
    if len(messageBytes) > MaxDataLen {
        panic("coinbase message too long")
    }
    buf := new(bytes.Buffer)
    // first stack item
    binary.Write(buf, binary.LittleEndian, uint8(4))
    binary.Write(buf, binary.LittleEndian, uint32(486604799))
    // second stack item
    binary.Write(buf, binary.LittleEndian, uint8(1))
    binary.Write(buf, binary.LittleEndian, uint8(4))
    // secret message
    binary.Write(buf, binary.LittleEndian, uint8(len(messageBytes)))
    binary.Write(buf, nil, messageBytes)

    return buf.Bytes()
}

// Create the coinbase transaction output script.
func CoinbaseTxScriptPubKey(pubKey []byte) []byte {
    if len(pubKey) > MaxDataLen {
        panic("coinbase pubkey too long")
    }
    buf := new(bytes.Buffer)
    // pubkey to stack
    binary.Write(buf, binary.LittleEndian, uint8(len(pubKey)))
    binary.Write(buf, nil, pubKey)
    // checksig
    binary.Write(buf, binary.LittleEndian, uint8(OpChecksig))

    return buf.Bytes()
}

// Shortcut for the hash of a genesis block if the block itself is not desired.
func GenesisHash(genValue uint64, timestamp time.Time, coinbaseMsg string,
        pubKey []byte, difficulty float64) bitcoin.Hash {
    block := Genesis(genValue, timestamp, coinbaseMsg, pubKey, difficulty)
    return bitcoin.Sha256d(block.Header())
}

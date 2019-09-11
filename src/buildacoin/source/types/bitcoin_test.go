package types

import (
    "buildacoin/bitcoin"
    "encoding/hex"
    "testing"
    "time"
)

func TestGenesisHash(t *testing.T) {
    expected := "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a" +
            "8ce26f"

    timestamp := time.Unix(1231006505, 0)

    merkleString := "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b" +
            "7afdeda33b"
    merkleBytes, err := hex.DecodeString(merkleString)
    if err != nil {
        t.Fatal(err.Error())
    }
    merkle, err := bitcoin.HashFromBytes(merkleBytes, bitcoin.BigEndian)

    block := bitcoin.NewBlock(1, 0x1d00ffff, 2083236893,
            bitcoin.Hash{}, timestamp).SetMerkleRoot(merkle)

    headerHash := bitcoin.Sha256d(block.Header())
    flippedHash, err := bitcoin.HashFromBytes(headerHash.Bytes(),
            bitcoin.FlipEndian)
    if err != nil {
        t.Fatal(err.Error())
    }
    actual := hex.EncodeToString(flippedHash.Bytes())

    if expected != actual {
        t.Fatalf("block hash mismatch:\nexpected\n%s\nactual\n%s\n", expected,
                actual)
    }
}

package altcoins

import (
    "buildacoin/bitcoin"
    "encoding/hex"
    "testing"
    "time"
)

func TestLitecoinGenesis(t *testing.T) {
    expected, err := bitcoin.HashFromHex("12a765e31ffd4059bada1e25190f6e98c" +
            "99d9714d334efa41a195a7e7e04bfe2")
    if err != nil {
        t.Fatal(err.Error())
    }
    pubKey, err := hex.DecodeString("040184710fa689ad5023690c80f3a49c8f13f8" +
            "d45b8c857fbcbc8bc4a8e4d3eb4b10f4d4604fa08dce601aaf0f470216fe1b" +
            "51850b4acf21b179c45070ac7b03a9")
    if err != nil {
        t.Fatal(err.Error())
    }

    gBlock := bitcoin.NewBlock(1, 0x1e0ffff0, 2084524493, bitcoin.Hash{},
            time.Unix(1317972665, 0))
    gTx := GenesisTx(50 * bitcoin.Coin, "NY Times 05/Oct/2011 Steve Jobs, " +
            "Apple’s Visionary, Dies at 56", pubKey)
    gBlock.AddTx(gTx)

    actual := bitcoin.Sha256d(gBlock.Header())

    if !actual.Equals(expected) {
        t.Fatal("hash mismatch: expected " + expected.String() + " / actual " +
                actual.String())
    }
}

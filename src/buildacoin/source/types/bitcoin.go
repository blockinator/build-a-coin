package types

import (
    "buildacoin/altcoins"
    "buildacoin/bitcoin"
    "crypto/rand"
    "encoding/hex"
    "errors"
    "math/big"
    "strconv"
    "strings"
    "time"
)

var (
    ErrOverflow error = errors.New("value overflows its type")
    ErrBadFloat error = errors.New("bad float value")
)

const maxCoins = (1 << 63 - 1) / bitcoin.Coin

type coinsType struct {}
func (tt coinsType) Produce(inputs ...string) (string, error) {
    if len(inputs) != 1 {
        return "", ErrWrongArity
    }
    input := inputs[0]
    val, err := strconv.ParseFloat(input, 64)
    if err != nil {
        return "", err
    }
    if val < 0.0 || val > float64(maxCoins) {
        return "", ErrOverflow
    }
    satoshis := int64(float64(bitcoin.Coin) * val)
    return strconv.FormatInt(satoshis, 10), nil
}

type difficultyType struct {}
func (tt difficultyType) Produce(inputs ...string) (string, error) {
    if len(inputs) != 1 {
        return "", ErrWrongArity
    }
    input := inputs[0]
    val, err := strconv.ParseFloat(input, 64)
    if err != nil {
        return "", err
    }
    target := bitcoin.Target(val)
    return strconv.FormatUint(uint64(target), 10), nil
}

var ErrBadPubkeyLen error = errors.New("bad length for public key")
type pubkeyType struct {}
func (tt pubkeyType) Produce(inputs ...string) (string, error) {
    if len(inputs) != 1 {
        return "", ErrWrongArity
    }
    input := inputs[0]
    if strings.HasPrefix(input, "0x") {
        input = input[2:]
    }
    val, err := hex.DecodeString(input)
    if err != nil {
        return "", err
    }
    if len(val) != bitcoin.PubkeyLen && len(val) != bitcoin.CompPubkeyLen {
        return "", ErrBadPubkeyLen
    }
    return hex.EncodeToString(val), nil
}

var ErrNoEntropy error = errors.New("insufficient entropy to create " +
        "random variable")

func doRandomBytes(length int) (string, error) {
    bytes := make([]byte, length)
    n, err := rand.Read(bytes)
    if n < len(bytes) || err != nil {
        return "", ErrNoEntropy
    }
    return hex.EncodeToString(bytes), nil
}

type randomPubkeyType struct {}
func (tt randomPubkeyType) Produce(inputs ...string) (string, error) {
    return doRandomBytes(bitcoin.PubkeyLen)
}

type randomHashType struct{}
func (tt randomHashType) Produce(inputs ...string) (string, error) {
    return doRandomBytes(bitcoin.HashSize)
}

type genesisMerkleRootType struct{}
func (tt genesisMerkleRootType) Produce(inputs ...string) (string, error) {
    if len(inputs) != 4 {
        return "", ErrWrongArity
    }
    coinbaseStr := inputs[2]
    pubkeyBytes, err := hex.DecodeString(inputs[3])
    if err != nil {
        return "", err
    }
    coinbaseValue, err := strconv.ParseUint(inputs[1], 10, 64)
    if err != nil {
        return "", err
    }

    genesis := altcoins.Genesis(coinbaseValue, time.Time{}, coinbaseStr,
            pubkeyBytes, 1.0)

    merkleBytes := genesis.MerkleRoot().Bytes()
    // flip endianness, because the result is a hash string, where bitcoin uses
    // big endian
    merkle, err := bitcoin.HashFromBytes(merkleBytes, bitcoin.FlipEndian)
    if err != nil {
        panic("bytes from a hash aren't a valid hash!?")
    }

    return hex.EncodeToString(merkle.Bytes()), nil
}

type genesisBlockHashType struct{}
func (tt genesisBlockHashType) Produce(inputs ...string) (string, error) {
    if len(inputs) != 5 {
        return "", ErrWrongArity
    }

    timestampUnix, err := strconv.ParseUint(inputs[1], 10, 32)
    if err != nil {
        return "", err
    }
    timestamp := time.Unix(int64(timestampUnix), 0)

    diffBits64, err := strconv.ParseUint(inputs[2], 10, 32)
    if err != nil {
        return "", err
    }
    diffBits := uint32(diffBits64)

    nonce64, err := strconv.ParseUint(inputs[3], 10, 32)
    if err != nil {
        return "", err
    }
    nonce := uint32(nonce64)

    merkleBytes, err := hex.DecodeString(inputs[4])
    if err != nil {
        return "", err
    }
    merkle, err := bitcoin.HashFromBytes(merkleBytes, bitcoin.BigEndian)
    if err != nil {
        return "", err
    }

    // the block itself

    block := bitcoin.NewBlock(1, diffBits, nonce,
            bitcoin.Hash{}, timestamp).SetMerkleRoot(merkle)

    headerHash := bitcoin.Sha256d(block.Header())
    flippedHash, err := bitcoin.HashFromBytes(headerHash.Bytes(), bitcoin.FlipEndian)
    if err != nil {
        panic("bytes from a hash aren't a valid hash!?")
    }

    return hex.EncodeToString(flippedHash.Bytes()), nil
}

type coinsMaxType struct{}
func (tt coinsMaxType) Produce(inputs ...string) (string, error) {
    if len(inputs) != 3 {
        return "", ErrWrongArity
    }

    initReward, err := strconv.ParseInt(inputs[1], 10, 64)
    if err != nil {
        return "", err
    }

    halvingBlocks, err := strconv.ParseInt(inputs[2], 10, 32)
    if err != nil {
        return "", err
    }

    totalCoins := int64(initReward * halvingBlocks * bitcoin.GeoFactor)
    return strconv.FormatInt(totalCoins, 10), nil
}

type doubleCoinsType struct{}
func (tt doubleCoinsType) Produce(inputs ...string) (string, error) {
    if len(inputs) != 2 {
        return "", ErrWrongArity
    }

    floatoshis, ok := new(big.Rat).SetString(inputs[1])
    if !ok {
        return "", ErrBadFloat
    }
    perCoin := new(big.Rat).SetInt64(bitcoin.Coin)

    result := new(big.Rat).Quo(floatoshis, perCoin)

    return result.FloatString(bitcoin.Precision), nil
}

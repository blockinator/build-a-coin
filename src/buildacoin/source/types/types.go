package types

// A type constrains template inputs to legal values for a field
type Type interface {
    // Take a user input string and produce a conforming source code string or
    // an error.
    Produce(input ...string) (string, error)
}

var (
    // accepts all inputs and does not modify them.  The empty string is an
    // alias for this type.
    Literal literalType
    Byte byteType
    // Hex values less than 128 (0x80)
    SevenBit sevenBitType
    Uint16 uint16Type
    Uint32 uint32Type
    Uint64 uint64Type
    Int32 int32Type
    Int64 int64Type
    Double doubleType
    // accepts all strings not containing the characters ' " and \
    Str strType
    // accepts all strings consisting entirely of unaccented latin alphabet
    // characters
    StrAlpha strAlphaType
    // accepts strAlpha inputs and produces all lowercase strings
    StrAlphaLower strAlphaLowerType
    // accepts strAlpha inputs and produces all uppercase strings
    StrAlphaUpper strAlphaLowerType
    // produces a random 32-bit unsigned integer
    RandomUint32 randomUint32Type
    // accepts floating representations of coin values and returns integer
    // representations of value units (satoshis)
    Coins coinsType
    // takes a floating difficulty value and returns the 32 bit compressed
    // target value
    Difficulty difficultyType
    // accepts hex-encoded compressed and uncompressed public keys
    Pubkey pubkeyType
    // produces a random hex-encoded uncompressed public key
    RandomPubkey randomPubkeyType
    // produces a random hex-encoded hash value
    RandomHash randomHashType
    // produces a unix timestamp of the current time
    UnixtimeCurrent unixtimeCurrentType
    // computes the merkle root of a genesis block from the genesis message and
    // the coinbase tx output pubkey
    GenesisMerkleRoot genesisMerkleRootType
    // computes the hash of a genesis block from the genesis timestamp nonce
    // difficulty bits and merkle hash
    GenesisHash genesisBlockHashType
    // the most coins that will ever exist based on initial reward and halving
    // time
    CoinsMax coinsMaxType
    // conversion to double from satoshis
    DoubleCoins doubleCoinsType
    // Mapping from the string name of a type to the type object itself
    Map = map[string]Type {
        "literal": Literal,
        "": Literal,
        "byte": Byte,
        "7bit": SevenBit,
        "uint16": Uint16,
        "uint32": Uint32,
        "uint64": Uint64,
        "int32": Int32,
        "int64": Int64,
        "double": Double,
        "str": Str,
        "str-alpha": StrAlpha,
        "str-alpha-lower": StrAlphaLower,
        "str-alpha-upper": StrAlphaUpper,
        "random-uint32": RandomUint32,
        "coins": Coins,
        "difficulty": Difficulty,
        "pubkey": Pubkey,
        "random-pubkey": RandomPubkey,
        "random-hash": RandomHash,
        "unixtime-current": UnixtimeCurrent,
        "genesis-merkle-root": GenesisMerkleRoot,
        "genesis-block-hash": GenesisHash,
        "coins-max": CoinsMax,
        "double-coins": DoubleCoins,
    }
)

package bitcoin

import (
    "bytes"
    "code.google.com/p/go.crypto/scrypt"
    "crypto/sha256"
    "encoding/hex"
    "errors"
)

const (
    // Length of bitcoin hash values in bytes
    HashSize = 32
    // Descriptive endianness constant for little endian operations
    LittleEndian = false
    // Descriptive endianness constant for big endian operations
    BigEndian = true
    // Descriptive endianness constant for flipping byte order
    FlipEndian = true
)
var (
    // Error when a hash is constructed (not calculated) from the wrong number
    // of bytes
    HashLengthError = errors.New("incorrect length for hash")
)

// A bitcoin hash value as used in numerous places in the protocol
type Hash [HashSize]byte

// Construct (not compute) a hash value from a sequence of bytes of a hash.
func HashFromBytes(input []byte, bigEndian bool) (Hash, error) {
    var output Hash
    if len(input) != HashSize {
        return output, HashLengthError
    }
    for ii := 0; ii < len(output); ii++ {
        outputIdx := ii
        if bigEndian {
            outputIdx = len(output)-1-ii
        }
        output[outputIdx] = input[ii]
    }
    return output, nil
}
// Construct (not compute) a hash value from a hex string of a hash.
func HashFromHex(input string) (Hash, error) {
    if len(input) < 2 {
        return Hash{}, HashLengthError
    }
    if input[0] == '0' && input[1] == 'x' {
        input = input[2:]
    }
    if hex.DecodedLen(len(input)) != HashSize {
        return Hash{}, HashLengthError
    }
    hashBytes, err := hex.DecodeString(input)
    if err != nil {
        return Hash{}, err
    }
    return HashFromBytes(hashBytes, BigEndian)
}

// Get the byte sequence that represents this hash.
func (tt Hash) Bytes() []byte {
    return tt[:]
}

// Get the hex string that represents this hash.
func (tt Hash) String() string {
    buf := make([]byte, len(tt))
    for ii := 0; ii < len(tt); ii++ {
        buf[len(tt)-1-ii] = tt[ii]
    }
    return hex.EncodeToString(buf)
}

func (tt Hash) Equals(other Hash) bool {
    for idx, val := range tt {
        if other[idx] != val {
            return false
        }
    }
    return true
}

// Hashers compute (not construct) hashes from arbitrary byte sequences.
type Hasher interface {
    Hash(input []byte) Hash
}

// Compute a Scrypt hash with Litecoin proof of work constants.
func Scrypt(input []byte) Hash {
    hashBytes, err := scrypt.Key(input, input, 1024, 1, 1, 32)
    if err != nil {
        panic("impossible flow, this is a bug: " + err.Error())
    }
    hash, err := HashFromBytes(hashBytes, LittleEndian)
    if err != nil {
        panic("impossible flow, this is a bug: " + err.Error())
    }
    return hash
}

type scryptHasher struct{}
func (tt scryptHasher) Hash(input []byte) Hash {
    return Scrypt(input)
}
// Hasher for Litecoin proof of work Scrypt
var HashScrypt scryptHasher

// Compute a double SHA-256 hash.
func Sha256d(input []byte) Hash {
    sha := sha256.New()
    sha.Write(input)
    intermediate := sha.Sum(nil)
    sha.Reset()
    sha.Write(intermediate)
    hash, err := HashFromBytes(sha.Sum(nil), LittleEndian)
    if err != nil {
        panic("impossible flow, this is a bug: " + err.Error())
    }
    return hash
}

type sha256dHasher struct{}
func (tt sha256dHasher) Hash(input []byte) Hash {
    return Sha256d(input)
}
// Hasher for double SHA-256
var HashSha256d sha256dHasher

// Produce a merkle tree from a list of hashes.
func MerkleTree(inputHashes []Hash, hasher Hasher) []Hash {
    // don't mutate the input
    hashes := make([]Hash, len(inputHashes), len(inputHashes) * 2)
    copy(hashes, inputHashes)
    // each row of the tree is half the length of the previous row until
    // reaching the root row of one hash
    for ii := len(hashes); ii > 1; ii /= 2 {
        // if a row has an odd number of hashes, duplicate the last hash for a
        // clean halving
        if ii % 2 != 0 {
            hashes = append(hashes, hashes[len(hashes)-1])
            ii++
        }
        // for each pair of hashes in the current row, compute the hash of the
        // concatenation of the two hashes and append the resulting hash to the
        // new tree row
        newRow := make([]Hash, 0, 8)
        for jj := ii; jj > 0; jj -= 2 {
            hashA := hashes[len(hashes)-jj]
            hashB := hashes[len(hashes)-(jj-1)]
            hashC := hasher.Hash(bytes.Join(
                    [][]byte{ hashA.Bytes(), hashB.Bytes() }, nil))
            newRow = append(newRow, hashC)
        }
        // add the new tree row to the tree
        hashes = append(hashes, newRow...)
    }

    return hashes
}

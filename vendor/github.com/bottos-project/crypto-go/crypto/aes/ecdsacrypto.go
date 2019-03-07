package aes

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "errors"
    "fmt"
    "io"
    "math/big"
    "github.com/bottos-project/crypto-go/crypto/secp256k1"
)

var (
	secp256k1_N, _  = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1_halfN = new(big.Int).Div(secp256k1_N, big.NewInt(2))

    Big1   = big.NewInt(1)
    Big2   = big.NewInt(2)
    Big3   = big.NewInt(3)
    Big0   = big.NewInt(0)
    Big32  = big.NewInt(32)
    Big256 = big.NewInt(256)
    Big257 = big.NewInt(257)
)

const (
    HashLength    = 32
    UUIDLength = 20
)

type UUID [UUIDLength]byte

var Reader io.Reader = &randEntropy{}

type randEntropy struct {
}

func BytesToUUID(b []byte) UUID {
    var a UUID
    a.SetBytes(b)
    return a
}

func (a UUID) Bytes() []byte { return a[:] }

func (a *UUID) SetBytes(b []byte) {
    if len(b) > len(a) {
            b = b[len(b)-UUIDLength:]
    }
    copy(a[UUIDLength-len(b):], b)
}


func (*randEntropy) Read(bytes []byte) (n int, err error) {
	readBytes := GetEntropyCSPRNG(len(bytes))
	copy(bytes, readBytes)
	return len(bytes), nil
}

func GetEntropyCSPRNG(n int) []byte {
	mainBuff := make([]byte, n)
	_, err := io.ReadFull(rand.Reader, mainBuff)
	if err != nil {
		panic("reading from crypto/rand failed: " + err.Error())
	}
	return mainBuff
}

func S256() elliptic.Curve {
    return secp256k1.S256()
}

// Keccak256 calculates and returns the Keccak256 hash of the input data.
func Keccak256(data ...[]byte) []byte {
	d := /*sha3.*/NewKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

// ToECDSACRYPTO creates a private key with the given D value.
func ToECDSACRYPTO(d []byte) (*ecdsa.PrivateKey, error) {
	return toECDSACRYPTO(d, true)
}

// ToECDSACRYPTOUnsafe blindly converts a binary blob to a private key. It should almost
// never be used unless you are sure the input is valid and want to avoid hitting
// errors due to bad origin encoding (0 prefixes cut off).
func ToECDSACRYPTOUnsafe(d []byte) *ecdsa.PrivateKey {
	priv, _ := toECDSACRYPTO(d, false)
	return priv
}

// toECDSACRYPTO creates a private key with the given D value. The strict parameter
// controls whether the key's length should be enforced at the curve size or
// it can also accept legacy encodings (0 prefixes).
func toECDSACRYPTO(d []byte, strict bool) (*ecdsa.PrivateKey, error) {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = S256()
	if strict && 8*len(d) != priv.Params().BitSize {
		return nil, fmt.Errorf("invalid length, need %d bits", priv.Params().BitSize)
	}
	priv.D = new(big.Int).SetBytes(d)

	// The priv.D must < N
	if priv.D.Cmp(secp256k1_N) >= 0 {
		return nil, fmt.Errorf("invalid private key, >=N")
	}
	// The priv.D must not be zero or negative.
	if priv.D.Sign() <= 0 {
		return nil, fmt.Errorf("invalid private key, zero or negative")
	}

	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(d)
	if priv.PublicKey.X == nil {
		return nil, errors.New("invalid private key")
	}
	return priv, nil
}

// FromECDSACRYPTO exports a private key into a binary dump.
func FromECDSACRYPTO(priv *ecdsa.PrivateKey) []byte {
	if priv == nil {
		return nil
	}
	return PaddedBigBytes(priv.D, priv.Params().BitSize/8)
}

// FromECDSACRYPTOPub exports a private key into a binary dump.
func FromPriECDSACRYPTOPub(priv *ecdsa.PrivateKey) []byte {
	if priv == nil {
		return nil
	}
	return elliptic.Marshal(secp256k1.S256(), priv.X, priv.Y)
}

func ToECDSACRYPTOPub(pub []byte) *ecdsa.PublicKey {
	if len(pub) == 0 {
		return nil
	}
	x, y := elliptic.Unmarshal(S256(), pub)
	return &ecdsa.PublicKey{Curve: S256(), X: x, Y: y}
}

func FromECDSACRYPTOPub(pub *ecdsa.PublicKey) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	return elliptic.Marshal(S256(), pub.X, pub.Y)
}

func PubkeyToUUID(p ecdsa.PublicKey) UUID {

	pubBytes := FromECDSACRYPTOPub(&p)
	return BytesToUUID(Keccak256(pubBytes[1:])[12:])
}

func zeroBytes(bytes []byte) {
	for i := range bytes {
		bytes[i] = 0
	}
}

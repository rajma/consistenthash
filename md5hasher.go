package consistenthash

import (
	"crypto/md5"
	"math/big"
)

//Hasher hashes the key and returns the hashed integer value
type Hasher interface {
	GetHash(identifer string) *big.Int
}

// MD5Hasher provides the MD5 hash implementation of a
// hashing a string
type MD5Hasher struct {
	//Possible to re-use the hash implementation since it has
	// a Reset function. Will have to worry about locking so
	// essentially it will become a single threaded operation
}

// GetHash returns a hashed value
func (mdh *MD5Hasher) GetHash(key string) *big.Int {
	h := md5.New()
	h.Write([]byte(key))
	hv := h.Sum(nil)
	val := big.NewInt(0)
	val = val.SetBytes(hv)
	return val
}

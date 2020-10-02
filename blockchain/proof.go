package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
)

// Take the data from the block
// create a counter (nonce) which starts at 0
// create a hash of the data plus the counter
// check the hash to see if it meets a set of requirements
// Requirements:
// The First few bytes must contain 0s
//
// choose Sum256 means using [32]byte
// to store the checksum result is perfect
// 32 * 8 = 256

// must >= 1
const Difficulty = 16

type ProofOfWork struct {
	Block  *Block
	Target *[32]byte
}

func NewProofOfWork(b *Block) *ProofOfWork {
	return &ProofOfWork{b, New32BytesArrWithRsh(Difficulty)}
}

func (pow *ProofOfWork) ToBytesWithNonce(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToBytes(int64(nonce)),
			ToBytes(int64(Difficulty)),
		},
		[]byte{},
	)
}

func (pow *ProofOfWork) Run() (int, []byte) {
	nonce := 0
	var checksum [32]byte

	for nonce < math.MaxInt64 {
		// give it a try with nonce
		data := pow.ToBytesWithNonce(nonce)
		checksum = sha256.Sum256(data)

		// todo: flush
		fmt.Printf("\t%x\n", checksum)

		if Cmp32BytesArr(&checksum, pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}

	return nonce, checksum[:]
}

func (pow *ProofOfWork) Validate() bool {
	data := pow.ToBytesWithNonce(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	return Cmp32BytesArr(&hash, pow.Target) == -1
}

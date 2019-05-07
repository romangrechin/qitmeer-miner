/**
	HLC FOUNDATION
	james
 */
package hlc

import (
	"math"
	"github.com/noxproject/nox/common/hash"
)

func nextPowerOfTwo(n int) int {
	// Return the number if it's already a power of 2.
	if n&(n-1) == 0 {
		return n
	}

	// Figure out and return the next power of two.
	exponent := uint(math.Log2(float64(n))) + 1
	return 1 << exponent // 2^exponent
}
func hashMerkleBranches(left *hash.Hash, right *hash.Hash) *hash.Hash {
	// Concatenate the left and right nodes.
	var h [hash.HashSize * 2]byte
	copy(h[:hash.HashSize], left[:])
	copy(h[hash.HashSize:], right[:])

	// TODO, add an abstract layer of hash func
	// TODO, double sha256 or other crypto hash
	newHash := hash.DoubleHashH(h[:])
	return &newHash
}
func (h *BlockHeader)BuildMerkleTreeStore() []*hash.Hash {
	// If there's an empty stake tree, return totally zeroed out merkle tree root
	// only.
	transactions := make([]hash.Hash,0)
	for i:=0;i<len(h.Transactions);i++{
		transactions = append(transactions,h.Transactions[i].Hash)
	}
	if len(transactions) == 0 {
		merkles := make([]*hash.Hash, 1)
		merkles[0] = &hash.Hash{}
		h.TxRoot = * merkles[len(merkles)-1]
		return merkles
	}

	// Calculate how many entries are required to hold the binary merkle
	// tree as a linear array and create an array of that size.
	nextPoT := nextPowerOfTwo(len(transactions))
	arraySize := nextPoT*2 - 1
	merkles := make([]*hash.Hash, arraySize)

	// Create the base transaction hashes and populate the array with them.
	for i, txHashFull := range transactions {
		//hash1 := common.Reverse(txHashFull[:])
		//var hash2 hash.Hash
		//copy(hash2[:],hash1[:])
		newHash := txHashFull
		merkles[i] = &newHash
	}

	// Start the array offset after the last transaction and adjusted to the
	// next power of two.
	offset := nextPoT
	for i := 0; i < arraySize-1; i += 2 {
		switch {
		// When there is no left child node, the parent is nil too.
		case merkles[i] == nil:
			merkles[offset] = nil

			// When there is no right child, the parent is generated by
			// hashing the concatenation of the left child with itself.
		case merkles[i+1] == nil:
			newHash := hashMerkleBranches(merkles[i], merkles[i])
			merkles[offset] = newHash

			// The normal case sets the parent node to the hash of the
			// concatentation of the left and right children.
		default:
			newHash := hashMerkleBranches(merkles[i], merkles[i+1])
			merkles[offset] = newHash
		}
		offset++
	}
	h.TxRoot = * merkles[len(merkles)-1]
	return merkles
}
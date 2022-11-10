package evm

import (
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
	"math/big"
)

// bigIntToByteArr converts Big Integer to fixed length binary byte array
func bigIntToByteArr(num *big.Int) []byte {

	localByteArr := make([]byte, 32)
	num.FillBytes(localByteArr)
	log.Info("(bigIntToByteArr) localByteArr=", localByteArr)

	return localByteArr
}

// calcMemGasCost calculates variable memory-expansion related gas cost
func calcMemGasCost(memory []byte) int {
	memorySizeWord := (len(memory) + 31) / 32
	memoryCost := (memorySizeWord*memorySizeWord)/512 + (3 * memorySizeWord)
	return memoryCost
}

// calcHash calculates KECCAK256 hash of memory (byte array)
func calcHash(memory []byte) string {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(memory)
	bs := hash.Sum(nil)
	return hex.EncodeToString(bs)
}

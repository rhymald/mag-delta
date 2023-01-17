package blockchain

import (
  "math/big"
  // "golang.org/x/crypto/bcrypt"
  "fmt"
  // "time"
  "encoding/binary"
  "bytes"
  "crypto/sha512"
  "math"
)

const Difficulty = 20

type PoW struct {
  Block *Block
  Target *big.Int
}

func NewProof(b *Block) *PoW {
  target := big.NewInt( 1 )
  target.Lsh(target, uint(512-Difficulty))
  return &PoW{Block: b, Target: target}
}

func InitData(pow *PoW, nonce int) []byte { return bytes.Join( [][]byte{ pow.Block.Data, pow.Block.Prev, BigToHex(int64(nonce)), BigToHex(int64(Difficulty)) }, []byte{} ) }

func BigToHex(num int64) []byte {
  buff := new(bytes.Buffer)
  _ = binary.Write(buff, binary.BigEndian, num)
  return buff.Bytes()
}

func Run(pow *PoW) (int, []byte) {
  var intHash big.Int
  var hash [64]byte
  nonce := 0
  for nonce < math.MaxInt64 {
    data := InitData(pow, nonce)
    hash = sha512.Sum512(data)
    fmt.Printf("\r%x", hash)
    // hash, _ = bcrypt.GenerateFromPassword( sum[:] , Difficulty)
    intHash.SetBytes(hash[:])
    if intHash.Cmp(pow.Target) == -1 {
      break
    } else {
      nonce++
    }
  }
  fmt.Println()
  // fin, _ := bcrypt.GenerateFromPassword( hash[:], Difficulty )
  return nonce, hash[:]
}

func Validate(pow *PoW) bool {
  var intHash big.Int
  data := InitData(pow, pow.Block.Nonce)
  hash := sha512.Sum512(data)
  intHash.SetBytes(hash[:])
  return intHash.Cmp(pow.Target) == -1
}

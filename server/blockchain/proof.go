package blockchain

import (
  "math/big"
  "encoding/binary"
  "bytes"
  "crypto/sha512"
  "crypto/rand"
  // "fmt"
)

type pow struct {
  Block *block
  Target *big.Int
  Mod *big.Int
}

func newProof(b *block, diff int) *pow {
  forearget := big.NewInt( 1 )
  forearget.Lsh(forearget, uint(512-diff))
  backarget := big.NewInt( 1 )
  backarget.Lsh(backarget, uint(diff))
  return &pow{Block: b, Target: forearget, Mod: backarget}
}

func initData(pow *pow, nonce int64) []byte { return bytes.Join( [][]byte{ pow.Block.Data, pow.Block.Prev, bigToHex(nonce), bigToHex(int64(takeDiff(pow.Block.Namespace, pow.Block.Time))) }, []byte{} ) }

func bigToHex(num int64) []byte {
  buff := new(bytes.Buffer)
  _ = binary.Write(buff, binary.BigEndian, num)
  return buff.Bytes()
}

func run(pow *pow) (int64, []byte, int) {
  var intHash big.Int
  var mod big.Int
  var hash [64]byte
  var counter int = 0
  nonce := new(big.Int)
  for {
    randy := make([]byte, 64)
    _,_ = rand.Read(randy)
    nonce.SetBytes(randy)
    data := initData(pow, int64(nonce.Uint64()))
    hash = sha512.Sum512(data)
    intHash.SetBytes(hash[:])
    mod.Mod(&intHash,pow.Mod)
    // fmt.Printf("\r%08b..%08b\r", hash[:3],hash[61:64])
    counter++
    if intHash.Cmp(pow.Target) == -1 && mod.Cmp(big.NewInt(0)) == 0 { break }
  }
  return int64(nonce.Uint64()), hash[:], counter
}

func validate(pow *pow) bool {
  var intHash big.Int
  var mod big.Int
  data := initData(pow, pow.Block.Nonce)
  hash := sha512.Sum512(data)
  intHash.SetBytes(hash[:])
  mod.Mod(&intHash,pow.Mod)
  return intHash.Cmp(pow.Target) == -1 && mod.Cmp(big.NewInt(0)) == 0
}

package blockchain

import (
  "math/big"
  "encoding/binary"
  "bytes"
  "crypto/sha512"
  "crypto/rand"
)

type pow struct {
  Block *block
  Target *big.Int
}

func newProof(b *block, diff int) *pow {
  target := big.NewInt( 1 )
  target.Lsh(target, uint(512-diff*4))
  return &pow{Block: b, Target: target}
}

func initData(pow *pow, nonce int64) []byte { return bytes.Join( [][]byte{ pow.Block.Data, pow.Block.Prev, bigToHex(nonce), bigToHex(int64(Diff[pow.Block.Namespace])) }, []byte{} ) }

func bigToHex(num int64) []byte {
  buff := new(bytes.Buffer)
  _ = binary.Write(buff, binary.BigEndian, num)
  return buff.Bytes()
}

func run(pow *pow) (int64, []byte) {
  var intHash big.Int
  var hash [64]byte
  nonce := new(big.Int)
  for {
    randy := make([]byte, 64)
    _,_ = rand.Read(randy)
    nonce.SetBytes(randy)
    data := initData(pow, int64(nonce.Uint64()))
    hash = sha512.Sum512(data)
    intHash.SetBytes(hash[:])
    if intHash.Cmp(pow.Target) == -1 { break }
  }
  return int64(nonce.Uint64()), hash[:]
}

func validate(pow *pow) bool {
  var intHash big.Int
  data := initData(pow, pow.Block.Nonce)
  hash := sha512.Sum512(data)
  intHash.SetBytes(hash[:])
  return intHash.Cmp(pow.Target) == -1
}

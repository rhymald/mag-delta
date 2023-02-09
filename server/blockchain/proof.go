package blockchain

import (
  "math/big"
  // "golang.org/x/crypto/bcrypt"
  // "fmt"
  // "time"
  "encoding/binary"
  "bytes"
  "crypto/sha512"
  "crypto/rand"
  // "math"
)

var Diff map[string]int = map[string]int{
  "/": 12,
  "/Players": 8,
  "/NPC": 4,
}

// const InitialDiff = 32
// const PlayerDiff = 24 // playable players
// const PhenomenaeDiff = 20 // nature, weather and objects
// const NPCDiff = 16 // nonplayable player
// const SessionDiff = 12 // for player sessions events: cast, ealth, regen, move, etc.
// const LifecycleDiff = 8 // for spawn/death and drop/loot

type pow struct {
  Block *block
  Target *big.Int
}

func newProof(b *block, diff int) *pow {
  target := big.NewInt( 1 )
  target.Lsh(target, uint(512-diff))
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
  for {//nonce.Cmp(big.NewInt(math.MaxInt64)) == -1 {
    randy := make([]byte, 64)
    _,_ = rand.Read(randy)
    nonce.SetBytes(randy)
    data := initData(pow, int64(nonce.Uint64()))
    hash = sha512.Sum512(data)
    // fmt.Printf("\r\u001b[38;5;229m\u001b[1m[%.10x] \u001b[0m", hash)
    // hash, _ = bcrypt.GenerateFromPassword( sum[:] , Difficulty)
    intHash.SetBytes(hash[:])
    if intHash.Cmp(pow.Target) == -1 { break }
  }
  // fmt.Println()
  // fin, _ := bcrypt.GenerateFromPassword( hash[:], Difficulty )
  return int64(nonce.Uint64()), hash[:]
}

func validate(pow *pow) bool {
  var intHash big.Int
  data := initData(pow, pow.Block.Nonce)
  hash := sha512.Sum512(data)
  intHash.SetBytes(hash[:])
  return intHash.Cmp(pow.Target) == -1
}

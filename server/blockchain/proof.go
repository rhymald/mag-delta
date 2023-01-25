package blockchain

import (
  "math/big"
  // "golang.org/x/crypto/bcrypt"
  // "fmt"
  // "time"
  "encoding/binary"
  "bytes"
  "crypto/sha512"
  "math"
)

const PlayerDiff = 24 // playable players
const PhenomenaeDiff = 20 // nature, weather and objects
const NPCDiff = 16 // nonplayable player
const SessionDiff = 12 // for player sessions events: cast, ealth, regen, move, etc.
const LifecycleDiff = 8 // for spawn/death and drop/loot

type pow struct {
  Block *block
  Target *big.Int
}

func newProof(b *block) *pow {
  target := big.NewInt( 1 )
  target.Lsh(target, uint(512-PlayerDiff))
  return &pow{Block: b, Target: target}
}

func initData(pow *pow, nonce int) []byte { return bytes.Join( [][]byte{ pow.Block.Data, pow.Block.Prev, bigToHex(int64(nonce)), bigToHex(int64(PlayerDiff)) }, []byte{} ) }

func bigToHex(num int64) []byte {
  buff := new(bytes.Buffer)
  _ = binary.Write(buff, binary.BigEndian, num)
  return buff.Bytes()
}

func run(pow *pow) (int, []byte) {
  var intHash big.Int
  var hash [64]byte
  nonce := 0
  for nonce < math.MaxInt64 {
    data := initData(pow, nonce)
    hash = sha512.Sum512(data)
    // fmt.Printf("\r%x", hash)
    // hash, _ = bcrypt.GenerateFromPassword( sum[:] , Difficulty)
    intHash.SetBytes(hash[:])
    if intHash.Cmp(pow.Target) == -1 {
      break
    } else {
      nonce++
    }
  }
  // fmt.Println()
  // fin, _ := bcrypt.GenerateFromPassword( hash[:], Difficulty )
  return nonce, hash[:]
}

func validate(pow *pow) bool {
  var intHash big.Int
  data := initData(pow, pow.Block.Nonce)
  hash := sha512.Sum512(data)
  intHash.SetBytes(hash[:])
  return intHash.Cmp(pow.Target) == -1
}

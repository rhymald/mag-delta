package blockchain

import (
  "bytes"
  "golang.org/x/crypto/bcrypt"
  "crypto/sha512"
  "encoding/json"
  "encoding/base64"
  "time"
  "fmt"
  "rhymald/mag-delta/player"
  "rhymald/mag-delta/funcs"
)

type BlockChain struct {
  Blocks []*Block
}

type Block struct {
  Time int64
  Hash []byte
  Data []byte
  Prev []byte
}

func (block *Block) CalculateHash() {
  info := bytes.Join( [][]byte{ block.Data, block.Prev }, []byte{} )
  sum := sha512.Sum512(info)
  hash, err := bcrypt.GenerateFromPassword( sum[:] , Difficulty)
  if err != nil { fmt.Println(err) }
  block.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
  block := &Block{Hash: []byte{}, Data: []byte(data), Prev: prevHash, Time: time.Now().UnixNano() }
  block.CalculateHash()
  return block
}

func AddBlock(chain *BlockChain, player player.Player) {
  player.Physical.Health.Current = 0
  player.Nature.Pool.Dots = []funcs.Dot{}
  datastring := ToJson(player)
  prevBlock := chain.Blocks[len(chain.Blocks)-1] // last block
  new := CreateBlock(datastring, prevBlock.Hash)
  chain.Blocks = append(chain.Blocks, new)
}

func Genesis() *Block {
  // text, _ := bcrypt.GenerateFromPassword( []byte("Hello, artifical world!") , 10)
  return CreateBlock( "Hello, artifical world!", []byte("His Allmightyness, Energy of Unpredictable Activity"))
}
func InitBlockChain() *BlockChain { return &BlockChain{[]*Block{Genesis()}} }

func ToJson(thing player.Player) string {
  // fmt.Println("  ─────────────────────────────────────────────────────────────────────────────────────────────────────")
  b, err := json.Marshal(thing)
  if err != nil { fmt.Println(err) ; return "" }
  // fmt.Println(string(b))
  encoded := base64.StdEncoding.EncodeToString(b)
  // fmt.Println(encoded)
  // fmt.Println("   ──── ──── ──── ──── ──── ──── ──── ──── ──── ──── ──── ──── ──── ──── ──── ──── ──── ──── ──── ──── ")
  _ = FromJson(encoded, thing)
  return encoded
}

func FromJson(code string, thing player.Player) player.Player {
  copy := &thing
  decoded, _ := base64.StdEncoding.DecodeString(code)
  // fmt.Println(string(decoded))
  err := json.Unmarshal(decoded, copy)
  if err != nil { fmt.Println(err) ; return player.Player{} }
  // fmt.Printf("%+v\n", *copy)
  // fmt.Println("  ─────────────────────────────────────────────────────────────────────────────────────────────────────")
  return *copy
}

func ListBlocks(chain *BlockChain) {
  fmt.Println(" ─┼─┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────")
  for i, each := range chain.Blocks {
    fmt.Printf("  ┼─┼─ %d ─── ─── ───\n  │ Hash\t%s\n  │ Time\t%d\n  │ Data\t%.100s\n  │ Parent\t%s\n", i, string(each.Hash), each.Time, each.Data, each.Prev)
  }
  fmt.Println(" ─┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────────")
}

package blockchain

import (
  // "bytes"
  // "crypto/sha512"
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
  Nonce int
}
//
// func (block *Block) CalculateHash() {
//   info := bytes.Join( [][]byte{ block.Data, block.Prev }, []byte{} )
//   sum := sha512.Sum512(info)
//   hash, err := bcrypt.GenerateFromPassword( sum[:] , Difficulty)
//   if err != nil { fmt.Println(err) }
//   block.Hash = hash[:]
// }

func CreateBlock(data string, prevHash []byte) *Block {
  block := &Block{Hash: []byte{}, Data: []byte(data), Prev: prevHash, Time: time.Now().UnixNano(), Nonce: 0 }
  // block.CalculateHash()
  pow := NewProof(block)
  nonce, hash := Run(pow)
  block.Hash = hash[:]
  block.Nonce = nonce
  return block
}

func AddBlock(chain *BlockChain, player player.Player) {
  player.Physical.Health.Current = 0
  player.Nature.Pool.Dots = []funcs.Dot{}
  player.Busy = false
  datastring := ToJson(player)
  prevBlock := chain.Blocks[len(chain.Blocks)-1] // last block
  new := CreateBlock(datastring, prevBlock.Hash)
  if len(chain.Blocks) != 0 {
    if datastring == string(chain.Blocks[len(chain.Blocks)-1].Data) { return }
  }
  chain.Blocks = append(chain.Blocks, new)
}

func Genesis() *Block {
  // text, _ := bcrypt.GenerateFromPassword( []byte("Hello, artifical world!") , 10)
  return CreateBlock( "Hello, artifical world!", []byte{})
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
  // fmt.Println(" ─┼─┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────")
  for i, each := range chain.Blocks {
    fmt.Printf("  │ %x\n", each.Prev)
    fmt.Printf(" ─┼─── %d ", i)
    fmt.Printf(" ─── Time %d", each.Time)
    fmt.Printf(" ─── Nonce %d", each.Nonce)
    fmt.Printf(" ─── Valid %v\n", Validate(NewProof(each)))
    fmt.Printf("  │ Data: %s\n", each.Data)
    fmt.Printf("  │ %x\n", string(each.Hash))
    // pow := NewProof(each)
  }
  fmt.Println("  │ ")
  // fmt.Println(" ─┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────────")
}

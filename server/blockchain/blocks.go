package blockchain

import (
  "bytes"
  // "crypto/sha512"
  "encoding/json"
  "encoding/base64"
  "time"
  "fmt"
  "rhymald/mag-delta/player"
  // "rhymald/mag-delta/funcs"
  "encoding/gob"
)


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

func Genesis() *Block {
  // text, _ := bcrypt.GenerateFromPassword( []byte("Hello, artifical world!") , 10)
  return CreateBlock( "Hello, artifical world!", []byte{})
}

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

// func ListBlocks(chain *BlockChain) {
//   // fmt.Println(" ─┼─┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────")
//   for i, each := range chain.Blocks {
//     fmt.Printf("  │ %x\n", each.Prev)
//     fmt.Printf(" ─┼─── %d ", i)
//     fmt.Printf(" ─── Time %d", each.Time)
//     fmt.Printf(" ─── Nonce %d", each.Nonce)
//     fmt.Printf(" ─── Valid %v\n", Validate(NewProof(each)))
//     fmt.Printf("  │ Data: %s\n", each.Data)
//     fmt.Printf("  │ %x\n", string(each.Hash))
//     // pow := NewProof(each)
//   }
//   fmt.Println("  │ \n")
//   // fmt.Println(" ─┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────────")
// }

func Serialize(b *Block) []byte {
  var res bytes.Buffer
  encoder := gob.NewEncoder(&res)
  err := encoder.Encode(b)
  if err != nil { fmt.Println(err) }
  return res.Bytes()
}

func Deserialize(data []byte) *Block {
  var block Block
  decoder := gob.NewDecoder(bytes.NewReader(data))
  err := decoder.Decode(&block)
  if err != nil { fmt.Println(err) }
  return &block
}

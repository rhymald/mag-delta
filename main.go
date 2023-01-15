package main

import (
  "fmt"
  "math"
  "rhymald/mag-delta/client/plot"
  "rhymald/mag-delta/client"
  "rhymald/mag-delta/player"
  "rhymald/mag-delta/act"
)
import (
  "bytes"
  "golang.org/x/crypto/bcrypt"
  "encoding/json"
  "encoding/base64"
  "time"
)

var You player.Player
var Target player.Player
var Action string
var Keys chan string = make(chan string)
var BC *BlockChain = InitBlockChain()

func init() {
  fmt.Println("\n\t\t  ", plot.Bar("Initializing...",8), "\n")
  player.PlayerBorn(&You,0)
  BC.AddBlock(ToJson(You))
  player.FoeSpawn(&Target,0)
}

func main() {
  fmt.Println("\n\t\t", plot.Bar("Successfully login",1),"\n")
  client.PlayerStatus(You, Target)
  fmt.Println("\n\t     ",plot.Bar("Press [Enter] to continue",8),"\n")
  fmt.Scanln()
  plot.ShowMenu(" ")
  client.PlayerStatus(You, Target)
  client.UI(Keys, You, Target)
  grow := 1/math.Phi/math.Phi/math.Phi
  for {
    key := actions()
    if Target.Physical.Health.Current <= 0 { grow = grow*math.Cbrt(math.Phi) ; player.FoeSpawn(&Target, grow) ; plot.ShowMenu(key)}// ; PlayerStatus(You, Target)}
  }
}

func actions() string {
  Action, _ := <-Keys
  key := Action
  if Action=="a" { act.Jinx(&You, &Target) ; Action = "" }
  return string(key)
}

// Stateful database SERVER!!!

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
  info := bytes.Join([][]byte{ block.Data, block.Prev }, []byte{})
  hash, _ := bcrypt.GenerateFromPassword(info, 7)
  block.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
  block := &Block{Hash: []byte{}, Data: []byte(data), Prev: prevHash, Time: time.Now().UnixNano() }
  block.CalculateHash()
  return block
}

func (chain *BlockChain) AddBlock(data string) {
  prevBlock := chain.Blocks[len(chain.Blocks)-1] // last block
  new := CreateBlock(data, prevBlock.Prev)
  chain.Blocks = append(chain.Blocks, new)
}

func Genesis() *Block { return CreateBlock("Hello World", []byte{}) }
func InitBlockChain() *BlockChain { return &BlockChain{[]*Block{Genesis()}} }

func ToJson(thing player.Player) string {
  fmt.Println("  ─────────────────────────────────────────────────────────────────────────────────────────────────────")
  b, err := json.Marshal(thing)
  if err != nil {
    fmt.Println(err)
    return ""
  }
  fmt.Println(string(b))
  encoded := base64.StdEncoding.EncodeToString(b)
  fmt.Println(encoded)
  fmt.Println("   ──── ──── ──── ──── ──── ──── ──── ──── ──── ──── ──── ──── ──── ──── ──── ──── ──── ──── ──── ──── ")
  _ = FromJson(encoded, thing)
  return encoded
}

func FromJson(code string, thing player.Player) player.Player {
  copy := &thing
  decoded, _ := base64.StdEncoding.DecodeString(code)
  fmt.Println(string(decoded))
  err := json.Unmarshal(decoded, copy)
  if err != nil {
    fmt.Println(err)
    return player.Player{}
  }
  fmt.Printf("%+v\n", *copy)
  fmt.Println("  ─────────────────────────────────────────────────────────────────────────────────────────────────────")
  return *copy
}

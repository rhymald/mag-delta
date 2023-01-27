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

type block struct {
  Time int64
  Namespace string
  Hash []byte
  Data []byte
  Prev []byte
  Nonce int64
}

func createBlock(data string, ns string, prevHash []byte, diff int) *block {
  block := &block{Hash: []byte{}, Data: []byte(data), Prev: prevHash, Time: time.Now().UnixNano(), Nonce: 0, Namespace: ns }
  pow := newProof(block, diff)
  nonce, hash := run(pow)
  block.Hash = hash[:]
  block.Nonce = nonce
  return block
}

func genesis() *block {
  return createBlock(base64.StdEncoding.EncodeToString([]byte("GENESIS BLOCK: ThickCat Concensus Protocol initialized. Hello, artifical World!")), "Initial", []byte{}, Diff["Initial"])
}

func toJson(thing player.Player) string {
  b, err := json.Marshal(thing)
  if err != nil { fmt.Println(err) ; return "" }
  encoded := base64.StdEncoding.EncodeToString(b)
  return encoded
}

func fromJson(code string, thing player.Player) player.Player {
  copy := &thing
  decoded, _ := base64.StdEncoding.DecodeString(code)
  err := json.Unmarshal(decoded, copy)
  if err != nil { fmt.Println(err) ; return player.Player{} }
  return *copy
}

func serialize(b *block) []byte {
  var res bytes.Buffer
  encoder := gob.NewEncoder(&res)
  err := encoder.Encode(b)
  if err != nil { fmt.Println(err) }
  return res.Bytes()
}

func deserialize(data []byte) *block {
  var block block
  decoder := gob.NewDecoder(bytes.NewReader(data))
  err := decoder.Decode(&block)
  if err != nil { fmt.Println(err) }
  return &block
}

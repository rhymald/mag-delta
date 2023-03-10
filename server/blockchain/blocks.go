package blockchain

import (
  "bytes"
  // "crypto/sha512"
  "encoding/base64"
  "time"
  "fmt"
  // "rhymald/mag-delta/player"
  // "rhymald/mag-delta/funcs"
  "encoding/gob"
)

type block struct {
  Time int64
  Namespace string
  Behind []byte
  Hash []byte
  Data []byte
  Prev []byte
  Nonce int64
}

func createBlock(data string, ns string, prevHash []byte, diff int, behind []byte, epoch int64) *block {
  block := &block{Hash: []byte{}, Data: []byte(data), Prev: prevHash, Behind: behind, Time: epoch, Nonce: 0, Namespace: ns }
  pow := newProof(block, diff)
  nonce, hash := run(pow)
  block.Hash = hash[:]
  block.Nonce = nonce
  return block
}

func genesis() *block {
  epoch := time.Now().UnixNano()-1317679200000000000
  return createBlock(base64.StdEncoding.EncodeToString([]byte("GENESIS BLOCK: ThickCat Concensus Protocol initialized. Hello, artifical World!")), "/", []byte{}, takeDiff("/", epoch), []byte{}, epoch)
}

func serialize(b *block) []byte {
  var res bytes.Buffer
  encoder := gob.NewEncoder(&res)
  err := encoder.Encode(b)
  if err != nil { fmt.Println(err) }
  return res.Bytes()
}

func Deserialize(data []byte) *block {
  var block block
  decoder := gob.NewDecoder(bytes.NewReader(data))
  err := decoder.Decode(&block)
  if err != nil { fmt.Println(err) }
  return &block
}

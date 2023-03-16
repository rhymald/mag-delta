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
  Transactions []*Transaction
  Namespace string
  Behind []byte
  Hash []byte
  Data []byte
  Prev []byte
  Nonce int64
}

func createBlock(data string, txs []*Transaction, ns string, prevHash []byte, diff int, behind []byte, epoch int64) (*block, int) {
  block := &block{Hash: []byte{}, Transactions: txs ,Data: []byte(data), Prev: prevHash, Behind: behind, Time: epoch, Nonce: 0, Namespace: ns }
  pow := newProof(block, diff)
  nonce, hash, counter := run(pow)
  block.Hash = hash[:]
  block.Nonce = nonce
  return block, counter
}

func genesis(coinbase *Transaction) *block {
  epoch := time.Now().UnixNano()-1317679200000000000
  genblock, _ := createBlock(base64.StdEncoding.EncodeToString([]byte("GENESIS BLOCK")), []*Transaction{coinbase},"/", []byte{}, takeDiff("/", epoch), []byte{}, epoch)
  return genblock
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
  if err != nil { fmt.Println("Block deserialise failed:", err) }
  return &block
}

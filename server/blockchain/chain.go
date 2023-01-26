package blockchain

import (
  "github.com/dgraph-io/badger"
  "fmt"
  "encoding/base64"
  "rhymald/mag-delta/player"
  "rhymald/mag-delta/funcs"
  "time"
)

type BlockChain struct {
  // Blocks []*Block
  Database *badger.DB
  LastHash []byte
}

type bcIterator struct {
  // Blocks []*Block
  Database *badger.DB
  Current []byte
}

func InitBlockChain(dbPath string) *BlockChain {
  // return &BlockChain{[]*Block{Genesis()}}
  var lastHash []byte
  opts := badger.DefaultOptions(dbPath)
  opts.Dir = dbPath
  opts.ValueDir = dbPath
  db, err := badger.Open(opts)
  if err != nil { fmt.Println(err) }
  // run writing query-connection
  err = db.Update(func(txn *badger.Txn) error {
    // if there is no last hash in db
    if _, err := txn.Get([]byte("Players")); err == badger.ErrKeyNotFound {
      fmt.Printf("Blockchain does not exist! Genereating...")
      genesis := genesis()
      fmt.Printf(" Writing...")
      err := txn.Set(genesis.Hash, serialize(genesis))
      if err != nil { fmt.Println(err) }
      err = txn.Set([]byte("Players"), genesis.Hash) // link to last block inside db
      fmt.Printf(" Genesis block provided!\n")
      lastHash = genesis.Hash
      return err
    } else { // if exists
      item, err := txn.Get([]byte("Players"))
      if err != nil { fmt.Println(err) }
      lastHash, err = item.ValueCopy([]byte("Players")) // ???
      return err
    }
  })
  if err != nil { fmt.Println(err) }
  return &BlockChain{LastHash: lastHash, Database: db}
}

func AddPlayer(chain *BlockChain, player player.Player) {
  player.Physical.Health.Current = 0
  player.Nature.Pool.Dots = []funcs.Dot{}
  player.Busy = false
  dataString := toJson(player)
  var lastHash []byte
  // run read only txn (connection query)
  err := chain.Database.View(func(txn *badger.Txn) error {
    item, err := txn.Get([]byte("Players"))
    if err != nil { fmt.Println(err) }
    lastHash, err = item.ValueCopy([]byte("Players"))
    return err
  })
  if err != nil { fmt.Println(err) }
  addBlock(chain, dataString, lastHash, []byte("Players"))
}

func addBlock(chain *BlockChain, data string, lastHash []byte, namespace []byte) {
  // Player clean
  // prevBlock := chain.Blocks[len(chain.Blocks)-1] // last block
  // new := CreateBlock(datastring, prevBlock.Hash)
  // if len(chain.Blocks) != 0 {
  //   if datastring == string(chain.Blocks[len(chain.Blocks)-1].Data) { return } !!!
  // }
  // chain.Blocks = append(chain.Blocks, new)
  // var lastHash []byte
  // // run read only txn (connection query)
  // err := chain.Database.View(func(txn *badger.Txn) error {
  //   item, err := txn.Get(namespace)
  //   if err != nil { fmt.Println(err) }
  //   lastHash, err = item.ValueCopy(namespace)
  //   return err
  // })
  // here +fetch
  new := createBlock(data, lastHash)
  err := chain.Database.Update(func(txn *badger.Txn) error {
    err := txn.Set(new.Hash, serialize(new))
    if err != nil { fmt.Println(err) }
    err = txn.Set((namespace), new.Hash) // link to last block inside db
    chain.LastHash = new.Hash
    return err
  })
  if err != nil { fmt.Println(err) }
}

func iterator(chain *BlockChain) *bcIterator { return &bcIterator{Current: chain.LastHash, Database: chain.Database} }

func deeper(iter *bcIterator) *block {
  var block *block
  err := iter.Database.View(func(txn *badger.Txn) error {
    item, err := txn.Get(iter.Current)
    if err != nil { fmt.Println(err) }
    encodedBlock, err := item.ValueCopy(iter.Current)
    block = deserialize(encodedBlock)
    return err
  })
  if err != nil { fmt.Println(err) }
  iter.Current = block.Prev // step back
  return block
}

func ListBlocks(chain *BlockChain, namespace string) {
  iter := iterator(chain)
  depth := 0
  next := &block{Time: time.Now().UnixNano()}
  // fmt.Println(" ─┼─┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────")
  for i:=0; i<10; i++ {
    each := deeper(iter)
    if each.Namespace != namespace { break }
    fmt.Printf("  │ %x\n", string(each.Hash))
    fmt.Printf(" ─┼─── %d'", -depth)
    fmt.Printf("%s", each.Namespace)
    fmt.Printf(" ─── Time %d", each.Time)
    fmt.Printf(" ─── Gape %0.3f s.", float64(each.Time-next.Time)/1000000000)
    fmt.Printf(" ─── Nonce %d", each.Nonce)
    fmt.Printf(" ─── Valid %v\n", validate(newProof(each)))
    decoded, _ := base64.StdEncoding.DecodeString(string(each.Data))
    fmt.Printf("  │ Data: \u001b[1m%s\n\u001b[0m", decoded)
    fmt.Printf("  │ %x\n", each.Prev)
    depth--
    next = each
    if len(each.Prev) == 0 { break }
    // pow := NewProof(each)
  }
  fmt.Println("  │ \n")
  // fmt.Println(" ─┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────────")
}

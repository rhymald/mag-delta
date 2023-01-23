package blockchain

import (
  "github.com/dgraph-io/badger"
  "fmt"
  "rhymald/mag-delta/player"
  "rhymald/mag-delta/funcs"
)

const (
  // dbPath = "./tmp/blocks"
  // runtime.Goexit() - shutdown app
)

type BlockChain struct {
  // Blocks []*Block
  Database *badger.DB
  LastHash []byte
}

type BCIterator struct {
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
    if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
      fmt.Printf("Blockchain does not exist! Genereating...")
      genesis := Genesis()
      fmt.Printf(" Writing...")
      err := txn.Set(genesis.Hash, Serialize(genesis))
      if err != nil { fmt.Println(err) }
      err = txn.Set([]byte("lh"), genesis.Hash) // link to last block inside db
      fmt.Printf(" Genesis block provided!\n")
      lastHash = genesis.Hash
      return err
    } else { // if exists
      item, err := txn.Get([]byte("lh"))
      if err != nil { fmt.Println(err) }
      lastHash, err = item.ValueCopy([]byte("lh")) // ???
      return err
    }
  })
  if err != nil { fmt.Println(err) }
  return &BlockChain{LastHash: lastHash, Database: db}
}

func AddBlock(chain *BlockChain, player player.Player) {
  // Player clean
  player.Physical.Health.Current = 0
  player.Nature.Pool.Dots = []funcs.Dot{}
  player.Busy = false
  datastring := ToJson(player)
  // prevBlock := chain.Blocks[len(chain.Blocks)-1] // last block
  // new := CreateBlock(datastring, prevBlock.Hash)
  // if len(chain.Blocks) != 0 {
  //   if datastring == string(chain.Blocks[len(chain.Blocks)-1].Data) { return } !!!
  // }
  // chain.Blocks = append(chain.Blocks, new)
  var lastHash []byte
  // run read only txn (connection query)
  err := chain.Database.View(func(txn *badger.Txn) error {
    item, err := txn.Get([]byte("lh"))
    if err != nil { fmt.Println(err) }
    lastHash, err = item.ValueCopy([]byte("lh"))
    return err
  })
  // here +fetch
  new := CreateBlock(datastring, lastHash)
  err = chain.Database.Update(func(txn *badger.Txn) error {
    err := txn.Set(new.Hash, Serialize(new))
    if err != nil { fmt.Println(err) }
    err = txn.Set([]byte("lh"), new.Hash) // link to last block inside db
    chain.LastHash = new.Hash
    return err
  })
  if err != nil { fmt.Println(err) }
}

func Iterator(chain *BlockChain) *BCIterator { return &BCIterator{Current: chain.LastHash, Database:  chain.Database} }

func Deeper(iter *BCIterator) *Block {
  var block *Block
  err := iter.Database.View(func(txn *badger.Txn) error {
    item, err := txn.Get(iter.Current)
    if err != nil { fmt.Println(err) }
    encodedBlock, err := item.ValueCopy(iter.Current)
    block = Deserialize(encodedBlock)
    return err
  })
  if err != nil { fmt.Println(err) }
  iter.Current = block.Prev // step back
  return block
}

func ListBlocks(chain *BlockChain) {
  iter := Iterator(chain)
  depth := 0
  // fmt.Println(" ─┼─┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────")
  for {//} i, each := range chain.Blocks {
    each := Deeper(iter)
    fmt.Printf("  │ %x\n", string(each.Hash))
    fmt.Printf(" ─┼─── %d ", -depth)
    fmt.Printf(" ─── Time %d", each.Time)
    fmt.Printf(" ─── Nonce %d", each.Nonce)
    fmt.Printf(" ─── Valid %v\n", Validate(NewProof(each)))
    fmt.Printf("  │ Data: %s\n", each.Data)
    fmt.Printf("  │ %x\n", each.Prev)
    depth--
    if len(each.Prev) == 0 { break }
    // pow := NewProof(each)
  }
  fmt.Println("  │ \n")
  // fmt.Println(" ─┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────────")
}

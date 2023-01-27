package blockchain

import (
  "fmt"
  "github.com/dgraph-io/badger"
  // "encoding/base64"
  // "rhymald/mag-delta/funcs"
  // "rhymald/mag-delta/player"
  // "time"
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

// func CreateContextAfter() error {
//   err = db.Update(func(txn *badger.Txn) error {
//     // if there is no last hash in db
//     if _, err := txn.Get([]byte("Initial")); err == badger.ErrKeyNotFound {
//       fmt.Printf("Blockchain does not exist! Genereating...")
//       genesis := genesis()
//       fmt.Printf(" Writing...")
//       err := txn.Set(genesis.Hash, serialize(genesis))
//       if err != nil { fmt.Println(err) }
//       err = txn.Set([]byte("Initial"), genesis.Hash) // link to last block inside db
//       fmt.Printf(" Genesis block provided!\n")
//       lastHash = genesis.Hash
//       return err
//     } else { // if exists
//       item, err := txn.Get([]byte("Initial"))
//       if err != nil { fmt.Println(err) }
//       lastHash, err = item.ValueCopy([]byte("Initial")) // ???
//       return err
//     }
//   })
//   return err
// }

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
    if _, err := txn.Get([]byte("Initial")); err == badger.ErrKeyNotFound {
      fmt.Printf("Blockchain does not exist! Genereating...")
      genesis := genesis()
      fmt.Printf(" Writing...")
      err := txn.Set(genesis.Hash, serialize(genesis))
      if err != nil { fmt.Println(err) }
      err = txn.Set([]byte("Initial"), genesis.Hash) // link to last block inside db
      fmt.Printf(" Genesis block provided!\n")
      lastHash = genesis.Hash
      return err
    } else { // if exists
      item, err := txn.Get([]byte("Initial"))
      if err != nil { fmt.Println(err) }
      lastHash, err = item.ValueCopy([]byte("Initial")) // ???
      return err
    }
  })
  if err != nil { fmt.Println(err) }
  return &BlockChain{LastHash: lastHash, Database: db}
}

// upodate context
func addBlock(chain *BlockChain, data string, lastHash []byte, namespace []byte) {
  new := createBlock(data, string(namespace), lastHash, Diff[string(namespace)])
  // get prev data
  var prevData []byte
  err := chain.Database.View(func(txn *badger.Txn) error {
    item, err := txn.Get(lastHash)
    if err != nil { fmt.Println(err) }
    prevData, err = item.ValueCopy(lastHash)
    return err
  })
  if err != nil { fmt.Println(err) }
  prevBlock := deserialize(prevData)
  // if == no write
  if data == string(*&prevBlock.Data) { return }
  // end if
  err = chain.Database.Update(func(txn *badger.Txn) error {
    err := txn.Set(new.Hash, serialize(new))
    if err != nil { fmt.Println(err) }
    err = txn.Set((namespace), new.Hash)
    if string(namespace) == "Players[]" { // auto create subNSs
      rowName := fmt.Sprintf("Players[%.8X]", new.Hash)
      err = txn.Set( []byte(rowName), new.Hash)
      rowName = fmt.Sprintf("Session[%.8X]", new.Hash)
      err = txn.Set( []byte(rowName), new.Hash)
    }
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

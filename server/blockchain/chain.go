package blockchain

import (
  "fmt"
  "github.com/dgraph-io/badger"
  "time"
  "encoding/base64"
  "sync"
)

type BlockChain struct {
  Epoch int64
  Database *badger.DB
  LastHash map[string][]byte
  sync.Mutex
}

type bcIterator struct {
  Database *badger.DB
  Current []byte
}

func FindByPrefixes(chain *BlockChain, prefix []byte) [][]byte {
  var playerList [][]byte
  chain.Database.View( func(txn *badger.Txn) error {
    iterator := txn.NewIterator(badger.DefaultIteratorOptions)
    defer iterator.Close()
    for iterator.Seek(prefix); iterator.ValidForPrefix(prefix); iterator.Next() {
      item := iterator.Item()
      key := item.Key()
      err := item.Value(func (v []byte) error {
        playerList = append(playerList, []byte(fmt.Sprintf("\u001b[1m%s\u001b[0m %x", key, v)))
        return nil
      })
      if err != nil { return err }
    }
    return nil
  })
  if len(playerList) == 0 { playerList = append(playerList, prefix) }
  return playerList
}

func InitBlockChain(dbPath string) *BlockChain {
  var lastHash []byte
  var epoch int64
  opts := badger.DefaultOptions(dbPath)
  opts.Dir = dbPath
  opts.ValueDir = dbPath
  db, err := badger.Open(opts)
  if err != nil { fmt.Println(err) }
  err = db.Update(func(txn *badger.Txn) error {
    if _, err := txn.Get([]byte("/")); err == badger.ErrKeyNotFound {
      fmt.Printf("Blockchain does not exist! Genereating...")
      genesis := genesis()
      epoch = genesis.Time
      fmt.Printf(" Writing...")
      err := txn.Set(genesis.Hash, serialize(genesis))
      if err != nil { fmt.Println(err) }
      err = txn.Set([]byte("/"), genesis.Hash) // link to last block inside db
      fmt.Printf(" Genesis block provided!\n")
      lastHash = genesis.Hash
      return err
    } else { // if exists
      item, err := txn.Get([]byte("/"))
      if err != nil { fmt.Println(err) }
      lastHash, err = item.ValueCopy([]byte("/")) // ???
      return err
    }
  })
  if err != nil { fmt.Println(err) }
  lasts := make(map[string][]byte)
  lasts["/"] = lastHash
  return &BlockChain{LastHash: lasts, Database: db, Epoch: epoch}
}

func iterator(chain *BlockChain, namespace string) *bcIterator {
  chain.Lock() ; current := chain.LastHash[namespace] ; chain.Unlock()
  return &bcIterator{Current: current, Database: chain.Database}
}

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
  if len(block.Behind)>0 { iter.Current = block.Behind }
  return block
}

func ListBlocks(chain *BlockChain, namespace string) {
  // var rows []string
  iter := iterator(chain, namespace)
  depth := 0
  next := &block{Time: time.Now().UnixNano(), Namespace: namespace}
  fmt.Println("════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════")
  for i:=0; i<10; i++ {
    each := deeper(iter)
    if each.Namespace != next.Namespace { fmt.Printf("────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────\n") }
    fmt.Printf("\u001b[1m%x\u001b[0m\n", string(each.Hash))
    if len(each.Behind)>0 { fmt.Printf("\u001b[1mTriggered by\u001b[0m %x\n", string(each.Behind)) }
    fmt.Printf("   %d'", -depth)
    fmt.Printf("\u001b[1m%s\u001b[0m", each.Namespace)
    fmt.Printf(" \u001b[1mTime\u001b[0m %d", each.Time)
    fmt.Printf(" \u001b[1mGape\u001b[0m %0.3fs.", float64(each.Time-next.Time)/1000000000)
    fmt.Printf(" \u001b[1mNonce\u001b[0m %d", each.Nonce)
    fmt.Printf(" \u001b[1mValid\u001b[0m %v\n", validate(newProof(each, Diff[each.Namespace])))
    decoded, _ := base64.StdEncoding.DecodeString(string(each.Data))
    fmt.Printf("\u001b[1mData\u001b[0m %s\n", decoded)
    depth--
    if len(each.Prev) == 0 { break } else { fmt.Printf("%x\n", each.Prev) }
    next = each
  }
  fmt.Println("════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════")
  playerList := FindByPrefixes(chain, []byte("/"))
  fmt.Println(" ─────── ──── ───────── ─ ─────── Metadata info ──────────────────────────────────────────────────────────────────────────────────────────────────────────────────")
  for _, each := range playerList { fmt.Println(string(each)) }
  fmt.Println(" ─────── ──── ───────── ─ ─────── ────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────\n")
}

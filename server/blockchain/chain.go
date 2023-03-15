package blockchain

import (
  "fmt"
  "github.com/dgraph-io/badger"
  "time"
  "encoding/base64"
  "sync"
  "os"
  "errors"
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

func FindByPrefixes(chain *BlockChain, prefix []byte) [][2][]byte {
  var playerList [][2][]byte
  chain.Database.View( func(txn *badger.Txn) error {
    iterator := txn.NewIterator(badger.DefaultIteratorOptions)
    defer iterator.Close()
    for iterator.Seek(prefix); iterator.ValidForPrefix(prefix); iterator.Next() {
      item := iterator.Item()
      key := item.Key()
      err := item.Value(func (v []byte) error {
        playerList = append(playerList, [2][]byte{ key, v })
        return nil
      })
      if err != nil { return err }
    }
    return nil
  })
  if len(playerList) == 0 { playerList = append(playerList, [2][]byte{ prefix, []byte{} }) }
  return playerList
}

func ensureDir(path string) error {
  err := os.Mkdir(path, 0755)
  if err == nil { return nil }
  if os.IsExist(err) {
    info, err := os.Stat(path)
    if err == nil { return nil }
    if !info.IsDir() { errors.New("invalid path, already exists, but not a dir") }
    return nil
  }
  return err
}

func InitBlockChain(dbPath string) *BlockChain {
  var lastHash []byte
  var epoch int64
  ensureDir(dbPath)
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
      err = db.View(func(txn *badger.Txn) error {
        item, err := txn.Get(lastHash)
        if err != nil { fmt.Println(err) }
        encodedBlock, err := item.ValueCopy(lastHash)
        block := Deserialize(encodedBlock)
        epoch = block.Time
        return err
      })
      return err
    }
  })
  if err != nil { fmt.Println(err) }
  lasts := make(map[string][]byte)
  lasts["/"] = lastHash
  chain := &BlockChain{LastHash: lasts, Database: db, Epoch: epoch}
  playerList := FindByPrefixes(chain, []byte("/Players"))
  for _, identity := range playerList { lasts[string(identity[0])] = identity[1] }
  playerList = FindByPrefixes(chain, []byte("/Session"))
  for _, identity := range playerList { lasts[string(identity[0])] = identity[1] }
  return chain
}

func iterator(chain *BlockChain, meta string) *bcIterator {
  chain.Lock() ; current := chain.LastHash[meta] ; chain.Unlock()
  return &bcIterator{Current: current, Database: chain.Database}
}

func deeper(iter *bcIterator, triggers bool) *block {
  var block *block
  err := iter.Database.View(func(txn *badger.Txn) error {
    item, err := txn.Get(iter.Current)
    if err != nil { fmt.Println(err) }
    encodedBlock, err := item.ValueCopy(iter.Current)
    block = Deserialize(encodedBlock)
    return err
  })
  if err != nil { fmt.Println(err) }
  iter.Current = block.Prev // step back
  if len(block.Behind)>0 && triggers { iter.Current = block.Behind }
  return block
}

func ListBlocks(chain *BlockChain, meta string, extended bool) {
  // var rows []string
  iter := iterator(chain, meta)
  depth := 0
  next := &block{Time: time.Now().UnixNano()-1317679200000000000-chain.Epoch, Namespace: meta}
  if !extended {
    fmt.Println("════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════")
    for {
      each := deeper(iter, false)
      if each.Namespace != next.Namespace { fmt.Printf("────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────\n") }
      fmt.Printf("\u001b[1m%x\u001b[0m\n", string(each.Hash))
      if len(each.Behind)>0 { fmt.Printf("\u001b[7mTriggered by \u001b[0m%x\n", string(each.Behind)) }
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
  } else {
    playerList := FindByPrefixes(chain, []byte("/Players"))
    // sessionsList := FindByPrefixes(chain, []byte("/Session"))
    fmt.Println(" ─────── ──── ───────── ─ ─────── Metadata info ──────────────────────────────────────────────────────────────────────────────────────────────────────────────────")
    for _, each := range playerList { fmt.Println(string(each[0]), fmt.Sprintf("%x", each[1])) }
    // for _, each := range sessionsList { fmt.Println(string(each[0]), fmt.Sprintf("%x", each[1])) }
    fmt.Println(" ─────── ──── ───────── ─ ─────── ────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────")
  }
}

func Gather_Blocks(chain *BlockChain, meta string) [][]byte {
  var buffer [][]byte 
  iter := iterator(chain, meta)
  next := &block{}
  for {
    each := deeper(iter, false)
    equalns := each.Namespace == meta && next.Namespace == meta //for same ns
    deadend := false ; 
    if len(meta) >= len(each.Namespace) { deadend = each.Namespace == meta[:len(each.Namespace)] && len(each.Prev) != 0 } // for "/Players"-born ns
    if len(meta) < len(each.Namespace) { deadend = each.Namespace[:len(meta)] == meta && len(each.Prev) != 0 } // for "/Players"-born ns
    if equalns || deadend { buffer = append(buffer, each.Data) ; fmt.Println(string(each.Data)) } else { break } 
    next = each
  }
  return buffer
}
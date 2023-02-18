package blockchain

import (
  "fmt"
  "github.com/dgraph-io/badger"
  // "encoding/base64"
  // "rhymald/mag-delta/funcs"
  // "rhymald/mag-delta/player"
  "time"
  "encoding/base64"
)

type BlockChain struct {
  // Blocks []*Block
  Database *badger.DB
  LastHash map[string][]byte
}

type bcIterator struct {
  // Blocks []*Block
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
    if _, err := txn.Get([]byte("/")); err == badger.ErrKeyNotFound {
      fmt.Printf("Blockchain does not exist! Genereating...")
      genesis := genesis()
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
  return &BlockChain{LastHash: lasts, Database: db}
}

// upodate context
func AddBlock(chain *BlockChain, data string, namespace string) []byte {
  // namespace := "/"
  // ^ immitation
  lastHash := chain.LastHash[namespace]
  new := createBlock(data, namespace, lastHash, Diff[namespace])
  var prevData []byte
  err := chain.Database.View(func(txn *badger.Txn) error {
    item, err := txn.Get(lastHash)
    if err != nil { fmt.Println(err) }
    prevData, err = item.ValueCopy(lastHash)
    return err
  })
  if err != nil { fmt.Println(err) }
  prevBlock := deserialize(prevData)
  if data == string(*&prevBlock.Data) { return []byte{} }
  err = chain.Database.Update(func(txn *badger.Txn) error {
    err := txn.Set(new.Hash, serialize(new))
    if err != nil { fmt.Println(err) }
    // update context
    // err = txn.Set([]byte(namespace), new.Hash)
    // chain.LastHash[namespace] = new.Hash
    // ^ updated context
    return err
  })
  if err != nil { fmt.Println(err) }
  return new.Hash[:]
}

func iterator(chain *BlockChain, namespace string) *bcIterator { return &bcIterator{Current: chain.LastHash[namespace], Database: chain.Database} }

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

func ListBlocks(chain *BlockChain, namespace string) []string {
  // var rows []string
  var playerIDs []string
  iter := iterator(chain, namespace)
  depth := 0
  next := &block{Time: time.Now().UnixNano(), Namespace: namespace}
  fmt.Println("════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════")
  for i:=0; i<10; i++ {
    each := deeper(iter)
    if each.Namespace != next.Namespace { fmt.Printf("────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────\n") }
    fmt.Printf("\u001b[1m%x\u001b[0m\n", string(each.Hash))
    fmt.Printf("   %d'", -depth)
    fmt.Printf("\u001b[1m%s\u001b[0m", each.Namespace)
    fmt.Printf(" \u001b[1mTime\u001b[0m %d", each.Time)
    if each.Namespace == "/Players" {
      // rows = append(rows, fmt.Sprintf("/Players/%.8X = %x", each.Hash, each.Hash))
      playerIDs = append(playerIDs, fmt.Sprintf("/Players/%.8X", each.Hash))
    }
    fmt.Printf(" \u001b[1mGape\u001b[0m %0.3fs.", float64(each.Time-next.Time)/1000000000)
    fmt.Printf(" \u001b[1mNonce\u001b[0m %d", each.Nonce)
    fmt.Printf(" \u001b[1mValid\u001b[0m %v\n", validate(newProof(each, Diff[each.Namespace])))
    decoded, _ := base64.StdEncoding.DecodeString(string(each.Data))
    fmt.Printf("\u001b[1mData\u001b[0m %s\n", decoded)
    // fmt.Printf("%x\n", each.Prev)
    // if each.Namespace != namespace { break }
    depth--
    if len(each.Prev) == 0 { break } else { fmt.Printf("%x\n", each.Prev) }
    next = each
    // pow := NewProof(each)
  }
  fmt.Println("════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════")
  playerList := FindByPrefixes(chain, []byte("/"))
  if namespace == "/Players" {
    fmt.Println("\n ─── Metadata info ────────────── ")
    // for _, link := range playerIDs { fmt.Println(link) }
    for _, each := range playerList { fmt.Println(string(each)) }
    fmt.Println(" ──────────────────────────────── \n")
  }
  if len(playerIDs) == 0 { playerIDs = append(playerIDs, "/Players") }
  return playerIDs
}

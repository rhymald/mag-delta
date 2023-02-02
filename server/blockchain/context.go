package blockchain

import(
  "time"
  "fmt"
  "rhymald/mag-delta/player"
  "github.com/dgraph-io/badger"
  "encoding/base64"
  // "rhymald/mag-delta/funcs"
)

var Diff map[string]int = map[string]int{
  "Initial": 12,
  "Players[]": 8,
  "NPC": 4,
}

func ListBlocks(chain *BlockChain, namespace string) []string {
  var rows []string
  var playerIDs []string
  iter := iterator(chain)
  depth := 0
  next := &block{Time: time.Now().UnixNano(), Namespace: namespace}
  fmt.Println(" ─┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────")
  for i:=0; i<10; i++ {
    each := deeper(iter)
    if each.Namespace != next.Namespace { fmt.Printf(" ─┼─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────── \n") }
    fmt.Printf("  │ \u001b[1m%x\u001b[0m\n", string(each.Hash))
    fmt.Printf(" ─┼─── %d'", -depth)
    fmt.Printf("\u001b[1m%s\u001b[0m", each.Namespace)
    fmt.Printf(" ─── \u001b[1mTime\u001b[0m %d", each.Time)
    if each.Namespace == "Players[]" {
      rows = append(rows, fmt.Sprintf("Players[%.8X] = %x", each.Hash, each.Hash))
      playerIDs = append(playerIDs, fmt.Sprintf("Players[%.8X]", each.Hash))
    }
    fmt.Printf(" ─── \u001b[1mGape\u001b[0m %0.3fs.", float64(each.Time-next.Time)/1000000000)
    fmt.Printf(" ─── \u001b[1mNonce\u001b[0m %d", each.Nonce)
    fmt.Printf(" ─── \u001b[1mValid\u001b[0m %v\n", validate(newProof(each, Diff[each.Namespace])))
    decoded, _ := base64.StdEncoding.DecodeString(string(each.Data))
    fmt.Printf("  │ \u001b[1mData\u001b[0m %s\n", decoded)
    // fmt.Printf("  │ %x\n", each.Prev)
    // if each.Namespace != namespace { break }
    depth--
    if len(each.Prev) == 0 { break } else { fmt.Printf("  │ %x\n", each.Prev) }
    next = each
    // pow := NewProof(each)
  }
  // fmt.Println("  │ \n")
  fmt.Println(" ─┼───────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────\n")
  playerList := FindByPrefixes(chain, []byte("Session[00"))
  if namespace == "Players[]" {
    fmt.Println("    ─────────────────── Born blocks: ─────────────────────────────────────────────────────────────────────────────────────────────────────────")
    for _, link := range rows { fmt.Println(link) }
    for _, each := range playerList { fmt.Println(string(each)) }
    fmt.Println("    ──────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────")
  }
  fmt.Println()
  if len(playerIDs) == 0 { playerIDs = append(playerIDs, "Players[]") }
  return playerIDs
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
        playerList = append(playerList, []byte(fmt.Sprintf("%s = %x", key, v)))
        return nil
      })
      if err != nil { return err }
    }
    return nil
  })
  if len(playerList) == 0 { playerList = append(playerList, prefix) }
  return playerList
}

func AddPlayer(chain *BlockChain, player player.BasicStats) {
  // dummy := player.Player{}
  // dummy.ID = player.ID
  dataString := toJson(player)
  var lastHash []byte
  // run read only txn (connection query)
  err := chain.Database.View(func(txn *badger.Txn) error {
    item, err := txn.Get([]byte("Players[]"))
    if err != nil { fmt.Println(err) }
    if err == badger.ErrKeyNotFound {
      fmt.Println(err)
      fmt.Printf("Context \"Players\" does not exist! Genereating...")
      err = chain.Database.View(func(txn *badger.Txn) error {
        item, err := txn.Get([]byte("Initial"))
        if err != nil { fmt.Println(err) }
        lastHash, err = item.ValueCopy([]byte("Initial"))
        return err
      })
    } else {
      lastHash, err = item.ValueCopy([]byte("Players[]")) // here!
    }
    return err
  })
  if err != nil { fmt.Println(err) }
  addBlock(chain, dataString, lastHash, []byte("Players[]"))
}

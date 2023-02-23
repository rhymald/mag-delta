package blockchain

import(
  "fmt"
  "github.com/dgraph-io/badger"
)

var Diff map[string]int = map[string]int{
  "/": 4,
  "/Players": 3,
  "/NPC": 2,
  "/Session": 1,
}

// upodate context
func AddBlock(chain *BlockChain, data string, namespace string, behind []byte) []byte {
  chain.Lock()
  lastHash := chain.LastHash[namespace]
  epoch := chain.Epoch
  chain.Unlock()
  new := createBlock(data, namespace, lastHash, takeDiff(namespace), behind, epoch)
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
    return err
  })
  if err != nil { fmt.Println(err) }
  return new.Hash[:]
}

func takeDiff(ns string) int {
  for diff, _ := range Diff {
    trigger := diff == ns[:len(diff)]
    if trigger && diff != "/" {
      // fmt.Printf("\r    U-USED DIF-F-FICULTY: %s = %d\r", diff, Diff[diff])
      return Diff[diff]
    }
  }
  return Diff["/"]
}

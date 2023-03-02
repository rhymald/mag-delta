package blockchain

import(
  "fmt"
  "github.com/dgraph-io/badger"
  "time"
  // "math"
)

var Diff map[string]int = map[string]int{
  "/": 0,
  "/Players": 1,
  "/NPC": 2,
  "/Session": 3,
}

// upodate context
func AddBlock(chain *BlockChain, data string, namespace string, behind []byte) []byte {
  chain.Lock()
  lastHash := chain.LastHash[namespace]
  epoch := time.Now().UnixNano()-1317679200000000000-chain.Epoch
  chain.Unlock()
  new := createBlock(data, namespace, lastHash, takeDiff(namespace, epoch), behind, epoch)
  var prevData []byte
  err := chain.Database.View(func(txn *badger.Txn) error {
    item, err := txn.Get(lastHash)
    if err != nil { fmt.Println(err) }
    prevData, err = item.ValueCopy(lastHash)
    return err
  })
  if err != nil { fmt.Println(err) }
  prevBlock := Deserialize(prevData)
  if data == string(*&prevBlock.Data) { return []byte{} }
  err = chain.Database.Update(func(txn *badger.Txn) error {
    err := txn.Set(new.Hash, serialize(new))
    if err != nil { fmt.Println(err) }
    return err
  })
  if err != nil { fmt.Println(err) }
  return new.Hash[:]
}

func takeDiff(ns string, epoch int64) int {
  maxdiff := 16.0// math.Log2(float64(epoch)+1)/math.Log2(1000*math.Phi)
  for diff, _ := range Diff {
    // fmt.Println(diff, ns)
    trigger := false
    if len(ns)>len(diff) { trigger = diff == ns[:len(diff)] }
    if trigger && diff != "/" {
      // fmt.Printf("\r    U-USED DIF-F-FICULTY: %s = %d %+d\r", diff, int(maxdiff), -Diff[diff])
      return int(maxdiff)-Diff[diff]
    }
  }
  // fmt.Printf("\r    U-USED DIF-F-FICULTY: %s = %d %+d\r", "/", int(maxdiff), -Diff["/"])
  return int(maxdiff)-Diff["/"]
}

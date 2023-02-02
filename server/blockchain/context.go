package blockchain

import(
  "fmt"
  "rhymald/mag-delta/player"
  "github.com/dgraph-io/badger"
  "encoding/base64"
  "encoding/json"
  // "rhymald/mag-delta/funcs"
)

func AddPlayer(chain *BlockChain, player player.BasicStats) {
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

func toJson(thing player.BasicStats) string {
  b, err := json.Marshal(thing)
  if err != nil { fmt.Println(err) ; return "" }
  encoded := base64.StdEncoding.EncodeToString(b)
  return encoded
}

func fromJson(code string, thing player.BasicStats) player.BasicStats {
  copy := &thing
  decoded, _ := base64.StdEncoding.DecodeString(code)
  err := json.Unmarshal(decoded, copy)
  if err != nil { fmt.Println(err) ; return player.BasicStats{} }
  return *copy
}

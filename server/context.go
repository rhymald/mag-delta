package server

import(
  "fmt"
  "rhymald/mag-delta/player"
  "rhymald/mag-delta/server/blockchain"
  "github.com/dgraph-io/badger"
  "encoding/base64"
  "encoding/json"
  "crypto/sha512"
  "encoding/binary"

)

// func UpdPlayerStats(chain *blockchain.BlockChain, player player.BasicStats) {
//   dataString := toJson(player)
//   in_bytes := make([]byte, 8)
//   binary.LittleEndian.PutUint64(in_bytes, uint64(player.ID.Born))
//   hsum := sha512.Sum512(in_bytes)
//   id := fmt.Sprintf("/Players/%s", fmt.Sprintf("%.5X", hsum))
//   var lastHash []byte
//   // run read only txn (connection query)
//   err := chain.Database.View(func(txn *badger.Txn) error {
//     item, err := txn.Get([]byte(id))
//     if err != nil { fmt.Println(err) }
//     if err == badger.ErrKeyNotFound {
//       fmt.Printf("Such player does not existhere, update node! How did you log in, cheater?!")
//     } else {
//       lastHash, err = item.ValueCopy([]byte(id)) // here!
//     }
//     return err
//   })
//   if err != nil { fmt.Println(err) }
//   blockchain.AddBlock(chain, dataString)
// }
//
//
func AddPlayer(chain *blockchain.BlockChain, player player.BasicStats) {
  dataString := toJson(player)
  lasthash := []byte{}
  err := chain.Database.View(func(txn *badger.Txn) error {
    item, err := txn.Get([]byte("/Players"))
    if err == badger.ErrKeyNotFound {
      err = chain.Database.View(func(txn *badger.Txn) error {
        item, err := txn.Get([]byte("/"))
        chain.LastHash["/Players"], err = item.ValueCopy([]byte("/"))
        return err
      })
    } else {
      chain.LastHash["/Players"], err = item.ValueCopy([]byte("/Players"))
    }
    return err
  })
  if err != nil { fmt.Println(err) }
  lasthash = blockchain.AddBlock(chain, dataString, "/Players")
  if len(lasthash) == 0 {return}
  // propagate dependant contexts
  in_bytes := make([]byte, 8)
  binary.LittleEndian.PutUint64(in_bytes, uint64(player.ID.Born))
  hsum := sha512.Sum512(in_bytes)
  statsid := fmt.Sprintf("/Players/%s/Attributes", fmt.Sprintf("%.8X", hsum))
  stateid := fmt.Sprintf("/Players/%s/Session", fmt.Sprintf("%.8X", hsum))
  // creating subcontexts
  err = chain.Database.Update(func(txn *badger.Txn) error {
    err = txn.Set([]byte("/Players"), lasthash)
    chain.LastHash["/Players"] = lasthash
    if err != nil { fmt.Println(err) }
    err = txn.Set([]byte(statsid), lasthash)
    chain.LastHash[statsid] = lasthash
    if err != nil { fmt.Println(err) }
    err = txn.Set([]byte(stateid), lasthash)
    chain.LastHash[stateid] = lasthash
    if err != nil { fmt.Println(err) }
    return err
  })
  if err != nil { fmt.Println(err) }
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

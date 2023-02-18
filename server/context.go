package server

import(
  "fmt"
  "rhymald/mag-delta/player"
  "rhymald/mag-delta/server/blockchain"
  "github.com/dgraph-io/badger"
  "encoding/base64"
  "encoding/json"
)

func UpdPlayerStats(chain *blockchain.BlockChain, person player.Player) {
  dataString := toJson(person.Basics)
  id := player.GetID(person)
  statsid := fmt.Sprintf("/Players/%s", id)
  stateid := fmt.Sprintf("/Session/%s", id)
  lasthash := blockchain.AddBlock(chain, dataString, statsid)
  if len(lasthash) == 0 {return}
  err := chain.Database.Update(func(txn *badger.Txn) error {
    err := txn.Set([]byte(statsid), lasthash)
    chain.LastHash[statsid] = lasthash
    if err != nil { fmt.Println(err) }
    err = txn.Set([]byte(stateid), lasthash)
    chain.LastHash[stateid] = lasthash
    if err != nil { fmt.Println(err) }
    return err
  })
  if err != nil { fmt.Println(err) }
}

func AddPlayer(chain *blockchain.BlockChain, person player.Player) {
  dataString := toJson(person.Basics)
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
  id := player.GetID(person)
  statsid := fmt.Sprintf("/Players/%s", id)
  stateid := fmt.Sprintf("/Session/%s", id)
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

package server

import(
  "rhymald/mag-delta/server/blockchain"
  "github.com/dgraph-io/badger"
  "fmt"
  "rhymald/mag-delta/player"
)

func UpdPlayerStatE(chain *blockchain.BlockChain, person player.Player) {
  dataString := toJson(person.Status)
  pid, sid := player.GetID(person)
  statsid := fmt.Sprintf("/Players/%s", pid)
  stateid := fmt.Sprintf("/Session/%s/%s", pid, sid)
  anchor := fmt.Sprintf("/Session/%s", pid)
  err := chain.Database.View(func(txn *badger.Txn) error {
    _, err := txn.Get([]byte(stateid))
    if err == badger.ErrKeyNotFound {
      err = chain.Database.View(func(txn *badger.Txn) error {
        item, err := txn.Get([]byte(statsid))
        chain.LastHash[stateid], err = item.ValueCopy([]byte(statsid))
        return err
      })
    }
    return err
  })
  lasthash := blockchain.AddBlock(chain, dataString, stateid)
  if len(lasthash) == 0 {return}
  err = chain.Database.Update(func(txn *badger.Txn) error {
    err = txn.Set([]byte(anchor), lasthash)
    chain.LastHash[anchor] = lasthash
    if err != nil { fmt.Println(err) }
    err = txn.Set([]byte(stateid), lasthash)
    chain.LastHash[stateid] = lasthash
    if err != nil { fmt.Println(err) }
    return err
  })
  if err != nil { fmt.Println(err) }
}

func UpdPlayerStats(chain *blockchain.BlockChain, person player.Player) {
  dataString := toJson(person.Basics)
  pid, sid := player.GetID(person)
  stateid := fmt.Sprintf("/Session/%s/%s", pid, sid)
  statsid := fmt.Sprintf("/Players/%s", pid)
  anchor := fmt.Sprintf("/Session/%s", pid)
  lasthash := blockchain.AddBlock(chain, dataString, statsid)
  if len(lasthash) == 0 {return}
  err := chain.Database.Update(func(txn *badger.Txn) error {
    err := txn.Set([]byte(statsid), lasthash)
    chain.LastHash[statsid] = lasthash
    if err != nil { fmt.Println(err) }
    err = txn.Set([]byte(stateid), lasthash)
    chain.LastHash[stateid] = lasthash
    if err != nil { fmt.Println(err) }
    err = txn.Set([]byte(anchor), lasthash)
    chain.LastHash[anchor] = lasthash
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
  pid, sid := player.GetID(person)
  stateid := fmt.Sprintf("/Session/%s/%s", pid, sid)
  statsid := fmt.Sprintf("/Players/%s", pid)
  anchor := fmt.Sprintf("/Session/%s", pid)
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
    err = txn.Set([]byte(anchor), lasthash)
    chain.LastHash[anchor] = lasthash
    if err != nil { fmt.Println(err) }
    return err
  })
  if err != nil { fmt.Println(err) }
}

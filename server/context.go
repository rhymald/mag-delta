package server

import(
  "rhymald/mag-delta/server/blockchain"
  "github.com/dgraph-io/badger"
  "fmt"
  "rhymald/mag-delta/player"
)

func AssumePlayer(chain *blockchain.BlockChain, id string) player.Player {
  // supposed to use as login
  dummy := player.Player{}
  var statsJson []byte
  var stateJson []byte
  chain.Lock()
  statsat := chain.LastHash[fmt.Sprintf("/Players/%s", id)]
  stateat := chain.LastHash[fmt.Sprintf("/Session/%s", id)]
  chain.Unlock()
  err := chain.Database.View(func(txn *badger.Txn) error {
    item, err := txn.Get([]byte(statsat))
    statsJson, err = item.ValueCopy([]byte(statsat))
    return err
  })
  if err != nil { fmt.Println(err) } else { dummy.Basics = statsFromJson(string(statsJson), dummy.Basics) }
  err = chain.Database.View(func(txn *badger.Txn) error {
    item, err := txn.Get([]byte(stateat))
    stateJson, err = item.ValueCopy([]byte(stateat))
    return err
  })
  if err != nil { fmt.Println(err) } else { dummy.Status = stateFromJson(string(stateJson), dummy.Status) }
  return dummy
}

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
        chain.Lock()
        chain.LastHash[stateid], err = item.ValueCopy([]byte(statsid))
        chain.Unlock()
        return err
      })
    }
    return err
  })
  lasthash := blockchain.AddBlock(chain, dataString, stateid, []byte{})
  if len(lasthash) == 0 {return}
  err = chain.Database.Update(func(txn *badger.Txn) error {
    err = txn.Set([]byte(anchor), lasthash)
    if err != nil { fmt.Println(err) }
    err = txn.Set([]byte(stateid), lasthash)
    if err != nil { fmt.Println(err) }
    chain.Lock()
    chain.LastHash[anchor] = lasthash
    chain.LastHash[stateid] = lasthash
    chain.Unlock()
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
  chain.Lock()
  trigger := chain.LastHash[stateid]
  chain.Unlock()
  // ???
  lasthash := blockchain.AddBlock(chain, dataString, statsid, trigger)
  if len(lasthash) == 0 {return}
  err := chain.Database.Update(func(txn *badger.Txn) error {
    err := txn.Set([]byte(statsid), lasthash)
    if err != nil { fmt.Println(err) }
    err = txn.Set([]byte(stateid), lasthash)
    if err != nil { fmt.Println(err) }
    err = txn.Set([]byte(anchor), lasthash)
    if err != nil { fmt.Println(err) }
    chain.Lock()
    chain.LastHash[statsid] = lasthash
    chain.LastHash[stateid] = lasthash
    chain.LastHash[anchor] = lasthash
    chain.Unlock()
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
        chain.Lock()
        chain.LastHash["/Players"], err = item.ValueCopy([]byte("/"))
        chain.Unlock()
        return err
      })
    } else {
      chain.Lock()
      chain.LastHash["/Players"], err = item.ValueCopy([]byte("/Players"))
      chain.Unlock()
    }
    return err
  })
  if err != nil { fmt.Println(err) }
  lasthash = blockchain.AddBlock(chain, dataString, "/Players", []byte{})
  if len(lasthash) == 0 {return}
  pid, sid := player.GetID(person)
  stateid := fmt.Sprintf("/Session/%s/%s", pid, sid)
  statsid := fmt.Sprintf("/Players/%s", pid)
  anchor := fmt.Sprintf("/Session/%s", pid)
  // creating subcontexts
  err = chain.Database.Update(func(txn *badger.Txn) error {
    err = txn.Set([]byte("/Players"), lasthash)
    if err != nil { fmt.Println(err) }
    err = txn.Set([]byte(statsid), lasthash)
    if err != nil { fmt.Println(err) }
    err = txn.Set([]byte(stateid), lasthash)
    if err != nil { fmt.Println(err) }
    err = txn.Set([]byte(anchor), lasthash)
    if err != nil { fmt.Println(err) }
    chain.Lock()
    chain.LastHash["/Players"] = lasthash
    chain.LastHash[statsid] = lasthash
    chain.LastHash[stateid] = lasthash
    chain.LastHash[anchor] = lasthash
    chain.Unlock()
    return err
  })
  if err != nil { fmt.Println(err) }
}

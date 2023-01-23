package main

import (
  "fmt"
  "math"
  "rhymald/mag-delta/client/plot"
  "rhymald/mag-delta/server/blockchain"
  "rhymald/mag-delta/client"
  "rhymald/mag-delta/player"
  "rhymald/mag-delta/act"
  "os"
  "os/exec"
  "time"
)

const DBPath = "./cache"

var You player.Player
var Target player.Player
var Action string
var Keys chan string = make(chan string)
var StatChain *blockchain.BlockChain = blockchain.InitBlockChain(DBPath)

func init() {
  fmt.Println("\n\t\t  ", plot.Bar("Initializing...",8), "\n")
  player.PlayerBorn(&You,0)
  go func() { for { blockchain.AddBlock(StatChain, You) }}()
  client.PlayerStatus(You, Target)
  player.FoeSpawn(&Target,0)
  fmt.Println("\n\t\t", plot.Bar("Successfully login",1),"\n")
  client.PlayerStatus(You, Target)
  fmt.Println("\n\t     ",plot.Bar("Press [Enter] to continue",8),"\n")
  fmt.Scanln()
}

func main() {
  defer os.Exit(0)
  defer StatChain.Database.Close()
  plot.ShowMenu(" ")
  client.PlayerStatus(You, Target)
  grow := 1/math.Phi/math.Phi/math.Phi
  go func() {
    exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
    exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
    var b = make([]byte, 1)
    for {
      os.Stdin.Read(b)
      Keys <- string(b)
      plot.ShowMenu(string(b))
      if string(b) != "?" { client.PlayerStatus(You, Target) }
      time.Sleep( time.Millisecond * time.Duration( 128 ))
    }
  }()
  go func () {
    for {
      Action, _ := <-Keys
      switch Action {
        case "a":
          go func(){ act.Jinx(&You, &Target) }()
          Action = ""
        case "?":
          go func(){ time.Sleep( time.Millisecond * time.Duration( 128 )) ; blockchain.ListBlocks(StatChain) }()
          Action = ""
        default:
          time.Sleep( time.Millisecond * time.Duration( 128 ))
      }
      // if Action=="a" { go func(){ act.Jinx(&You, &Target) }() ; Action = "" }
      // if Action=="?" { go func(){ time.Sleep( time.Millisecond * time.Duration( 128 )) ; blockchain.ListBlocks(StatChain) }() ; Action = "" }
      if Target.Physical.Health.Current <= 0 { grow = grow*math.Cbrt(math.Phi) ; player.FoeSpawn(&Target, grow) ; plot.ShowMenu(Action)}// ; PlayerStatus(You, Target)}
      time.Sleep( time.Millisecond * time.Duration( 128 ))
    }
  }()
  for { time.Sleep( time.Millisecond * time.Duration( 4096 )) }
}

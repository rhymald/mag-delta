package main

import (
  "fmt"
  "math"
  "rhymald/mag-delta/client/plot"
  "rhymald/mag-delta/server/blockchain"
  // "rhymald/mag-delta/server"
  "rhymald/mag-delta/client"
  "rhymald/mag-delta/player"
  "rhymald/mag-delta/act"
  "rhymald/mag-delta/funcs"
  "os"
  "os/exec"
  "time"
)

const DBPath = "./cache"
const CacheAddr = "local"

var You player.Player
var Target player.Player
var StatChain *blockchain.BlockChain = blockchain.InitBlockChain(DBPath)

var Frame plot.LogFrame = plot.CleanFrame()
var Keys chan string = make(chan string)

func init() {
  fmt.Println("\n\t\t  ", plot.Bar("Initializing...",8), "\n")
  // player.PlayerBorn(&You,1024) ; blockchain.AddPlayer(StatChain, You)
  // player.PlayerBorn(&You,6) ; blockchain.AddPlayer(StatChain, You)
  player.PlayerBorn(&You,0,&Frame.Player) //; server.AddPlayer(StatChain, You.Basics)
  // go func() { for { server.UpdPlayerStats(StatChain, You.Basics) } }()
  client.PlayerStatus(You, Target)
  fmt.Println("\n\t\t", plot.Bar("Successfully login",1),"\n")
  player.FoeSpawn(&Target,0,&Frame.Foe)
  client.PlayerStatus(You, Target)
  fmt.Println("\n\t     ",plot.Bar("Press [Enter] to continue",8),"\n")
  fmt.Scanln()
}

func main() {
  defer os.Exit(0)
  defer StatChain.Database.Close()
  plot.ShowMenu(" ")
  grow := math.Cbrt(math.Phi)-1
  var b = make([]byte, 1)
  go func() {
    exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
    exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
    for {
      os.Stdin.Read(b)
      Keys <- string(b)
    }
  }()
  go func () {
    for {
      Action, _ := <-Keys
      switch Action {
      case "e":
          go func(){ act.Jinx(&You, &Target, &Frame) }()
          Action = " "
        default:
      }
      if Target.Status.Health <= 0 { player.PlayerEmpower(&You, 0) ; player.FoeSpawn(&Target, (funcs.Vector(You.Basics.Streams.Cre,You.Basics.Streams.Alt,You.Basics.Streams.Des)/math.Sqrt(3)-1)+grow, &Frame.Foe) }
    }
  }()
  for {
    plot.Clean()
    plot.ShowMenu(string(b))
    if string(b) != "?" {
      client.PlayerStatus(You, Target) ; plot.Frame(Frame)
      time.Sleep( time.Millisecond * time.Duration( 128 ))
    } else {
      _ = blockchain.ListBlocks(StatChain, "/")
      time.Sleep( time.Millisecond * time.Duration( 2048 ))
    }
  }
}

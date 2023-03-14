package main

import (
  "fmt"
  // "math"
  "rhymald/mag-delta/client/plot"
  "rhymald/mag-delta/server/blockchain"
  "rhymald/mag-delta/server"
  "rhymald/mag-delta/client"
  "rhymald/mag-delta/player"
  // "rhymald/mag-delta/balance"
  "rhymald/mag-delta/act"
  // "rhymald/mag-delta/funcs"
  "os"
  "os/exec"
  "time"
  "flag"
)

var (
  // arguments: 
  DBPath string //= "./cache"
  playerID string
  reborn bool = false
  // objects: 
  You player.Player
  Target player.Player
  StatChain *blockchain.BlockChain //= blockchain.InitBlockChain(DBPath)
  // CLI (TBDeprecated): 
  Frame plot.LogFrame = plot.CleanFrame()
  Keys chan string = make(chan string)
)

// read args, open or generate chain
func connect() string { //read app args and connect to db
  id := flag.String("p", "[no id defined]", "Player ID to login")
  flag.StringVar(&DBPath, "d", "cache", "Directory for cache")
  flag.BoolVar(&reborn, "n", false, "Create new player")
  help := flag.Bool("h", false, "Show this help")
  flag.Parse()
  if *help { fmt.Println("Application usage: keys") ; flag.PrintDefaults() ; fmt.Println(); os.Exit(1)}
  fmt.Println("\n\t\t", plot.Bar("  Initializing... ",0))
  StatChain = blockchain.InitBlockChain(DBPath)
  return *id
}

// get player from chain or create new
func init() {
  id := connect()
  fmt.Printf("Login with player ID: \u001b[1m")
  if reborn {fmt.Println(("[generating new...]"))} else {fmt.Println((id))}
  fmt.Printf("\u001b[0m")
  if reborn {
    playerID = player.PlayerBorn(&You, 0, 1, &Frame.Player)
    server.AddPlayer(StatChain, You)
    fmt.Println("\n\t", plot.Bar("  Successfully created new player  ",4))
  } else {
    if len(id) == 14 {
      You = server.AssumePlayer(StatChain, id, &Frame.Player)
      player.Live(&You, &Frame.Player)
      if You.Basics.ID.Born == 0 {
        fmt.Println("\n\t\t", plot.Bar("   Login failed   ",6))
        fmt.Println("\t\t", plot.Bar("  No such player  ",6))
        playerID = fmt.Sprintf("/Players")
      } else {
        fmt.Println("\n\t\t", plot.Bar("Successfully login",1))
        playerID = fmt.Sprintf("/Session/%s", id)
      }
    } else {
      fmt.Println("\n\t\t", plot.Bar(" Invalid playerId ",6))
      playerID = fmt.Sprintf("/Players")
    }
  }
  if You.Basics.ID.Born != 0 {
    go func() { for { server.UpdPlayerStats(StatChain, You) } }()
    go func() { for { server.UpdPlayerStatE(StatChain, You) } }()
    client.PlayerStatus(You)
    player.FoeSpawn(&Target, 1, &Frame.Foe)
    client.PlayerStatus(Target)
  }
  fmt.Println("\t     ",plot.Bar("Press [Enter] to continue",0))
  fmt.Scanln()
}

// interactive ui
func main() {
  defer os.Exit(0)
  defer StatChain.Database.Close()
  plot.ShowMenu(" ")
  var b = make([]byte, 1)
  go func() {
    exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
    exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
    for {
      os.Stdin.Read(b)
      Keys <- string(b)
    }
  }()
  go func() {
    for {
      Action, _ := <-Keys
      switch Action {
        case "e":
          go func(){ act.Fractal_Jinx(&You, &Target, &Frame) }()
          Action = " "
        case "+":
          go func(){ player.PlayerEmpower(&You, 0, &Frame.Player) }()
          Action = " "
          // player.FoeSpawn(&Target, 1024, &Frame.Foe)
        // default:
      }
    }
  }()
  for {
    plot.Clean()
    plot.ShowMenu(string(b))
    if string(b) == "/" {
      blockchain.ListBlocks(StatChain, playerID, true)
      time.Sleep( time.Millisecond * time.Duration( 2048 ))
    } else if string(b) == "?" || string(b) == "+" {
      blockchain.ListBlocks(StatChain, playerID, false)
      time.Sleep( time.Millisecond * time.Duration( 2048 ))
    } else {
      client.PlayerStatus(Target) ; plot.Frame(Frame) ; client.PlayerStatus(You)
      time.Sleep( time.Millisecond * time.Duration( 128 ))
    }
  }
}

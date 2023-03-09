package main

import (
  "fmt"
  "math"
  "rhymald/mag-delta/client/plot"
  "rhymald/mag-delta/server/blockchain"
  "rhymald/mag-delta/server"
  "rhymald/mag-delta/client"
  "rhymald/mag-delta/player"
  "rhymald/mag-delta/act"
  "rhymald/mag-delta/funcs"
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
  for x:=0 ;x<len(funcs.Physical); x++ { fmt.Printf(" x[%s]%.3f ", funcs.Physical[x], math.Pow(math.Log2(float64(x)+1), 2) ) } ; fmt.Println()
  for x:=0 ;x<3; x++ { fmt.Printf(" ^[%s]%.3f ", funcs.Elements[x], math.Pow(math.Sqrt(math.Log2(float64(x)+2))-1, 2)+1 ) } ; fmt.Println()
  id := flag.String("p", "[no id defined]", "Player ID to login")
  flag.StringVar(&DBPath, "d", "cache", "Directory for cache")
  flag.BoolVar(&reborn, "n", false, "Create new player")
  help := flag.Bool("h", false, "Show this help")
  flag.Parse()
  if *help { fmt.Println("Application usage: keys") ; flag.PrintDefaults() ; fmt.Println(); os.Exit(1)}
  fmt.Println("\n\t\t", plot.Bar("  Initializing... ",0), "\n")
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
    playerID = player.PlayerBorn(&You, 0, &Frame.Player)
    server.AddPlayer(StatChain, You)
    fmt.Println("\n\t", plot.Bar("  Successfully created new player  ",4),"\n")
  } else {
    if len(id) == 14 {
      You = server.AssumePlayer(StatChain, id, &Frame.Player)
      player.Live(&You, &Frame.Player)
      if You.Basics.ID.Born == 0 {
        fmt.Println("\n\t\t", plot.Bar("   Login failed   ",6))
        fmt.Println("\t\t", plot.Bar("  No such player  ",6),"\n")
        playerID = fmt.Sprintf("/Players")
      } else {
        fmt.Println("\n\t\t", plot.Bar("Successfully login",1),"\n")
        playerID = fmt.Sprintf("/Session/%s", id)
      }
    } else {
      fmt.Println("\n\t\t", plot.Bar(" Invalid playerId ",6),"\n")
      playerID = fmt.Sprintf("/Players")
    }
  }
  if You.Basics.ID.Born != 0 {
    go func() { for { server.UpdPlayerStats(StatChain, You) } }()
    go func() { for { server.UpdPlayerStatE(StatChain, You) } }()
    client.PlayerStatus(You)
    player.FoeSpawn(&Target,0,&Frame.Foe)
    client.PlayerStatus(You, Target)
  }
  fmt.Println("\n\t     ",plot.Bar("Press [Enter] to continue",0),"\n")
  fmt.Scanln()
}

// interactive ui
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
  go func() {
    for {
      Action, _ := <-Keys
      switch Action {
      case "e":
          go func(){ act.Jinx(&You, &Target, &Frame) }()
          Action = " "
        default:
      }
      if Target.Status.Health <= 0 { 
        player.PlayerEmpower(&You, 0, &Frame.Player) 
        _, stats := funcs.ReStr(You.Basics.Streams[0])
        player.FoeSpawn(&Target, (funcs.Vector(stats[0],stats[1],stats[2])/math.Sqrt(3)-1)+grow, &Frame.Foe)
      }
    }
  }()
  for {
    plot.Clean()
    plot.ShowMenu(string(b))
    if string(b) == "/" {
      blockchain.ListBlocks(StatChain, playerID, true)
      time.Sleep( time.Millisecond * time.Duration( 2048 ))
    } else if string(b) == "?" {
      blockchain.ListBlocks(StatChain, playerID, false)
      time.Sleep( time.Millisecond * time.Duration( 2048 ))
    } else {
      client.PlayerStatus(You, Target) ; plot.Frame(Frame)
      time.Sleep( time.Millisecond * time.Duration( 128 ))
    }
  }
}

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

// func debug() {
//   // player fetch
//   var ef player.Player 
//   var ad player.Player 
//   var bc player.Player 
//   _, _, _ = player.PlayerBorn(&ef, 3, &Frame.Player), player.PlayerBorn(&ad, 1, &Frame.Player), player.PlayerBorn(&bc, 2, &Frame.Player)
//   fmt.Println("Different:", player.Fetch_Stats(bc.Basics,ad.Basics))
//   fmt.Println("Same:     ", player.Fetch_Stats(bc.Basics,bc.Basics))
//   bc.Basics.ID.Born, ef.Basics.ID.Born = ad.Basics.ID.Born, ad.Basics.ID.Born  
//   gh := player.Fetch_Stats(bc.Basics,ad.Basics)
//   fmt.Println("Different 2:", gh) ; fmt.Println("Sum 2:      ", player.Grow_Stats(ef.Basics,gh))
//   player.TakeAll_Stats(&ef, []player.BasicStats{gh, gh, gh})
//   fmt.Println("Cascade 2:  ", ef)
//   // picker
//   fmt.Print(funcs.PickXFrom(17, 39)); fmt.Print(funcs.PickXFrom(4, 9)); fmt.Print(funcs.PickXFrom(17, 3)) ; fmt.Println(funcs.PickXFrom(4, 6))
//   // list all elements
//   for x:=0 ;x<len(funcs.Physical); x++ { fmt.Printf(" x[%s]%.3f ", funcs.Physical[x], math.Pow(math.Log2(float64(x)+1), 2) ) } ; fmt.Println()
//   for x:=0 ;x<3; x++ { fmt.Printf(" ^[%s]%.3f ", funcs.Elements[x], math.Pow(math.Sqrt(math.Log2(float64(x)+2))-1, 2)+1 ) } ; fmt.Println()
//   // streams count randomizer
//   a,b,c,d,e := 0,0,0,0,0
//   for x:=int64(0); x<1000; x++ { 
//     aaa := balance.BasicStats_StreamsCountAndModifier(funcs.Epoch())
//     if aaa == 2 {a++} else if aaa == 3 {b++} else if aaa == 4 {c++} else if aaa == 5 {d++} else {e++}
//   }
//   fmt.Println("    ||:",a,"\t|||:",b,"\t||||:",c,"\t|||||:",d,"\terr:",e)
// }

// read args, open or generate chain
func connect() string { //read app args and connect to db
  // debug()
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
        default:
      }
      if Target.Status.Health <= 0 { 
        player.PlayerEmpower(&You, 0, &Frame.Player) 
        player.FoeSpawn(&Target, 1024, &Frame.Foe)
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
      client.PlayerStatus(Target) ; plot.Frame(Frame) ; client.PlayerStatus(You)
      time.Sleep( time.Millisecond * time.Duration( 128 ))
    }
  }
}

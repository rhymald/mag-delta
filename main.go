package main

import (
  "fmt"
  "math"
  "time"
  "math/rand"
  "rhymald/mag-delta/plot"
  "rhymald/mag-delta/funcs"
  // "rhymald/mag-delta/balance"
  "rhymald/mag-delta/player"
  "os"
  "os/exec"
  // "github.com/rhymald/mag-gamma"
)

// type Player struct {
//   // Physical
//   Health struct {
//     Current float64
//     Max float64
//   }
//   // Energetical
//   Nature struct {
//     Resistance float64
//     Stream funcs.Stream
//     Pool struct {
//       Max float64
//       Dots []funcs.Dot
//     }
//   }
// }

var You player.Player
var Target player.Player
var Action string
var Keys chan string = make(chan string)


func init() {
  fmt.Println("\n\t\t", plot.Bar("  Initializing...  ",8), "\n")
  player.PlayerBorn(&You,100)
  player.FoeSpawn(&Target,0)

}

func main() {
  fmt.Println("\n    ",plot.Bar("Successfully login. Pres [Enter] to continue.",8),"\n")
  fmt.Scanln()
  plot.ShowMenu(" ")
  PlayerStatus(You, Target)
  UI(Keys)
  grow := 1.0
  for {
    // fmt.Print("\033[H\033[2J")
    // PlayerStatus(You, Target)
    // fmt.Printf("Do: \n")
    // fmt.Scanln(&Action)
    Action, _ := <-Keys
    key := Action
    if Action=="a" { Jinx(&You, &Target) ; Action = "" }
    if Target.Health.Current == 0 { grow+=grow*0.1+1 ; player.FoeSpawn(&Target, grow) ; plot.ShowMenu(key)}// ; PlayerStatus(You, Target)}
  }
}

// ███████████████████
// █▓▒░server side░▒▓█
// ███████████████████


// +Punch(Da) +Sting(Ad) - [physicals]
func Jinx(caster *player.Player, target *player.Player) {
  damage := 0.0
  dotsForConsume := int(*&caster.Nature.Pool.Max / (math.Pi + *&caster.Nature.Stream.Cre))
  pause := 1024 / float64(dotsForConsume)
  reach := 1024.0 // between
  dotCounter := 0
  for i:=0.0; i<float64(dotsForConsume); i+=1 {
    if len(*&caster.Nature.Pool.Dots) == 0 { break }// fmt.Printf("\n█▓▒░ DEBUG[Cast][Jinx]: Out of energy\n") ; break}
    _, w := MinusDot(&(*&caster.Nature.Pool.Dots))
    damage += w
    dotCounter++
    time.Sleep( time.Millisecond * time.Duration( pause ))
  }
  if CastFailed(dotsForConsume,dotCounter) {
    fmt.Printf("DEBUG[Cast][Jinx]: cast failed                                       \n") ; return
  } else {
    fmt.Printf("DEBUG[Cast][Jinx][From]: %0.1f damage sent                           \n", damage)
    go func(){
      time.Sleep( time.Millisecond * time.Duration( reach )) // immitation
      *&target.Health.Current += -damage*caster.Nature.Stream.Des/target.Nature.Resistance
      fmt.Printf("DEBUG[Cast][Jinx][ To ]: %0.1f damage received                       \n", damage*caster.Nature.Stream.Des/target.Nature.Resistance)
      if *&target.Health.Current < 0 { *&target.Health.Current = 0 }
    }()
  }
}

func CastFailed(need int, got int) bool { return funcs.Rand() >= math.Sqrt(float64(got)/float64(need)) }
func MinusDot(pool *[]funcs.Dot) (string, float64) {
  index := rand.New(rand.NewSource(time.Now().UnixNano())).Intn( len(*pool) )
  buffer := *pool
  ddelement := buffer[index].Element
  ddweight := buffer[index].Weight
  buffer[index] = buffer[len(buffer)-1]
  *pool = buffer[:len(buffer)-1]
  return ddelement, ddweight
}


// ███████████████████
// █▓▒░client side░▒▓█
// ███████████████████

func UI(Keys chan string) {
  // Keys := make(chan string)
  go func(Keys chan string) {
    // disable input buffering
    exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
    // do not display entered characters on the screen
    exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
    var b = make([]byte, 1)
    for {
      os.Stdin.Read(b)
      Keys <- string(b)
      plot.ShowMenu(string(b))
      PlayerStatus(You, Target)
      time.Sleep( time.Millisecond * time.Duration( 1 ))
    }
  }(Keys)
}

func PlayerStatus(players ...player.Player) {
  it, foe, compare := players[0], player.Player{}, len(players) > 1
  if players[1].Health.Current <= 0 { compare = false }
  if compare { foe = players[1] }
  playerTuple := [][]string{}
  fmt.Println(plot.Color("Player status",0),"[comparing to a foe]:")
  line := ""
  if compare {
    line = fmt.Sprintf("Health|Max: %0.0f|Current: %0.0f|Rate: %3.0f%%|[%3.0f%%]", it.Health.Max, it.Health.Current, 100*it.Health.Current/it.Health.Max,100*foe.Health.Current/foe.Health.Max)
  } else {
    line = fmt.Sprintf("Health|Max: %0.0f|Current: %0.0f|Rate: %1.0f%%", it.Health.Max, it.Health.Current, 100*it.Health.Current/it.Health.Max)
  }
  playerTuple = plot.AddRow(line, playerTuple)
  if compare {
    line = fmt.Sprintf(
      " \n %s \n[%s]|Creation\n  %0.3f \n [%0.3f]|Alteration\n  %0.3f \n [%0.3f]|Destruction\n  %0.3f \n [%0.3f]|Resistance\n  %0.3f \n [%0.3f]",
      it.Nature.Stream.Element,
      foe.Nature.Stream.Element,
      it.Nature.Stream.Cre,
      foe.Nature.Stream.Cre,
      it.Nature.Stream.Alt,
      foe.Nature.Stream.Alt,
      it.Nature.Stream.Des,
      foe.Nature.Stream.Des,
      it.Nature.Resistance,
      foe.Nature.Resistance,
    )
  } else {
    line = fmt.Sprintf(
      "Element\n%s|Creation\n%0.3f|Alteration\n%0.3f|Destruction\n%0.3f",
      it.Nature.Stream.Element,
      it.Nature.Stream.Cre,
      it.Nature.Stream.Alt,
      it.Nature.Stream.Des,
    )
  }
  playerTuple = plot.AddRow(line,playerTuple)
  if compare {
    line = fmt.Sprintf("Pool|Max: %0.0f|Current: %d|Rate: %1.0f%%|[%0.0f]", it.Nature.Pool.Max, len(it.Nature.Pool.Dots), 100*float64(len(it.Nature.Pool.Dots))/float64(it.Nature.Pool.Max), foe.Nature.Pool.Max )
  } else {
    line = fmt.Sprintf("Pool|Max: %0.0f|Current: %d|Rate: %1.0f%%", it.Nature.Pool.Max, len(it.Nature.Pool.Dots), 100*float64(len(it.Nature.Pool.Dots))/float64(it.Nature.Pool.Max) )
  }
  playerTuple = plot.AddRow(line,playerTuple)
  plot.Table(playerTuple, false)
  fmt.Println()
}

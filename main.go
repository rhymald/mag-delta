package main

import (
  "fmt"
  "math"
  "time"
  "math/rand"
  "rhymald/mag-delta/plot"
  "rhymald/mag-delta/funcs"
  "os"
  "os/exec"
  // "github.com/rhymald/mag-gamma"
)

type Player struct {
  // Physical
  Health struct {
    Current float64
    Max float64
  }
  // Energetical
  Nature struct {
    Resistance float64
    Stream funcs.Stream
    Pool struct {
      Max float64
      Dots []funcs.Dot
    }
  }
}

var You Player
var Target Player
var Action string
var Keys chan string = make(chan string)


func init() {
  plot.Bar("Initializing...    ")
  PlayerBorn(&You,0)
  FoeSpawn(&Target,0)
}

func main() {
  KeyListener(Keys)
  fmt.Println("[Go!..]")
  grow := 0.0
  for {
    // fmt.Print("\033[H\033[2J")
    PlayerStatus(You, Target)

    // fmt.Printf("Do: \n")
    fmt.Scanln(&Action)
    if Action=="a" { Jinx(&You, &Target) ; Action = "" }

    PlayerStatus(You, Target)
    if Target.Health.Current == 0 { grow += 1 ; FoeSpawn(&Target, grow) }
  }
}

// ███████████████████
// █▓▒░server side░▒▓█
// ███████████████████

func Jinx(caster *Player, target *Player) {
  damage := 0.0
  dotsForConsume := int(*&caster.Nature.Pool.Max / (math.Pi + *&caster.Nature.Stream.Cre))
  pause := 1024 / float64(dotsForConsume)
  reach := 1024.0 // between
  fmt.Printf("█▓▒░ DEBUG[Cast][Jinx]: %v needs %d dots \n", caster, dotsForConsume)
  for i:=0.0; i<float64(dotsForConsume); i+=1 {
    if len(*&caster.Nature.Pool.Dots) == 0 { fmt.Printf("█▓▒░ DEBUG[Cast][Jinx]: Out of energy\n") ; break}
    _, w := MinusDot(&(*&caster.Nature.Pool.Dots))
    damage += w
    time.Sleep( time.Millisecond * time.Duration( pause ))
  }
  fmt.Printf("█▓▒░ DEBUG[Cast][Jinx]: %0.1f damage sent\n", damage)
  go func(){
    time.Sleep( time.Millisecond * time.Duration( reach )) // immitation
    *&target.Health.Current += -damage*caster.Nature.Stream.Des/target.Nature.Resistance
    fmt.Printf("█▓▒░ DEBUG[Cast][Jinx]: %0.1f damage received\n", damage*caster.Nature.Stream.Des/target.Nature.Resistance)
    if *&target.Health.Current < 0 { *&target.Health.Current = 0 }
  }()
}
func MinusDot(pool *[]funcs.Dot) (string, float64) {
  index := rand.New(rand.NewSource(time.Now().UnixNano())).Intn( len(*pool) )
  buffer := *pool
  ddelement := buffer[index].Element
  ddweight := buffer[index].Weight
  buffer[index] = buffer[len(buffer)-1]
  *pool = buffer[:len(buffer)-1]
  return ddelement, ddweight
}

func PlayerBorn(player *Player, mean float64){
  mean += math.Sqrt(3)
  playerTuple := [][]string{}
  buffer := Player{}
  fmt.Println("Player creation start:")
  buffer.Health.Max = (mean/10+1)*(mean/10+1)*50 // from db
  buffer.Health.Current = math.Sqrt(buffer.Health.Max+1)-1 //from db
  current := fmt.Sprintf("Health|Max: %0.0f|Current: %0.0f|Rate: %1.0f%%", buffer.Health.Max, buffer.Health.Current, 100*buffer.Health.Current/buffer.Health.Max)
  playerTuple = plot.AddRow(current, playerTuple)
  buffer.Nature.Stream.Cre = 1+funcs.Rand()
  buffer.Nature.Stream.Alt = 1+funcs.Rand()
  buffer.Nature.Stream.Des = 1+funcs.Rand()
  stabilizer := mean/funcs.Vector(buffer.Nature.Stream.Cre, buffer.Nature.Stream.Alt, buffer.Nature.Stream.Des)
  buffer.Nature.Stream.Cre *= stabilizer
  buffer.Nature.Stream.Alt *= stabilizer
  buffer.Nature.Stream.Des *= stabilizer
  buffer.Nature.Stream.Element = "Common"
  buffer.Nature.Resistance = 3 / (1/buffer.Nature.Stream.Cre + 1/buffer.Nature.Stream.Alt + 1/buffer.Nature.Stream.Des)
  row := fmt.Sprintf(
    "Element\n%s|Creation\n%0.3f|Alteration\n%0.3f|Destruction\n%0.3f|Resistance\n%0.3f",
    // "Element\n%s|Creation\n%0.3f|Alteration\n%0.3f|Destruction\n%0.3f",
    buffer.Nature.Stream.Element,
    buffer.Nature.Stream.Cre,
    buffer.Nature.Stream.Alt,
    buffer.Nature.Stream.Des,
    buffer.Nature.Resistance,
  )
  playerTuple = plot.AddRow(row,playerTuple)
  thickness := math.Pi / ( 1/buffer.Nature.Stream.Des + 1/buffer.Nature.Stream.Alt + 1/buffer.Nature.Stream.Cre)
  buffer.Nature.Pool.Max = math.Sqrt( thickness *1024 + 1024) - 1
  playerTuple = plot.AddRow( fmt.Sprintf("Pool|Max: %0.0f|Current: %d|Rate: %1.0f%%", buffer.Nature.Pool.Max, len(buffer.Nature.Pool.Dots), 100*float64(len(buffer.Nature.Pool.Dots))/float64(buffer.Nature.Pool.Max) ) ,playerTuple)
  plot.Table(playerTuple, false)
  *player = buffer
  go func(){ Regeneration(&(*&player.Nature.Pool.Dots), &(*&player.Health.Current), *&player.Nature.Pool.Max, *&player.Health.Max, *&player.Nature.Stream) }()
}

func FoeSpawn(foe *Player, mean float64) {
  mean += math.Sqrt(3)
  playerTuple := [][]string{}
  buffer := Player{}
  fmt.Println("Foe spawning start:")
  buffer.Health.Max = (mean/10+1)*(mean/10+1)*50 // from db
  buffer.Health.Current = buffer.Health.Max //from db
  current := fmt.Sprintf("Health|||Rate: %1.0f%%", 100*buffer.Health.Current/buffer.Health.Max)
  playerTuple = plot.AddRow(current, playerTuple)
  buffer.Nature.Stream.Cre = 1+funcs.Rand()
  buffer.Nature.Stream.Alt = 1+funcs.Rand()
  buffer.Nature.Stream.Des = 1+funcs.Rand()
  stabilizer := mean/funcs.Vector(buffer.Nature.Stream.Cre, buffer.Nature.Stream.Alt, buffer.Nature.Stream.Des)
  buffer.Nature.Stream.Cre *= stabilizer
  buffer.Nature.Stream.Alt *= stabilizer
  buffer.Nature.Stream.Des *= stabilizer
  buffer.Nature.Stream.Element = "Common"
  buffer.Nature.Resistance = 3 / (1/buffer.Nature.Stream.Cre + 1/buffer.Nature.Stream.Alt + 1/buffer.Nature.Stream.Des)
  row := fmt.Sprintf(
    "Element\n%s|Creation\n%0.3f|Alteration\n%0.3f|Destruction\n%0.3f|Resistance\n%0.3f",
    buffer.Nature.Stream.Element,
    math.Sqrt(mean*mean/3),
    math.Sqrt(mean*mean/3),
    math.Sqrt(mean*mean/3),
    buffer.Nature.Resistance,
  )
  playerTuple = plot.AddRow(row,playerTuple)
  buffer.Nature.Pool.Max = math.Sqrt(buffer.Nature.Stream.Cre*1024 + 1024) - 1
  playerTuple = plot.AddRow( fmt.Sprintf("Pool|Max: %0.0f", buffer.Nature.Pool.Max ) ,playerTuple)
  plot.Table(playerTuple, false)
  *foe = buffer
  go func(){ Negeneration(&(*&foe.Health.Current), *&foe.Health.Max, *&foe.Nature.Stream) }()
}

func Regeneration(pool *[]funcs.Dot, health *float64, max float64, maxhp float64, stream funcs.Stream) {
  for {
    if max-float64(len(*pool))<1 { time.Sleep( time.Millisecond * time.Duration( 4096 )) } else {
      weight := math.Pow( math.Log2( 1+funcs.Vector(stream.Cre,stream.Des,stream.Alt) ), 2)
      dot := funcs.Dot{ Element: stream.Element, Weight: weight*(funcs.Rand()*0.5+0.75) }
      pause := 256.0
      heal := 1.0
      time.Sleep( time.Millisecond * time.Duration( pause ))
      //block
      *pool = append(*pool, dot )
      if *health <= 0 { fmt.Println("█▓▒░ FATAL[Player][Regeneration]: YOU ARE DEAD") ; break }
      if *health < maxhp { *health += heal } else { *health = maxhp }
      //unblock
    }
  }
}

func Negeneration(health *float64, maxhp float64, stream funcs.Stream) {
  for {
    pause := 256.0
    heal := 1.0
    //block
    time.Sleep( time.Millisecond * time.Duration( pause ))
    if *health <= 0 { fmt.Printf("█▓▒░ DEBUG[Fee][Regeneration]: Foe died\n") ; break }
    if *health < maxhp { *health += heal } else { *health = maxhp }
    //unblock
  }
}

// ███████████████████
// █▓▒░client side░▒▓█
// ███████████████████

func KeyListener(Keys chan string) {
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
    }
  }(Keys)
  for {
    stdin, _ := <-Keys
    fmt.Print("\033[H\033[2J")
    fmt.Println("█▓▒░ DEBUG[Keys pressed]:", stdin)
    PlayerStatus(You, Target)
    time.Sleep( time.Millisecond * time.Duration( 1 ))
  }
}

func PlayerStatus(players ...Player) {
  it, foe, compare := players[0], Player{}, len(players) > 1
  if players[1].Health.Current <= 0 { compare = false }
  if compare { foe = players[1] }
  playerTuple := [][]string{}
  fmt.Println("Player status [comparing to a foe]:")
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
}

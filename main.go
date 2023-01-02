package main

import (
  "fmt"
  "math"
  "time"
  "math/rand"
  "rhymald/mag-delta/plot"
  "rhymald/mag-delta/funcs"
)

type Player struct {
  // Physical
  Health struct {
    Current float64
    Max float64
  }
  // Energetical
  Nature struct {
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

func init() {
  fmt.Println("[Initializing...]")
  You = PlayerBorn(8)
  Target = FoeSpawn(13)
  go func(){
    go func(){ Regeneration(&You.Nature.Pool.Dots, &You.Health.Current, You.Nature.Pool.Max, You.Health.Max, You.Nature.Stream) }()
    go func(){ Negeneration(&Target.Health.Current, Target.Health.Max, Target.Nature.Stream) }()
  }()
}

func main() {
  fmt.Println("[Go!..]")
  for {
    fmt.Printf("Do: ")
    fmt.Scanln(&Action)
    fmt.Print("\033[H\033[2J")
    if Action=="a" {Jinx(&You, &Target) ; Action = ""}
    PlayerStatus(You, Target)
  }
}

func Jinx(caster *Player, target *Player) {
  damage := 0.0
  dotsForConsume := int(*&caster.Nature.Pool.Max / (math.Pi + *&caster.Nature.Stream.Cre))
  for i:=0.0; i<float64(dotsForConsume); i+=1 {
    if len(*&caster.Nature.Pool.Dots) == 0 {break}
    _, w := MinusDot(&(*&caster.Nature.Pool.Dots))
    damage += w + caster.Nature.Stream.Des
  }
  *&target.Health.Current += -damage
  if *&target.Health.Current < 0 { *&target.Health.Current = 0 }
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

func PlayerBorn(mean float64) Player {
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
  row := fmt.Sprintf(
    "Element\n%s|Creation\n%0.3f|Alteration\n%0.3f|Destruction\n%0.3f",
    buffer.Nature.Stream.Element,
    buffer.Nature.Stream.Cre,
    buffer.Nature.Stream.Alt,
    buffer.Nature.Stream.Des,
  )
  playerTuple = plot.AddRow(row,playerTuple)
  thickness := math.Pi / ( 1/buffer.Nature.Stream.Des + 1/buffer.Nature.Stream.Alt + 1/buffer.Nature.Stream.Cre)
  buffer.Nature.Pool.Max = math.Sqrt( thickness *1024 + 1024) - 1
  playerTuple = plot.AddRow( fmt.Sprintf("Pool|Max: %0.0f|Current: %d|Rate: %1.0f%%", buffer.Nature.Pool.Max, len(buffer.Nature.Pool.Dots), 100*float64(len(buffer.Nature.Pool.Dots))/float64(buffer.Nature.Pool.Max) ) ,playerTuple)
  plot.PlotTable(playerTuple, false)
  return buffer
}

func FoeSpawn(mean float64) Player {
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
  row := fmt.Sprintf(
    "Element\n%s|Creation\n%0.3f|Alteration\n%0.3f|Destruction\n%0.3f",
    buffer.Nature.Stream.Element,
    math.Sqrt(mean*mean/3),
    math.Sqrt(mean*mean/3),
    math.Sqrt(mean*mean/3),
  )
  playerTuple = plot.AddRow(row,playerTuple)
  buffer.Nature.Pool.Max = math.Sqrt(buffer.Nature.Stream.Cre*1024 + 1024) - 1
  playerTuple = plot.AddRow( fmt.Sprintf("Pool|Max: %0.0f", buffer.Nature.Pool.Max ) ,playerTuple)
  plot.PlotTable(playerTuple, false)
  return buffer
}

func PlayerStatus(players ...Player) {
  it, foe, compare := players[0], Player{}, len(players) > 1
  if compare { foe = players[1] }
  playerTuple := [][]string{}
  fmt.Println("Player status [comparing to a foe]:")
  line := ""
  if compare {
    line = fmt.Sprintf("Health|Max: %0.0f|Current: %0.0f|Rate: %3.0f%%\n  [%3.0f%%]", it.Health.Max, it.Health.Current, 100*it.Health.Current/it.Health.Max,100*foe.Health.Current/foe.Health.Max)
  } else {
    line = fmt.Sprintf("Health|Max: %0.0f|Current: %0.0f|Rate: %1.0f%%", it.Health.Max, it.Health.Current, 100*it.Health.Current/it.Health.Max)
  }
  playerTuple = plot.AddRow(line, playerTuple)
  if compare {
    line = fmt.Sprintf(
      " \n %s \n[%s]|Creation\n %0.3f \n[%0.3f]|Alteration\n %0.3f \n[%0.3f]|Destruction\n %0.3f \n[%0.3f]",
      it.Nature.Stream.Element,
      foe.Nature.Stream.Element,
      it.Nature.Stream.Cre,
      foe.Nature.Stream.Cre,
      it.Nature.Stream.Alt,
      foe.Nature.Stream.Alt,
      it.Nature.Stream.Des,
      foe.Nature.Stream.Des,
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
    line = fmt.Sprintf("Pool|Max: %0.0f\n  [%0.0f]|Current: %d|Rate: %1.0f%%", it.Nature.Pool.Max, foe.Nature.Pool.Max, len(it.Nature.Pool.Dots), 100*float64(len(it.Nature.Pool.Dots))/float64(it.Nature.Pool.Max) )
  } else {
    line = fmt.Sprintf("Pool|Max: %0.0f|Current: %d|Rate: %1.0f%%", it.Nature.Pool.Max, len(it.Nature.Pool.Dots), 100*float64(len(it.Nature.Pool.Dots))/float64(it.Nature.Pool.Max) )
  }
  playerTuple = plot.AddRow(line,playerTuple)
  plot.PlotTable(playerTuple, false)
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
      if *health <= 0 { fmt.Println("YOU ARE DEAD") ; break }
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
    if *health <= 0 { fmt.Printf("[Hint: Foe is DEAD] ") ; break }
    if *health < maxhp { *health += heal } else { *health = maxhp }
    //unblock
  }
}

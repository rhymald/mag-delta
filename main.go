package main

import (
  "fmt"
  "math"
  "math/rand"
  "crypto/sha512"
  "encoding/binary"
  "time"
  "strings"
)

type Dot struct {
  Weight float64
  Element float64
}

type Player struct {
  // Physical
  Health struct {
    Current float64
    Max float64
  }
  // Energetical
  Nature struct {
    Stream struct {
      Length float64
      Width float64
      Power float64
      Element string
    }
    Pool struct {
      Max float64
      Dots []Dot
    }
  }
}

var You Player
var Target Player

func init() {
  fmt.Println("[Initializing...]")
  You = PlayerBorn(1)
  return
}

func main() {
  fmt.Println("[Go!..]")
  FoeSpawn(4)
  return
}

func PlayerBorn(mean float64) Player {
  playerTuple := [][]string{}
  buffer := Player{}
  fmt.Println("Player creation start:")
  buffer.Health.Max = (mean/100+1)*(mean/100+1)*50 // from db
  buffer.Health.Current = math.Sqrt(buffer.Health.Max+1)-1 //from db
  // current := fmt.Sprintf("Health|Current: %0.0f|Max: %0.0f|Rate: %1.0f%%", buffer.Health.Current, buffer.Health.Max, 100*buffer.Health.Current/buffer.Health.Max)
  current := fmt.Sprintf("Health|Max: %0.0f|Current: %0.0f|Rate: %1.0f%%", buffer.Health.Max, buffer.Health.Current, 100*buffer.Health.Current/buffer.Health.Max)
  playerTuple = AddRow(current, playerTuple)
  buffer.Nature.Stream.Length  = 1+Rand()
  buffer.Nature.Stream.Width   = 1+Rand()
  buffer.Nature.Stream.Power   = 1+Rand()
  stabilizer := mean/Vector(buffer.Nature.Stream.Length, buffer.Nature.Stream.Width, buffer.Nature.Stream.Power)
  buffer.Nature.Stream.Length *= stabilizer
  buffer.Nature.Stream.Width  *= stabilizer
  buffer.Nature.Stream.Power  *= stabilizer
  buffer.Nature.Stream.Element = "Common"
  playerTuple = AddRow("Element|Creation|Alteration|Destruction",playerTuple)
  row := fmt.Sprintf(
    "%s|%0.3f|%0.3f|%0.3f",
    buffer.Nature.Stream.Element,
    buffer.Nature.Stream.Length,
    buffer.Nature.Stream.Width,
    buffer.Nature.Stream.Power,
  )
  playerTuple = AddRow(row,playerTuple)
  thickness := math.Pi / ( 1/buffer.Nature.Stream.Power + 1/buffer.Nature.Stream.Width + 1/buffer.Nature.Stream.Length)
  buffer.Nature.Pool.Max = math.Sqrt( thickness *1024 + 1024) - 1
  playerTuple = AddRow( fmt.Sprintf("Pool|Max: %0.0f|Current: %d|Rate: %1.0f%%", buffer.Nature.Pool.Max, len(buffer.Nature.Pool.Dots), 100*float64(len(buffer.Nature.Pool.Dots))/float64(buffer.Nature.Pool.Max) ) ,playerTuple)
  PlotTable(playerTuple)
  return buffer
}

func FoeSpawn(mean float64) Player {
  playerTuple := [][]string{}
  buffer := Player{}
  fmt.Println("Foe spawning start:")
  buffer.Health.Max = (mean/100+1)*(mean/100+1)*50 // from db
  buffer.Health.Current = buffer.Health.Max //from db
  current := fmt.Sprintf("Health|||Rate: %1.0f%%", 100*buffer.Health.Current/buffer.Health.Max)
  playerTuple = AddRow(current, playerTuple)
  buffer.Nature.Stream.Length  = 1+Rand()
  buffer.Nature.Stream.Width   = 1+Rand()
  buffer.Nature.Stream.Power   = 1+Rand()
  stabilizer := mean/Vector(buffer.Nature.Stream.Length, buffer.Nature.Stream.Width, buffer.Nature.Stream.Power)
  buffer.Nature.Stream.Length *= stabilizer
  buffer.Nature.Stream.Width  *= stabilizer
  buffer.Nature.Stream.Power  *= stabilizer
  buffer.Nature.Stream.Element = "Common"
  playerTuple = AddRow("Element|Creation|Alteration|Destruction",playerTuple)
  row := fmt.Sprintf(
    "%s|%0.3f|%0.3f|%0.3f",
    buffer.Nature.Stream.Element,
    buffer.Nature.Stream.Length,
    buffer.Nature.Stream.Width,
    buffer.Nature.Stream.Power,
  )
  playerTuple = AddRow(row,playerTuple)
  buffer.Nature.Pool.Max = math.Sqrt(buffer.Nature.Stream.Length*1024 + 1024) - 1
  playerTuple = AddRow( fmt.Sprintf("Pool|Max: %0.0f", buffer.Nature.Pool.Max ) ,playerTuple)
  PlotTable(playerTuple)
  return buffer
}

func Rand() float64 {
  x := (time.Now().UnixNano())
  in_bytes := make([]byte, 8)
  binary.LittleEndian.PutUint64(in_bytes, uint64(x))
  hsum := sha512.Sum512(in_bytes)
  sum  := binary.BigEndian.Uint64(hsum[:])
  return rand.New(rand.NewSource( int64(sum) )).Float64()
}
func Vector(props ...float64) float64 {
  sum := 0.0
  for _, each := range props { sum += each*each }
  return math.Sqrt(sum)
}

func PlotTable(tuple [][]string) {
  // tuple := tuple
  maxs := make([]int, len(tuple[0]))
  for j, y := range tuple {
    for i, _ := range y {
      if j == 0 {
        for c := 0; c>len(maxs); c++ { maxs[i] = 2+len(tuple[0][i]) }
      }
      maxs[i] = int(math.Max(2+float64(len(tuple[j][i])), float64(maxs[i])))
    }
  }
  //Head:
  fmt.Printf(" ╔")
  for i, wid := range maxs {
    for counter:=0 ;counter < wid; counter++ {
      fmt.Printf("═")
    }
    if i+1 == len(maxs) {fmt.Printf("╗\n")} else {fmt.Printf("╤")} //╤
  }
  //String:
  for I, _ := range tuple {
    fmt.Printf(" ║")
    for i, wid := range maxs {
      fmt.Printf(" ")
      fmt.Printf("%s", tuple[I][i])
      for counter:=0 ;counter < wid-1-len(tuple[I][i]); counter++ {
        fmt.Printf(" ")
      }
      if i+1 == len(maxs) {fmt.Printf("║\n")} else {
        if tuple[I][i+1]=="" || tuple[I][i+1]==" " {fmt.Printf("│")} else {fmt.Printf("│")}
      }
    }
    if I+1 == len(tuple) {
      //Footer:
      fmt.Printf(" ╚")
      for i, wid := range maxs {
        for counter:=0 ;counter < wid; counter++ {
          fmt.Printf("═")
        }
        if i+1 == len(maxs) {fmt.Printf("╝\n")} else {fmt.Printf("╧")} //╧
      }
    } else {
      //Delimiter:
      fmt.Printf(" ╟")
      for i, wid := range maxs {
        for counter:=0 ;counter < wid; counter++ {
          fmt.Printf("─")
        }
        if i+1 == len(maxs) {fmt.Printf("╢\n")} else {fmt.Printf("┼")} //┼
      }
    }
  }
}
func AddRow(row string, tuple [][]string) [][]string {
  buffer := strings.Split(row, "|")
  if len(tuple)==0 { return [][]string{buffer} }
  if len(buffer) > len(tuple[0]) {
    for l, _ := range tuple { for count := 0; count < len(buffer)-len(tuple[0]); count++ { tuple[l] = append(tuple[l], " ") } }
  } else if len(buffer) < len(tuple[0]) {
    for count := 0; count <= len(tuple[0])-len(buffer); count++ { buffer = append(buffer, " ") }
  }
  tuple = append(tuple, buffer)
  return tuple
}

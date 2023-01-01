package main

import (
  "fmt"
  "math"
  "math/rand"
  "crypto/sha512"
  "encoding/binary"
  "time"
  "strings"
  "golang.org/x/term"
)

type Dot struct {
  Weight float64
  Element string
}

type Stream struct {
  Cre float64
  Alt float64
  Des float64
  Element string
}

type Player struct {
  // Physical
  Health struct {
    Current float64
    Max float64
  }
  // Energetical
  Nature struct {
    Stream Stream
    Pool struct {
      Max float64
      Dots []Dot
    }
  }
}

var You Player
var Target Player
var Action string

func init() {
  fmt.Println("[Initializing...]")
  You = PlayerBorn(1)
  go func(){
    go func(){ Regeneration(&You.Nature.Pool.Dots, &You.Health.Current, You.Nature.Pool.Max, You.Health.Max, You.Nature.Stream) }()
  }()
  return
}

func main() {
  fmt.Println("[Go!..]")
  FoeSpawn(4)
  for {
    fmt.Scanf("Any move?..", &Action)
    PlayerStatus(You)
  }
  return
}

func PlayerBorn(mean float64) Player {
  playerTuple := [][]string{}
  buffer := Player{}
  fmt.Println("Player creation start:")
  buffer.Health.Max = (mean/10+1)*(mean/10+1)*50 // from db
  buffer.Health.Current = math.Sqrt(buffer.Health.Max+1)-1 //from db
  // current := fmt.Sprintf("Health|Current: %0.0f|Max: %0.0f|Rate: %1.0f%%", buffer.Health.Current, buffer.Health.Max, 100*buffer.Health.Current/buffer.Health.Max)
  current := fmt.Sprintf("Health|Max: %0.0f|Current: %0.0f|Rate: %1.0f%%", buffer.Health.Max, buffer.Health.Current, 100*buffer.Health.Current/buffer.Health.Max)
  playerTuple = AddRow(current, playerTuple)
  buffer.Nature.Stream.Cre  = 1+Rand()
  buffer.Nature.Stream.Alt   = 1+Rand()
  buffer.Nature.Stream.Des   = 1+Rand()
  stabilizer := mean/Vector(buffer.Nature.Stream.Cre, buffer.Nature.Stream.Alt, buffer.Nature.Stream.Des)
  buffer.Nature.Stream.Cre *= stabilizer
  buffer.Nature.Stream.Alt  *= stabilizer
  buffer.Nature.Stream.Des  *= stabilizer
  buffer.Nature.Stream.Element = "Common"
  // playerTuple = AddRow("Element|Creation|Alteration|Destruction",playerTuple)
  row := fmt.Sprintf(
    "Element\n%s|Creation\n%0.3f|Alteration\n%0.3f|Destruction\n%0.3f",
    buffer.Nature.Stream.Element,
    buffer.Nature.Stream.Cre,
    buffer.Nature.Stream.Alt,
    buffer.Nature.Stream.Des,
  )
  playerTuple = AddRow(row,playerTuple)
  thickness := math.Pi / ( 1/buffer.Nature.Stream.Des + 1/buffer.Nature.Stream.Alt + 1/buffer.Nature.Stream.Cre)
  buffer.Nature.Pool.Max = math.Sqrt( thickness *1024 + 1024) - 1
  playerTuple = AddRow( fmt.Sprintf("Pool|Max: %0.0f|Current: %d|Rate: %1.0f%%", buffer.Nature.Pool.Max, len(buffer.Nature.Pool.Dots), 100*float64(len(buffer.Nature.Pool.Dots))/float64(buffer.Nature.Pool.Max) ) ,playerTuple)
  PlotTable(playerTuple, false)
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
  buffer.Nature.Stream.Cre  = 1+Rand()
  buffer.Nature.Stream.Alt   = 1+Rand()
  buffer.Nature.Stream.Des   = 1+Rand()
  stabilizer := mean/Vector(buffer.Nature.Stream.Cre, buffer.Nature.Stream.Alt, buffer.Nature.Stream.Des)
  buffer.Nature.Stream.Cre *= stabilizer
  buffer.Nature.Stream.Alt  *= stabilizer
  buffer.Nature.Stream.Des  *= stabilizer
  buffer.Nature.Stream.Element = "Common"
  // playerTuple = AddRow("Element|Creation|Alteration|Destruction",playerTuple)
  row := fmt.Sprintf(
    "Element\n%s|Creation\n%0.3f|Alteration\n%0.3f|Destruction\n%0.3f",
    buffer.Nature.Stream.Element,
    math.Sqrt(mean*mean/3),
    math.Sqrt(mean*mean/3),
    math.Sqrt(mean*mean/3),
    // buffer.Nature.Stream.Length,
    // buffer.Nature.Stream.Width,
    // buffer.Nature.Stream.Power,
  )
  playerTuple = AddRow(row,playerTuple)
  buffer.Nature.Pool.Max = math.Sqrt(buffer.Nature.Stream.Cre*1024 + 1024) - 1
  playerTuple = AddRow( fmt.Sprintf("Pool|Max: %0.0f", buffer.Nature.Pool.Max ) ,playerTuple)
  PlotTable(playerTuple, false)
  return buffer
}

func PlayerStatus(it Player) {
  playerTuple := [][]string{}
  fmt.Println("Player status:")
  current := fmt.Sprintf("Health|Max: %0.0f|Current: %0.0f|Rate: %1.0f%%", it.Health.Max, it.Health.Current, 100*it.Health.Current/it.Health.Max)
  playerTuple = AddRow(current, playerTuple)
  // playerTuple = AddRow("Element|Creation|Alteration|Destruction",playerTuple)
  row := fmt.Sprintf(
    "Element\n%s|Creation\n%0.3f|Alteration\n%0.3f|Destruction\n%0.3f",
    it.Nature.Stream.Element,
    it.Nature.Stream.Cre,
    it.Nature.Stream.Alt,
    it.Nature.Stream.Des,
  )
  playerTuple = AddRow(row,playerTuple)
  playerTuple = AddRow( fmt.Sprintf("Pool|Max: %0.0f|Current: %d|Rate: %1.0f%%", it.Nature.Pool.Max, len(it.Nature.Pool.Dots), 100*float64(len(it.Nature.Pool.Dots))/float64(it.Nature.Pool.Max) ) ,playerTuple)
  PlotTable(playerTuple, false)
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

func PlotTable(tuple [][]string, stretch bool) {
  // tuple := tuple
  maxs := make([]int, len(tuple[0]))
  for j, y := range tuple {
    for i, _ := range y {
      if j == 0 {
        for c := 0; c>len(maxs); c++ { maxs[i] = 2+ MaxInCell(tuple[0][i]) }
      }
      maxs[i] = int(math.Max(2+float64( MaxInCell(tuple[j][i]) ), float64(maxs[i])))
    }
  }
  if stretch {
    for e, _ := range maxs { maxs[e] = int( math.Log2(float64(maxs[e])+2)/math.Log2(1.1459) ) }
    sums := 0
    for _, each := range maxs { sums+=each }
    termWigth, _, _ := term.GetSize(0)
    modificator := float64(termWigth - len(maxs) - 2) / float64(sums)
    for e, _ := range maxs { maxs[e] = int( float64(maxs[e])*modificator ) }
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
  for I, row := range tuple {
    PlotRow(row, maxs)
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
func FindDelim(row []string) ([]int, int) {
  buffer := make([]int, len(row))
  max := 0
  for i, cell := range row { buffer[i] = len(strings.Split(cell, "\n")) ; max = int(math.Max( float64(buffer[i]), float64(max) )) }
  return buffer, max
}
func PlotRow(row []string, widths []int) {
  _, max := FindDelim(row)
  for linenum:=0; linenum<max; linenum++ {
    fmt.Printf(" ║")
    for i, wid := range widths {
      fmt.Printf(" ")
      cell := strings.Split(row[i], "\n")
      if len(cell) < max { for count:=0; count<max-len(cell); count++ { cell = append(cell, string(" ")) } }
      toprint := cell[linenum]
      fmt.Printf("%s", toprint)
      for counter:=0 ;counter < wid-1-len(toprint); counter++ {
        fmt.Printf(" ")
      }
      if i+1 == len(widths) {fmt.Printf("║\n")} else {
        if row[i+1]=="" || row[i+1]==" " {fmt.Printf("│")} else {fmt.Printf("│")}
      }
    }
  }
}
func MaxInCell(cell string) int {
  lines, max := strings.Split(cell, "\n"), 0
  for _, each := range lines { max = int(math.Max( float64(max), float64(len(each)) )) }
  return max
}

func Regeneration(pool *[]Dot, health *float64, max float64, maxhp float64, stream Stream) {
  for {
    if max-float64(len(*pool))<1  { time.Sleep( time.Millisecond * time.Duration( 4096 )) ; return }
    weight := math.Pow( math.Log2( 1+Vector(stream.Cre,stream.Des,stream.Alt) ), 2)
    dot := Dot{ Element: stream.Element, Weight: weight }
    pause := 1024
    heal := 1.0
    time.Sleep( time.Millisecond * time.Duration( pause ))
    //block
    *pool = append(*pool, dot )
    if *health < maxhp { *health += heal }
    //unblock
  }
}

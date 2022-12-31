package main

import (
  "fmt"
  "math"
  "math/rand"
  "crypto/sha512"
  "encoding/binary"
  "time"
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
var Foe Player

func init() {
  fmt.Println("Initializing...")
  You = PlayerBorn()
  return
}

func main() {
  fmt.Println("Here we go!")
  return
}

func PlayerBorn() Player {
  buffer := Player{}
  fmt.Println("Player creation start,")
  buffer.Health.Max = 1000 // from db
  buffer.Health.Current = math.Sqrt(buffer.Health.Max+1) //from db
  fmt.Printf("\tHealth: current %0.1f/%0.1f max,\n", buffer.Health.Current, buffer.Health.Max)
  buffer.Nature.Stream.Length  = 1+Rand()
  buffer.Nature.Stream.Width   = 1+Rand()
  buffer.Nature.Stream.Power   = 1+Rand()
  buffer.Nature.Stream.Element = "Common"
  fmt.Printf(
    "\tStreams: %s  cre %0.3f  alt %0.3f  des %0.3f,\n",
    buffer.Nature.Stream.Element,
    buffer.Nature.Stream.Length,
    buffer.Nature.Stream.Width,
    buffer.Nature.Stream.Power,
  )
  buffer.Nature.Pool.Max = math.Sqrt(buffer.Nature.Stream.Length*1024 + 1024) - 1
  fmt.Printf("\tPool: current %d/%0.0f max,\n", len(buffer.Nature.Pool.Dots), buffer.Nature.Pool.Max)
  fmt.Println("Player created.")
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

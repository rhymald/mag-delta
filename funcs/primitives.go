package funcs

import (
  "math"
  "math/rand"
  "crypto/sha512"
  "encoding/binary"
  "time"
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

func LogF(n float64) float64 { return math.Log10(1+math.Abs(n))/math.Log10(math.Phi) }

func ChancedRound(a float64) int {
  b,l:=math.Ceil(a),math.Floor(a)
  c:=math.Abs(math.Abs(a)-math.Abs(math.Min(b, l)))
  if a<0 {c = 1-c}
  if Rand() < c {return int(b)} else {return int(l)}
  return 0
}

func ChancedRand(i int) float64 {
  counter, randy := 0, 0.0
  for {
    if counter >= i {break}
    randy += Rand()+Rand()
    counter++
  }
  return randy
}

func MeanStream(strs []Stream) Stream {
  mean := Stream{}
  for _, each := range strs { mean.Alt += 1/each.Alt ; mean.Cre += 1/each.Cre ; mean.Des += 1/each.Des }
  mean.Alt, mean.Cre, mean.Des = float64(len(strs))/mean.Alt, float64(len(strs))/mean.Cre, float64(len(strs))/mean.Des
  return mean
}

package balance

import (
  "math"
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
    Resistance float64
    Stream funcs.Stream
    Pool struct {
      Max float64
      Dots []funcs.Dot
    }
  }
}

func Regeneration_DefaultTimeout() float64 {return 1024*math.Pi}
func Regeneration_DotWeight_FromStream(stream funcs.Stream) funcs.Dot {return funcs.Dot{Element:stream.Element,Weight:(math.Pow(math.Log2(1+funcs.Vector(stream.Cre,stream.Des,stream.Alt)),2))*(funcs.Rand()*0.5+0.75)} }
func Regeneration_TimeoutMilliseconds_FromWeightPool(w float64, curr float64, max float64) float64 {return Regeneration_DefaultTimeout()/(math.Sqrt(max-curr+1)-1)}
func Regeneration_Heal_FromWeight(w float64) float64 {return math.Sqrt(w+1)-1}

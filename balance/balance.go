package balance

import (
  "math"
  "rhymald/mag-delta/funcs"
)

func BasicStats_MaxHP_FromNormale(norm float64) float64 {return (norm*32+16)}
func BasicStats_Resistance_FromStream(str funcs.Stream) float64 { return math.Pi/(1/str.Cre+1/str.Alt+1/str.Des) }
func BasicStats_MaxPool_FromStream(str funcs.Stream) float64 { return math.Sqrt(funcs.Vector(str.Cre,str.Alt,str.Des))*32 }
func BasicStats_Stream_FromNormaleWithElement(norm float64, element string) funcs.Stream {
  norm *= math.Sqrt(3)
  cre, alt, des := 5+funcs.Rand()+funcs.Rand(), 5+funcs.Rand()+funcs.Rand(), 5+funcs.Rand()+funcs.Rand()
  stabilizer := norm/funcs.Vector(cre,alt,des)
  cre, alt, des = cre*stabilizer, alt*stabilizer, des*stabilizer
  return funcs.Stream{Cre: cre, Alt: alt, Des: des, Element: element}
}

func Regeneration_TimeoutMilliseconds_FromWeightPool(w float64, curr float64, max float64) float64 {return Regeneration_DefaultTimeout()/(math.Sqrt(max-curr+1)-1)}
func Regeneration_Heal_FromWeight(w float64) float64 {return math.Log2(w+2)-1 }
func Regeneration_DefaultTimeout() float64 {return 1024*math.Pi}
func Regeneration_DotWeight_FromStream(stream funcs.Stream) funcs.Dot {
  w := math.Pow(math.Log2(1+funcs.Vector(stream.Cre,stream.Des,stream.Alt)),2)
  return funcs.Dot{Element:stream.Element,Weight:w/math.Pi*(funcs.Rand()*0.5+0.75)}
}

func Cast_Common_Failed(need int, got int) bool { return funcs.Rand() >= math.Sqrt(float64(got)/float64(need)) }
func Cast_Common_TimePerString(str funcs.Stream) float64 { return Regeneration_DefaultTimeout()/funcs.LogF(str.Alt) }
func Cast_Common_ExecutionRapidity(str funcs.Stream) float64 { return 1/math.Phi+funcs.LogF(str.Des)/math.Pi }
func Cast_Common_DotsPerString(str funcs.Stream) int { return funcs.ChancedRound(funcs.LogF(funcs.Vector(str.Cre,str.Alt,str.Des)*1024)) }

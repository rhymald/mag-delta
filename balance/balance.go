package balance

import (
  "math"
  "rhymald/mag-delta/funcs"
)

func ReStr(stream funcs.Stream) (string, [3]float64) { return funcs.ReStr(stream) } 

func BasicStats_MaxHP_FromBody(body funcs.Stream) float64 { _, stats := ReStr(body) ; return (math.Pi/(1/stats[0]+1/stats[1]+1/stats[2]))*64+16 }
func BasicStats_Resistance_FromStream(str funcs.Stream) float64 { _, stats := ReStr(str) ; return math.Pi/(1/stats[0]+1/stats[1]+1/stats[2]) }
func BasicStats_MaxPool_FromStream(str funcs.Stream) float64 { _, stats := ReStr(str) ; return math.Sqrt(funcs.Vector(stats[0],stats[1],stats[2]))*32 }
func BasicStats_Stream_FromNormaleWithElement(norm float64, element string) funcs.Stream {
  norm *= math.Sqrt(3)
  cre, alt, des := 5+funcs.Rand()+funcs.Rand(), 5+funcs.Rand()+funcs.Rand(), 5+funcs.Rand()+funcs.Rand()
  stabilizer := norm/funcs.Vector(cre,alt,des)
  cre, alt, des = cre*stabilizer, alt*stabilizer, des*stabilizer
  return funcs.Stream{ element : [3]int{ funcs.FloorRound(cre*1000), funcs.FloorRound(alt*1000), funcs.FloorRound(des*1000) }}
}

func Regeneration_TimeoutMilliseconds_FromWeightPool(w int, curr float64, max float64) float64 { return Regeneration_DefaultTimeout()/math.Sqrt( 2/( 1/(max-curr+1) + 1/(curr+1) )) }
func Regeneration_Heal_FromBody(body funcs.Stream) float64 { _, stats := ReStr(body) ; return math.Sqrt(math.Log10(1+funcs.Vector(stats[0],stats[1],stats[2])))*(funcs.Rand()*0.2+0.9) }
func Regeneration_DefaultTimeout() float64 {return 1024*math.Pi}
func Regeneration_DotWeight_FromStream(stream funcs.Stream) funcs.Dot {
  elem, stats := ReStr(stream)
  w := funcs.Log(math.Cbrt(stats[0]*stats[1]*stats[2]))
  return map[string]int{ elem : funcs.ChancedRound(w*(funcs.Rand()*0.2+0.9)) }
}

// Here!
func Cast_Common_Failed(need int, got int) bool { return funcs.Rand() >= math.Sqrt(float64(got)/float64(need)) }
func Cast_Common_TimePerString(str funcs.Stream) float64 { _, stats := ReStr(str) ; return Regeneration_DefaultTimeout()/math.Log2(stats[1]+1) }
func Cast_Common_ExecutionRapidity(str funcs.Stream) float64 { _, stats := ReStr(str) ; return math.Log10(stats[2]+10) }
func Cast_Common_DotsPerString(str funcs.Stream) int { _, stats := ReStr(str) ; return funcs.ChancedRound(10* math.Sqrt(math.Log2(1+funcs.Vector(stats[0],stats[1],stats[2]))) ) }

func StreamStructure2(a float64, b float64, c float64, t float64) bool { if a > b && b*math.Sqrt(t) > c && a/b > 1 && a/b < t { return true } ; return false }
func StreamStructure3(a float64, b float64, c float64, t float64) bool { if ( StreamStructure2(a,b,c,t) || StreamStructure2(a,c,b,t) ) && math.Max(math.Max(a/b,a/c),b/c)<math.Cbrt(t)*math.Cbrt(t) && math.Max(b/c,c/b) < math.Sqrt(t) { return true } ; return false }
func StreamAffinity2(a float64, b float64, t float64) float64 { return math.Pow(math.Log2(t/(a/b))/math.Log2(t), 2) }
func StreamAffinity3(a float64, b float64, c float64, t float64) float64 { ab, ca := math.Max(a,b)/math.Min(a,b), math.Max(a,c)/math.Min(a,c) ; return math.Pow(math.Log2(t/(2/(1/ab+1/ca)))/math.Log2(t), 2)}
func StreamAbilities_FromStream(str funcs.Stream) map[string]float64 {
  rate := math.Phi // resonating coefficient: bigger = more effect, - must be >1
  buffer := make(map[string]float64)
  _, stats := ReStr(str) 
  // Antibarrier (enchantment, poisoned weapon) = +AddDamage, +ticks, - if D>C close to each other
  if StreamStructure2(stats[2],stats[0],stats[1],rate) { buffer["Dc"] = StreamAffinity2(stats[2],stats[0],rate) }
  // Permanent debuff (hard to clean, need restore, not just cancel - canceling is stopping it) = +Speed, +effectiveness, - if D>A close to each other
  if StreamStructure2(stats[2],stats[1],stats[0],rate) { buffer["Da"] = StreamAffinity2(stats[2],stats[1],rate) }
  // Pulsing damage = +efectiveness, +damage, +speed, - if D>(A=C) when ac close to each other
  if StreamStructure3(stats[2],stats[1],stats[0],rate) { buffer["Dac"] = StreamAffinity3(stats[2],stats[1],stats[0],rate) ; buffer["Dca"] = buffer["Dac"] }
  // Smooth damaging conditions (easy to clean) = +time, +damage : A>D
  if StreamStructure2(stats[1],stats[2],stats[0],rate) { buffer["Ad"] = StreamAffinity2(stats[1],stats[2],rate) }
  // Smooth buff (easy to rip-off) = +time, +edfectiveness : A>C
  if StreamStructure2(stats[1],stats[0],stats[2],rate) { buffer["Ac"] = StreamAffinity2(stats[1],stats[0],rate) }
  // Permanent buff trigger = +effectiveness, +chance, +speed : A>(D=C)
  if StreamStructure3(stats[1],stats[0],stats[2],rate) { buffer["Adc"] = StreamAffinity3(stats[1],stats[2],stats[0],rate) ; buffer["Acd"] = buffer["Adc"] }
  // Shield (barrier) = +amount, +time : C>D
  if StreamStructure2(stats[0],stats[2],stats[1],rate) { buffer["Cd"] = StreamAffinity2(stats[0],stats[2],rate) }
  // Heal recovery, restoration = +efectiveness, +speed : C>A
  if StreamStructure2(stats[0],stats[1],stats[2],rate) { buffer["Ca"] = StreamAffinity2(stats[0],stats[1],rate) }
  // Conjuration local shadows, wells, self-regenerating energy shields = +volume, +activity, +efectiveness : C>(A=D)
  if StreamStructure3(stats[0],stats[2],stats[1],rate) { buffer["Cad"] = StreamAffinity3(stats[0],stats[1],stats[2],rate) ; buffer["Cda"] = buffer["Cad"] }
  return buffer
} 
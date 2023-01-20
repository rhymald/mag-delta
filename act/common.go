package act

import (
  "rhymald/mag-delta/funcs"
  "rhymald/mag-delta/balance"
  "rhymald/mag-delta/player"
  "math/rand"
  // "math"
  "time"
  "fmt"
)

// +Punch(Da) +Sting(Ad) - [physicals]
func Jinx(caster *player.Player, target *player.Player) {
  if *&caster.Busy { fmt.Printf("DEBUG[Cast][Jinx]: busy ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\n") ; return }
  dotsForConsume := balance.Cast_Common_DotsPerString(caster.Nature.Stream) //Cre
  pause := 1/float64(dotsForConsume) * balance.Cast_Common_TimePerString(caster.Nature.Stream) //Alt
  reach := 1024.0 / balance.Cast_Common_ExecutionRapidity(caster.Nature.Stream) // Des
  damage := 0.0
  dotCounter := 0
  *&caster.Busy = true
  for i:=0; i<dotsForConsume; i++ {
    if len(*&caster.Nature.Pool.Dots) == 0 { break }
    _, w := MinusDot(&(*&caster.Nature.Pool.Dots))
    damage += w
    dotCounter++
    time.Sleep( time.Millisecond * time.Duration( pause ))
  }
  *&caster.Busy = false
  if balance.Cast_Common_Failed(dotsForConsume,dotCounter) {
    fmt.Printf("DEBUG[Cast][Jinx]: cast failed ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\n") ; return
  } else {
    fmt.Printf("DEBUG[Cast][Jinx][From]: %0.1f damage sent for %.0f ms░░░░░░░░░░░░░░░░░░░░░░░░░\n", damage, pause*float64(dotsForConsume))
    go func(){
      time.Sleep( time.Millisecond * time.Duration( reach )) // immitation
      *&target.Physical.Health.Current += -damage*(caster.Nature.Stream.Des/target.Nature.Resistance)
      fmt.Printf("DEBUG[Cast][Jinx][ To ]: %0.1f damage received after %.0f ms ░░░░░░░░░░░░░░░░░░░\n", damage*caster.Nature.Stream.Des/target.Nature.Resistance, reach)
      if *&target.Physical.Health.Current < 0 { *&target.Physical.Health.Current = 0 }
    }()
  }
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

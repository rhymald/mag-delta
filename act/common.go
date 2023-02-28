package act

import (
  "rhymald/mag-delta/funcs"
  "rhymald/mag-delta/balance"
  "rhymald/mag-delta/player"
  "rhymald/mag-delta/client/plot"
  "math/rand"
  // "math"
  "time"
  "fmt"
)

// +Punch(Da) +Sting(Ad) - [physicals]
func Jinx(caster *player.Player, target *player.Player, logs *plot.LogFrame) {
  castId := time.Now().UnixNano() % 1000
  if *&caster.Attributes.Busy { plot.AddAction(logs, fmt.Sprintf("CAST [Jinx %3d] Fail: player is busy", castId)) ; return }
  if !(*&caster.Attributes.Login) || !(*&target.Attributes.Login) { plot.AddAction(logs, fmt.Sprintf("CAST [Jinx %3d] Fail: no player / no target", castId)) ; return }
  dotsForConsume := balance.Cast_Common_DotsPerString(caster.Basics.Streams) //Cre
  pause := 1/float64(dotsForConsume) * balance.Cast_Common_TimePerString(caster.Basics.Streams) //Alt
  reach := 1024.0 / balance.Cast_Common_ExecutionRapidity(caster.Basics.Streams) // Des
  damage := 0.0
  dotCounter := 0
  plot.AddAction(logs, fmt.Sprintf("CAST [Jinx %3d]:      start casting of %d dots ", castId, dotsForConsume))
  *&caster.Attributes.Busy = true
  for i:=0; i<dotsForConsume; i++ {
    if len(*&caster.Status.Pool) == 0 { break }
    _, w := MinusDot(&(*&caster.Status.Pool))
    damage += w
    dotCounter++
    time.Sleep( time.Millisecond * time.Duration( pause ))
  }
  *&caster.Attributes.Busy = false
  if balance.Cast_Common_Failed(dotsForConsume,dotCounter) {
    plot.AddAction(logs, fmt.Sprintf("CAST [Jinx %3d] Fail: cast of %d/%d dots failed ", castId, dotCounter, dotsForConsume)) ; return
  } else {
    plot.AddAction(logs, fmt.Sprintf("CAST [Jinx %3d] From: %0.1f damage as %d sent for %.0f ms", castId, damage, dotsForConsume, pause*float64(dotsForConsume)))
    go func(){
      time.Sleep( time.Millisecond * time.Duration( reach )) // immitation
      *&target.Status.Health += -damage*(caster.Basics.Streams.Des/(target.Attributes.Resistances[caster.Basics.Streams.Element]))
      plot.AddAction(logs, fmt.Sprintf("CAST [Jinx %3d] To:   %0.1f damage received after %.0f ms ", castId, damage*(caster.Basics.Streams.Des/(target.Attributes.Resistances[caster.Basics.Streams.Element])), reach))
      if *&target.Status.Health < 0 { *&target.Status.Health = 0 }
    }()
  }
}

// + MinusStamina
func MinusDot(pool *[]funcs.Dot) (string, float64) {
  index := rand.New(rand.NewSource(time.Now().UnixNano())).Intn( len(*pool) )
  buffer := *pool
  ddelement := buffer[index].Element
  ddweight := buffer[index].Weight
  buffer[index] = buffer[len(buffer)-1]
  *pool = buffer[:len(buffer)-1]
  return ddelement, ddweight
}

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
  action := funcs.Action{ Time: (time.Now().UnixNano()/1000000) , Kind: "Jinx" }
  castId := fmt.Sprintf("%s#%3d", action.Kind, action.Time%1000)
  if *&caster.Attributes.Busy { plot.AddAction(logs, fmt.Sprintf("%s Fail: player is busy", castId)) ; return }
  if !(*&caster.Attributes.Login) || !(*&target.Attributes.Login) { plot.AddAction(logs, fmt.Sprintf("%s Fail: no player / no target", castId)) ; return }
  dotsForConsume := balance.Cast_Common_DotsPerString(caster.Basics.Streams) //Cre
  pause := 1/float64(dotsForConsume) * balance.Cast_Common_TimePerString(caster.Basics.Streams) //Alt
  reach := 1024.0 / balance.Cast_Common_ExecutionRapidity(caster.Basics.Streams) // Des
  damage := 0.0
  dotCounter := 0
  plot.AddAction(logs, fmt.Sprintf("%s:      start casting of %d dots ", castId, dotsForConsume))
  *&caster.Attributes.Busy = true
  action.By = append(action.By, 0)
  for i:=0; i<dotsForConsume; i++ {
    if len(*&caster.Status.Pool) == 0 { break }
    _, w, index := MinusDot(&(*&caster.Status.Pool))
    damage += w
    dotCounter++
    action.With = append(action.With, index)
    time.Sleep( time.Millisecond * time.Duration( pause ))
  }
  *&caster.Attributes.Busy = false
  // actions logging for anticheat? 
  affection := action
  cpid, csid := player.GetID(*caster)
  tpid, tsid := player.GetID(*target)
  action.From = fmt.Sprintf("%s/%s", cpid, csid) ; action.To = fmt.Sprintf("%s/%s", tpid, tsid)
  affection.To = action.From ; affection.From = action.To
  *&caster.Status.ActionLog = append(*&caster.Status.ActionLog, action)
  *&target.Status.ActionLog = append(*&target.Status.ActionLog, affection)
  // maybe move it to spawnchain, friendly fire and self leave in playchain
  if balance.Cast_Common_Failed(dotsForConsume,dotCounter) {
    plot.AddAction(logs, fmt.Sprintf("%s Fail: cast of %d/%d dots failed ", castId, dotCounter, dotsForConsume)) ; return
  } else {
    plot.AddAction(logs, fmt.Sprintf("%s From: %0.1f damage as %d sent for %.0f ms", castId, damage, dotsForConsume, pause*float64(dotsForConsume)))
    go func(){
      time.Sleep( time.Millisecond * time.Duration( reach )) // immitation
      *&target.Status.Health += -damage*(caster.Basics.Streams.Des/(target.Attributes.Resistances[caster.Basics.Streams.Element]))
      plot.AddAction(logs, fmt.Sprintf("%s To:   %0.1f damage received after %.0f ms ", castId, damage*(caster.Basics.Streams.Des/(target.Attributes.Resistances[caster.Basics.Streams.Element])), reach))
      if *&target.Status.Health < 0 { *&target.Status.Health = 0 }
      // +exp?
    }()
  }
}

// + MinusStamina
func MinusDot(pool *[]funcs.Dot) (string, float64, int) {
  index := rand.New(rand.NewSource(time.Now().UnixNano())).Intn( len(*pool) )
  buffer := *pool
  ddelement := buffer[index].Element
  ddweight := buffer[index].Weight
  buffer[index] = buffer[len(buffer)-1]
  *pool = buffer[:len(buffer)-1]
  return ddelement, ddweight, index
}

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

// +Punch(Da) +Sting(Ad) +Block(C) - [physicals]
func Jinx(caster *player.Player, target *player.Player, logs *plot.LogFrame) {
  action := funcs.Action{ Time: (funcs.Epoch()/1000000) , Kind: "Jinx" }
  castId := fmt.Sprintf("%s#%3d", action.Kind, action.Time%1000)
  if *&caster.Attributes.Busy { plot.AddAction(logs, fmt.Sprintf("%s Fail: player is busy", castId)) ; return }
  if !(*&caster.Attributes.Login) || !(*&target.Attributes.Login) { plot.AddAction(logs, fmt.Sprintf("%s Fail: no player / no target", castId)) ; return }

  minstreams := balance.BasicStats_StreamsCountAndModifier(caster.Basics.ID.Born) 
  picked := funcs.PickXFrom(minstreams, len(caster.Basics.Streams))
  reach := 0.0 // Des
  totalDotsNeeded := 0
  damage := 0.0
  dotCounter := 0
  totalpause := 0.0
  
  *&caster.Attributes.Busy = true
  for _, each := range picked {
    dotsForConsume := balance.Cast_Common_DotsPerString(caster.Basics.Streams[each], minstreams) //Cre
    reach += 1024.0 / balance.Cast_Common_ExecutionRapidity(caster.Basics.Streams[each]) / float64(minstreams) // Des
    totalDotsNeeded += dotsForConsume
    pause := 1/float64(dotsForConsume) * balance.Cast_Common_TimePerString(caster.Basics.Streams[each]) // alt
    totalpause += pause * float64(dotsForConsume)
    plot.AddAction(logs, fmt.Sprintf("%s:      stream #%d demands %d dots ", castId, each, dotsForConsume))
    action.By = append(action.By, 0) // only 1 stream yet
    for i:=0; i<dotsForConsume; i++ {
      if len(*&caster.Status.Pool) == 0 { break }
      _, w, index := MinusDot(&(*&caster.Status.Pool))
      damage += float64(w)
      dotCounter++
      action.With = append(action.With, index)
      time.Sleep( time.Millisecond * time.Duration( pause ))
    }
  }
  *&caster.Attributes.Busy = false

  // actions logging for anticheat? 
  // TBRefactored for less duplication
  affection := action
  cpid, csid := player.GetID(*caster)
  tpid, tsid := player.GetID(*target)
  action.From = fmt.Sprintf("%s/%s", cpid, csid) ; action.To = fmt.Sprintf("%s/%s", tpid, tsid)
  affection.To = action.From ; affection.From = action.To
  *&caster.Status.ActionLog = append(*&caster.Status.ActionLog, action)
  *&target.Status.ActionLog = append(*&target.Status.ActionLog, affection)
  // maybe move it to spawnchain, friendly fire and self leave in playchain
  if balance.Cast_Common_Failed(totalDotsNeeded, dotCounter) {
    plot.AddAction(logs, fmt.Sprintf("%s Fail: cast of %d/%d dots failed ", castId, dotCounter, totalDotsNeeded)) ; return
  } else {
    plot.AddAction(logs, fmt.Sprintf("%s From: %0.1f damage as %d sent for %.0f ms", castId, damage, totalDotsNeeded, totalpause))
    // calculate effectiveness here
    go func(){
      elem, stats := funcs.ReStr(caster.Basics.Streams[0])
      time.Sleep( time.Millisecond * time.Duration( reach )) // immitation
      *&target.Status.Health += funcs.ChancedRound( -damage * (stats[2]+caster.Attributes.Resistances[funcs.Elements[0]]) / (target.Attributes.Resistances[elem]+target.Attributes.Resistances[funcs.Elements[0]]) * 1000/target.Attributes.Vitality)
      plot.AddAction(logs, fmt.Sprintf("%s To:   %0.1f damage received after %.0f ms ", castId, damage*((stats[2]+caster.Attributes.Resistances[funcs.Elements[0]])/(target.Attributes.Resistances[elem]+target.Attributes.Resistances[funcs.Elements[0]])), reach))
      if *&target.Status.Health < 0 { *&target.Status.Health = 0 }
      // +exp?
    }()
  }
}

// + MinusStamina
func MinusDot(pool *[]funcs.Dot) (string, int, int) {
  index := rand.New(rand.NewSource(time.Now().UnixNano())).Intn( len(*pool) )
  buffer := *pool
  ddelement, ddweight := "", 0 
  for elem, weig := range buffer[index] { ddelement = elem ; ddweight = weig }
  buffer[index] = buffer[len(buffer)-1]
  *pool = buffer[:len(buffer)-1]
  return ddelement, ddweight, index
}

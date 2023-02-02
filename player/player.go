package player

import (
  "rhymald/mag-delta/funcs"
  "rhymald/mag-delta/balance"
  "math"
  "fmt"
  _ "encoding/json"
  "time"
)

// need restruct: separate block of hp+mp+exp
type Player struct {
  Basics BasicStats `json:"Basics"`
  // For BlockChain
  // ID struct {
  //   NPC bool `json:"NPC"`
  //   Name string `json:"Name,omitempty"`
  //   Date int64 `json:"Date,omitempty"`
  //   Last int64 `json:"Last,omitempty"`
  // } `json:"ID"`
  // For actions
  Busy bool `json:"Busy"`
  // Physical
  Physical struct {
    Health struct {
      Current float64 `json:"Current"`
      Max float64 `json:"Max"`
    } `json:"Health"`
    // Body funcs.Stream `json:"Body"`
  } `json:"Physical"`
  // Energetical
  Nature struct {
    Resistance float64 `json:"Resistances,omitempty"`
    // Stream funcs.Stream `json:"Stream"`
    Pool struct {
      Max float64 `json:"Max"`
      Dots []funcs.Dot `json:"Dots,omitempty"`
    } `json:"Pool"`
  } `json:"Nature"`
}

type BasicStats struct {
  ID struct {
    NPC bool `json:"NPC"`
    Name string `json:"Name,omitempty"`
    Date int64 `json:"Date,omitempty"`
    Last int64 `json:"Last,omitempty"`
  } `json:"ID,omitempty"`
  Body funcs.Stream `json:"Body"`
  Streams funcs.Stream `json:"Streams"`
  Items []funcs.Stream `json:"Items,omitempty"`
} // ^ stored in stats/spawn chains
type CharStatus struct { // +heats +exp +consumables
  XYZ [3]float64 `json:"XYZ"`
  Health float64 `json:"Health"`
  Barrier map[string]float64 `json:"Barrier,omitempty"`
  Dots []funcs.Dot `json:"Dots,omitempty"`
} // ^ stored in status chain
type CharAttributes struct { // +states
  Busy bool `json:"Busy,omitempty"`
  Vitality float64 `json:"Vitality,omitempty"`
  Resistances map[string]float64 `json:"Resistances,omitempty"`
  Poolsize float64 `json:"Poolsize,omitempty"`
} // ^ calculated when login

func PlayerEmpower(player *Player, mean float64){ // immitation
  buffer := *player
  buffer.Basics.ID.Last = time.Now().UnixNano()
  *player = buffer
}

func PlayerBorn(player *Player, mean float64){
  buffer := Player{}
  buffer.Basics.ID.NPC = false
  buffer.Basics.ID.Date = time.Now().UnixNano()
  buffer.Basics.ID.Last = time.Now().UnixNano()
  buffer.Basics.Body = balance.BasicStats_Stream_FromNormaleWithElement(2, "Physical")
  buffer.Physical.Health.Max = balance.BasicStats_MaxHP_FromBody(buffer.Basics.Body) // from db
  buffer.Physical.Health.Current = math.Sqrt(buffer.Physical.Health.Max+1)-1 //from db
  buffer.Basics.Streams = balance.BasicStats_Stream_FromNormaleWithElement(1+mean, "Common")
  buffer.Nature.Resistance = balance.BasicStats_Resistance_FromStream(buffer.Basics.Streams)
  buffer.Nature.Pool.Max = balance.BasicStats_MaxPool_FromStream(buffer.Basics.Streams)
  *player = buffer
  go func(){ Regeneration(&(*&player.Nature.Pool.Dots), &(*&player.Physical.Health.Current), *&player.Nature.Pool.Max, *&player.Physical.Health.Max, *&player.Basics.Streams, *&player.Basics.Body) }()
}

func FoeSpawn(foe *Player, mean float64) {
  buffer := Player{}
  buffer.Basics.ID.NPC = true
  buffer.Basics.ID.Date = time.Now().UnixNano()
  buffer.Basics.ID.Last = time.Now().UnixNano()
  buffer.Basics.Body = balance.BasicStats_Stream_FromNormaleWithElement(2, "Physical")
  buffer.Physical.Health.Max = balance.BasicStats_MaxHP_FromBody(buffer.Basics.Body) // from db
  buffer.Physical.Health.Current = buffer.Physical.Health.Max / math.Sqrt2 //from db
  buffer.Basics.Streams = balance.BasicStats_Stream_FromNormaleWithElement(1+mean, "Common")
  buffer.Nature.Resistance = balance.BasicStats_Resistance_FromStream(buffer.Basics.Streams)
  buffer.Nature.Pool.Max = balance.BasicStats_MaxPool_FromStream(buffer.Basics.Streams)
  *foe = buffer
  go func(){ Negeneration(&(*&foe.Physical.Health.Current), *&foe.Physical.Health.Max, *&foe.Nature.Pool.Max, *&foe.Basics.Body) }()
}

func Regeneration(pool *[]funcs.Dot, health *float64, max float64, maxhp float64, stream funcs.Stream, body funcs.Stream) {
  for {
    if max-float64(len(*pool))<1 { time.Sleep( time.Millisecond * time.Duration( balance.Regeneration_DefaultTimeout() )) } else {
      dot := balance.Regeneration_DotWeight_FromStream(stream)
      pause := balance.Regeneration_TimeoutMilliseconds_FromWeightPool(dot.Weight, float64(len(*pool)), max)
      heal := balance.Regeneration_Heal_FromBody(body)
      time.Sleep( time.Millisecond * time.Duration( pause ))
      //block
      if *health >= maxhp {
        fmt.Printf("DEBUG[Player][Regeneration]: ░░░░░░░░░ for %0.3fs +%s %0.3f'e ░░░░░░░░░░░░░░░░░░\r", pause/1000, dot.Element, dot.Weight)
      } else {
        fmt.Printf("DEBUG[Player][Regeneration]: %+0.3f'hp for %0.3fs +%s %0.3f'e ░░░░░░░░░░░░░░░░░░\r", heal, pause/1000, dot.Element, dot.Weight)
      }
      *pool = append(*pool, dot )
      if *health <= 0 { fmt.Printf("DEBUG[Player][Regeneration]: ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ You are Died ░░░░░░░░░\n") ; break }
      if *health < maxhp { *health += heal } else { *health = maxhp }
      //unblock
    }
  }
}

func Negeneration(health *float64, maxhp float64, maxe float64, body funcs.Stream) {
  for {
    if maxhp<=*health { time.Sleep( time.Millisecond * time.Duration( balance.Regeneration_DefaultTimeout() )) } else {
      dot := balance.Regeneration_DotWeight_FromStream(body)
      pause := balance.Regeneration_TimeoutMilliseconds_FromWeightPool(dot.Weight, funcs.Log(maxe), maxe)
      heal := balance.Regeneration_Heal_FromBody(body)
      time.Sleep( time.Millisecond * time.Duration( pause ))
      //block
      if *health < maxhp { fmt.Printf("DEBUG[ NPC  ][Regeneration]: %+0.3f'hp for %0.3fs ░░░░░░░░░░░░░░░░░░░░░░░░░\r", heal, pause/1000) }
      if *health <= 0 { fmt.Printf("DEBUG[ NPC  ][Regeneration]: Foe died ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\n") ; break }
      if *health < maxhp { *health += heal } else { *health = maxhp }
      //unblock
    }
  }
}

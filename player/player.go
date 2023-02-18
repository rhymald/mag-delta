package player

import (
  "rhymald/mag-delta/funcs"
  "rhymald/mag-delta/balance"
  "math"
  "fmt"
  _ "encoding/json"
  "time"
  "crypto/sha512"
  "encoding/binary"
)

// need restruct: separate block of hp+mp+exp
type Player struct {
  Basics BasicStats `json:"Basics"` // from stats chain
  Status CharStatus `json:"Status"` // from state chain
  Attributes CharAttributes `json:"Attributes,omitepmty"`
}

type BasicStats struct {
  ID struct {
    NPC bool `json:"NPC"`
    Name string `json:"Name,omitempty"`
    Born int64 `json:"Born,omitempty"`
    Last int64 `json:"Last,omitempty"`
  } `json:"ID,omitempty"`
  Body funcs.Stream `json:"Body"`
  Streams funcs.Stream `json:"Streams"`
  Items []funcs.Stream `json:"Items,omitempty"`
} // ^ stored in stats/spawn chains
type CharStatus struct { // +heats +exp +consumables
  Health float64 `json:"Health"`
  Pool []funcs.Dot `json:"Pool,omitempty"`
  Focus []string `json:"Focus,omitempty"` // 0 for target, 1+ other
  XYZ [3]float64 `json:"XYZ"` // not used yet
  Barrier map[string]float64 `json:"Barrier,omitempty"` // max hp mods per element
} // ^ stored in status chain
type CharAttributes struct { // +states
  Busy bool `json:"Busy,omitempty"`
  Vitality float64 `json:"Vitality,omitempty"`
  Resistances map[string]float64 `json:"Resistances,omitempty"`
  Poolsize float64 `json:"Poolsize,omitempty"`
} // ^ calculated when login

func PlayerEmpower(player *Player, mean float64){ // immitation
  buffer := *player
  buffer.Basics.ID.Last = funcs.Epoch()
  *player = buffer
}

func GetID(player Player) string {
  in_bytes := make([]byte, 8)
  binary.LittleEndian.PutUint64(in_bytes, uint64(player.Basics.ID.Born))
  id := sha512.Sum512(in_bytes)
  streang := fmt.Sprintf("%X", id)
  final := fmt.Sprintf("%v-%v-%v-%v", streang[:4], streang[4:13], streang[13:14], streang[14:21])
  return final
}

func CalculateAttributes_FromBasics(player *Player){
  buffer := *player
  buffer.Attributes.Vitality = balance.BasicStats_MaxHP_FromBody(buffer.Basics.Body) // from db
  resists := make(map[string]float64)
  resists[buffer.Basics.Streams.Element] = balance.BasicStats_Resistance_FromStream(buffer.Basics.Streams)
  buffer.Attributes.Resistances = resists
  buffer.Attributes.Poolsize = balance.BasicStats_MaxPool_FromStream(buffer.Basics.Streams)
  *player = buffer
}

func PlayerBorn(player *Player, mean float64, logger *string){
  buffer := Player{}
  buffer.Basics.ID.NPC = false
  buffer.Basics.ID.Born = funcs.Epoch()
  buffer.Basics.ID.Last = funcs.Epoch()
  buffer.Basics.Body = balance.BasicStats_Stream_FromNormaleWithElement(2, "Physical")
  buffer.Basics.Streams = balance.BasicStats_Stream_FromNormaleWithElement(1+mean, "Common")
  CalculateAttributes_FromBasics(&buffer)
  buffer.Status.Health = math.Sqrt(buffer.Attributes.Vitality+1)-1 //from db
  *player = buffer
  go func(){ Regeneration(&(*&player.Status.Pool), &(*&player.Status.Health), *&player.Attributes.Poolsize, *&player.Attributes.Vitality, *&player.Basics.Streams, *&player.Basics.Body, logger) }()
}

func FoeSpawn(foe *Player, mean float64, logger *string) { // old, new+ template Stream{}
  buffer := Player{}
  buffer.Basics.ID.NPC = true
  buffer.Basics.ID.Born = funcs.Epoch()
  buffer.Basics.ID.Last = funcs.Epoch()
  buffer.Basics.Body = balance.BasicStats_Stream_FromNormaleWithElement(2, "Physical")
  buffer.Basics.Streams = balance.BasicStats_Stream_FromNormaleWithElement(1+mean, "Common")
  CalculateAttributes_FromBasics(&buffer)
  buffer.Status.Health = buffer.Attributes.Vitality / math.Sqrt2
  *foe = buffer
  go func(){ Negeneration(&(*&foe.Status.Health), *&foe.Attributes.Vitality, *&foe.Attributes.Poolsize, *&foe.Basics.Body, logger) }()
}

func Regeneration(pool *[]funcs.Dot, health *float64, max float64, maxhp float64, stream funcs.Stream, body funcs.Stream, logger *string) {
  for {
    if max-float64(len(*pool))<1 { time.Sleep( time.Millisecond * time.Duration( balance.Regeneration_DefaultTimeout() )) } else {
      dot := balance.Regeneration_DotWeight_FromStream(stream)
      pause := balance.Regeneration_TimeoutMilliseconds_FromWeightPool(dot.Weight, float64(len(*pool)), max)
      heal := balance.Regeneration_Heal_FromBody(body)
      time.Sleep( time.Millisecond * time.Duration( pause ))
      //block
      if *health >= maxhp {
        *logger = fmt.Sprintf("          for %0.3fs +%s %0.3f'e ", pause/1000, dot.Element, dot.Weight)
      } else {
        *logger = fmt.Sprintf("%+0.3f'hp for %0.3fs +%s %0.3f'e ", heal, pause/1000, dot.Element, dot.Weight)
      }
      *pool = append(*pool, dot )
      if *health <= 0 { *logger = fmt.Sprintf("You are Died") ; break }
      if *health < maxhp { *health += heal } else { *health = maxhp }
      //unblock
    }
  }
}

func Negeneration(health *float64, maxhp float64, maxe float64, body funcs.Stream, logger *string) {
  for {
    if maxhp<=*health { time.Sleep( time.Millisecond * time.Duration( balance.Regeneration_DefaultTimeout() )) } else {
      dot := balance.Regeneration_DotWeight_FromStream(body)
      pause := balance.Regeneration_TimeoutMilliseconds_FromWeightPool(dot.Weight, funcs.Log(maxe), maxe)
      heal := balance.Regeneration_Heal_FromBody(body)
      time.Sleep( time.Millisecond * time.Duration( pause ))
      //block
      if *health < maxhp { *logger = fmt.Sprintf("Dummy: %+0.3f'hp for %0.3fs ", heal, pause/1000) }
      if *health <= 0 { *logger = fmt.Sprintf("Dummy: Foe died ") ; break }
      if *health < maxhp { *health += heal } else { *health = maxhp }
      //unblock
    }
  }
}

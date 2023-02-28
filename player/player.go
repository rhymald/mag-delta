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
  XYZ [3]float64 `json:"XYZ"` // not used yet
  Health float64 `json:"Health"`
  Pool []funcs.Dot `json:"Pool,omitempty"`
  InFocus []string `json:"InFocus,omitempty"` // 0 for target, 1+ other
  Barrier map[string]float64 `json:"Barrier,omitempty"` // max hp mods per element
} // ^ stored in status chain
type CharAttributes struct { // +states
  Login bool `json:"Login,omitempty"`
  Busy bool `json:"Busy,omitempty"`
  Vitality float64 `json:"Vitality,omitempty"`
  Resistances map[string]float64 `json:"Resistances,omitempty"`
  Poolsize float64 `json:"Poolsize,omitempty"`
} // ^ calculated when login

func PlayerEmpower(player *Player, mean float64, logger *string){ // immitation
  *&player.Attributes.Login = false
  time.Sleep( time.Millisecond * time.Duration( balance.Regeneration_DefaultTimeout() ))
  *&player.Basics.ID.Last = funcs.Epoch()
  CalculateAttributes_FromBasics(player)
  Live(player, logger)
}

func Live(player *Player, logger *string) {
  *&player.Attributes.Login = true
  if player.Basics.ID.NPC {
    go func(){ Negeneration(&(*&player.Status.Health), *&player.Attributes.Login, *&player.Attributes.Vitality, *&player.Attributes.Poolsize, *&player.Basics.Body, logger) }()
  } else {
    go func(){ Regeneration(&(*&player.Status.Pool), &(*&player.Status.Health), *&player.Attributes.Login, *&player.Attributes.Poolsize, *&player.Attributes.Vitality, *&player.Basics.Streams, *&player.Basics.Body, logger) }()
  }
  // + start calm
  // + start stamina
}

func GetID(player Player) (string, string) {
  in_bytes := make([]byte, 8)
  binary.LittleEndian.PutUint64(in_bytes, uint64(player.Basics.ID.Born))
  pid := fmt.Sprintf("%X", sha512.Sum512(in_bytes))
  binary.LittleEndian.PutUint64(in_bytes, uint64(player.Basics.ID.Last))
  sid := fmt.Sprintf("%X", sha512.Sum512(in_bytes))
  pstring := fmt.Sprintf("%X", pid)
  sstring := fmt.Sprintf("%X", sid)
  pfinal := fmt.Sprintf("%v-%v", pstring[:4], pstring[119:128])
  sfinal := fmt.Sprintf("%v-%v", sstring[:1], sstring[121:128])
  return pfinal, sfinal
}

func CalculateAttributes_FromBasics(player *Player){
  buffer := *player
  buffer.Attributes.Login = false
  buffer.Attributes.Vitality = balance.BasicStats_MaxHP_FromBody(buffer.Basics.Body) // from db
  resists := make(map[string]float64)
  resists[buffer.Basics.Streams.Element] = balance.BasicStats_Resistance_FromStream(buffer.Basics.Streams)
  buffer.Attributes.Resistances = resists
  buffer.Attributes.Poolsize = balance.BasicStats_MaxPool_FromStream(buffer.Basics.Streams)
  *player = buffer
}

func PlayerBorn(player *Player, mean float64, logger *string) string {
  buffer := Player{}
  buffer.Basics.ID.NPC = false
  buffer.Basics.ID.Born = funcs.Epoch()
  buffer.Basics.ID.Last = funcs.Epoch()
  buffer.Basics.Body = balance.BasicStats_Stream_FromNormaleWithElement(2, "Physical")
  buffer.Basics.Streams = balance.BasicStats_Stream_FromNormaleWithElement(1+mean, "Common")
  CalculateAttributes_FromBasics(&buffer)
  buffer.Status.Health = math.Sqrt(buffer.Attributes.Vitality+1)-1 //from db
  *player = buffer
  Live(player, logger)
  // go func(){ Regeneration(&(*&player.Status.Pool), &(*&player.Status.Health), *&player.Attributes.Poolsize, *&player.Attributes.Vitality, *&player.Basics.Streams, *&player.Basics.Body, logger) }()
  pid, _ := GetID(buffer)
  return fmt.Sprintf("/Session/%s", pid)
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
  Live(foe, logger)
  // go func(){ Negeneration(&(*&foe.Status.Health), *&foe.Attributes.Vitality, *&foe.Attributes.Poolsize, *&foe.Basics.Body, logger) }()
}

func Regeneration(pool *[]funcs.Dot, health *float64, alive bool, max float64, maxhp float64, stream funcs.Stream, body funcs.Stream, logger *string) {
  for {
    if !(alive) { *logger = fmt.Sprintf("Not logged in yet. Stop Regeneration") ; break }
    if max-float64(len(*pool))<1 { time.Sleep( time.Millisecond * time.Duration( balance.Regeneration_DefaultTimeout() )) } else {
      dot := balance.Regeneration_DotWeight_FromStream(stream)
      pause := balance.Regeneration_TimeoutMilliseconds_FromWeightPool(dot.Weight, float64(len(*pool)), max)
      heal := balance.Regeneration_Heal_FromBody(body)
      time.Sleep( time.Millisecond * time.Duration( pause ))
      // +break logout
      if *health <= 0 { *logger = fmt.Sprintf("You are Died") ; break }
      if !(alive) { *logger = fmt.Sprintf("Logged out") ; break }
      //block
      *pool = append(*pool, dot )
      if *health >= maxhp {
        *logger = fmt.Sprintf("          for %0.3fs +%s %0.3f'e ", pause/1000, dot.Element, dot.Weight)
        } else {
          *logger = fmt.Sprintf("%+0.3f'hp for %0.3fs +%s %0.3f'e ", heal, pause/1000, dot.Element, dot.Weight)
        }
      if *health < maxhp { *health += heal } else { *health = maxhp }
      //unblock
    }
  }
}

func Negeneration(health *float64, alive bool, maxhp float64, maxe float64, body funcs.Stream, logger *string) {
  for {
    if !(alive) { *logger = fmt.Sprintf("Not spawned yet. Stop Regeneration") ; break }
    if maxhp<=*health { time.Sleep( time.Millisecond * time.Duration( balance.Regeneration_DefaultTimeout() )) } else {
      dot := balance.Regeneration_DotWeight_FromStream(body)
      pause := balance.Regeneration_TimeoutMilliseconds_FromWeightPool(dot.Weight, funcs.Log(maxe), maxe)
      heal := balance.Regeneration_Heal_FromBody(body)
      time.Sleep( time.Millisecond * time.Duration( pause ))
      // +break unspawn
      if *health <= 0 { *logger = fmt.Sprintf("Dummy: Foe died ") ; break }
      if !(alive) { *logger = fmt.Sprintf("Unspawned") ; break }
      //block
      if *health < maxhp { *logger = fmt.Sprintf("Dummy: %+0.3f'hp for %0.3fs ", heal, pause/1000) }
      if *health < maxhp { *health += heal } else { *health = maxhp }
      //unblock
    }
  }
}

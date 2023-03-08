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
  ActionLog []funcs.Action `json:"ActionLog,omitempty"`
  XYZ [3]int `json:"XYZ"` // not used yet
  Health int `json:"Health"`
  Pool []funcs.Dot `json:"Pool,omitempty"`
  // InFocus []string `json:"InFocus,omitempty"` // 0 for target, 1+ other
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
  elem, _ := funcs.ReStr(buffer.Basics.Streams)
  resists[elem] = balance.BasicStats_Resistance_FromStream(buffer.Basics.Streams)
  buffer.Attributes.Resistances = resists
  buffer.Attributes.Poolsize = balance.BasicStats_MaxPool_FromStream(buffer.Basics.Streams)
  *player = buffer
}

func PlayerBorn(player *Player, mean float64, logger *string) string {
  buffer := Player{}
  buffer.Basics.ID.NPC = false
  buffer.Basics.ID.Born = funcs.Epoch()
  buffer.Basics.ID.Last = buffer.Basics.ID.Born
  buffer.Basics.Body = balance.BasicStats_Stream_FromNormaleWithElement(2, funcs.Physical[1])
  buffer.Basics.Streams = balance.BasicStats_Stream_FromNormaleWithElement(1+mean, funcs.Elements[0])
  CalculateAttributes_FromBasics(&buffer)
  buffer.Status.Health = int(1000/math.Sqrt(buffer.Attributes.Vitality)) //from db
  *player = buffer
  Live(player, logger)
  pid, _ := GetID(buffer)
  return fmt.Sprintf("/Session/%s", pid)
}

func FoeSpawn(foe *Player, mean float64, logger *string) { // old, new+ template Stream{}
  buffer := Player{}
  buffer.Basics.ID.NPC = true
  buffer.Basics.ID.Born = funcs.Epoch()
  buffer.Basics.ID.Last = buffer.Basics.ID.Born
  buffer.Basics.Body = balance.BasicStats_Stream_FromNormaleWithElement(2, funcs.Physical[2])
  buffer.Basics.Streams = balance.BasicStats_Stream_FromNormaleWithElement(1+mean, funcs.Elements[0])
  CalculateAttributes_FromBasics(&buffer)
  buffer.Status.Health = 618
  *foe = buffer
  Live(foe, logger)
}

func Regeneration(pool *[]funcs.Dot, health *int, alive bool, max float64, vitality float64, stream funcs.Stream, body funcs.Stream, logger *string) {
  for {
    if !(alive) { *logger = fmt.Sprintf("Not logged in yet. Stop Regeneration") ; break }
    if max-float64(len(*pool))<1 { 
      time.Sleep( time.Millisecond * time.Duration( balance.Regeneration_DefaultTimeout() )) 
      *logger = fmt.Sprintf("Energy full, regeneration paused. ")
    } else {
      elem, _ := funcs.ReStr(stream)
      dot := balance.Regeneration_DotWeight_FromStream(stream)
      pause := balance.Regeneration_TimeoutMilliseconds_FromWeightPool(dot[elem], float64(len(*pool)), max)
      hpause := balance.Regeneration_TimeoutMilliseconds_FromWeightPool(dot[elem], float64(len(*pool)), max)
      heal := balance.Regeneration_Heal_FromBody(body) * pause/hpause
      time.Sleep( time.Millisecond * time.Duration( pause ))
      // +break logout
      if *health <= 0 { *logger = fmt.Sprintf("You are Died") ; break }
      if !(alive) { *logger = fmt.Sprintf("Logged out") ; break }
      //block
      *pool = append(*pool, dot )
      for elem, weig := range dot {
        if *health >= 1000 {
          *logger = fmt.Sprintf("          for %0.3fs +%d'%s", pause/1000, weig, elem)
        } else {
          *logger = fmt.Sprintf("%d'hp for %0.3fs +%d'%s", funcs.ChancedRound(heal), pause/1000, weig, elem)
        }
      }
      if *health < 1000 { *health += funcs.ChancedRound( heal*1000/vitality ) } else { *health = 1000 }
      //unblock
    }
  }
}

func Negeneration(health *int, alive bool, vitality float64, maxe float64, body funcs.Stream, logger *string) {
  for {
    if !(alive) { *logger = fmt.Sprintf("Not spawned yet. Stop Regeneration") ; break }
    if 1000 <= *health { 
      time.Sleep( time.Millisecond * time.Duration( balance.Regeneration_DefaultTimeout() )) 
      *logger = fmt.Sprintf("Health full, regeneration paused. ")
    } else {
      dot := balance.Regeneration_DotWeight_FromStream(body)
      elem, _ := funcs.ReStr(body)
      pause := balance.Regeneration_TimeoutMilliseconds_FromWeightPool(dot[elem], funcs.Log(maxe), maxe)
      heal := balance.Regeneration_Heal_FromBody(body)
      time.Sleep( time.Millisecond * time.Duration( pause ))
      // +break unspawn
      if *health <= 0 { *logger = fmt.Sprintf("Dummy: Foe died ") ; break }
      if !(alive) { *logger = fmt.Sprintf("Unspawned") ; break }
      //block
      if *health < 1000 { *logger = fmt.Sprintf("Dummy: %d'hp for %0.3fs ", funcs.ChancedRound(heal), pause/1000) }
      if *health < 1000 { *health += funcs.ChancedRound( heal*1000/vitality ) } else { *health = 1000 }
      //unblock
    }
  }
}

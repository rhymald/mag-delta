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
    Description string `json:"Name,omitempty"`
    Born int64 `json:"Born,omitempty"`
    Last int64 `json:"Last,omitempty"`
  } `json:"ID,omitempty"`
  Body funcs.Stream `json:"Body"`
  Streams []funcs.Stream `json:"Streams,omitempty"`
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
  go func(){ Regeneration(&(*&player.Status.Pool), &(*&player.Status.Health), &(*&player.Attributes.Login), *&player.Attributes.Poolsize, *&player.Attributes.Vitality, *&player.Basics.Streams[0], *&player.Basics.Body, *&player.Basics.ID.NPC, logger) }()
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
  resists, pool := make(map[string]float64), 0.0
  for _, each := range buffer.Basics.Streams {
    elem, _ := funcs.ReStr(each)
    if elem != funcs.Elements[0] { resists[elem] += balance.BasicStats_Resistance_FromStream(each) } 
    resists[funcs.Elements[0]] += math.Pow(math.Log2(balance.BasicStats_Resistance_FromStream(each)+2), 3)
    pool += balance.BasicStats_MaxPool_FromStream(each)
  }
  buffer.Attributes.Resistances = resists
  buffer.Attributes.Poolsize = pool
  *player = buffer
}

func PlayerBorn(player *Player, mean float64, logger *string) string {
  buffer := Player{}
  buffer.Basics.ID.NPC = false
  buffer.Basics.ID.Description = "Oh, that's me!"
  buffer.Basics.ID.Born = funcs.Epoch()
  buffer.Basics.ID.Last = buffer.Basics.ID.Born
  buffer.Basics.Body = balance.BasicStats_Stream_FromNormaleWithElement(2, funcs.Physical[1])
  strc := balance.BasicStats_StreamsCountAndModifier(buffer.Basics.ID.Born)
  for i:=0; i<strc; i++ { buffer.Basics.Streams = append(buffer.Basics.Streams, balance.BasicStats_Stream_FromNormaleWithElement(mean/float64(strc), funcs.Elements[0])) }
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
  buffer.Basics.ID.Description = "Training dummy, human-sized"
  buffer.Basics.ID.Born = funcs.Epoch()
  buffer.Basics.ID.Last = buffer.Basics.ID.Born
  buffer.Basics.Body = balance.BasicStats_Stream_FromNormaleWithElement(2, funcs.Physical[2])
  buffer.Basics.Streams = append(buffer.Basics.Streams, balance.BasicStats_Stream_FromNormaleWithElement(mean, funcs.Elements[0]))
  CalculateAttributes_FromBasics(&buffer)
  buffer.Status.Health = 618
  *foe = buffer
  Live(foe, logger)
}

func Regeneration(pool *[]funcs.Dot, health *int, alive *bool, max float64, vitality float64, stream funcs.Stream, body funcs.Stream, npc bool, logger *string) {
  for {
    if *alive != true { *logger = fmt.Sprintf("Not logged in yet. Stop Regeneration") ; break }
    trigger := max-float64(len(*pool))<1
    if npc { trigger = 1000 <= *health }
    if trigger { 
      time.Sleep( time.Millisecond * time.Duration( balance.Regeneration_DefaultTimeout() )) 
      *logger = fmt.Sprintf("Energy full, regeneration paused. ")
    } else {
      elem, _ := funcs.ReStr(stream)
      dot := balance.Regeneration_DotWeight_FromStream(stream)
      // if npc { dot = balance.Regeneration_DotWeight_FromStream(body) }
      // hpause := balance.Regeneration_TimeoutMilliseconds_FromWeightPool(dot[elem], float64(len(*pool)), max)
      curen := float64(len(*pool)) ; if npc { curen = max/2 }
      pause := balance.Regeneration_TimeoutMilliseconds_FromWeightPool(dot[elem], curen, max)
      heal := funcs.ChancedRound( balance.Regeneration_Heal_FromBody(body) * 1000/vitality )
      time.Sleep( time.Millisecond * time.Duration( pause ))
      // +break logout
      if *health <= 0 { *logger = fmt.Sprintf("You are died / Foe is dead") ; break }
      if *alive != true { *logger = fmt.Sprintf("Logged out / Monster escaped") ; break }
      //block
      if !npc { *pool = append(*pool, dot ) }
      for elem, weig := range dot {
        if *health >= 1000 {
          *logger = fmt.Sprintf("           for %0.3fs +%d'%s", pause/1000, weig, elem)
        } else {
          if npc {
            *logger = fmt.Sprintf("+%0.1f%% hp/s for %0.3fs ", float64(heal)/10/(pause/1000), pause/1000)
          } else {
            *logger = fmt.Sprintf("+%0.1f%% hp/s for %0.3fs +%d'%s", float64(heal)/10/(pause/1000), pause/1000, weig, elem)
          }
        }
      }
      if *health < 1000 { *health += heal } else { *health = 1000 }
      //unblock
    }
  }
}

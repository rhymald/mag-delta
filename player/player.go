package player

import (
  "rhymald/mag-delta/funcs"
  "rhymald/mag-delta/balance"
  "math"
  "fmt"
  "time"
)

type Player struct {
  Busy bool
  // Physical
  Physical struct {
    Health struct {
      Current float64
      Max float64
    }
    Body funcs.Stream
  }
  // Energetical
  Nature struct {
    Resistance float64
    Stream funcs.Stream
    Pool struct {
      Max float64
      Dots []funcs.Dot
    }
  }
}

func PlayerBorn(player *Player, mean float64){
  buffer := Player{}
  buffer.Physical.Body = balance.BasicStats_Stream_FromNormaleWithElement(2, "Common")
  buffer.Physical.Health.Max = balance.BasicStats_MaxHP_FromBody(buffer.Physical.Body) // from db
  buffer.Physical.Health.Current = math.Sqrt(buffer.Physical.Health.Max+1)-1 //from db
  buffer.Nature.Stream = balance.BasicStats_Stream_FromNormaleWithElement(1+mean, "Common")
  buffer.Nature.Resistance = balance.BasicStats_Resistance_FromStream(buffer.Nature.Stream)
  buffer.Nature.Pool.Max = balance.BasicStats_MaxPool_FromStream(buffer.Nature.Stream)
  *player = buffer
  go func(){ Regeneration(&(*&player.Nature.Pool.Dots), &(*&player.Physical.Health.Current), *&player.Nature.Pool.Max, *&player.Physical.Health.Max, *&player.Nature.Stream, *&player.Physical.Body) }()
}

func FoeSpawn(foe *Player, mean float64) {
  buffer := Player{}
  buffer.Physical.Body = balance.BasicStats_Stream_FromNormaleWithElement(2, "Common")
  buffer.Physical.Health.Max = balance.BasicStats_MaxHP_FromBody(buffer.Physical.Body) // from db
  buffer.Physical.Health.Current = buffer.Physical.Health.Max / math.Sqrt2 //from db
  buffer.Nature.Stream = balance.BasicStats_Stream_FromNormaleWithElement(1+mean, "Common")
  buffer.Nature.Resistance = balance.BasicStats_Resistance_FromStream(buffer.Nature.Stream)
  buffer.Nature.Pool.Max = balance.BasicStats_MaxPool_FromStream(buffer.Nature.Stream)
  *foe = buffer
  go func(){ Negeneration(&(*&foe.Physical.Health.Current), *&foe.Physical.Health.Max, *&foe.Nature.Pool.Max, *&foe.Physical.Body) }()
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
      pause := balance.Regeneration_TimeoutMilliseconds_FromWeightPool(dot.Weight, 0, maxe)
      heal := balance.Regeneration_Heal_FromBody(body)
      time.Sleep( time.Millisecond * time.Duration( pause ))
      //block
      if *health < maxhp { fmt.Printf("\rDEBUG[ NPC  ][Regeneration]: %+0.3f'hp for %0.3fs ░░░░░░░░░░░░░░░░░░░░░░░░░\r", heal, pause/1000) }
      if *health <= 0 { fmt.Printf("DEBUG[ NPC  ][Regeneration]: Foe died ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░\n") ; break }
      if *health < maxhp { *health += heal } else { *health = maxhp }
      //unblock
    }
  }
}

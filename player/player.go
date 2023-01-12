package player

import (
  "rhymald/mag-delta/funcs"
  "rhymald/mag-delta/balance"
  "math"
  "fmt"
  "time"
)

type Player struct {
  // Physical
  Health struct {
    Current float64
    Max float64
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
  // mean += math.Sqrt(3)
  // playerTuple := [][]string{}
  buffer := Player{}
  // fmt.Println(plot.Color("Player in game:",0))
  buffer.Health.Max = balance.BasicStats_MaxHP_FromNormale(mean) // from db
  buffer.Health.Current = math.Sqrt(buffer.Health.Max+1)-1 //from db
  // current := fmt.Sprintf("Health|Max: %0.0f|Current: %0.0f|Rate: %1.0f%%", buffer.Health.Max, buffer.Health.Current, 100*buffer.Health.Current/buffer.Health.Max)
  // playerTuple = plot.AddRow(current, playerTuple)
  buffer.Nature.Stream = balance.BasicStats_Stream_FromNormaleWithElement(mean, "Common")
  buffer.Nature.Resistance = balance.BasicStats_Resistance_FromStream(buffer.Nature.Stream)
  // row := fmt.Sprintf(
  //   "Element\n%s|Creation\n%0.3f|Alteration\n%0.3f|Destruction\n%0.3f|Resistance\n%0.3f",
  //   // "Element\n%s|Creation\n%0.3f|Alteration\n%0.3f|Destruction\n%0.3f",
  //   buffer.Nature.Stream.Element,
  //   buffer.Nature.Stream.Cre,
  //   buffer.Nature.Stream.Alt,
  //   buffer.Nature.Stream.Des,
  //   buffer.Nature.Resistance,
  // )
  // playerTuple = plot.AddRow(row,playerTuple)
  buffer.Nature.Pool.Max = balance.BasicStats_MaxPool_FromStream(buffer.Nature.Stream)
  // playerTuple = plot.AddRow( fmt.Sprintf("Pool|Max: %0.0f|Current: %d|Rate: %1.0f%%", buffer.Nature.Pool.Max, len(buffer.Nature.Pool.Dots), 100*float64(len(buffer.Nature.Pool.Dots))/float64(buffer.Nature.Pool.Max) ) ,playerTuple)
  // plot.Table(playerTuple, false)
  *player = buffer
  go func(){ Regeneration(&(*&player.Nature.Pool.Dots), &(*&player.Health.Current), *&player.Nature.Pool.Max, *&player.Health.Max, *&player.Nature.Stream) }()
}

func FoeSpawn(foe *Player, mean float64) {
  // mean += math.Sqrt(3)
  // playerTuple := [][]string{}
  buffer := Player{}
  // fmt.Println(plot.Color("Foe spawned:",0))
  buffer.Health.Max = balance.BasicStats_MaxHP_FromNormale(mean) // from db
  buffer.Health.Current = buffer.Health.Max / math.Sqrt2 //from db
  // current := fmt.Sprintf("Health|||Rate: %1.0f%%", 100*buffer.Health.Current/buffer.Health.Max)
  // playerTuple = plot.AddRow(current, playerTuple)
  buffer.Nature.Stream = balance.BasicStats_Stream_FromNormaleWithElement(mean, "Common")
  buffer.Nature.Resistance = balance.BasicStats_Resistance_FromStream(buffer.Nature.Stream)
  // row := fmt.Sprintf(
  //   "Element\n%s|Creation\n%0.3f|Alteration\n%0.3f|Destruction\n%0.3f|Resistance\n%0.3f",
  //   buffer.Nature.Stream.Element,
  //   math.Sqrt(mean*mean/3),
  //   math.Sqrt(mean*mean/3),
  //   math.Sqrt(mean*mean/3),
  //   buffer.Nature.Resistance,
  // )
  // playerTuple = plot.AddRow(row,playerTuple)
  buffer.Nature.Pool.Max = balance.BasicStats_MaxPool_FromStream(buffer.Nature.Stream)
  // playerTuple = plot.AddRow( fmt.Sprintf("Pool|Max: %0.0f", buffer.Nature.Pool.Max ) ,playerTuple)
  // plot.Table(playerTuple, false)
  *foe = buffer
  go func(){ Negeneration(&(*&foe.Health.Current), *&foe.Health.Max, *&foe.Nature.Pool.Max, *&foe.Nature.Stream) }()
}

func Regeneration(pool *[]funcs.Dot, health *float64, max float64, maxhp float64, stream funcs.Stream) {
  for {
    if max-float64(len(*pool))<1 { time.Sleep( time.Millisecond * time.Duration( balance.Regeneration_DefaultTimeout() )) } else {
      dot :=   balance.Regeneration_DotWeight_FromStream(stream)
      pause := balance.Regeneration_TimeoutMilliseconds_FromWeightPool(dot.Weight, float64(len(*pool)), max)
      heal :=  balance.Regeneration_Heal_FromWeight(dot.Weight)
      time.Sleep( time.Millisecond * time.Duration( pause ))
      //block
      if *health >= maxhp {
        fmt.Printf("DEBUG[Player][Regeneration]:           for %0.3fs +%s %0.3f'e \r", pause/1000, dot.Element, dot.Weight)
      } else {
        fmt.Printf("DEBUG[Player][Regeneration]: %+0.3f'hp for %0.3fs +%s %0.3f'e \r", heal, pause/1000, dot.Element, dot.Weight)
      }
      *pool = append(*pool, dot )
      if *health <= 0 { fmt.Printf("DEBUG[Player][Regeneration]: %s\n", "You are Died") ; break }
      if *health < maxhp { *health += heal } else { *health = maxhp }
      //unblock
    }
  }
}

func Negeneration(health *float64, maxhp float64, maxe float64, stream funcs.Stream) {
  for {
    dot :=   balance.Regeneration_DotWeight_FromStream(stream)
    pause := balance.Regeneration_TimeoutMilliseconds_FromWeightPool(dot.Weight, 0, maxe)
    heal :=  balance.Regeneration_Heal_FromWeight(dot.Weight)
    //block
    if *health < maxhp { fmt.Printf("\rDEBUG[ NPC  ][Regeneration]: %+0.3f'hp for %0.3fs                      \r", heal, pause/1000) }
    time.Sleep( time.Millisecond * time.Duration( pause ))
    if *health <= 0 { fmt.Printf("DEBUG[ NPC  ][Regeneration]: %s\n", "Foe died") ; break }
    if *health < maxhp { *health += heal } else { *health = maxhp }
    //unblock
  }
}

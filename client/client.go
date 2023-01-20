package client

import (
  "fmt"
  // "time"
  "rhymald/mag-delta/client/plot"
  "rhymald/mag-delta/player"
)

func PlayerStatus(players ...player.Player) {
  it, foe, compare := players[0], player.Player{}, len(players) > 1
  if players[1].Physical.Health.Current <= 0 || players[0].ID.NPC { compare = false }
  if compare { foe = players[1] }
  playerTuple := [][]string{}
  fmt.Println(plot.Color("Player status",0),"[comparing to a foe]:")
  line := ""
  if compare {
    line = fmt.Sprintf(
      "Health|Max: %0.0f|Current: %0.0f|Rate: %3.0f%%|[%3.0f%%]",
      it.Physical.Health.Max,
      it.Physical.Health.Current,
      100*it.Physical.Health.Current/it.Physical.Health.Max,
      100*foe.Physical.Health.Current/foe.Physical.Health.Max,
    )
  } else {
    line = fmt.Sprintf(
      "Health|Max: %0.0f|Current: %0.0f|Rate: %1.0f%%",
      it.Physical.Health.Max,
      it.Physical.Health.Current,
      100*it.Physical.Health.Current/it.Physical.Health.Max,
    )
  }
  playerTuple = plot.AddRow(line, playerTuple)
  if compare {
    line = fmt.Sprintf(
      " \nPhysical|Size\n  %0.3f \n [%0.3f]|Endurance\n  %0.3f \n [%0.3f]|Strength\n  %0.3f \n [%0.3f]",
      it.Physical.Body.Cre,
      foe.Physical.Body.Cre,
      it.Physical.Body.Alt,
      foe.Physical.Body.Alt,
      it.Physical.Body.Des,
      foe.Physical.Body.Des,
    )
  } else {
    line = fmt.Sprintf(
      " \nPhysical|Size\n%0.3f|Endurance\n%0.3f|Strength\n%0.3f",
      it.Physical.Body.Cre,
      it.Physical.Body.Alt,
      it.Physical.Body.Des,
    )
  }
  playerTuple = plot.AddRow(line,playerTuple)
  if compare {
    line = fmt.Sprintf(
      " \n %s \n[%s]|Creation\n  %0.3f \n [%0.3f]|Alteration\n  %0.3f \n [%0.3f]|Destruction\n  %0.3f \n [%0.3f]|Resistance\n  %0.3f \n [%0.3f]",
      it.Nature.Stream.Element,
      foe.Nature.Stream.Element,
      it.Nature.Stream.Cre,
      foe.Nature.Stream.Cre,
      it.Nature.Stream.Alt,
      foe.Nature.Stream.Alt,
      it.Nature.Stream.Des,
      foe.Nature.Stream.Des,
      it.Nature.Resistance,
      foe.Nature.Resistance,
    )
  } else {
    line = fmt.Sprintf(
      "Element\n%s|Creation\n%0.3f|Alteration\n%0.3f|Destruction\n%0.3f",
      it.Nature.Stream.Element,
      it.Nature.Stream.Cre,
      it.Nature.Stream.Alt,
      it.Nature.Stream.Des,
    )
  }
  playerTuple = plot.AddRow(line,playerTuple)
  if compare {
    line = fmt.Sprintf("Pool|Max: %0.0f|Current: %d|Rate: %1.0f%%|[%0.0f]", it.Nature.Pool.Max, len(it.Nature.Pool.Dots), 100*float64(len(it.Nature.Pool.Dots))/float64(it.Nature.Pool.Max), foe.Nature.Pool.Max )
  } else {
    line = fmt.Sprintf("Pool|Max: %0.0f|Current: %d|Rate: %1.0f%%", it.Nature.Pool.Max, len(it.Nature.Pool.Dots), 100*float64(len(it.Nature.Pool.Dots))/float64(it.Nature.Pool.Max) )
  }
  playerTuple = plot.AddRow(line,playerTuple)
  plot.Table(playerTuple, false)
  fmt.Println()
}

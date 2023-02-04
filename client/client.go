package client

import (
  "fmt"
  // "time"
  "rhymald/mag-delta/client/plot"
  "rhymald/mag-delta/player"
)

func PlayerStatus(players ...player.Player) {
  it, foe, compare := players[0], player.Player{}, len(players) > 1
  if players[1].Status.Health <= 0 || players[0].Basics.ID.NPC { compare = false }
  if compare { foe = players[1] }
  playerTuple := [][]string{}
  fmt.Println(plot.Color("Player status",0),"[comparing to a foe]:")
  line := ""
  if compare {
    line = fmt.Sprintf(
      "Health|Max: %0.0f|Current: %0.0f|Rate: %3.0f%%|[%3.0f%%]",
      it.Attributes.Vitality,
      it.Status.Health,
      100*it.Status.Health/it.Attributes.Vitality,
      100*foe.Status.Health/foe.Attributes.Vitality,
    )
  } else {
    line = fmt.Sprintf(
      "Health|Max: %0.0f|Current: %0.0f|Rate: %1.0f%%",
      it.Attributes.Vitality,
      it.Status.Health,
      100*it.Status.Health/it.Attributes.Vitality,
    )
  }
  playerTuple = plot.AddRow(line, playerTuple)
  if compare {
    line = fmt.Sprintf(
      " \nPhysical|Size\n  %0.3f \n [%0.3f]|Endurance\n  %0.3f \n [%0.3f]|Strength\n  %0.3f \n [%0.3f]",
      it.Basics.Body.Cre,
      foe.Basics.Body.Cre,
      it.Basics.Body.Alt,
      foe.Basics.Body.Alt,
      it.Basics.Body.Des,
      foe.Basics.Body.Des,
    )
  } else {
    line = fmt.Sprintf(
      " \nPhysical|Size\n%0.3f|Endurance\n%0.3f|Strength\n%0.3f",
      it.Basics.Body.Cre,
      it.Basics.Body.Alt,
      it.Basics.Body.Des,
    )
  }
  playerTuple = plot.AddRow(line,playerTuple)
  if compare {
    line = fmt.Sprintf(
      " \n %s \n[%s]|Creation\n  %0.3f \n [%0.3f]|Alteration\n  %0.3f \n [%0.3f]|Destruction\n  %0.3f \n [%0.3f]|Resistance\n  %0.3f \n [%0.3f]",
      it.Basics.Streams.Element,
      foe.Basics.Streams.Element,
      it.Basics.Streams.Cre,
      foe.Basics.Streams.Cre,
      it.Basics.Streams.Alt,
      foe.Basics.Streams.Alt,
      it.Basics.Streams.Des,
      foe.Basics.Streams.Des,
      it.Attributes.Resistances[it.Basics.Streams.Element],
      foe.Attributes.Resistances[foe.Basics.Streams.Element],
    )
  } else {
    line = fmt.Sprintf(
      "Element\n%s|Creation\n%0.3f|Alteration\n%0.3f|Destruction\n%0.3f",
      it.Basics.Streams.Element,
      it.Basics.Streams.Cre,
      it.Basics.Streams.Alt,
      it.Basics.Streams.Des,
    )
  }
  playerTuple = plot.AddRow(line,playerTuple)
  if compare {
    line = fmt.Sprintf("Pool|Max: %0.0f|Current: %d|Rate: %1.0f%%|[%0.0f]", it.Attributes.Poolsize, len(it.Status.Pool), 100*float64(len(it.Status.Pool))/float64(it.Attributes.Poolsize), foe.Attributes.Poolsize )
  } else {
    line = fmt.Sprintf("Pool|Max: %0.0f|Current: %d|Rate: %1.0f%%", it.Attributes.Poolsize, len(it.Status.Pool), 100*float64(len(it.Status.Pool))/float64(it.Attributes.Poolsize) )
  }
  playerTuple = plot.AddRow(line,playerTuple)
  plot.Table(playerTuple, false)
  fmt.Println()
}

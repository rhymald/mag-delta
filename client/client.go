package client

import (
  "fmt"
  "rhymald/mag-delta/client/plot"
  "rhymald/mag-delta/player"
  "rhymald/mag-delta/funcs"
  "rhymald/mag-delta/balance"
)
var ElemList []string = []string{"Common", "Air", "Fire", "Earth", "Water", "Void", "Mallom", "Noise", "Resonance",}
var PhysList []string = []string{"Ghosty", "Flesh", "Wooden", "Stone", "Forged"}

func elemTotStr(e string) string {
  for i, elem := range funcs.Elements { if e == elem { return ElemList[i] } } 
  for i, phys := range funcs.Physical { if e == phys { return PhysList[i] } }
  return "ERROR" 
}

func PlayerStatus(players ...player.Player) {
  it, foe, compare := players[0], player.Player{}, len(players) > 1
  if players[1].Status.Health <= 0 || players[0].Basics.ID.NPC { compare = false }
  if compare { foe = players[1] }

  fmt.Print("  Health ")
  plot.Baaar( it.Status.Health/it.Attributes.Vitality, 50, "right" )
  fmt.Print("\n  ")
  plot.Baaar( float64(len(it.Status.Pool))/it.Attributes.Poolsize, 50, "fade" )
  fmt.Print(" Energy\n")
  fmt.Println()
  if compare {
    fmt.Print("  Dummy ")
    plot.Baaar( foe.Status.Health/foe.Attributes.Vitality, 51, "up" )
    fmt.Print("\n\n")  
  }

  playerTuple := [][]string{}
  fmt.Println(plot.Color("\nPlayer status",0),"[comparing to a foe]:")
  line := ""
  if compare {
    line = fmt.Sprintf(
      "Physical\n    %s\n   [%s]|Toughness\n   %0.3f \n  [%0.3f]|Agility\n   %0.3f \n  [%0.3f]|Strength\n   %0.3f \n  [%0.3f]",
      elemTotStr(it.Basics.Body.Element),
      elemTotStr(foe.Basics.Body.Element),
      it.Basics.Body.Cre,
      foe.Basics.Body.Cre,
      it.Basics.Body.Alt,
      foe.Basics.Body.Alt,
      it.Basics.Body.Des,
      foe.Basics.Body.Des,
    )
  } else {
    line = fmt.Sprintf(
      "Physical\n    %s|Toughness\n%0.3f|Agility\n%0.3f|Strength\n%0.3f",
      elemTotStr(it.Basics.Body.Element),
      it.Basics.Body.Cre,
      it.Basics.Body.Alt,
      it.Basics.Body.Des,
    )
  }
  playerTuple = plot.AddRow(line,playerTuple)
  yourAbilities := balance.StreamAbilities_FromStream(it.Basics.Streams)
  foeAbilities := balance.StreamAbilities_FromStream(foe.Basics.Streams)
  if compare {
    line = fmt.Sprintf(
      "Energy \n    %s \n   [%s]|Resistance\n   %0.3f \n  [%0.3f]|Creation\n   %0.3f \n  [%0.3f]|Alteration\n   %0.3f \n  [%0.3f]|Destruction\n   %0.3f \n  [%0.3f]",
      elemTotStr(it.Basics.Streams.Element),
      elemTotStr(foe.Basics.Streams.Element),
      it.Attributes.Resistances[it.Basics.Streams.Element],
      foe.Attributes.Resistances[foe.Basics.Streams.Element],
      it.Basics.Streams.Cre,
      foe.Basics.Streams.Cre,
      it.Basics.Streams.Alt,
      foe.Basics.Streams.Alt,
      it.Basics.Streams.Des,
      foe.Basics.Streams.Des,
    )
  } else {
    line = fmt.Sprintf(
      "Element\n    %s|Resistance\n%0.3f|Creation\n%0.3f|Alteration\n%0.3f|Destruction\n%0.3f",
      elemTotStr(it.Basics.Streams.Element),
      it.Attributes.Resistances[it.Basics.Streams.Element],
      it.Basics.Streams.Cre,
      it.Basics.Streams.Alt,
      it.Basics.Streams.Des,
    )
  }
  playerTuple = plot.AddRow(line,playerTuple)
  if compare {
    line = fmt.Sprintf(
      " |Creation|  %+5.1f%% \n [%+5.1f%%]|  %+5.1f%% \n [%+5.1f%%]|  %+5.1f%% \n [%+5.1f%%]",
      yourAbilities["Cad"]*100,
      foeAbilities["Cad"]*100,
      yourAbilities["Ca"]*100,
      foeAbilities["Ca"]*100,
      yourAbilities["Cd"]*100,
      foeAbilities["Cd"]*100,
    )
  } else {
    line = fmt.Sprintf(
      " |Creation|%+5.1f%%|%+5.1f%%|%+5.1f%%",
      yourAbilities["Cad"]*100,
      yourAbilities["Ca"]*100,
      yourAbilities["Cd"]*100,
    )
  }
  playerTuple = plot.AddRow(line,playerTuple)
  if compare {
    line = fmt.Sprintf(
      " |Alteration|  %+5.1f%% \n [%+5.1f%%]|  %+5.1f%% \n [%+5.1f%%]|  %+5.1f%% \n [%+5.1f%%]",
      yourAbilities["Ac"]*100,
      foeAbilities["Ac"]*100,
      yourAbilities["Acd"]*100,
      foeAbilities["Acd"]*100,
      yourAbilities["Ad"]*100,
      foeAbilities["Ad"]*100,
    )
  } else {
    line = fmt.Sprintf(
      " |Alteration|%+5.1f%%|%+5.1f%%|%+5.1f%%",
      yourAbilities["Ac"]*100,
      yourAbilities["Acd"]*100,
      yourAbilities["Ad"]*100,
    )
  }
  playerTuple = plot.AddRow(line,playerTuple)
  if compare {
    line = fmt.Sprintf(
      " |Destruction|  %+5.1f%% \n [%+5.1f%%]|  %+5.1f%% \n [%+5.1f%%]|  %+5.1f%% \n [%+5.1f%%]",
      yourAbilities["Dc"]*100,
      foeAbilities["Dc"]*100,
      yourAbilities["Da"]*100,
      foeAbilities["Da"]*100,
      yourAbilities["Dac"]*100,
      foeAbilities["Dac"]*100,
    )
  } else {
    line = fmt.Sprintf(
      " |Creation|%+5.1f%%|%+5.1f%%|%+5.1f%%",
      yourAbilities["Dc"]*100,
      yourAbilities["Da"]*100,
      yourAbilities["Dac"]*100,
    )
  }
  playerTuple = plot.AddRow(line,playerTuple)
  plot.Table(playerTuple, false)
  fmt.Println()
}


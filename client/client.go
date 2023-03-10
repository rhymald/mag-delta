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

func PlayerStatus(player player.Player) {
  // health and mana bars
  if player.Basics.ID.NPC {
    fmt.Println(plot.Color("\nTarget status",0),"[",player.Basics.ID.Description ,"]:\n")
  } else {
    fmt.Println(plot.Color("\nCharacter status",0),"[",player.Basics.ID.Description ,"]:\n")
  }
  elem, stats := funcs.ReStr(player.Basics.Streams[0]) // npc has only 1 sream
  line, playerTuple := "", [][]string{}
  attrs := []string{"Ca", "Cd", "Cad", "Ac", "Ad", "Acd", "Dc", "Da", "Dca"}
  fmt.Print("  Health ")
  plot.Baaar( float64(player.Status.Health)/1000 , 50, "right" )
  fmt.Printf(" | Vitality: %.0f\n  ", player.Attributes.Vitality)
  if !player.Basics.ID.NPC {
    plot.Baaar( float64(len(player.Status.Pool))/player.Attributes.Poolsize, 50, "up" )
    fmt.Printf(" Energy | Pool size: %.0f\n", player.Attributes.Poolsize)
  } ; fmt.Println()
  elem, stats = funcs.ReStr(player.Basics.Body)
  line = fmt.Sprintf(
    "Physical\n%s|Toughness\n%.0f|Agility\n%.0f|Strength\n%.0f",
    elemTotStr(elem),
    stats[0]*1000,
    stats[1]*1000,
    stats[2]*1000,
  )
  playerTuple = plot.AddRow(line,playerTuple)
  plot.Table(playerTuple, false) ; line, playerTuple = "", [][]string{}
  fmt.Print("    Resistances:")
  for el, resist := range player.Attributes.Resistances { if resist != 0 { fmt.Printf("    %s: %.3f", el, resist) } } ; fmt.Println()
  line = "Element|Create|Alterate|Destroy|Ca|Cd|Cad|Ac|Ad|Acd|Dc|Da|Dca"
  playerTuple = plot.AddRow(line,playerTuple)
  for index, every := range player.Basics.Streams {
  // if player.Basics.ID.NPC {
    elem, stats = funcs.ReStr(every)
    line = fmt.Sprintf(
      "%d'%s|%.0f|%.0f|%.0f",
      index+1, elemTotStr(elem),
      stats[0]*1000,
      stats[1]*1000,
      stats[2]*1000,
    )
    abilities := balance.StreamAbilities_FromStream(every)
    for _, each := range attrs {
      if abilities[each] != 0 {
        line = fmt.Sprintf("%s|%0.1f%%", line, abilities[each])
      } else {
        line = fmt.Sprintf("%s|---", line)
      }
    }
    playerTuple = plot.AddRow(line,playerTuple)
  }
  plot.Table(playerTuple, false)
  fmt.Println()
}
  // stats tuple preparation
  // playerTuple := [][]string{}
  // fmt.Println(plot.Color("\nPlayer status",0),"[comparing to a foe]:")
  // line := ""
  // if compare {
  //   itelem, itstats := funcs.ReStr(it.Basics.Body)
  //   foelem, foestats := funcs.ReStr(foe.Basics.Body)  
  //   line = fmt.Sprintf(
  //     "Physical\n    %s\n   [%s]|Toughness\n   %.3f \n  [%.3f]|Agility\n   %.3f \n  [%.3f]|Strength\n   %.3f \n  [%.3f]",
  //     elemTotStr(itelem),
  //     elemTotStr(foelem),
  //     itstats[0],
  //     foestats[0],
  //     itstats[1],
  //     foestats[1],
  //     itstats[2],
  //     foestats[2],
  //   )
  // } else {
  //   elem, stats := funcs.ReStr(it.Basics.Body)
  //   line = fmt.Sprintf(
  //     "Physical\n    %s|Toughness\n%.3f|Agility\n%.3f|Strength\n%.3f",
  //     elemTotStr(elem),
  //     stats[0],
  //     stats[1],
  //     stats[2],
  //   )
  // }
  // playerTuple = plot.AddRow(line,playerTuple)
  // yourAbilities := balance.StreamAbilities_FromStream(it.Basics.Streams[0])
  // foeAbilities := make(map[string]float64)
  // if compare { foeAbilities = balance.StreamAbilities_FromStream(foe.Basics.Streams[0]) }
  // if compare {
  //   itelem, itstats := funcs.ReStr(it.Basics.Streams[0])
  //   foelem, foestats := funcs.ReStr(foe.Basics.Streams[0])  
  //   line = fmt.Sprintf(
  //     "Energy \n    %s\n   [%s]|Affinity (Resist)\n %.3f vs [%.3f] \n[%.3f] vs %.3f|Creation\n   %.3f \n  [%.3f]|Alteration\n   %.3f \n  [%.3f]|Destruction\n   %.3f \n  [%.3f]",
  //     elemTotStr(itelem),
  //     elemTotStr(foelem),
  //     it.Attributes.Resistances[itelem]+balance.Cast_Common_Equalator(),
  //     it.Attributes.Resistances[foelem]+balance.Cast_Common_Equalator(),
  //     foe.Attributes.Resistances[foelem]+balance.Cast_Common_Equalator(),
  //     foe.Attributes.Resistances[itelem]+balance.Cast_Common_Equalator(),
  //     itstats[0],
  //     foestats[0],
  //     itstats[1],
  //     foestats[1],
  //     itstats[2],
  //     foestats[2],
  //   )
  // } else {
  //   elem, stats := funcs.ReStr(it.Basics.Streams[0])
  //   line = fmt.Sprintf(
  //     "Element\n    %s|Affinity (Resist)\n%.3f|Creation\n%.3f|Alteration\n%.3f|Destruction\n%.3f",
  //     elemTotStr(elem),
  //     it.Attributes.Resistances[elem]+it.Attributes.Resistances[funcs.Elements[0]],
  //     stats[0],
  //     stats[1],
  //     stats[2],
  //   )
  // }
  // playerTuple = plot.AddRow(line,playerTuple)
  // if compare {
  //   line = fmt.Sprintf(
  //     " |Creation|  %5.1f%% \n [%5.1f%%]|  %5.1f%% \n [%5.1f%%]|  %5.1f%% \n [%5.1f%%]",
  //     yourAbilities["Cad"]*100,
  //     foeAbilities["Cad"]*100,
  //     yourAbilities["Ca"]*100,
  //     foeAbilities["Ca"]*100,
  //     yourAbilities["Cd"]*100,
  //     foeAbilities["Cd"]*100,
  //   )
  // } else {
  //   line = fmt.Sprintf(
  //     " |Creation|%5.1f%%|%5.1f%%|%5.1f%%",
  //     yourAbilities["Cad"]*100,
  //     yourAbilities["Ca"]*100,
  //     yourAbilities["Cd"]*100,
  //   )
  // }
  // playerTuple = plot.AddRow(line,playerTuple)
  // if compare {
  //   line = fmt.Sprintf(
  //     " |Alteration|  %5.1f%% \n [%5.1f%%]|  %5.1f%% \n [%5.1f%%]|  %5.1f%% \n [%5.1f%%]",
  //     yourAbilities["Ac"]*100,
  //     foeAbilities["Ac"]*100,
  //     yourAbilities["Acd"]*100,
  //     foeAbilities["Acd"]*100,
  //     yourAbilities["Ad"]*100,
  //     foeAbilities["Ad"]*100,
  //   )
  // } else {
  //   line = fmt.Sprintf(
  //     " |Alteration|%5.1f%%|%5.1f%%|%5.1f%%",
  //     yourAbilities["Ac"]*100,
  //     yourAbilities["Acd"]*100,
  //     yourAbilities["Ad"]*100,
  //   )
  // }
  // playerTuple = plot.AddRow(line,playerTuple)
  // if compare {
  //   line = fmt.Sprintf(
  //     " |Destruction|  %5.1f%% \n [%5.1f%%]|  %5.1f%% \n [%5.1f%%]|  %5.1f%% \n [%5.1f%%]",
  //     yourAbilities["Dc"]*100,
  //     foeAbilities["Dc"]*100,
  //     yourAbilities["Da"]*100,
  //     foeAbilities["Da"]*100,
  //     yourAbilities["Dac"]*100,
  //     foeAbilities["Dac"]*100,
  //   )
  // } else {
  //   line = fmt.Sprintf(
  //     " |Creation|%5.1f%%|%5.1f%%|%5.1f%%",
  //     yourAbilities["Dc"]*100,
  //     yourAbilities["Da"]*100,
  //     yourAbilities["Dac"]*100,
  //   )
  // }
  // playerTuple = plot.AddRow(line,playerTuple)


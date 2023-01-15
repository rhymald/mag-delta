package main

import (
  "fmt"
  "math"
  "rhymald/mag-delta/client/plot"
  "rhymald/mag-delta/client"
  "rhymald/mag-delta/player"
  "rhymald/mag-delta/act"
)

var You player.Player
var Target player.Player
var Action string
var Keys chan string = make(chan string)

func init() {
  fmt.Println("\n\t\t  ", plot.Bar("Initializing...",8), "\n")
  player.PlayerBorn(&You,0)
  player.FoeSpawn(&Target,0)
}

func main() {
  fmt.Println("\n\t\t", plot.Bar("Successfully login",1),"\n")
  client.PlayerStatus(You, Target)
  fmt.Println("\n\t     ",plot.Bar("Press [Enter] to continue",8),"\n")
  fmt.Scanln()
  plot.ShowMenu(" ")
  client.PlayerStatus(You, Target)
  client.UI(Keys, You, Target)
  grow := 1/math.Phi/math.Phi/math.Phi
  for {
    key := actions()
    if Target.Physical.Health.Current <= 0 { grow = grow*math.Cbrt(math.Phi) ; player.FoeSpawn(&Target, grow) ; plot.ShowMenu(key)}// ; PlayerStatus(You, Target)}
  }
}

func actions() string {
  Action, _ := <-Keys
  key := Action
  if Action=="a" { act.Jinx(&You, &Target) ; Action = "" }
  return string(key)
}

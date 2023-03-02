package plot

import (
  "fmt"
)

type LogFrame struct {
  Player string
  Actions []string
  Foe string
  Size int
}

func CleanFrame() LogFrame {
  buffer := LogFrame{}
  buffer.Player = ""
  buffer.Size = 32
  buffer.Actions = append(buffer.Actions," ")
  buffer.Actions = append(buffer.Actions,"Welcome!")
  buffer.Actions = append(buffer.Actions,"Here you can find a list of actions you have made.")
  buffer.Actions = append(buffer.Actions,"Just press [E] key to attack the dummy,")
  buffer.Actions = append(buffer.Actions,"     Press [?] key to get the chain tree,")
  buffer.Actions = append(buffer.Actions,"        Or [/] key to get the list of players.")
  buffer.Actions = append(buffer.Actions," ")
  buffer.Foe = ""
  return buffer
}

func Frame(frame LogFrame){
  fmt.Println()
  fmt.Println(" ─┼──[Player status]─────────────────────────────────────────────")
  fmt.Println("    " ,frame.Player)
  fmt.Println(" ─┼─────[Actions]────────────────────────────────────────────────")
  for x:=0 ; x<(len(frame.Actions)) ; x++ {
    // if x%2 == 1 {fmt.Printf(" ")}
    fmt.Printf("    %s\n" ,frame.Actions[x])
  }
  fmt.Println(" ─┼───────[Foes]─────────────────────────────────────────────────")
  fmt.Println("    " ,frame.Foe)
  fmt.Println(" ─┼────────[End]─────────────────────────────────────────────────")
  fmt.Println()
}

func AddAction(frame *LogFrame, action string){
  *&frame.Actions = append(*&frame.Actions, action)
  if len(*&frame.Actions) > *&frame.Size {
    buffer := *&frame.Actions
    *&frame.Actions = buffer[(len(*&frame.Actions)-*&frame.Size):*&frame.Size]
  }
}

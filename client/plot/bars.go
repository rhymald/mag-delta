package plot 

import(
  "fmt"
  "math"
)

func Baaar(current float64, wid int, mode string) {
  if wid <= 0 { wid = 53 }
  if 1 < current { current = 1 }
  var growup []string = []string{"◦", "▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"} // for health
  var grow []string   = []string{"◦", "▏", "▎", "▍", "▌", "▋", "▊", "▉", "█"} // for barrier
  var fade []string   = []string{"◦", "░", "▒", "▓", "█"} // for mana
  pick := grow ; if mode == "right" { pick = grow } else if mode == "fade" { pick = fade  } else if mode == "up" { pick = growup }
  wid *= (len(pick)-1)
  rate, counter := int( math.Floor(float64(wid) * current ) ), 0
  if math.Abs(1 - current) <= float64(len(pick)-1)/float64(wid-1) { rate = wid }
  for {
    counter += (len(pick)-1)
    if counter-rate <= 0 { fmt.Print(pick[(len(pick)-1)]) } else if counter-rate < (len(pick)-1) && counter-rate > 0 { fmt.Print(pick[rate%(len(pick)-1)]) } else { fmt.Print(pick[0]) }
    if counter >= wid {break}
  }
}
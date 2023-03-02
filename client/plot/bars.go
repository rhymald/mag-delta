package plot 

import(
  "fmt"
  "math"
)

func Baaar(current float64, wid int, mode string) {
  if wid <= 0 { wid = 53 }
  if 1 < current { current = 1 }
  var growup [11]string = [11]string{"◦", "▁", "▁", "▂", "▃", "▄", "▅", "▆", "▇", "▇", "█"} // for health
  var grow [11]string   = [11]string{"◦", "▏", "▏", "▎", "▍", "▌", "▋", "▊", "▉", "▉", "█"} // for barrier
  var fade [11]string   = [11]string{"◦", "░", "░", "░", "░", "▒", "▒", "▒", "▓", "▓", "█"} // for mana
  wid *= 10
  rate, counter := int( math.Floor(float64(wid) * current ) ), 0
  if math.Abs(1 - current) <= 10/float64(wid-1) { rate = wid }
  pick := grow ; if mode == "right" { pick = grow } else if mode == "fade" { pick = fade  } else if mode == "up" { pick = growup }
  for {
    counter += 10
    if counter-rate <= 0 { fmt.Print(pick[10]) } else if counter-rate < 10 && counter-rate > 0 { fmt.Print(pick[rate%10]) } else { fmt.Print(pick[0]) }
    if counter >= wid {break}
  }
}
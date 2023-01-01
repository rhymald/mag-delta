package plot

import(
  "fmt"
  "math"
  "strings"
  "golang.org/x/term"
)

func PlotTable(tuple [][]string, stretch bool) {
  // tuple := tuple
  maxs := make([]int, len(tuple[0]))
  for j, y := range tuple {
    for i, _ := range y {
      if j == 0 {
        for c := 0; c>len(maxs); c++ { maxs[i] = 2+ MaxInCell(tuple[0][i]) }
      }
      maxs[i] = int(math.Max(2+float64( MaxInCell(tuple[j][i]) ), float64(maxs[i])))
    }
  }
  if stretch {
    for e, _ := range maxs { maxs[e] = int( math.Log2(float64(maxs[e])+2)/math.Log2(1.1459) ) }
    sums := 0
    for _, each := range maxs { sums+=each }
    termWigth, _, _ := term.GetSize(0)
    modificator := float64(termWigth - len(maxs) - 2) / float64(sums)
    for e, _ := range maxs { maxs[e] = int( float64(maxs[e])*modificator ) }
  }
  //Head:
  fmt.Printf(" ╔")
  for i, wid := range maxs {
    for counter:=0 ;counter < wid; counter++ {
      fmt.Printf("═")
    }
    if i+1 == len(maxs) {fmt.Printf("╗\n")} else {fmt.Printf("╤")} //╤
  }
  //String:
  for I, row := range tuple {
    PlotRow(row, maxs)
    if I+1 == len(tuple) {
      //Footer:
      fmt.Printf(" ╚")
      for i, wid := range maxs {
        for counter:=0 ;counter < wid; counter++ {
          fmt.Printf("═")
        }
        if i+1 == len(maxs) {fmt.Printf("╝\n")} else {fmt.Printf("╧")} //╧
      }
    } else {
      //Delimiter:
      fmt.Printf(" ╟")
      for i, wid := range maxs {
        for counter:=0 ;counter < wid; counter++ {
          fmt.Printf("─")
        }
        if i+1 == len(maxs) {fmt.Printf("╢\n")} else {fmt.Printf("┼")} //┼
      }
    }
  }
}
func AddRow(row string, tuple [][]string) [][]string {
  buffer := strings.Split(row, "|")
  if len(tuple)==0 { return [][]string{buffer} }
  if len(buffer) > len(tuple[0]) {
    for l, _ := range tuple { for count := 0; count < len(buffer)-len(tuple[0]); count++ { tuple[l] = append(tuple[l], " ") } }
  } else if len(buffer) < len(tuple[0]) {
    for count := 0; count <= len(tuple[0])-len(buffer); count++ { buffer = append(buffer, " ") }
  }
  tuple = append(tuple, buffer)
  return tuple
}
func FindDelim(row []string) ([]int, int) {
  buffer := make([]int, len(row))
  max := 0
  for i, cell := range row { buffer[i] = len(strings.Split(cell, "\n")) ; max = int(math.Max( float64(buffer[i]), float64(max) )) }
  return buffer, max
}
func PlotRow(row []string, widths []int) {
  _, max := FindDelim(row)
  for linenum:=0; linenum<max; linenum++ {
    fmt.Printf(" ║")
    for i, wid := range widths {
      fmt.Printf(" ")
      cell := strings.Split(row[i], "\n")
      if len(cell) < max { for count:=0; count<max-len(cell); count++ { cell = append(cell, string(" ")) } }
      toprint := cell[linenum]
      fmt.Printf("%s", toprint)
      for counter:=0 ;counter < wid-1-len(toprint); counter++ {
        fmt.Printf(" ")
      }
      if i+1 == len(widths) {fmt.Printf("║\n")} else {
        if row[i+1]=="" || row[i+1]==" " {fmt.Printf("│")} else {fmt.Printf("│")}
      }
    }
  }
}
func MaxInCell(cell string) int {
  lines, max := strings.Split(cell, "\n"), 0
  for _, each := range lines { max = int(math.Max( float64(max), float64(len(each)) )) }
  return max
}

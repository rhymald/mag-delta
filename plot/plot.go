package plot

import(
  "fmt"
  "math"
  "strings"
  "golang.org/x/term"
)

func Table(tuple [][]string, stretch bool) {
  // tuple := tuple
  maxs := make([]int, len(tuple[0]))
  for j, y := range tuple {
    for i, _ := range y {
      if j == 0 {
        for c := 0; c>len(maxs); c++ { maxs[i] = 2+ maxInCell(tuple[0][i]) }
      }
      maxs[i] = int(math.Max(2+float64( maxInCell(tuple[j][i]) ), float64(maxs[i])))
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
    plotRow(row, maxs)
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
    for l, _ := range tuple { for { if len(buffer) == len(tuple[0]) {break} ; tuple[l] = append(tuple[l], " ") } }
  }
  if len(buffer) < len(tuple[0])  {
    for { if len(buffer) == len(tuple[0]) {break} ; buffer = append(buffer, " ") }
  }
  tuple = append(tuple, buffer)
  return tuple
}

func findDelim(row []string) ([]int, int) {
  buffer := make([]int, len(row))
  max := 0
  for i, cell := range row { buffer[i] = len(strings.Split(cell, "\n")) ; max = int(math.Max( float64(buffer[i]), float64(max) )) }
  return buffer, max
}

func plotRow(row []string, widths []int) {
  _, max := findDelim(row)
  for linenum:=0; linenum<max; linenum++ {
    fmt.Printf(" ║")
    for i, wid := range widths {
      fmt.Printf(" ")
      cell := []string{" "}
      if row[i] != " " {cell = strings.Split(row[i], "\n")}
      if len(cell) < max { for { if max == len(cell) {break} ; cell = append(cell, string(" ")) } }
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

func maxInCell(cell string) int {
  lines, max := strings.Split(cell, "\n"), 0
  for _, each := range lines { max = int(math.Max( float64(max), float64(len(each)) )) }
  return max
}

func Bar(text string) string { return fmt.Sprintf("░▒▓█%s%s%s█▓▒░",R,text,E[0]) }

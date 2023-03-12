package player

import (
	"fmt"
	"testing"
  "rhymald/mag-delta/client/plot"
  // "rhymald/mag-delta/player"
)

// func debug() {
//   // player fetch
//   var ef player.Player 
//   var ad player.Player 
//   var bc player.Player 
//   _, _, _ = player.PlayerBorn(&ef, 3, &Frame.Player), player.PlayerBorn(&ad, 1, &Frame.Player), player.PlayerBorn(&bc, 2, &Frame.Player)
//   fmt.Println("Different:", player.Fetch_Stats(bc.Basics,ad.Basics))
//   fmt.Println("Same:     ", player.Fetch_Stats(bc.Basics,bc.Basics))
//   bc.Basics.ID.Born, ef.Basics.ID.Born = ad.Basics.ID.Born, ad.Basics.ID.Born  
//   gh := player.Fetch_Stats(bc.Basics,ad.Basics)
//   fmt.Println("Different 2:", gh) ; fmt.Println("Sum 2:      ", player.Grow_Stats(ef.Basics,gh))
//   player.TakeAll_Stats(&ef, []player.BasicStats{gh, gh, gh})
//   fmt.Println("Cascade 2:  ", ef)
//   // picker
//   fmt.Print(funcs.PickXFrom(17, 39)); fmt.Print(funcs.PickXFrom(4, 9)); fmt.Print(funcs.PickXFrom(17, 3)) ; fmt.Println(funcs.PickXFrom(4, 6))
//   // list all elements
//   for x:=0 ;x<len(funcs.Physical); x++ { fmt.Printf(" x[%s]%.3f ", funcs.Physical[x], math.Pow(math.Log2(float64(x)+1), 2) ) } ; fmt.Println()
//   for x:=0 ;x<3; x++ { fmt.Printf(" ^[%s]%.3f ", funcs.Elements[x], math.Pow(math.Sqrt(math.Log2(float64(x)+2))-1, 2)+1 ) } ; fmt.Println()
//   // streams count randomizer
//   a,b,c,d,e := 0,0,0,0,0
//   for x:=int64(0); x<1000; x++ { 
//     aaa := balance.BasicStats_StreamsCountAndModifier(funcs.Epoch())
//     if aaa == 2 {a++} else if aaa == 3 {b++} else if aaa == 4 {c++} else if aaa == 5 {d++} else {e++}
//   }
//   fmt.Println("    ||:",a,"\t|||:",b,"\t||||:",c,"\t|||||:",d,"\terr:",e)
// }


func TestFetch_Stats(t *testing.T) {
	Frame := plot.CleanFrame()
  var ad Player 
  var bc Player 
  _, _ = PlayerBorn(&ad, 2, 1, &Frame.Player), PlayerBorn(&bc, 2, 6, &Frame.Player)
  bc.Basics.ID.Born = ad.Basics.ID.Born
	diff, same := Fetch_Stats(bc.Basics,ad.Basics), Fetch_Stats(bc.Basics,bc.Basics)
	if diff.ID.Description[:6] == "ERROR!" { t.Errorf("%sFAIL Fetch_Stats%s(diff): %s", plot.E[2], plot.E[0], diff.ID.Description) } else { t.Logf("%sSUCCESS Fetch_Stats%s(diff):\n%+v\n", plot.E[1], plot.E[0], diff) }
	if same.ID.Description[:6] == "ERROR!" { t.Errorf("%sFAIL Fetch_Stats%s(same): %s", plot.E[2], plot.E[0], same.ID.Description) } else { t.Logf("%sSUCCESS Fetch_Stats%s(same):\n%+v\n", plot.E[1], plot.E[0], same) }
	fmt.Println("   --------------------------------------------------------------------------------------------------------------")
}

func TestGrow_Stats(t *testing.T) {
	Frame := plot.CleanFrame()
  var ef Player 
  var ad Player 
  var bc Player 
  _, _, _ = PlayerBorn(&ef, 3, 3, &Frame.Player), PlayerBorn(&ad, 3, 1, &Frame.Player), PlayerBorn(&bc, 3, 2, &Frame.Player)
  bc.Basics.ID.Born, ef.Basics.ID.Born = ad.Basics.ID.Born, ad.Basics.ID.Born  
  gh := Fetch_Stats(bc.Basics,ad.Basics)
	if gh.ID.Description[:6] == "ERROR!" { t.Errorf("%sFAIL%s Fetch_Stats(2-1): %s", plot.B, plot.E[0], gh.ID.Description) } else { t.Logf("%sSUCCESS%s Fetch_Stats(2-1):\n%+v\n", plot.B, plot.E[0], gh) }
	grow := Grow_Stats(ef.Basics,gh)
	if grow.ID.Description[:6] == "ERROR!" { t.Errorf("%sFAIL Grow_Stats%s(3+(2-1)): %s", plot.E[2], plot.E[0], grow.ID.Description) } else { t.Logf("%sSUCCESS Grow_Stats%s(3+(2-1)):\n%+v\n", plot.E[1], plot.E[0], grow) }
	fmt.Println("   --------------------------------------------------------------------------------------------------------------")
}

func TestTakeAll_Stats(t *testing.T) {
	Frame := plot.CleanFrame()
  var ef Player 
  var ad Player 
  var bc Player 
  _, _, _ = PlayerBorn(&ef, 4, 3, &Frame.Player), PlayerBorn(&ad, 4, 1, &Frame.Player), PlayerBorn(&bc, 4, 2, &Frame.Player)
  bc.Basics.ID.Born, ef.Basics.ID.Born = ad.Basics.ID.Born, ad.Basics.ID.Born  
  gh := Fetch_Stats(bc.Basics,ad.Basics)
	if gh.ID.Description[:6] == "ERROR!" { t.Errorf("%sFAIL%s Fetch_Stats(2-1): %s", plot.B, plot.E[0], gh.ID.Description) } else { t.Logf("%sSUCCESS%s Fetch_Stats(2-1):\n%+v\n", plot.B, plot.E[0], gh) }
  err := TakeAll_Stats(&ef, []BasicStats{gh, gh, gh}) 
	if err != nil { t.Errorf("%sFAIL TakeAll_Stats%s(3+(2-1)x3): %s", plot.E[2], plot.E[0], err) } else { t.Logf("%sSUCCESS TakeAll_Stats%s(3+(2-1)x3):\n%+v\n%+v\n", plot.E[1], plot.E[0], ef.Basics, ef.Attributes) }
	fmt.Println("=================================================================================================================")
}
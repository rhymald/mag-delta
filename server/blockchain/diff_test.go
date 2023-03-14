package blockchain

import (
	"time"
	"testing"
	"fmt"
	"math"
  "rhymald/mag-delta/client/plot"
)

type Result struct {
	Count int
	Time int
}

func TestCreateBlock(t *testing.T){
	begin := time.Now().UnixNano()
	stopper, trials, diffMax := int64(1000000000*60*60*10), 256, 20 // 1 hour
	results := make([]Result, diffMax)
	for diff := 0 ; diff <= diffMax ; diff++ {
		hashrate := 0
		prev := 0.0
		for trial := 0 ; trial < trials ; trial++ {
			date := time.Now().UnixNano()
			data := fmt.Sprintf("Data here: %d", date)
			block, hashcount := createBlock(data, "NS", []byte("Data"), diff, []byte("Last"), date)
			hashrate += hashcount
			results[diff].Time += int(time.Now().UnixNano() - date)
			results[diff].Count += 1
			now := float64(results[diff].Time/results[diff].Count/1000)/1000
			t.Logf("%08b..%08b Difficult: %d Block #%d time: %d ms. ", block.Hash[:3], block.Hash[61:64], diff, trial, (int(time.Now().UnixNano() - date))/1000000)
			if stopper <= time.Now().UnixNano() - begin || math.Abs(math.Log2(prev/now)) < math.Log2(1+1/float64(trials)) { t.Logf("%sFinished diff: %d, mean time: %.3f ms, hashrate: %.2f MH/sec%s", plot.E[1], diff, now, float64(hashrate)/now/1000, plot.E[0]) ; break}			
			prev = now 
		}
		if stopper <= time.Now().UnixNano() - begin {break}
	}
}
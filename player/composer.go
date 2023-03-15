package player 

// + split -s and +s
// + mean

import (
	"rhymald/mag-delta/funcs"
	"rhymald/mag-delta/balance"
	"time"
	"fmt"
	"errors"
)

func Fetch_Stats(last BasicStats, prev BasicStats) (BasicStats, error) {
	buffer := BasicStats{}
	// ID block
	if len(last.ID.Description) < 6 || len(prev.ID.Description) < 0 { buffer.ID.Description = last.ID.Description } else if last.ID.Description[:6] == "ERROR!" || prev.ID.Description[:6] == "ERROR!" { 
		buffer.ID.Description = "ERROR! A state is already ruined" ; return buffer, errors.New(buffer.ID.Description) 
	} else { buffer.ID.Description = last.ID.Description }
	if last.ID.NPC == prev.ID.NPC { buffer.ID.NPC = last.ID.NPC } else { buffer.ID.Description = "ERROR! Foreign stats: ID/NPC" ; return buffer, errors.New(buffer.ID.Description) } 
	if last.ID.Born == prev.ID.Born { buffer.ID.Born = last.ID.Born } else { buffer.ID.Description = "ERROR! Foreign stats: ID/Born" ; return buffer, errors.New(buffer.ID.Description) }
	if last.ID.Last >= prev.ID.Last { buffer.ID.Last = last.ID.Last } else { buffer.ID.Description = "ERROR! Newer stats: ID/Last" ; return buffer, errors.New(buffer.ID.Description) }
	// BODY 
	lelem, _ := funcs.ReStr(last.Body) 
	pelem, _ := funcs.ReStr(prev.Body) 
	sameElement := lelem == pelem
	if sameElement {
		stream := make(map[string][3]int)
		stream[lelem] = [3]int{ last.Body[lelem][0]-prev.Body[pelem][0], last.Body[lelem][1]-prev.Body[pelem][1], last.Body[lelem][2]-prev.Body[pelem][2] }
		buffer.Body = stream
	} else { buffer.ID.Description = "ERROR! Wrong body stream element" ; return buffer, errors.New(buffer.ID.Description) }
	// STREAMS
	for x:=0; x<len(last.Streams)&&x<len(prev.Streams); x++ {
		lelem, _ = funcs.ReStr(last.Streams[x]) 
		pelem, _ = funcs.ReStr(prev.Streams[x]) 
		sameElement = lelem == pelem
		inherited := pelem == funcs.Elements[0]
		if sameElement || inherited {
			stream := make(map[string][3]int)
			stream[lelem] = [3]int{ last.Streams[x][lelem][0]-prev.Streams[x][pelem][0], last.Streams[x][lelem][1]-prev.Streams[x][pelem][1], last.Streams[x][lelem][2]-prev.Streams[x][pelem][2] }
			buffer.Streams = append(buffer.Streams, stream)
		} else { buffer.ID.Description = "ERROR! Wrong stream element #"+fmt.Sprint(x) ; return buffer, errors.New(buffer.ID.Description) }
	}
	if len(last.Streams)-len(prev.Streams) > 0 {
		for x:=0; x<len(last.Streams)-len(prev.Streams); x++ { buffer.Streams = append(buffer.Streams, last.Streams[x+len(prev.Streams)]) }
	} else if len(last.Streams)-len(prev.Streams) < 0 { buffer.ID.Description = "ERROR! Wrong stream count (old bigger than new): "+fmt.Sprint(len(last.Streams)-len(prev.Streams)) ; return buffer, errors.New(buffer.ID.Description) }
	// ITEMS - tbd not implemented totally, add DROP / LOOT prefix to description
	// NEGATIVE check
	negativeCheckFail := false
	for _, each := range buffer.Streams {
		belem, _ := funcs.ReStr(each) 
		negativeCheckFail = negativeCheckFail || each[belem][0] < 0 || each[belem][1] < 0 || each[belem][2] < 0	
	}
	if negativeCheckFail { buffer.ID.Description = "ERROR! Negative check failed, stream grow can't be negative (body can)" ; return buffer, errors.New(buffer.ID.Description) }
	return buffer, nil
}

func Grow_Stats(last BasicStats, diff BasicStats) BasicStats {
	buffer := BasicStats{}
	// ID block
	if last.ID.Description[:6] == "ERROR!" || diff.ID.Description[:6] == "ERROR!" { buffer.ID.Description = "ERROR! Diff or state is already ruined" ; return buffer } else { buffer.ID.Description = diff.ID.Description }
	if last.ID.NPC == diff.ID.NPC { buffer.ID.NPC = diff.ID.NPC } else { buffer.ID.Description = "ERROR! Foreign stats: ID/NPC" ; return buffer } 
	if last.ID.Born == diff.ID.Born { buffer.ID.Born = diff.ID.Born } else { buffer.ID.Description = "ERROR! Foreign stats: ID/Born" ; return buffer }
	if last.ID.Last <= diff.ID.Last { buffer.ID.Last = diff.ID.Last } else { buffer.ID.Description = "ERROR! Fetch in older: ID/Last" ; return buffer }
	// BODY 
	lelem, _ := funcs.ReStr(last.Body) 
	delem, _ := funcs.ReStr(diff.Body) 
	sameElement := lelem == delem
	if sameElement {
		stream := make(map[string][3]int)
		stream[delem] = [3]int{ last.Body[lelem][0]+diff.Body[delem][0], last.Body[lelem][1]+diff.Body[delem][1], last.Body[lelem][2]+diff.Body[delem][2] }
		buffer.Body = stream
	} else { buffer.ID.Description = "ERROR! Wrong body stream element" ; return buffer }
	// Streams
	for x:=0; x<len(last.Streams)&&x<len(diff.Streams); x++ {
		lelem, _ = funcs.ReStr(last.Streams[x]) 
		delem, _ = funcs.ReStr(diff.Streams[x]) 
		sameElement = lelem == delem
		inherited := lelem == funcs.Elements[0]
		if sameElement || inherited {
			stream := make(map[string][3]int)
			stream[delem] = [3]int{ last.Streams[x][lelem][0]+diff.Streams[x][delem][0], last.Streams[x][lelem][1]+diff.Streams[x][delem][1], last.Streams[x][lelem][2]+diff.Streams[x][delem][2] }
			buffer.Streams = append(buffer.Streams, stream)
		} else { buffer.ID.Description = "ERROR! Wrong stream element #"+fmt.Sprint(x) ; return buffer }
	}
	if len(diff.Streams)-len(last.Streams) > 0 {
		for x:=0; x<len(diff.Streams)-len(last.Streams); x++ { buffer.Streams = append(buffer.Streams, diff.Streams[x+len(last.Streams)]) }
	} else if len(diff.Streams)-len(last.Streams) < 0 { buffer.ID.Description = "ERROR! Wrong stream count (old bigger than new): "+fmt.Sprint(len(diff.Streams)-len(last.Streams)) ; return buffer }
	// ITEMS - tbd not implemented totally, add DROP / LOOT prefix to description
	// NEGATIVE check
	belem, _ := funcs.ReStr(buffer.Body) 
	negativeCheckFail := buffer.Body[belem][0] <= 0 || buffer.Body[belem][1] <= 0 || buffer.Body[belem][2] <= 0
	for _, each := range buffer.Streams {
		belem, _ = funcs.ReStr(each) 
		negativeCheckFail = negativeCheckFail || each[belem][0] <= 0 || each[belem][1] <= 0 || each[belem][2] <= 0	
	}
	if negativeCheckFail { buffer.ID.Description = "ERROR! Negative check failed, streams and body can't be negative" }
	return buffer
}

func TakeAll_Stats(base *Player, improves []BasicStats) error {
	wasLogged := *&base.Attributes.Login
  *&base.Attributes.Login = false // to stop regen
  if !wasLogged { time.Sleep( time.Millisecond * time.Duration( balance.Regeneration_DefaultTimeout() )) }
	for index, each := range improves { if each.ID.Description[:6] == "ERROR!" { return errors.New(fmt.Sprintf("ERROR: Composing stats failed - some stats are invalid: %d", index)) } }
	buffer := *&base.Basics
	for index, each := range improves {
		buffer = Grow_Stats(buffer, each)
		if buffer.ID.Description[:6] == "ERROR!" { return errors.New(fmt.Sprintf("ERROR: Applying stats failed - final invalid: %d", index)) }
	}
	*&base.Basics = buffer
	// CalculateAttributes_FromBasics(base)
	return nil
}
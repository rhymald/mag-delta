package player 

// + split -s and +s
// + mean
// + sum

import (
	"rhymald/mag-delta/funcs"
)

// ID struct {
// 	NPC bool `json:"NPC"`
// 	Description string `json:"Name,omitempty"`
// 	Born int64 `json:"Born,omitempty"`
// 	Last int64 `json:"Last,omitempty"`
// } `json:"ID,omitempty"`
// Body funcs.Stream `json:"Body"`
// Streams []funcs.Stream `json:"Streams,omitempty"`
// Items []funcs.Stream `json:"Items,omitempty"`

func Fetch_Stats(last BasicStats, prev BasicStats) BasicStats {
	buffer := BasicStats{}
	// ID block
	if last.ID.Description[:6] == "ERROR!" || prev.ID.Description[:6] == "ERROR!" { buffer.ID.Description = "ERROR! A state is already ruined" ; return buffer } else { buffer.ID.Description = last.ID.Description }
	if last.ID.NPC == prev.ID.NPC { buffer.ID.NPC = last.ID.NPC } else { buffer.ID.Description = "ERROR! Foreign stats: ID/NPC" ; return buffer } 
	if last.ID.Born == prev.ID.Born { buffer.ID.Born = last.ID.Born } else { buffer.ID.Description = "ERROR! Foreign stats: ID/Born" ; return buffer }
	if last.ID.Last > prev.ID.Last { buffer.ID.Last = last.ID.Last } else { buffer.ID.Description = "ERROR! Newer stats: ID/Last" ; return buffer }
	// BODY 
	lelem, _ := funcs.ReStr(last.Body) 
	pelem, _ := funcs.ReStr(prev.Body) 
	sameElement := lelem == pelem
	if sameElement {
		stream := make(map[string][3]int)
		stream[lelem] = [3]int{ last.Body[lelem][0]-prev.Body[pelem][0], last.Body[lelem][1]-prev.Body[pelem][1], last.Body[lelem][2]-prev.Body[pelem][2] }
		buffer.Body = stream
	} else { buffer.ID.Description = "ERROR! Wrong body sream element" ; return buffer }
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
		} else { buffer.ID.Description = "ERROR! Wrong stream stream element #"+string(x) ; return buffer }
	}
	if len(last.Streams)-len(prev.Streams) > 0 {
		for x:=0; x<len(last.Streams)-len(prev.Streams); x++ { buffer.Streams = append(buffer.Streams, last.Streams[x+len(prev.Streams)]) }
	} else if len(last.Streams)-len(prev.Streams) < 0 { buffer.ID.Description = "ERROR! Wrong stream count (old bigger than new): "+string(len(last.Streams)-len(prev.Streams)) ; return buffer }
	// ITEMS - tbd not implemented totally, add DROP / LOOT prefix to description
	return buffer
}
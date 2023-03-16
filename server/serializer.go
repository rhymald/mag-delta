package server

import(
  "fmt"
  "encoding/base64"
  "encoding/json"
  "rhymald/mag-delta/player"
)

func toJson(thing interface{}) string {
  b, err := json.Marshal(thing)
  if err != nil { fmt.Println(err) ; return "" }
  encoded := base64.StdEncoding.EncodeToString(b)
  return encoded
}

func statsFromJson(code string, thing player.BasicStats) player.BasicStats {
  copy := &thing
  decoded, _ := base64.StdEncoding.DecodeString(code)
  err := json.Unmarshal(decoded, copy)
  if err != nil { fmt.Println("Stats read failed:", err, string(decoded)) ; return thing }
  return *copy
}

func stateFromJson(code string, thing player.CharStatus) player.CharStatus {
  copy := &thing
  decoded, _ := base64.StdEncoding.DecodeString(code)
  err := json.Unmarshal(decoded, copy)
  if err != nil { fmt.Println("Stats read failed:", err, string(decoded)) ; return thing }
  return *copy
}

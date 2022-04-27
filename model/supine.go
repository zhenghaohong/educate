package model
// {"versions":"0.0.1","mac":"94B5552C2C14","cmd":1,"type":2,"data":{"radar_cnt":5}}

// 仰卧起坐
type SupineResponse struct {
	Versions string `json:"versions"`
	Mac      string `json:"mac"`
	Cmd      int    `json:"cmd"`
	Type     int    `json:"type"`
	Data     struct {
		RadarCnt int `json:"radar_cnt"`  // 仰卧起坐次数
	}
}
package model

// {"versions":"0.0.1","mac":"94B55526BF14","cmd":1,"type":3,"data":{"respiratory":258}}
// 硬件返回的数据
type LungCapacityResponse struct {
	Versions string `json:"versions"`
	Mac      string `json:"mac"`
	Cmd      int    `json:"cmd"`
	Type     int    `json:"type"`
	Data     struct {
		Respiratory int `json:"respiratory"`
	}	
}





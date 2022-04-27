package model

// weight/up/94B55525F1EC

// 接收mqtt
type HeightWeightReceive struct {
	Versions string `json:"versions"`
	Mac      string `json:"mac"`
	Cmd      int    `json:"cmd"`
	Type     int    `json:"type"`
	Data     struct {
		Weight []int `json:"weight"`
	}
}


// 返回给websocket
type HeightWeightResponse struct {
	Versions string `json:"versions"`
	Mac      string `json:"mac"`
	Cmd      int    `json:"cmd"`
	Type     int    `json:"type"`
	Data     struct {
		Weight float64 `json:"weight"`
		Height int `json:"height"`
	}
}

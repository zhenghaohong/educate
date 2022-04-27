package model

// websocket order
// {
// 	"order":1 		  // 1 开始  0 结束考试
// 	"mac":"x001",   // 设备编号
// 	"type":"jump"  //  项目类型
//   }

type WebSocketOrder struct {
	Cmd  int    `json:"cmd"`
	Mac  string `json:"mac"`
	Type int    `json:"type"`
	Data struct {
		// Laser1 int `json:"laser1"`
		// Laser2 int `json:"laser2"`
		StudentNum string `json:"studentNum"`
	}
	// Data interface{} `json:"data"`
}

type WebSocketRespone struct {
	Cmd  int    `json:"cmd"`
	Mac  string `json:"mac"`
	Type int    `json:"type"`
	Data struct {
		Laser1 int `json:"laser1"`
		Laser2 int `json:"laser2"`
	}
}

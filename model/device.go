package model

// {
// "mac":"004D098030C4", //设备mac
// "cmd":1,                        //备用命令
// "type":1,                        //设备类型
// "data"
//  {
//     "state":1   ,//1:开始 0:结束
//     "unit_time":   n, // 倒计时内就不发送了
//     "unit_value":  Y, // 距离
// }
// }
type JumpConfigController struct {
	Request
	Data struct {
		State     int `json:"state"`
		UnitTime  int `json:"unit_time"`
		UnitValue int `json:"unit_value"`
	}
}

// {
// 	"versions" : "0.0.1",
// 	"mac" : "004D098030C4",
// 	"cmd" : 1,
// 	"type" : 1,
// 	"data" : {
// 	  "laser1" : 0,
// 	  "laser2" : 16,
// 	  "radar1" : 0
// 	}
//   }


// {
// 	"mac":"004D098030C4", //设备mac
// 	"cmd":1,                        //备用命令
// 	"type":1,                        //设备类型
// 	"data"
// 	 {
// 		"state":0   ,//1:开始 0:结束
// 	}
// 	}

type DeviceResponseData struct {
	Response
	Data struct {
		State int `json:"state"`
	}
}

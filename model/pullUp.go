package model

//  {"versions":"0.0.1","mac":"94B55526BF14","cmd":1,"type":3,"data":{"radar_cnt":258}}
type PullUpResponse struct {
	Versions string `json:"versions"`
	Mac      string `json:"mac"`
	Cmd      int    `json:"cmd"`
	Type     int    `json:"type"`
	Data     struct {
		RadarCnt int `json:"radar_cnt"` // 引体向上次数
	}
}

// "state":1   ,//1:开始 0:结束
//     "farthest_distance":   n,
//     "shortest_distance":   y,
//     " last_time": x

// 配置指令 开始
type PullUpStartRequest struct {
	Versions string `json:"versions"`
	Mac      string `json:"mac"`
	Cmd      int    `json:"cmd"`
	Type     int    `json:"type"`
	Data     struct {
		State            int `json:"state"`
		FarthestDistance int `json:"farthest_distance"`
		ShortestDistance int `json:"shortest_distance"`
		LastTime         int `json:"last_time"`
	}
}

// 配置指令 结束
type PullUpStopRequest struct {
	Versions string `json:"versions"`
	Mac      string `json:"mac"`
	Cmd      int    `json:"cmd"`
	Type     int    `json:"type"`
	Data     struct {
		State            int `json:"state"`
	}
}



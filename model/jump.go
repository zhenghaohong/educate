package model


type JumpRespone struct {
	Cmd  int    `json:"cmd"`
	Mac  string `json:"mac"`
	Type int    `json:"type"`
	Data struct {
		Laser1 int `json:"laser1"`
		Laser2 int `json:"laser2"`
	}
}


type JumpReceive struct {
	Response
	Data struct {
		// State  int `json:"state"`
		Laser1 int `json:"laser1"`
		Laser2 int `json:"laser2"`
		Radar1 int `json:"radar1"`
	}
}
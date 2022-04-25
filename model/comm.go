package model

// common Request represents the request data
type Request struct {
	Mac  string `json:"mac"`
	Cmd  int    `json:"cmd"`
	Type int    `json:"type"`
	// Data interface{} `json:"data"`
}

// common Response represents the response data
type Response struct {
	Mac  string `json:"mac"`
	Cmd  int    `json:"cmd"`
	Type int    `json:"type"`
	// Data interface{} `json:"data"`
}

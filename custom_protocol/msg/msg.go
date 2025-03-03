package msg

//go:generate go run go_study/custom_protocol/msg_general $GOFILE

type Message struct {
	Id      uint32  `msg_tag:"uint32,1,id"`
	Name    string  `msg_tag:"string,2,name"`
	Coins   float32 `msg_tag:"float32,3,coins"`
	Damoind float64 `msg_tag:"float64,4,damoind"`
}

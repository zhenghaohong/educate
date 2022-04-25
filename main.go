package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"phyFit/mqtt/db"
	"phyFit/mqtt/model"
	"phyFit/mqtt/public"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"golang.org/x/net/websocket"

	_ "phyFit/mqtt/db"
)

var cliPool = make(map[string]*websocket.Conn, 0)
var mu *sync.RWMutex
var broker = "mqtt.calfkaka.com"
var port = 1883
var UserName = "emqx"
var Password = "public"
var ClientId = "go_mqtt_client"
var WebSocketPort = ":8083"

var ()

var c MQTT.Client

func init() {
	opts := MQTT.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetUsername(UserName)
	opts.SetPassword(Password)
	opts.SetClientID(ClientId)

	opts.SetDefaultPublishHandler(f)
	c = MQTT.NewClient(opts)

}

func main() {
	mu = &sync.RWMutex{}

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.Subscribe("test/online/data", 0, nil); token.Wait() &&
		token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	// jump1: send start order to device for mqtt
	if token := c.Subscribe("ZA_jump/up/004D098030C4", 0, nil); token.Wait() &&
		token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	// jump1: receive data to device for mqtt
	if token := c.Subscribe("ZA_jump/down/004D098030C4", 0, nil); token.Wait() &&
		token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	if token := c.Subscribe("websocket/device/control", 0, nil); token.Wait() &&
		token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	go startWebSocket()
	select {}

}

func startWebSocket() {
	//接受websocket的路由地址
	http.Handle("/ws", websocket.Handler(WsHandle))
	if err := http.ListenAndServe(WebSocketPort, nil); err != nil {
		fmt.Printf("webSocket 失败:%+v", err)
		log.Fatal("ListenAndServe:", err)
	}
}

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	var (
		b []byte
	)

	// payLoadStr := string(msg.Payload())
	// fmt.Printf("payLoadStr:%+v\n", payLoadStr)

	var req model.Request
	err := json.Unmarshal(msg.Payload(), &req)
	if err != nil {
		fmt.Printf("device format json.Unmarshal err:%+v\n", err)
		// return
	}
	fmt.Println("接收数据", string(msg.Payload()))

	// websocket topic
	var webSocketOrder model.WebSocketOrder
	err = json.Unmarshal(msg.Payload(), &webSocketOrder)
	if err != nil {
		fmt.Printf("webSocketOrder format json.Unmarshal err:%+v\n", err)
	}

	var jumpReceive model.JumpReceive
	var jumpReceiveArrs []model.JumpReceive
	// if webSocketOrder.Order == 1 && webSocketOrder.Type == 1 { // jump start order 1是jump
	// jump data
	if req.Type == 1 { // jump

		err = json.Unmarshal(msg.Payload(), &jumpReceive)
		if err != nil {
			fmt.Printf("jumpReceive format json.Unmarshal err:%+v\n", err)
		} else {
			s, err := db.GetKey("jump")

			if err == nil { // 取到值
				fmt.Printf("取到数据 开始累加数据啦:%+v\n", s)
				jumpReceiveArrByte, err := json.Marshal(&jumpReceive) // 将结构体转换成json
				if err != nil {
					fmt.Printf("json.Marshal failed, err:%v\n", err)
					return
				}

				db.AddCache("jump", s+","+string(jumpReceiveArrByte)) // 先拼接在添加上去

				fmt.Println("jumpReceiveArrsStr:", string(jumpReceiveArrByte))

				s_arr := "[" + s + "]"
				err = json.Unmarshal([]byte(s_arr), &jumpReceiveArrs)
				fmt.Printf("累计存取数据啦:%+v\n", jumpReceiveArrs)
				if err == nil {
					fmt.Printf("len:%+v\n", len(jumpReceiveArrs))
					if len(jumpReceiveArrs) > 20 { // 大于20条数据删掉
						db.DeleteCache("jump")
					}
				} else {
					fmt.Printf("jumpReceiveArrs 错误 err:%+v\n", err)
				}

			} else {
				jumpReceiveArrsStr, _ := json.Marshal(jumpReceive)
				db.AddCache("jump", string(jumpReceiveArrsStr))
			}
		}

	}

	if req.Cmd == 2 && req.Type == 1 { // jump end order
		// 统计出最小成绩 发给客户端 然后删掉数据

		// 结束指令 发送保存的数据到websocket 然后删除掉它
		// 发送至mqtt device 并且删除redis key 为jump
		s, err := db.GetKey("jump")
		if err != nil { // 咩有取到key
			b = []byte("请先点开始测试")
			for _, v := range cliPool { // 重新写入websocket
				if v == nil {
					fmt.Println("结束")
					continue
				}
				_, err := v.Write(b)
				if err != nil {
					fmt.Printf("写入body 失败:%+v", err)
				} else {
					fmt.Println("写入成功")
				}
			}
		} else { // 取到key 拼成数组返回给客户端
			fmt.Println("取到key  拼成数组返回给客户端")
			s_arr := "[" + s + "]"
			_ = json.Unmarshal([]byte(s_arr), &jumpReceiveArrs)
			laser1, laser2 := GetJumpMin(jumpReceiveArrs)

			fmt.Printf("laser1:%+v\n", laser1)
			fmt.Printf("laser2:%+v\n", laser2)
			var websocketResponse model.WebSocketRespone
			websocketResponse.Cmd = 3
			websocketResponse.Type = 1
			websocketResponse.Mac = "jump001"
			websocketResponse.Data.Laser1 = laser1
			websocketResponse.Data.Laser2 = laser2

			websocketResponseByte, _ := json.Marshal(&websocketResponse)
			result_jump := string(websocketResponseByte)
			b = []byte(result_jump)
			for _, v := range cliPool { // 重新写入websocket
				if v == nil {
					fmt.Println("结束")
					continue
				}
				_, err := v.Write(b)
				if err != nil {
					fmt.Printf("写入body 失败:%+v", err)
				} else {
					fmt.Println("写入成功")
					// 删除redis key
					db.DeleteCache("jump")
				}
			}

		}

	}

	// 接收的数据全部发到websocket
	// b = msg.Payload()
	// for _, v := range cliPool {
	// 	if v == nil {
	// 		fmt.Println("结束")
	// 		continue
	// 	}
	// 	_, err := v.Write(b)
	// 	if err != nil {
	// 		fmt.Printf("写入body 失败:%+v", err)
	// 	} else {
	// 		fmt.Println("写入成功")
	// 	}
	// }
	// fmt.Printf("\n-----cliPool----%+v\n", cliPool)
}

func WsHandle(ws *websocket.Conn) {
	var err error
	t := fmt.Sprintf("%d", time.Now().UnixNano())
	cliPool[t] = ws

	var webSocketOrder model.WebSocketOrder
	for {
		var reply string
		// websocket接受信息
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("receive failed:", err)
			break
		}

		// 发送 至 mqtt

		// 请求头为 setUp
		fmt.Println("received from client: " + reply)
		//msg := "received:" + reply
		msg := reply
		fmt.Println("send to client:" + msg)
		fmt.Printf("\n----mqtt send to client----%+v\n", cliPool)

		// websocket 指令
		err = json.Unmarshal([]byte(reply), &webSocketOrder)
		if err != nil {
			fmt.Println("websocket 指令 json.Unmarshal err:", err)
		}
		switch webSocketOrder.Cmd {
		case 1:
			// start
			fmt.Println("websocket start order")
			public.WebsocketStartTopic("004D098030C4", c)
			public.WebsocketJumpStartOrderToTopic(c)

		case 2:
			// stop
			fmt.Println("websocket stop order")
			public.WebsocketStopTopic("004D098030C4", c)
			public.WebsocketJumpStopOrderToTopic(c)
		default:

		}
		//push message
		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"), " send failed:", err)
			break
		}
	}
	// mu.Lock()
	// delete(cliPool, t)
	// fmt.Println("删除cli 成功:", t)
	// mu.Unlock()

}

// var jumpReceiveArrs []model.JumpReceive
// laser1[] 最小值
// laser2[] 最小值
func GetJumpMin(jumpReceiveArrs []model.JumpReceive) (int, int) {
	var laser1_int_array = []int{}
	var laser2_int_array = []int{}
	for _, j := range jumpReceiveArrs {
		laser1_int_array = append(laser1_int_array, j.Data.Laser1)
		laser2_int_array = append(laser2_int_array, j.Data.Laser2)
	}
	laser1 := GetMin(laser1_int_array)
	laser2 := GetMin(laser2_int_array)
	return laser1, laser2
}

func GetMin(l []int) (min int) {
	min = l[0]
	for _, v := range l {
		if v <= 0 {
			continue
		}
		if v < min {
			min = v
		}
	}
	return
}

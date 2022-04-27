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
	"phyFit/mqtt/config"

	"phyFit/mqtt/tools/zlog"
	"phyFit/mqtt/util"
)

var cliPool = make(map[string]*websocket.Conn, 0)
var mu *sync.RWMutex

// var broker = "mqtt.calfkaka.com"
var broker = config.MQTTConfig.IP
var port = config.MQTTConfig.Port
var UserName = "emqx"
var Password = "public"
var ClientId = "go_mqtt_client"
var WebSocketPort = config.WsConfig.PORT




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
	zlog.Debugf("启动成功 :%+v",mu)
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

	// lung_capacity
	if token := c.Subscribe("Lung_capacity/up/94B55526BF14", 0, nil); token.Wait() &&
		token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}


	// 坐位体前屈
	if token := c.Subscribe("ZA_trunk/up/94B55526CDE4", 0, nil); token.Wait() &&
		token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// 身高体重
	if token := c.Subscribe("weight/up/94B55525F1EC", 0, nil); token.Wait() &&
		token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// 仰卧起坐
	if token := c.Subscribe("ZA_supine/up/94B5552C2C14", 0, nil); token.Wait() &&
		token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// 引体向上
	if token := c.Subscribe("ZA_chinning/up/94B5552C1D68", 0, nil); token.Wait() &&
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

	// var (
	// 	b []byte
	// )

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
	var (
		jumpReceive model.JumpReceive
		jumpReceiveArrs []model.JumpReceive

		lungCapacityResponse model.LungCapacityResponse
		trunkResponse model.TrunkResponse
		heightWeightReceive model.HeightWeightReceive // mqtt
		heightWeightResponse model.HeightWeightResponse // redis

		supineResponse model.SupineResponse // 仰卧起坐
		pullUpResponse model.PullUpResponse
	)

	// if webSocketOrder.Order == 1 && webSocketOrder.Type == 1 { // jump start order 1是jump
	// jump data
	if req.Type == 1 { // jump

		err = json.Unmarshal(msg.Payload(), &jumpReceive)
		if err != nil {
			fmt.Printf("jumpReceive format json.Unmarshal err:%+v\n", err)
		} else {
			s, err := db.GetKey("jump")

			if err == nil { // 取到值
				// fmt.Printf("取到数据 开始累加数据啦:%+v\n", s)
				jumpReceiveArrByte, err := json.Marshal(&jumpReceive) // 将结构体转换成json
				if err != nil {
					fmt.Printf("json.Marshal failed, err:%v\n", err)
					return
				}

				db.AddCache("jump", s+","+string(jumpReceiveArrByte)) // 先拼接在添加上去

				// fmt.Println("jumpReceiveArrsStr:", string(jumpReceiveArrByte))

				s_arr := "[" + s + "]"
				err = json.Unmarshal([]byte(s_arr), &jumpReceiveArrs)
				// fmt.Printf("累计存取数据啦:%+v\n", jumpReceiveArrs)
				if err == nil {
					// fmt.Printf("len:%+v\n", len(jumpReceiveArrs))
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
	// 仰卧起坐
	if req.Type == 2 {
		err = json.Unmarshal(msg.Payload(), &supineResponse)
		if err != nil {

		}else {
			s, err := db.GetKey("supine")
			if err == nil {
				fmt.Println("获取仰卧起坐数据成功",s)
				supineResponseByte, _ := json.Marshal(&supineResponse) // 将结构体转换成json
				db.AddCache("supine",string(supineResponseByte)) // 替换掉值
			}else{
				fmt.Println("获取不到仰卧起坐数据",s)
				supineResponseByte, _ := json.Marshal(&supineResponse) 
				db.AddCache("supine", string(supineResponseByte)) 
			}
		}

	}
	// 引体向上
	if req.Type == 8 {
		err = json.Unmarshal(msg.Payload(), &pullUpResponse)
		if err != nil {

		}else {
			s, err := db.GetKey("pullUp")
			if err == nil {
				fmt.Println("获取引体向上数据成功",s)
				pullUpResponseByte, _ := json.Marshal(&pullUpResponse) // 将结构体转换成json
				db.AddCache("pullUp",string(pullUpResponseByte)) // 替换掉值
			}else{
				fmt.Println("获取不到引体向上数据",s)
				pullUpResponseByte, _ := json.Marshal(&pullUpResponse) 
				db.AddCache("pullUp", string(pullUpResponseByte)) 
			}
		}

	}

	// 身高体重
	if req.Type == 3 {
		err = json.Unmarshal(msg.Payload(), &heightWeightReceive)
		fmt.Printf("heightWeightReceive:%+v\n", heightWeightReceive)
		if err != nil {

		}else {

			_, _ = db.GetKey("heightWeight")
				heightWeightResponse.Versions = "0.0.1"
				heightWeightResponse.Mac = "heightWeight01"
				heightWeightResponse.Cmd = 2
				heightWeightResponse.Type = 3
				heightWeightResponse.Data.Weight = 0
				heightWeightResponse.Data.Height = 0
				if len(heightWeightReceive.Data.Weight) > 0 {
					weight,height:= util.HexToString(heightWeightReceive.Data.Weight)
					fmt.Println("weight:",weight)
					fmt.Println("height:",height)
					heightWeightResponse.Data.Weight,_ = util.StringToFloat(weight)
					heightWeightResponse.Data.Height,_ = util.StringToInt(height)
				}
				heightWeightResponseByte, _ := json.Marshal(&heightWeightResponse) // 将结构体转换成json
				db.AddCache("heightWeight",string(heightWeightResponseByte)) // 替换掉值
			
		}


	}
	// 坐位体前屈
	if req.Type == 4 {
		err = json.Unmarshal(msg.Payload(), &trunkResponse)
		if err != nil {

		}else{
			s, err := db.GetKey("trunk")
			if err == nil {
				trunkResponseByte, _ := json.Marshal(&trunkResponse) // 将结构体转换成json
				db.AddCache("trunk", s+","+string(trunkResponseByte)) // 替换掉值
			}else{
				trunkResponseByte, _ := json.Marshal(&trunkResponse) 
				db.AddCache("trunk", string(trunkResponseByte)) 
			}
		}
	}
	// 肺活量
	if req.Type == 5 { 
		err = json.Unmarshal(msg.Payload(), &lungCapacityResponse)
		if err != nil {
			fmt.Printf("lungCapacityResponse format json.Unmarshal err:%+v\n", err)
			zlog.Errorf("lungCapacityResponse format json.Unmarshal err:%+v\n", err)
		}else{
			s, err := db.GetKey("lungCapacity")
			if err == nil { // 取到值 替换掉值
				lungCapacityResponseByte, _ := json.Marshal(&lungCapacityResponse) // 将结构体转换成json
				db.AddCache("lungCapacity", s+","+string(lungCapacityResponseByte)) // 替换掉值

				// websocket 收到结束指令去redis中取值

			} else { // 存key 值  studentNum 从websocket 获取
				// studentNum := ""
				lungCapacityResponseByte, _ := json.Marshal(&lungCapacityResponse) 
				db.AddCache("lungCapacity", string(lungCapacityResponseByte)) 
			}
		}
	}





	// 跑步
	if req.Type == 6 {}
	// 跳绳
	if req.Type == 7 { }




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
	var (
		webSocketOrder model.WebSocketOrder
		lungCapacityResponse model.LungCapacityResponse
		lungCapacityResponseArr []model.LungCapacityResponse


		trunkResponse model.TrunkResponse
		trunkResponseArr []model.TrunkResponse

		heightWeightResponse model.HeightWeightResponse
		// heightWeightResponseArr []model.HeightWeightResponse
	)
	
	for {
		var reply string
		// websocket接受信息
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("receive failed:", err)
			break
		}

		// 发送 至 mqtt

		fmt.Println("received from client: " + reply)
		//msg := "received:" + reply
		// msg := reply
		// fmt.Println("send to client:" + msg)
		fmt.Printf("\n----mqtt send to client----%+v\n", cliPool)

		// websocket 指令
		err = json.Unmarshal([]byte(reply), &webSocketOrder)
		if err != nil {
			fmt.Println("websocket 指令 json.Unmarshal err:", err)
		}
		fmt.Printf("\n----websocket 指令----%+v\n", webSocketOrder)
		switch webSocketOrder.Cmd {
		case 1: // start		
			switch webSocketOrder.Type {
			case 1: // jump
				// start
				fmt.Println("websocket jump start order")
				public.WebsocketStartTopic("004D098030C4", c)
				public.WebsocketJumpStartOrderToTopic(c)
			case 2: // 仰卧起坐 supine
			fmt.Println("websocket supine start order")
				public.SupineStartOrderToTopic("94B5552C2C14",c)

			case 3: // heightWeight

			case 4: // trunk 坐位体前屈


			case 5: // lungCapacity
				// start
				fmt.Println("websocket lungCapacity start order")
				// db.AddCache("lungCapacity", "无啊")
			case 8: // pullUp 引体向上
				public.PutUpStartOrderToTopic("94B5552C1D68",c)

			}


		case 2: // stop
			switch webSocketOrder.Type {
			case 1: // jump
					fmt.Println("websocket stop order")
					public.WebsocketStopTopic("004D098030C4", c)
					public.WebsocketJumpStopOrderToTopic(c)

						// 统计出最小成绩 发给客户端 然后删掉数据
						var jumpRespone model.JumpRespone
						var jumpResponeArr []model.JumpRespone
						jumpRespone.Cmd = 3
						jumpRespone.Type = 1
						jumpRespone.Mac = "jump001"
						// 结束指令 发送保存的数据到websocket 然后删除掉它
						s, err := db.GetKey("jump")
						fmt.Printf("\n----s----%+v\n", s)
						if err != nil { // 咩有取到key
							jumpRespone.Data.Laser1 = 0
							jumpRespone.Data.Laser2 = 0
							
							b,_:= json.Marshal(jumpRespone)
						
							if err = websocket.Message.Send(ws, string(b)); err != nil {
								fmt.Println(time.Now().Format("2006-01-02 15:04:05"), " send jumpRespone failed:", err)
								break
							}
							fmt.Println("jump stop send no data success")

					
						} else { // 取到key 拼成数组返回给客户端
							s_arr := "[" + s + "]"
							_ = json.Unmarshal([]byte(s_arr), &jumpResponeArr)

							fmt.Printf("jumpResponeArr:%+v\n", jumpResponeArr)
							laser1, laser2 := GetJumpMin(jumpResponeArr)
				
							fmt.Printf("laser1:%+v\n", laser1)
							fmt.Printf("laser2:%+v\n", laser2)
							
							
							jumpRespone.Data.Laser1 = laser1
							jumpRespone.Data.Laser2 = laser2
							// jumpResponeByte, _ := json.Marshal(&jumpRespone)
							// result_jump := string(jumpResponeByte)
							// b = []byte(result_jump)

							b,_:= json.Marshal(jumpRespone)
							
							if err = websocket.Message.Send(ws, string(b)); err != nil {
								fmt.Println(time.Now().Format("2006-01-02 15:04:05"), " send jumpRespone failed:", err)
								break
							}else{
								db.DeleteCache("jump")
							}					
				
						}
					
			case 2: // 仰卧起坐 supine
					public.SupineStopOrderToTopic("94B5552C2C14",c)			
					var supineRespone model.SupineResponse
					supineRespone.Cmd = 3
					supineRespone.Type = 2
					supineRespone.Mac = "supine"
					// 结束指令 发送保存的数据到websocket 然后删除掉它
					s, err := db.GetKey("supine")
					if err != nil {
						supineRespone.Data.RadarCnt = 0
						b,_ := json.Marshal(supineRespone)
						if err = websocket.Message.Send(ws, string(b)); err != nil {
							fmt.Println(time.Now().Format("2006-01-02 15:04:05"), " send failed:", err)
							break
						}
					}else{
						if err = websocket.Message.Send(ws, string(s)); err != nil {
							fmt.Println(time.Now().Format("2006-01-02 15:04:05"), " send failed:", err)
							break
						}else{
							db.DeleteCache("supine")
						}
					}


			case 3: // heightWeight
					heightWeightResponse.Versions = "0.0.1"
					trunkResponse.Mac = "heightWeight001"
					heightWeightResponse.Cmd = 2
					heightWeightResponse.Type = 3
					fmt.Println("websocket height stop order")
					s, err := db.GetKey("heightWeight")
					if err != nil {
						heightWeightResponse.Data.Weight = 0
						heightWeightResponse.Data.Height = 0
						b,_:= json.Marshal(heightWeightResponse)
						
						if err = websocket.Message.Send(ws, string(b)); err != nil {
							fmt.Println(time.Now().Format("2006-01-02 15:04:05"), " send trunkResponse failed:", err)
							break
						}
						fmt.Println("heightweight stop send no data success")

					}else{

						// b,_ := json.Marshal(heightWeightResponse)
						if err = websocket.Message.Send(ws, string(s)); err != nil {
							fmt.Println(time.Now().Format("2006-01-02 15:04:05"), " send trunkResponse failed:", err)
							break
						}else{
							db.DeleteCache("heightWeight")
						}
						fmt.Println("heightweight stop send no data success")
					}

			case 4: // trunk 坐位体前屈
					trunkResponse.Versions = "0.0.1"
					trunkResponse.Mac = "trunk"
					trunkResponse.Cmd = 2
					trunkResponse.Type = 4
					fmt.Println("websocket trunk stop order")
					s, err := db.GetKey("trunk")
					if err != nil {
						trunkResponse.Data.Respiratory = 0
						b,err:= json.Marshal(trunkResponse)
						
						if err = websocket.Message.Send(ws, string(b)); err != nil {
							fmt.Println(time.Now().Format("2006-01-02 15:04:05"), " send trunkResponse failed:", err)
							break
						}
						fmt.Printf("该学生还没有坐位体前屈数据")
					}else{
						s_arr := "[" + s + "]"
						_ = json.Unmarshal([]byte(s_arr), &trunkResponseArr)
						max1 := 0
						fmt.Printf("s_arr:%+v\n", s_arr)
						fmt.Printf("trunkResponseArr:%+v\n", trunkResponseArr)

						if len(trunkResponseArr) > 0 {
							max1 = GetTrunkMax(trunkResponseArr)
						} 
						fmt.Printf("trunk max1:%+v\n", max1)
						// {versions: "0.0.1",mac: "94B55526BF14",cmd: 1,type: 5,Data: {respiratory: 548}}
						
						trunkResponse.Data.Respiratory = max1
						b,err := json.Marshal(trunkResponse)
						if err != nil {
							fmt.Println("json.Marshal err:", err)
						}
					
						if err = websocket.Message.Send(ws, string(b)); err != nil {
							fmt.Println(time.Now().Format("2006-01-02 15:04:05"), " send trunkResponse failed:", err)
							break
						}else{
							db.DeleteCache("trunk")
						}

					}
			case 5: //肺活量
					// 读取key 并返回给客户端 
					// studentNum:="studentNum"
					lungCapacityResponse.Versions = "0.0.1"
					lungCapacityResponse.Mac = "lungCapacity"
					lungCapacityResponse.Cmd = 2
					lungCapacityResponse.Type = 5
					s, err := db.GetKey("lungCapacity")
					if err != nil {
						lungCapacityResponse.Data.Respiratory = 0
						b,err := json.Marshal(lungCapacityResponse)
						if err = websocket.Message.Send(ws, string(b)); err != nil {
							fmt.Println(time.Now().Format("2006-01-02 15:04:05"), " send failed:", err)
							break
						}
						fmt.Printf("请测试肺活量数据")
					}else {
						// 计算最大值 
						s_arr := "[" + s + "]"
						_ = json.Unmarshal([]byte(s_arr), &lungCapacityResponseArr)
						// max1 := 0
						fmt.Printf("s_arr:%+v\n", s_arr)
						fmt.Printf("lungCapacityResponseArr:%+v\n", lungCapacityResponseArr)

						// if len(lungCapacityResponseArr) > 0 {
						// 	max1 = GetLungCapacityMax(lungCapacityResponseArr)
						// } 
							
						// 发送所有记录 
						
							b,err := json.Marshal(lungCapacityResponseArr)
							if err = websocket.Message.Send(ws, string(b)); err != nil {
								fmt.Println(time.Now().Format("2006-01-02 15:04:05"), " send failed:", err)
								break
							}else{
								// 保存至mysql
								// 删除redis key
								db.DeleteCache("lungCapacity")
							}
						
						

						// fmt.Printf("max1:%+v\n", max1)
						// // {versions: "0.0.1",mac: "94B55526BF14",cmd: 1,type: 5,Data: {respiratory: 548}}
						
						// lungCapacityResponse.Data.Respiratory = max1
						// b,err := json.Marshal(lungCapacityResponse)
						// if err != nil {
						// 	fmt.Println("json.Marshal err:", err)
						// }
						// if err = websocket.Message.Send(ws, string(b)); err != nil {
						// 	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), " send failed:", err)
						// 	break
						// }else{
						// 	// 保存至mysql
						// 	// 删除redis key
						// 	db.DeleteCache("lungCapacity")
						// }
						
					}				
				case 8: //引体向上		
					public.PutUpStopOrderToTopic("94B5552C1D68",c)			
					var pullUpRespone model.PullUpResponse
					pullUpRespone.Cmd = 3
					pullUpRespone.Type = 8
					pullUpRespone.Mac = "pullUp"
					// 结束指令 发送保存的数据到websocket 然后删除掉它
					s, err := db.GetKey("pullUp")
					if err != nil {
						pullUpRespone.Data.RadarCnt = 0
						b,_ := json.Marshal(pullUpRespone)
						if err = websocket.Message.Send(ws, string(b)); err != nil {
							fmt.Println(time.Now().Format("2006-01-02 15:04:05"), " send failed:", err)
							break
						}
					}else{
						if err = websocket.Message.Send(ws, string(s)); err != nil {
							fmt.Println(time.Now().Format("2006-01-02 15:04:05"), " send failed:", err)
							break
						}else{
							db.DeleteCache("pullUp")
						}
					}

				default:

		}
		//push message
		// if err = websocket.Message.Send(ws, msg); err != nil {
		// 	fmt.Println(time.Now().Format("2006-01-02 15:04:05"), " send failed:", err)
		// 	break
		// }
	}
	// mu.Lock()
	// delete(cliPool, t)
	// fmt.Println("删除cli 成功:", t)
	// mu.Unlock()

	}
}


// var jumpReceiveArrs []model.JumpReceive
// laser1[] 最小值
// laser2[] 最小值
func GetJumpMin(jumpReceiveArrs []model.JumpRespone) (int, int) {
	var laser1_int_array = []int{}
	var laser2_int_array = []int{}
	for _, j := range jumpReceiveArrs {
		if j.Data.Laser1 > 0 {
			laser1_int_array = append(laser1_int_array, j.Data.Laser1)
		}
		if j.Data.Laser2 > 0 {
			laser2_int_array = append(laser2_int_array, j.Data.Laser2)
		}
		
	}
	fmt.Printf("laser1_int_array:%+v\n", laser1_int_array)
	fmt.Printf("laser2_int_array:%+v\n", laser2_int_array)
	laser1 := GetMin(laser1_int_array)
	laser2 := GetMin(laser2_int_array)
	fmt.Printf("laser1:%+v\n", laser1)
	fmt.Printf("laser2:%+v\n", laser2)
	return laser1, laser2
}


func GetTrunkMax(trunkResponseArrs []model.TrunkResponse) int {
	var trunk_int_array = []int{}
	for _, j := range trunkResponseArrs {
		trunk_int_array = append(trunk_int_array, j.Data.Respiratory)
	}
	max1 := GetMax(trunk_int_array)
	return max1
}


func GetLungCapacityMax(lungCapacityResponseArrs []model.LungCapacityResponse) int {
	var lungCapacity_int_array = []int{}
	for _, j := range lungCapacityResponseArrs {
		lungCapacity_int_array = append(lungCapacity_int_array, j.Data.Respiratory)
	}
	max1 := GetMax(lungCapacity_int_array)
	return max1
}

func GetMin(arr []int) (min int) {
	if len(arr) == 0 {
		return 0
	}
	min = arr[0]
	for _, v := range arr {
		if v > 0 {
			if v < min {
				min = v
			}
		}
	}
	return
}

func GetMax(l []int) (max int) {
	max = l[0]
	for _, v := range l {
		if v <= 0 {
			continue
		}
		if v > max {
			max = v
		}
	}
	return
}
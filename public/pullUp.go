package public


import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"

	"phyFit/mqtt/config"
	"phyFit/mqtt/util"
)

// PullUpFarthestDistance
// PullUpShortestDistance
// PullUpLastTime
  
// 引体向上  开始指令
func PutUpStartOrderToTopic(mac string, client MQTT.Client) {
	farthestDistance:= util.IntToString(config.DeviceConfigConfig.PullUpFarthestDistance)
	shortestDistance:= util.IntToString(config.DeviceConfigConfig.PullUpShortestDistance)
	lastTime:= util.IntToString(config.DeviceConfigConfig.PullUpLastTime)

	topic := "ZA_chinning/down/" + mac
	msg := `{"mac":"` + mac + `", ` +
		` "cmd":1, "type":8, "data":{"state":1 ,"farthest_distance":`+ farthestDistance + `,`+
		`"shortest_distance":`+ shortestDistance +`,`+
		`"last_time":`+lastTime+`}}`
	fmt.Printf("putUp start config : %s, msg: %s\n", topic, msg)
	token := client.Publish(topic, 0, false, []byte(msg))
	if token.Error() != nil {
		fmt.Printf("send to WebsocketStartTopic server start err : %+v", token.Error())
		return
	} else {
		fmt.Println("send to WebsocketStartTopic server start success ")
	}
	token.Wait()
	// time.Sleep(time.Second)
}

// 引体向上停止指令
func PutUpStopOrderToTopic(mac string, client MQTT.Client) {
	topic := "ZA_chinning/down/" + mac
	msg := `{"mac":"` + mac + `", ` +
		` "cmd":1, "type":8, "data":{"state":0}}`
	fmt.Printf("putUp stop config : %s, msg: %s\n", topic, msg)
	token := client.Publish(topic, 0, false, []byte(msg))
	if token.Error() != nil {
		fmt.Printf("send to WebsocketStartTopic server start err : %+v", token.Error())
		return
	} else {
		fmt.Println("send to WebsocketStartTopic server start success ")
	}
	token.Wait()
	// time.Sleep(time.Second)
}

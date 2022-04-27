package public
import (
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	"phyFit/mqtt/config"
	"phyFit/mqtt/util"

)

// SupineFarthestDistance
// SupineShortestDistance
// SupineLastTime
// 仰卧起坐开始指令
func SupineStartOrderToTopic(mac string, client MQTT.Client) {
	farthestDistance:= util.IntToString(config.DeviceConfigConfig.SupineConfigFarthestDistance)
	shortestDistance:= util.IntToString(config.DeviceConfigConfig.SupineConfigShortestDistance)
	lastTime:= util.IntToString(config.DeviceConfigConfig.SupineConfigLastTime)
	topic := "ZA_supine/down/" + mac
	msg := `{"mac":"` + mac + `", ` +
		` "cmd":1, "type":2, "data":{"state":1 ,`+
		`"farthest_distance":`+ farthestDistance + `,`+
		`"shortest_distance":`+shortestDistance+`,`+
		`"last_time":`+lastTime+`}}`
	fmt.Printf("supine start config : %s, msg: %s\n", topic, msg)
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

// 仰卧起坐停止指令
// SupineConfigFarthestDistance
// SupineConfigShortestDistance
// SupineConfigLastTime 

func SupineStopOrderToTopic(mac string, client MQTT.Client) {
	topic := "ZA_supine/down/" + mac
	msg := `{"mac":"` + mac + `", ` +
		` "cmd":1, "type":2, "data":{"state":0}}`
	fmt.Printf("supine stop config : %s, msg: %s\n", topic, msg)
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



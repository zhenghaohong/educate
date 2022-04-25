package public

import (
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// websocket Controller topic
func WebsocketStartTopic(mac string, client MQTT.Client) {
	topic := "ZA_jump/down/" + mac
	msg := `{"mac":"004D098030C4", "cmd":1, "type":1, "data":{"state":1 ,"unit_time": 11,"unit_value":  5}}`
	fmt.Printf("websocket send to mqtt to device topic: %s, msg: %s\n", topic, msg)
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

func WebsocketStopTopic(mac string, client MQTT.Client) {
	topic := "ZA_jump/down/" + mac
	msg := `{"mac":"004D098030C4", "cmd":1, "type":1, "data":{"state":0 }}`
	fmt.Printf("websocket send to mqtt to device stop topic: %s, msg: %s\n", topic, msg)
	token := client.Publish(topic, 0, false, []byte(msg))
	if token.Error() != nil {
		fmt.Printf("send to WebsocketStartTopic server stop err : %+v", token.Error())
		return
	} else {
		fmt.Println("send to WebsocketStartTopic server stop success ")
	}
	token.Wait()
	// time.Sleep(time.Second)
}

// websocket to websocket to mqtt
func WebsocketJumpStartOrderToTopic(client MQTT.Client) {
	topic := "websocket/device/control"
	msg := ` {"cmd":1,"mac":"x001","type":1}`
	fmt.Printf("jump Start Order : %s, msg: %s\n", topic, msg)
	token := client.Publish(topic, 0, false, []byte(msg))
	if token.Error() != nil {
		fmt.Printf("Jump server start order err : %+v", token.Error())
		return
	} else {
		fmt.Println("Jump  server start Order success ")
	}
	token.Wait()
	// time.Sleep(time.Second)
}

func WebsocketJumpStopOrderToTopic(client MQTT.Client) {
	topic := "websocket/device/control"
	msg := ` {"cmd":2,"mac":"x001","type":1}`
	fmt.Printf("jump End Order : %s, msg: %s\n", topic, msg)
	token := client.Publish(topic, 0, false, []byte(msg))
	if token.Error() != nil {
		fmt.Printf("Jump server stop order err : %+v", token.Error())
		return
	} else {
		fmt.Println("Jump  server stop Order success ")
	}
	token.Wait()
}

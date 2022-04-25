package public

import (
	"fmt"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func PublicToTopic(topic string, client MQTT.Client, msg string) {
	token := client.Publish(topic, 0, false, []byte(msg))
	if token.Error() != nil {
		fmt.Printf("send to topic server err : %+v", token.Error())
		return
	} else {
		fmt.Println("send to topic server success ")
	}
	token.Wait()
	time.Sleep(time.Second)
}

// device Controller topic
func DeviceControllerToTopic(topic string, client MQTT.Client, msg string) {
	topic = "JA_jump/up/" + topic
	token := client.Publish(topic, 0, false, []byte(msg))
	if token.Error() != nil {
		fmt.Printf("send to topic server err : %+v", token.Error())
		return
	} else {
		fmt.Println("send to topic server success ")
	}
	token.Wait()
	time.Sleep(time.Second)
}

func GetMac() string {
	return "123456789"
}

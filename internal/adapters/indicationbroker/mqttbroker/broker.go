package mqttbroker

import (
	"fmt"
	"log/slog"
	"encoding/json"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"gitlab.phystech.pro/kimmybaez/sensorstray/pkg/settings"
	"gitlab.phystech.pro/kimmybaez/sensorstray/internal/indications/dtos"
)

type indicationsMQTTBroker struct {
	client mqtt.Client
}

// Используем singleton т.к наше приложение и есть client у брокера и мы не должны определять множество клиентов у брокера
var brokerInstancePtr *indicationsMQTTBroker
var once sync.Once

func CreateNewIndicationMQTTBroker() *indicationsMQTTBroker {
	once.Do(func() {
		var brokerURI string = fmt.Sprintf("tcp://%s:%s", settings.GetSettings().BrokerHost, settings.GetSettings().BrokerPort)

		options := mqtt.NewClientOptions()
		options.AddBroker(brokerURI)
		options.SetUsername(settings.GetSettings().BrokerUser)
		options.SetPassword(settings.GetSettings().BrokerPassword)
		options.SetClientID("tray_client")

		client := mqtt.NewClient(options)

		// token - это по факту объект подключения, ожидаем подключение проверяем на наличие ошибок и обрабатываем их
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			slog.Error("Failed to connect to mqtt broker. Can't proceed", "Options", options, "Error", token.Error())
			panic(token.Error())
		}

		brokerInstancePtr = &indicationsMQTTBroker{client: client}
	})

	return brokerInstancePtr
}

func (broker *indicationsMQTTBroker) SendIndications(data dtos.IndicationsDTO) error {
	// Превращаем наше DTO в байт строку для отправки
	messageBody, err := json.Marshal(data)
	if err != nil {
		slog.Error("Failed to marshall indication structure into json", "Structure", data, "Error", err)
		return err
	}

	// Отправляем в брокер и  в случае ошибки обрабаитываем её
	token := broker.client.Publish("sensors", 0, false, messageBody)
	token.Wait()
	if token.Error() != nil {
		slog.Error("Failed to push message into broker", "Error", token.Error())
		return token.Error()
	}
	
	return nil
}
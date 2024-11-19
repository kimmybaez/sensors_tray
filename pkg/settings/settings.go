package settings

import (
	"log/slog"
	"os"
	"sync"
	"encoding/json"
)

type settings struct {
	BrokerHost string `json:"broker_host"`
	BrokerPort string `json:"broker_port"`
	BrokerUser string `json:"broker_user"`
	BrokerPassword string `json:"broker_password"`
}

// Будем использовать Singleton т.к нам необходимо четко быть увереннымми,что объект настроек инстанцирован лишь 1 раз
var insatancePtr *settings
var once sync.Once

func GetSettings() *settings {
	once.Do(func() {
		// Читаем файл config.json в масив байт file
		file, err := os.ReadFile("./config.json")
		if err != nil {
			slog.Error("Failed to open config file. Can't instantiate the application", "Error", err)
			panic(err)
		}

		var settingsInstance settings
		err = json.Unmarshal(file, &settingsInstance)
		if err != nil {
			slog.Error("Failed to unmarshalling file content into settings structure", "Content", string(file), "Error", err)
			panic(err)
		}

		insatancePtr = &settingsInstance
	})

	return insatancePtr
}

func GetIcon() []byte {
	icon, err := os.ReadFile("./assets/icon.svg")
	if err != nil {
		slog.Error("Failed to open icon file", "Error", err)
		return make([]byte, 0)
	}

	return icon
}
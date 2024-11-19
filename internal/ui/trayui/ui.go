package trayui

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/getlantern/systray"
	"gitlab.phystech.pro/kimmybaez/sensorstray/internal/indications/dtos"
	"gitlab.phystech.pro/kimmybaez/sensorstray/internal/indications/ports"
	"gitlab.phystech.pro/kimmybaez/sensorstray/pkg/settings"
)

type trayUIApp struct {
	broker ports.IndicationBroker
	indicationsCommand ports.Command[dtos.IndicationsDTO]
}

var trayAppPtr *trayUIApp
var once sync.Once

func CreateNewTrayUiApp(broker ports.IndicationBroker, command ports.Command[dtos.IndicationsDTO]) *trayUIApp {
	once.Do(func() {
		trayAppPtr = &trayUIApp{
			broker: broker,
			indicationsCommand: command,
		}
	})

	return trayAppPtr
}

func (app *trayUIApp) StartApp() {
	systray.SetTitle("Fiztech-Climatic")
	systray.SetTooltip("Fiztech-Climatic")
	systray.SetIcon(settings.GetIcon())

	indicationsButton := systray.AddMenuItem("Получить показатели", "Получить показатели")
	systray.AddSeparator()
	quitButton := systray.AddMenuItem("Выйти", "Выйти")

	go func() {
		for {
			select {
			case <-indicationsButton.ClickedCh:
				go app.GetIndications()
			case <-quitButton.ClickedCh:
				systray.Quit()
            	return
			}
		}
	}()
}

func (app *trayUIApp) CloseApp() {
	app.indicationsCommand.Close()
}

func (app *trayUIApp) GetIndications() {
	slog.Info("Get indication button pressed")

	go app.indicationsCommand.Execute()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	indications, err := app.indicationsCommand.GetResult(ctx)
	if err != nil {
		slog.Error("Error while getting inidication", "Error", err)
		return
	}
	
	err = app.broker.SendIndications(*indications)
	if err != nil {
		slog.Error("Error while sending indication to broker", "Error", err)
		return
	}
}
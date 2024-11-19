package serialcommands

import (
	"log/slog"
	"encoding/json"
	"context"

	"gitlab.phystech.pro/kimmybaez/sensorstray/internal/indications/dtos"
	"gitlab.phystech.pro/kimmybaez/sensorstray/internal/indications/ports"
)

/*
В структуре команды есть:
1. Communicator - порт с помощью которого мы общаемся с датчиком. Communicator - может считаться ресивером
2. ExecutionError - ошибка полученная при попытке исполнить команду
3. Канал с результатами
*/
type GetIndicationsCommand struct {
	communicator ports.Communicator
	executionError error
	result chan *dtos.IndicationsDTO
}

func CreateGetCommand(communicator ports.Communicator) *GetIndicationsCommand {
	return &GetIndicationsCommand{
		communicator: communicator,
		result: make(chan *dtos.IndicationsDTO, 1),
	}
}

func (command *GetIndicationsCommand) Execute() {
	// Отправляем команду комуникаторы и обрабатываем возможную ошибку
	resp, err := command.communicator.SendCommand("/indications\n")
	if err != nil {
		slog.Error("Failed to send /indications command", "Error", err)
		command.executionError = err
		return
	}

	// Переводим ответ от датчика из байтового масива в структуру, в случае ошибки обрабатываем её
	var data dtos.IndicationsDTO
	err = json.Unmarshal(resp, &data)
	if err != nil {
		slog.Error("Failed to unmarshall /indications command response", "Response", string(resp), "Error", err)
		command.executionError = err
		return
	}

	// Записываем в структуру в качестве результата
	command.result <- &data
}

func (command *GetIndicationsCommand) GetResult(ctx context.Context) (*dtos.IndicationsDTO, error) {
	// Если в процессе была записана какая-либо ошибка возвращаем её и пустой результат
	if command.executionError != nil {
		return nil, command.executionError
	}


	// Какой из каналов первый отправит результат тот и вернет, в случае если команда уже была успешно выполнена
	// Первым отправит канал result, если же нет и вышел таймаут контекста будет возвращен пустой результат и ошибка контекста
	select {
		case result := <-command.result:
			return result, nil
		case <-ctx.Done():
			return nil, ctx.Err()
	}

}

func (command *GetIndicationsCommand) Close() {
	command.communicator.Close()
}
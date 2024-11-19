package ports

import (
	"context"
	"gitlab.phystech.pro/kimmybaez/sensorstray/internal/indications/dtos"
)

type CommandDTO interface {
	dtos.IndicationsDTO
}

type Command[T CommandDTO] interface {
	Execute()
	GetResult(ctx context.Context) (*T, error)
	Close()
}

type Communicator interface {
	SendCommand(command string) ([]byte, error)
	Close()
}

type IndicationBroker interface {
	SendIndications(data dtos.IndicationsDTO) error
}
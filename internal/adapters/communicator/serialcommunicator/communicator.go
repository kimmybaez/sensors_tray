package serialcommunicator

import (
	"log/slog"

	"github.com/tarm/serial"
)

type SerialCommunicator struct {
	portName string
	portBaud int
	port *serial.Port
}

func CreateNewSerialCommunicator() (*SerialCommunicator, error) {
	config := serial.Config{Name: "/dev/ttyACM0", Baud: 9600}

	port, err := serial.OpenPort(&config)
	if err != nil {
		slog.Error("Failed to open serial port", "Port Name", config.Name, "Baud", config.Baud, "Error", err)
		return nil, err
	}

	return &SerialCommunicator{portName: config.Name, portBaud: config.Baud, port: port}, nil
}

func (communicator *SerialCommunicator) SendCommand(command string) ([]byte, error) {
	// Записываем в serial порт n байт, который получили путем перевода команды в байт строку
	_, err := communicator.port.Write([]byte(command))
	if err != nil {
		slog.Error("Failed to write command into serial port", "Port", communicator.port, "Baud", communicator.portBaud, "Command", command, "Error", err)
		return nil, err
	}

	// Создаем буфер куда запишем ответ из порта
	buffer := make([]byte, 128)

	// Читаем n байт из serial порта в буфер
	n, err := communicator.port.Read(buffer)
	if err != nil {
		slog.Error("Failed to read command response from serial port", "Port", communicator.portName, "Baud", communicator.portBaud, "Command", command, "Error", err)
		return nil, err
	}

	return buffer[:n], nil
}

func (communcator *SerialCommunicator) Close() {
	communcator.port.Close()
}
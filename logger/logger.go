package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogEntry содержит поля для записи лога
type LogEntry struct {
	Timestamp    string   // Время запроса
	Route        string   // Тип запроса (sendOne, sendBulk)
	User         string   // Имя пользователя
	IP           string   // IP-адрес клиента
	Sender       string   // Отправитель
	MSISDN       []string // Номера получателей
	Status       string   // Статус выполнения (success, failed)
	Errors       []string // Ошибки при отправке
	SuccessCount int      // Количество успешных отправок
	FailedCount  int      // Количество неудачных отправок
	Duration     string   // Длительность выполнения запроса
}

// Logger структура для логирования
type Logger struct {
	file *os.File
}

// NewLogger создаёт новый логгер с ротацией по дате
func NewLogger() (*Logger, error) {
	currentDate := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("logs/%s.log", currentDate)

	// Создаём папку logs, если её нет
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{file: file}, nil
}

// Log записывает запись лога в файл
func (l *Logger) Log(entry LogEntry) {
	log.SetOutput(l.file)
	log.Printf(
		"[Timestamp: %s] Route=%s, User=%s, IP=%s, Sender=%s, MSISDN=%v, Status=%s, Errors=%v, SuccessCount=%d, FailedCount=%d, Duration=%s\n",
		entry.Timestamp, entry.Route, entry.User, entry.IP, entry.Sender, entry.MSISDN, entry.Status, entry.Errors, entry.SuccessCount, entry.FailedCount, entry.Duration,
	)
}

// Close закрывает лог-файл
func (l *Logger) Close() {
	l.file.Close()
}

package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"smpp-sender/logger" // Импортируем наш логгер
	"smpp-sender/smpp"
	"time"
)

type SendOneRequest struct {
	Server   string `json:"server"`
	Username string `json:"username"`
	Password string `json:"password"`
	Sender   string `json:"sender"`
	Msisdn   string `json:"msisdn"`
	Message  string `json:"message"`
}

type SendOneResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func SendOneHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	var req SendOneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(SendOneResponse{
			Status: "error",
			Error:  "Invalid request body",
		})
		return
	}

	// Инициализация логгера
	log, err := logger.NewLogger()
	if err != nil {
		http.Error(w, "Failed to initialize logger", http.StatusInternalServerError)
		return
	}
	defer log.Close()

	// Обновление файла конфигурации
	configFile, err := os.OpenFile("config/config.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		http.Error(w, "Failed to open config file", http.StatusInternalServerError)
		return
	}
	defer configFile.Close()

	updatedConfig := map[string]string{
		"server":   req.Server,
		"username": req.Username,
		"password": req.Password,
		"sender":   req.Sender,
	}
	json.NewEncoder(configFile).Encode(updatedConfig)

	// Отправка SMS
	err = smpp.SendSMS(req.Server, req.Username, req.Password, req.Sender, req.Msisdn, req.Message)

	// Подготовка для логирования
	status := "success"
	errors := []string{}
	if err != nil {
		status = "failed"
		errors = append(errors, err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(SendOneResponse{
			Status: "error",
			Error:  "Failed to send message: " + err.Error(),
		})
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SendOneResponse{
			Status:  "success",
			Message: "Message sent successfully",
		})
	}

	duration := time.Since(startTime).String()

	// Логирование запроса и результата
	log.Log(logger.LogEntry{
		Timestamp:    time.Now().Format("2006-01-02 15:04:05"),
		Route:        "sendOne",
		User:         req.Username,
		IP:           r.RemoteAddr,
		Sender:       req.Sender,
		MSISDN:       []string{req.Msisdn},
		Status:       status,
		Errors:       errors,
		SuccessCount: 1,
		FailedCount:  len(errors),
		Duration:     duration,
	})
}

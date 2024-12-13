package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"smpp-sender/logger" // Импортируем наш логгер
	"smpp-sender/smpp"
	"time"
)

type SendBulkRequest struct {
	Server   string   `json:"server"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Sender   string   `json:"sender"`
	Msisdn   []string `json:"msisdn"`
	Message  string   `json:"message"`
}

type SendBulkResponse struct {
	Total        int      `json:"total"`
	SuccessCount int      `json:"success_count"`
	FailedCount  int      `json:"failed_count"`
	Errors       []string `json:"errors"`
}

func SendBulkHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	var req SendBulkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Инициализация логгера
	log, err := logger.NewLogger()
	if err != nil {
		http.Error(w, "Failed to initialize logger", http.StatusInternalServerError)
		return
	}
	defer log.Close()

	var (
		successCount int
		failedCount  int
		errors       []string
	)

	for _, msisdn := range req.Msisdn {
		err := smpp.SendSMS(req.Server, req.Username, req.Password, req.Sender, msisdn, req.Message)
		if err != nil {
			failedCount++
			errors = append(errors, fmt.Sprintf("Failed for %s: %v", msisdn, err))
			continue
		}
		successCount++
	}

	// Подготовка ответа
	response := SendBulkResponse{
		Total:        len(req.Msisdn),
		SuccessCount: successCount,
		FailedCount:  failedCount,
		Errors:       errors,
	}

	// Логирование запроса и результата
	duration := time.Since(startTime).String()
	status := "success"
	if failedCount > 0 {
		status = "partial"
	}

	log.Log(logger.LogEntry{
		Timestamp:    time.Now().Format("2006-01-02 15:04:05"),
		Route:        "sendBulk",
		User:         req.Username,
		IP:           r.RemoteAddr,
		Sender:       req.Sender,
		MSISDN:       req.Msisdn,
		Status:       status,
		Errors:       errors,
		SuccessCount: successCount,
		FailedCount:  failedCount,
		Duration:     duration,
	})

	// Отправка ответа клиенту
	w.Header().Set("Content-Type", "application/json")
	if failedCount > 0 {
		w.WriteHeader(http.StatusPartialContent) // 206 Partial Content для частичных ошибок
	} else {
		w.WriteHeader(http.StatusOK)
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Write(responseJSON)
}

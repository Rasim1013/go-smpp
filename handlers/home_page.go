package handlers

import (
	"encoding/json"
	"net/http"
)

func HomeHandlers(w http.ResponseWriter, r *http.Request) {
	// Пример данных для возврата
	response := map[string]string{
		"message": "THE BOOK!!",
	}

	// Установка заголовков
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Кодирование и отправка JSON
	json.NewEncoder(w).Encode(response)
}

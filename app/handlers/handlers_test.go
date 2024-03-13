package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	database "github.com/ddlifter/BashAPI/db"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func TestRunCommandError(t *testing.T) {
	// Подготавливаем фейковую базу данных
	db := database.Database()
	defer db.Close()

	// Создаем фейковый HTTP запрос
	req, _ := http.NewRequest("GET", "/command/9999", nil)
	vars := map[string]string{"id": "9999"}
	req = mux.SetURLVars(req, vars)

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	// Вызываем функцию RunCommand
	RunCommand(rr, req)

	// Проверяем статус код ответа
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}
func TestRunCommand(t *testing.T) {
	req, err := http.NewRequest("GET", "/command/555", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RunCommand)

	// Вызываем RunCommand параллельно несколько раз
	for i := 0; i < 5; i++ {
		go func() {
			handler.ServeHTTP(rr, req)
		}()
	}

	// Ждем завершения всех горутин
	for i := 0; i < 5; i++ {
		<-time.After(1 * time.Second) // Ждем 1 секунду между запросами
	}
}

func TestGetCommands(t *testing.T) {
	req, err := http.NewRequest("GET", "/commands", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetCommands)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, expectedContentType)
	}
}

func TestDeleteCommands(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/commands", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteCommands)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, expectedContentType)
	}
}

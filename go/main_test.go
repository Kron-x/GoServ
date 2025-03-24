package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
)

func TestHomeHandler(t *testing.T) {
	// Создаем тестовый запрос
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HomeHandler)

	// Вызываем обработчик напрямую
	handler.ServeHTTP(rr, req)

	// Проверяем статус код
	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получили %d", rr.Code)
	}

	// Проверяем Content-Type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/html; charset=UTF-8" {
		t.Errorf("Неверный Content-Type: %s", contentType)
	}
}

func TestNewDimensionHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/new-dimension", nil)
	rr := httptest.NewRecorder()
	
	http.HandlerFunc(NewDimensionHandler).ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("AboutHandler: неверный статус-код: %v", rr.Code)
	}
}

func TestSubmitTextHandler(t *testing.T) {
	// Подготовка JSON-тела запроса
	jsonBody := `{"text":"hello"}`

	req := httptest.NewRequest("POST", "/submit-text", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	SubmitTextHandler(rr, req) // Ваш обработчик для /submit-text

	// Проверки
	if rr.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получили %d", rr.Code)
	}

	// Проверка JSON-ответа
	expected := `{"message":"Вы ввели: hello"}`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("Ожидался ответ %s, получили %s", expected, rr.Body.String())
	}
}
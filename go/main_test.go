package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/username/GoServ/go/pkg"
)

// Тест для функции loadConfig
'''
func TestLoadConfig(t *testing.T) {
    // Создаем временный config.json

    config := server.LoadConfig()

    if config.Port != "5000" {
        t.Errorf("Expected port 5000, got %s", config.Port)
    }
    if config.ImagesDir != "images" {
        t.Errorf("Expected images_dir: ./images, got: %s", config.ImagesDir)
    }
    if config.LogFile != "logs/server.log" {
        t.Errorf("Expected log_file: ./logs/server.log, got: %s", config.LogFile)
    }
}
'''
// Тест для homeHandler
func TestHomeHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(server.HomeHandler)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expected := "<h1>Welcome to the Image Server</h1>"
    if rr.Body.String()[:len(expected)] != expected {
        t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
    }
}

// Тест для submitTextHandler
func TestSubmitTextHandler(t *testing.T) {
    requestBody := map[string]string{"text": "Hello, World!"}
    jsonBody, _ := json.Marshal(requestBody)

    req, err := http.NewRequest("POST", "/submit-text", bytes.NewBuffer(jsonBody))
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(server.SubmitTextHandler)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    var response map[string]string
    err = json.Unmarshal(rr.Body.Bytes(), &response)
    if err != nil {
        t.Fatal(err)
    }

    expected := "Вы ввели: Hello, World!"
    if response["message"] != expected {
        t.Errorf("handler returned unexpected message: got %v want %v", response["message"], expected)
    }
}

// Тест для NewDimensionHandler
func TestNewDimensionHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/new-dimension", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(server.NewDimensionHandler)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expected := "<h1>Вы попали в новое измерение, пристегнитесь</h1>"
    if rr.Body.String()[:len(expected)] != expected {
        t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
    }
}

// Тест для healthHandler
func TestHealthHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/health", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(server.HealthHandler)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expected := "OK"
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
    }
}
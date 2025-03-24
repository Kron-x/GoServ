package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
)

// логирования запросов
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

func SubmitTextHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Чтение JSON-тела запроса
	var request struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Утилизация текста (например, просто возвращаем его обратно)
	response := map[string]string{
		"message": "Вы ввели: " + request.Text,
	}

	// Отправка ответа
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	fmt.Fprintf(w, `
		<h1>Welcome to the Image Server</h1>
		<img src="/images/image1.jpg" alt="Image 1" width="300">
		<img src="/images/image2.jpg" alt="Image 2" width="300">
        <div style="position: fixed; bottom: 20px; right: 20px;">
			<button onclick="window.location.href='/new-dimension'">Перейти в новое измерение</button>
		</div>
        <div style="text-align: center; margin-top: 20px;">
			<input type="text" id="textInput" placeholder="Введите текст" style="padding: 10px; font-size: 16px;">
			<button onclick="submitText()" style="padding: 10px; font-size: 16px;">Отправить</button>
		</div>
		<div id="result" style="text-align: center; margin-top: 20px; font-size: 18px;"></div>
		<script>
			function submitText() {
				const input = document.getElementById("textInput");
				const result = document.getElementById("result");
				const text = input.value.trim();

				if (text) {
					fetch("/submit-text", {
						method: "POST",
						headers: {
							"Content-Type": "application/json",
						},
						body: JSON.stringify({ text: text }),
					})
					.then(response => response.json())
					.then(data => {
						result.textContent = "Ответ сервера: " + data.message;
						input.value = "";
					})
					.catch(error => {
						result.textContent = "Ошибка: " + error.message;
					});
				} else {
					result.textContent = "Пожалуйста, введите текст.";
				}
			}

			// Обработка нажатия Enter
			document.getElementById("textInput").addEventListener("keyup", function(event) {
				if (event.key === "Enter") {
					submitText();
				}
			});
		</script>
	`)
}

func NewDimensionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(w, `
		<h1>Вы попали в новое измерение, пристегнитесь</h1>
		<p>Текущая дата и время: %s</p>
	    `, currentTime)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

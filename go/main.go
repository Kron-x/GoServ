package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "log"
    "io/ioutil"
    "os/signal"
    "syscall"
    "time"
    "context"
)


type Config struct {
    Port      string `json:"port"`
    ImagesDir string `json:"images_dir"`
    LogFile   string `json:"log_file"`
}

func loadConfig() Config {
    file, err := ioutil.ReadFile("config.json")
    if err != nil {
        log.Fatalf("Failed to read config file: %v", err)
    }
    var config Config
    err = json.Unmarshal(file, &config)
    if err != nil {
        log.Fatalf("Failed to parse config file: %v", err)
    }
    return config
}

// логирования запросов
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

func submitTextHandler(w http.ResponseWriter, r *http.Request) {
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

func homeHandler(w http.ResponseWriter, r *http.Request) {
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

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func main() {
    config := loadConfig()

    // Настройка логгера
    logFile, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatalf("Failed to open log file: %v", err)
    }
    defer logFile.Close()
    log.SetOutput(logFile)

    log.Printf("Starting server on port %s", config.Port)

    http.Handle("/", loggingMiddleware(http.HandlerFunc(homeHandler)))
    http.Handle("/images/", loggingMiddleware(http.StripPrefix("/images/", http.FileServer(http.Dir(config.ImagesDir)))))
    http.HandleFunc("/new-dimension", NewDimensionHandler)
    http.HandleFunc("/health", healthHandler)
    http.HandleFunc("/submit-text", submitTextHandler)

    server := &http.Server{
        Addr: ":" + config.Port,
    }

    // Graceful shutdown
    go func() {
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server error: %v", err)
        }
    }()

    log.Println("Server is running...")

    // Ожидание сигнала для graceful shutdown
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

    <-stop
    log.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Server shutdown error: %v", err)
    }

    log.Println("Server stopped.")
}

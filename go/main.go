package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "log"
    "io/ioutil"
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

    http.Handle("/", loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html")
        fmt.Fprintf(w, `
            <h1>Welcome to the Image Server</h1>
            <img src="/images/image1.jpg" alt="Image 1" width="300">
            <img src="/images/image2.jpg" alt="Image 2" width="300">
        `)
    })))

    http.Handle("/images/", loggingMiddleware(http.StripPrefix("/images/", http.FileServer(http.Dir(config.ImagesDir)))))

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "OK")
    })

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

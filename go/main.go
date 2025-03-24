package main

import (
    "net/http"
    "os"
    "log"
    "os/signal"
    "syscall"
    "time"
    "context"
)

func main() {
    config := LoadConfig()

    // Настройка логгера
    logFile, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatalf("Failed to open log file: %v", err)
    }
    defer logFile.Close()
    log.SetOutput(logFile)

    log.Printf("Starting server on port %s", config.Port)

    http.Handle("/", LoggingMiddleware(http.HandlerFunc(HomeHandler)))
    http.Handle("/images/", LoggingMiddleware(http.StripPrefix("/images/", http.FileServer(http.Dir(config.ImagesDir)))))
    http.HandleFunc("/new-dimension", NewDimensionHandler)
    http.HandleFunc("/health", HealthHandler)
    http.HandleFunc("/submit-text", SubmitTextHandler)

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

package main

import (
	"budget-recommendations/internal/handlers"
	"budget-recommendations/internal/service"
	"budget-recommendations/internal/storage/postgres"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Получаем DSN из переменных окружения
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:password@localhost:5432/budget?sslmode=disable"
	}

	// Инициализируем хранилище
	storage, err := postgres.NewStorage(dsn)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer storage.Close()

	// Инициализируем сервис
	transactionService := service.NewTransactionService(storage)

	// Инициализируем обработчики
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Настраиваем роутер
	router := mux.NewRouter()

	// Middleware для логирования
	router.Use(loggingMiddleware)

	// Роуты
	router.HandleFunc("/health", transactionHandler.HealthCheck).Methods("GET")
	router.HandleFunc("/transactions", transactionHandler.CreateTransaction).Methods("POST")
	router.HandleFunc("/transactions", transactionHandler.GetTransactions).Methods("GET")
	router.HandleFunc("/transactions/{id}", transactionHandler.GetTransactionByID).Methods("GET")

	// Настраиваем сервер
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Ожидаем сигналов для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
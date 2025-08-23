package main

import (
	"context"
	"log"
	"myapp/internal/cache"
	"myapp/internal/config"
	"myapp/internal/handler"
	"myapp/internal/kafka"
	"myapp/internal/repository"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 1. Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Инициализация репозитория (БД)
	repo, err := repository.New(cfg.DBURL)
	if err != nil {
		log.Fatalf("Failed to init repository: %v", err)
	}
	defer repo.Close()

	// 3. Инициализация кэша
	orderCache := cache.New()

	// 4. Восстановление кэша из БД при старте
	log.Println("Restoring cache from database...")
	orders, err := repo.GetAllOrders(context.Background())
	if err != nil {
		log.Fatalf("Failed to restore cache: %v", err)
	}
	orderCache.Load(orders)
	log.Printf("Cache restored. Loaded %d orders.", len(orders))

	// 5. Запуск Kafka Consumer в отдельной горутине
	ctx, cancel := context.WithCancel(context.Background())
	consumer := kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaTopic, repo, orderCache)
	go consumer.Run(ctx)

	// 6. Настройка и запуск HTTP-сервера
	orderHandler := handler.NewOrderHandler(orderCache)
	mux := http.NewServeMux()
	mux.Handle("/order/", http.HandlerFunc(orderHandler.GetOrderByUID))
	// Отдаем статическую веб-страницу
	mux.Handle("/", http.FileServer(http.Dir("./web")))

	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: mux,
	}

	go func() {
		log.Printf("HTTP server starting on port %s", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	// 7. Грациозное завершение работы
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	cancel() // Сигнал для остановки Kafka consumer
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server exited properly")
}

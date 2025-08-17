package main

import (
	"database/sql"
	"net/http"

	"github.com/tvbondar/go-server/internal/handlers"
	"github.com/tvbondar/go-server/internal/repositories"
	"github.com/tvbondar/go-server/internal/usecases"

	_ "github.com/lib/pq"
)

func main() {
	// PostgreSQL
	db, err := sql.Open("postgres", "user=myuser password=pass dbname=orders_db sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	dbRepo := repositories.NewPostgresOrderRepository(db)
	cacheRepo := repositories.NewCacheOrderRepository()

	// Восстановление кэша
	cacheRepo.LoadFromDB(dbRepo)

	processUseCase := usecases.NewProcessOrderUseCase(dbRepo, cacheRepo)
	getUseCase := usecases.NewGetOrderUseCase(cacheRepo, dbRepo) // Создай GetOrderUseCase: сначала кэш, затем DB

	// Kafka
	go handlers.StartKafkaConsumer(processUseCase)

	// HTTP
	httpHandler := handlers.NewHTTPHandler(getUseCase)
	http.HandleFunc("/order/", httpHandler.GetOrder)
	http.ListenAndServe(":8081", nil)
}

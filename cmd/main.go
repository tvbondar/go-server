package main

import (
	"database/sql"
	"internal/handlers"
	"internal/repositories"
	"internal/usecases"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=myuser password=pass dbname=orders_db sslmode=disable")
	if err != nil {
		panic(err)
	}

	repo := repositories.NewPostgresOrderRepository(db)
	usecase := usecases.NewProcessOrderUseCase(repo)

	go handlers.KafkaConsumer(usecase)

	// Запуск HTTP-сервера
	// ...
}

package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/tvbondar/go-server/internal/handlers"
	"github.com/tvbondar/go-server/internal/repositories"
	"github.com/tvbondar/go-server/internal/usecases"
)

func main() {
	connStr := "user=postgres password=pass dbname=orders_db host=postgres port=5432 sslmode=disable"
	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal("Failed to open database connection", err)
			time.Sleep(2 * time.Second)
			continue
		}
		err = db.Ping()
		if err == nil {
			break
		}
		log.Fatal("Failed to ping database", err)
		db.Close()
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Failed to connect to database after retries", err)
	}
	defer db.Close()
	log.Println("Successfully connected to database")

	dbRepo := repositories.NewPostgresOrderRepository(db)
	cacheRepo := repositories.NewCacheOrderRepository()

	if err := cacheRepo.LoadFromDB(dbRepo); err != nil {
		log.Fatal(err)
	}

	processUseCase := usecases.NewProcessOrderUseCase(dbRepo, cacheRepo)
	getUseCase := usecases.NewGetOrderUseCase(cacheRepo, dbRepo)

	go handlers.StartKafkaConsumer(processUseCase)

	httpHandler := handlers.NewHTTPHandler(getUseCase)
	http.HandleFunc("/order/", httpHandler.GetOrder)
	http.Handle("/", http.FileServer(http.Dir("web")))
	log.Fatal(http.ListenAndServe(":8081", nil))
}

package main

import (
	"context"
	"fmt"
	"log"
	"main/controller"
	"main/repository"
	"main/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
    // Connect to both DBs
    writeDB, err := pgxpool.New(context.Background(), "postgres://demo:demo@localhost:5438/cqrs_write")
    if err != nil { log.Fatal(err) }
    readDB, err := pgxpool.New(context.Background(), "postgres://demo:demo@localhost:5439/cqrs_read")
    if err != nil { log.Fatal(err) }

    // wire repos & services
    writeRepo := repository.NewProductWriteRepository(writeDB)
    readRepo := repository.NewProductReadRepository(readDB)
    srv := service.NewProductService(writeRepo, readRepo)
    ctrl := controller.NewProductController(srv)

    // start outbox processor
    processor := service.NewOutboxProcessor(writeDB, readDB)
    processor.Start()

    // gin setup
    r := gin.Default()
    ctrl.RegisterRoutes(r)

    fmt.Println("Server running on :8080")
    r.Run(":8080")
}

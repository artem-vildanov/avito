package app

import (
	"avito/internal/api"
	"avito/internal/db"
	"fmt"
	"log"
	"net/http"
	"os"
)

type App struct {
	router  *api.Router
	storage *db.PostgresStorage
}

func RunApp() {
	app := App{}

	app.storage = db.NewPostgresStorage()
	//app.storage.Migrate()

	app.router = api.NewRouter(app.storage)

	println("running server...")
	err := http.ListenAndServe(
		fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")),
		app.router.Router)
	if err != nil {
		log.Fatalf("Failed to start server: %s", err.Error())
	}
}

package main

import (
	"fmt"
	"net/http"

	"github.com/blue-script/url-shortener/configs"
	"github.com/blue-script/url-shortener/internal/auth"
	"github.com/blue-script/url-shortener/internal/link"
	"github.com/blue-script/url-shortener/pkg/db"
)

func main() {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	router := http.NewServeMux()

	// Repositories
	linkRepository := link.NewLinkRepository(db)

	// Handler
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: conf,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Server is lintening on port 8081")
	server.ListenAndServe()
}

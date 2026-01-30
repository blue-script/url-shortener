package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/blue-script/url-shortener/configs"
	"github.com/blue-script/url-shortener/internal/auth"
	"github.com/blue-script/url-shortener/internal/link"
	"github.com/blue-script/url-shortener/internal/stat"
	"github.com/blue-script/url-shortener/internal/user"
	"github.com/blue-script/url-shortener/pkg/db"
	"github.com/blue-script/url-shortener/pkg/middleware"
)

func tickOperation(ctx context.Context) {
	ticker := time.NewTicker(200 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			fmt.Println("Tick")
		case <-ctx.Done():
			fmt.Println("Cancel")
			return
		}
	}
}

func main() {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	router := http.NewServeMux()

	// Repositories
	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)
	statRepository := stat.NewStatRepository(db)

	// Services
	authService := auth.NewAuthService(userRepository)

	// Handler
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
		StatRepository: statRepository,
		Config:         conf,
	})

	// Middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}

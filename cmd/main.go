package main

import (
	"fmt"
	"net/http"

	"github.com/blue-script/url-shortener/configs"
	"github.com/blue-script/url-shortener/internal/auth"
	"github.com/blue-script/url-shortener/internal/link"
	"github.com/blue-script/url-shortener/internal/stat"
	"github.com/blue-script/url-shortener/internal/user"
	"github.com/blue-script/url-shortener/pkg/db"
	"github.com/blue-script/url-shortener/pkg/event"
	"github.com/blue-script/url-shortener/pkg/middleware"
)

func App() http.Handler {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	router := http.NewServeMux()
	eventBus := event.NewEventBus()

	// Repositories
	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)
	statRepository := stat.NewStatRepository(db)

	// Services
	authService := auth.NewAuthService(userRepository)
	statService := stat.NewStatService(&stat.StatServiceDeps{
		StatRepository: statRepository,
		EventBus:       eventBus,
	})

	// Handler
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
		EventBus:       eventBus,
		Config:         conf,
	})
	stat.NewStatHandler(router, stat.StatHandlerDeps{
		StatRepository: statRepository,
		Config:         conf,
	})

	go statService.AddClick()

	// Middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	return stack(router)
}

func main() {
	app := App()

	server := http.Server{
		Addr:    ":8081",
		Handler: app,
	}

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}

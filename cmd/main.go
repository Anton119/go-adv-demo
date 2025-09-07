package main

import (
	"fmt"
	"go-adv-demo/configs"
	"go-adv-demo/internal/auth"
	"go-adv-demo/internal/link"
	"go-adv-demo/internal/user"
	"go-adv-demo/pkg/db"
	"go-adv-demo/pkg/middleware"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	router := http.NewServeMux()
	// repositories

	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)

	// services
	authService := auth.NewAuthService(userRepository)

	//handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})

	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
	})

	//Midlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	fmt.Printf("server is lisening on port: %v \n", server.Addr)
	server.ListenAndServe()
}

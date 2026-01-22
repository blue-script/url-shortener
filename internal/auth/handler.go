package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/blue-script/url-shortener/configs"
	"github.com/blue-script/url-shortener/pkg/req"
	"github.com/blue-script/url-shortener/pkg/res"
)

type AuthHandlerDeps struct {
	*configs.Config
}

type AuthHandler struct {
	*configs.Config
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](&w, r)
		if err != nil {
			log.Printf("Error: %v", err.Error())
			return
		}
		fmt.Println(body)
		data := LoginResponse{
			Token: "123",
		}
		res.Json(w, data, 200)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			log.Printf("Error: %v", err.Error())
			return
		}
		fmt.Println(body)
		data := LoginResponse{
			Token: "123",
		}
		res.Json(w, data, 200)
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/blue-script/url-shortener/internal/auth"
	"github.com/blue-script/url-shortener/internal/user"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDb() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	fmt.Println(os.Getenv("DSN"))

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	return db
}

func initData(db *gorm.DB) {
	db.Create(&user.User{
		Email:    "a77@a.com",
		Password: "$2a$10$hdOE8CGuuvaVZSdZOsjfiOCczhzVzHIRVvyCJUhnNPRKZbi57r0Da",
		Name:     "Jon",
	})
}

func removeData(db *gorm.DB) {
	db.Unscoped().
		Where("email = ?", "a77@a.com").
		Delete(&user.User{})
}

func TestLoginSuccess(t *testing.T) {
	db := initDb()
	initData(db)
	defer removeData(db)

	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "a77@a.com",
		Password: "1",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	fmt.Println(res.StatusCode)

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("Expected %d got %d", 200, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err.Error())
	}

	var resData auth.LoginResponse

	err = json.Unmarshal(body, &resData)
	if err != nil {
		t.Fatal(err.Error())
	}

	if resData.Token == "" {
		t.Fatal("Token empty")
	}
}

func TestLoginFail(t *testing.T) {
	db := initDb()
	initData(db)
	defer removeData(db)

	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "a@a.com",
		Password: "no_password",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 401 {
		t.Fatalf("Expected no %d got %d", 401, res.StatusCode)
	}
}

package auth_test

import (
	"testing"

	"github.com/blue-script/url-shortener/internal/auth"
	"github.com/blue-script/url-shortener/internal/user"
)

type MockUserRepository struct{}

func (repo *MockUserRepository) Create(u *user.User) (*user.User, error) {
	return &user.User{
		Email: "a@a.en",
	}, nil
}

func (repo *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	return nil, nil
}

func TestRegisterSuccess(t *testing.T) {
	const initialEmail = "a@a.com"
	authService := auth.NewAuthService(&MockUserRepository{})
	email, err := authService.Register("a@a.com", "1", "a_user")
	if err != nil {
		t.Fatal(err)
	}
	if email != initialEmail {
		t.Fatalf("Email %s don't match %s", email, initialEmail)
	}
}



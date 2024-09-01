package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/gistapp/api/server"
	"github.com/gistapp/api/storage"
	"github.com/gistapp/api/tests/mock"
	"github.com/gistapp/api/user"
	"github.com/gistapp/api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func InitServerUsers() *fiber.App {
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		if err := storage.Migrate(); err != nil {
			return nil
		}
		return nil
	}

	// Set up the server
	port := utils.Get("PORT")
	s := server.NewServer(fmt.Sprintf(":%s", port))

	// Set up routers

	auth_router := user.AuthRouter{
		Controller: &mock.MockAuthController{
			AuthService: &mock.MockAuthService{},
		},
	}

	user_router := user.UserRouter{
		Controller: &user.UserControllerImpl{},
	}

	// Initialize the server with the routers
	s.Setup(&auth_router, &user_router)
	return s.App
}

func TestRetreiveUser(t *testing.T) {
	t.Run("Retreive user", func(t *testing.T) {
		app := InitServerUsers()
		if app == nil {
			t.Fatal("Failed to initialize the application")
		}

		// Begin the sign-up process
		//token corresponds to "test@test.com user"
		auth_token := GetAuthToken(t, app)

		// Retrieve the user
		body, _ := utils.MakeRequest("GET", t, app, "/user/me", nil, map[string]string{
			"Authorization": "Bearer " + auth_token,
		})

		if body["email"] != "test@test.com" {
			t.Fatalf("Expected email to be test@test.com")
		}

		shouldHave := map[string]bool{
			"email":   true,
			"name":    true,
			"picture": true,
			"id":      true,
		}

		for key := range body {
			if !shouldHave[key] {
				t.Fatalf("Unexpected key %s", key)
			}
		}

		log.Info(body)

		DeleteAuthUser(t, auth_token)
	})
}

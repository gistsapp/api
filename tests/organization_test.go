package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/gistapp/api/auth"
	"github.com/gistapp/api/gists"
	"github.com/gistapp/api/server"
	"github.com/gistapp/api/storage"
	"github.com/gistapp/api/tests/mock"
	"github.com/gistapp/api/utils"
	"github.com/gofiber/fiber/v2"
)

var endpoint = "http://localhost:4000"

func Init() *fiber.App {
	// Check for command-line arguments
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
	gistRouter := gists.GistRouter{
		Controller: gists.GistController,
	}

	authRouter := auth.AuthRouter{
		Controller: &mock.MockAuthController{
			AuthService: &mock.MockAuthService{},
		},
	}

	// Initialize the server with the routers
	s.Setup(&gistRouter, &authRouter)
	return s.App
}

func Stop() {
	os.Exit(0)
}

func getAuthToken(t *testing.T, app *fiber.App) string {
	// Begin the sign-up process
	beginPayload := map[string]string{
		"email": "test@test.com",
	}
	respBody, _ := utils.MakeRequest(t, app, "/auth/local/begin", beginPayload)
	token := respBody["token"]

	// Verify the sign-up process
	verifyPayload := map[string]string{
		"email": "test@test.com",
		"token": token,
	}
	_, resp := utils.MakeRequest(t, app, "/auth/local/verify", verifyPayload)

	auth_token := resp.Cookies()[0].Value
	return auth_token
}

func TestCreateOrganization(t *testing.T) {
	t.Run("Create organization", func(t *testing.T) {
		app := Init()
		if app == nil {
			t.Fatal("Failed to initialize the application")
		}

		// Begin the sign-up process

		//TODO: continue that

		// auth_token := getAuthToken(t, app)

		// Create a new organization

	})
}

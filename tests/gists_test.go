package tests

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/gistapp/api/auth"
	"github.com/gistapp/api/gists"
	"github.com/gistapp/api/organizations"
	"github.com/gistapp/api/server"
	"github.com/gistapp/api/storage"
	"github.com/gistapp/api/tests/mock"
	"github.com/gistapp/api/utils"
	"github.com/gofiber/fiber/v2"
)

func InitServerGists() *fiber.App {
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
	gist_router := gists.GistRouter{
		Controller: gists.GistController,
	}

	auth_router := auth.AuthRouter{
		Controller: &mock.MockAuthController{
			AuthService: &mock.MockAuthService{},
		},
	}

	organization_router := organizations.OrganizationRouter{
		Controller: organizations.OrganizationControllerImpl{},
	}

	// Initialize the server with the routers
	s.Setup(&gist_router, &auth_router, &organization_router)
	return s.App
}

func TestCreateGists(t *testing.T) {
	t.Run("Create a new personal gist", func(t *testing.T) {
		app := InitServerGists()
		authToken := GetAuthToken(t, app)

		body, req := utils.MakeRequest(t, app, "/gists", map[string]string{
			"name":    "Test Gist",
			"content": "Test content",
		}, map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", authToken),
		})

		if req.StatusCode != 201 {
			t.Fatalf("Expected status code 201, got %d", req.StatusCode)
		}

		fmt.Println(body)

	})

	t.Run("Create a new organization gist", func(t *testing.T) {
		app := InitServerGists()
		auth_token := GetAuthToken(t, app)
		claims, _ := auth.AuthService.IsAuthenticated(auth_token)

		org_mod := organizations.OrganizationSQL{
			Name: sql.NullString{
				String: "Test Org",
				Valid:  true,
			},
		}

		org, err := org_mod.Save(claims.Pub)
		if err != nil {
			t.Fatalf("Failed to create organization: %v", err)
		}

		payload := map[string]string{
			"name":    "Test Gist",
			"content": "Test content",
			"org_id":  org.ID,
		}

		body, req := utils.MakeRequest(t, app, "/gists", payload, map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", auth_token),
		})

		if req.StatusCode != 201 {
			t.Fatalf("Expected status code 201, got %d", req.StatusCode)
		}

		fmt.Println(body)

		DeleteAuthUser(t, auth_token)

	})
}

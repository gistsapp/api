package tests

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
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

var endpoint = "http://localhost:4000"

func InitServerOrgs() *fiber.App {
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

func TestCreateOrganization(t *testing.T) {
	t.Run("Create organization", func(t *testing.T) {
		app := InitServerOrgs()
		if app == nil {
			t.Fatal("Failed to initialize the application")
		}

		// Begin the sign-up process
		//
		auth_token := GetAuthToken(t, app)
		fmt.Println(auth_token)
		//
		// // Create a new organization
		org_payload := map[string]string{
			"name": "Test Organization",
		}
		fmt.Println(org_payload)
		//
		body, _ := utils.MakeRequest(t, app, "/orgs", org_payload, map[string]string{
			"Authorization": "Bearer " + auth_token,
		})
		//
		if body["name"] != "Test Organization" {
			t.Errorf("Expected organization name to be 'Test Organization', got %s", body["name"])
		}

		// cleanup
		id, err := strconv.ParseInt(body["id"], 10, 32)

		if err != nil {
			t.Errorf("Failed to parse organization ID: %v", err)
		}

		org := organizations.OrganizationSQL{
			ID: sql.NullInt32{
				Int32: int32(id),
				Valid: true,
			},
		}
		if err = org.Delete(); err != nil {
			t.Errorf("Failed to delete organization: %v", err)
		}
		DeleteAuthUser(t, auth_token)
	})
}

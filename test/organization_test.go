package test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/gistapp/api/gists"
	"github.com/gistapp/api/organizations"
	"github.com/gistapp/api/server"
	"github.com/gistapp/api/storage"
	"github.com/gistapp/api/test/mock"
	"github.com/gistapp/api/user"
	"github.com/gistapp/api/utils"
	"github.com/gofiber/fiber/v2"
)

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

	auth_router := user.AuthRouter{
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
		//
		body, _ := utils.MakeRequest("POST", t, app, "/orgs", org_payload, map[string]string{
			"Authorization": "Bearer " + auth_token,
		}, []int{201})
		//
		if body["name"] != "Test Organization" {
			t.Errorf("Expected organization name to be 'Test Organization', got %s", body["name"])
		}

		// cleanup
		DeleteOrganization(t, body["id"])
		DeleteAuthUser(t, auth_token)
	})

}

func DeleteOrganization(t *testing.T, org_id string) {

	org := organizations.OrganizationSQL{
		ID: sql.NullString{
			String: org_id,
			Valid:  true,
		},
	}
	if err := org.Delete(); err != nil {
		t.Errorf("Failed to delete organization: %v", err)
		return
	}

}

func TestDeleteOrganization(t *testing.T) {
	t.Run("Delete organization", func(t *testing.T) {
		app := InitServerOrgs()
		if app == nil {
			t.Fatal("Failed to initialize the application")
		}

		auth_token := GetAuthToken(t, app)

		org_payload := map[string]string{
			"name": "Test Organization",
		}

		body, _ := utils.MakeRequest("POST", t, app, "/orgs", org_payload, map[string]string{
			"Authorization": "Bearer " + auth_token,
		}, []int{201}) //before previous test tests the creation, we should be pretty sure that the creation works

		id := body["id"]

		body, _ = utils.MakeRequest("DELETE", t, app, fmt.Sprintf("/orgs/%s", id), nil, map[string]string{
			"Authorization": "Bearer " + auth_token,
		}, []int{200})

		org_dto := organizations.OrganizationSQL{
			ID: sql.NullString{
				String: id,
				Valid:  true,
			},
		}

		_, err := org_dto.Get()

		if err == nil {
			t.Fatal("Organization was not deleted")
		}

		DeleteAuthUser(t, auth_token)
	})
}

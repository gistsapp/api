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
	"github.com/gistapp/api/test/factory"
	"github.com/gistapp/api/test/mock"
	"github.com/gistapp/api/user"
	"github.com/gistapp/api/utils"
	"github.com/go-faker/faker/v4"
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

	auth_router := user.AuthRouter{
		Controller: &mock.MockAuthController{
			//needs GetUser to fit the auth interface
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

	t.Run("Create a personal gist", func(t *testing.T) {
		app := InitServerGists()
		user_factory := factory.UserWithAuthFactory()
		bob := user_factory.Create()
		Client(t, app).Post("/gists").WithJson().LoginAs(bob).WithPayload(map[string]string{
			"name":    faker.Name(),
			"content": faker.Sentence(),
		}).Send().ExpectStatus(201).JSON()

		user_factory.Clean()
	})

	t.Run("Create a new organization gist", func(t *testing.T) {
		app := InitServerGists()
		auth_token := GetAuthToken(t, app)
		claims, _ := user.AuthService.IsAuthenticated(auth_token)

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

		body, req := utils.MakeRequest("POST", t, app, "/gists", payload, map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", auth_token),
		}, []int{201})

		if req.StatusCode != 201 {
			t.Fatalf("Expected status code 201, got %d", req.StatusCode)
		}

		fmt.Println(body)
		DeleteOrganization(t, org.ID)
		DeleteAuthUser(t, auth_token)

	})

	t.Run("Create a new gist with invalid payload", func(t *testing.T) {
		app := InitServerGists()
		user_factory := factory.UserWithAuthFactory()
		user := user_factory.Create()
		access_token, err := user.GetAccessToken()

		if err != nil {
			t.Fatalf("Failed to get access token: %v", err)
		}

		Client(t, app).Post("/gists").WithPayload(map[string]string{
			"zob": "test",
		}).WithHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", access_token),
		}).Send().ExpectStatus(400)

		err = factory.UserWithAuthFactory().Clean()

		if err != nil {
			t.Fatalf("Failed to clean up user: %v", err)
		}
	})

	t.Run("Verify anyone can access public gist", func(t *testing.T) {
		app := InitServerGists()
		user_factory := factory.UserWithAuthFactory()
		bob := user_factory.Create()
		access_token, err := bob.GetAccessToken()

		if err != nil {
			t.Fatalf("Failed to get access token: %v", err)
		}

		client := Client(t, app).Post("/gists").WithPayload(map[string]string{
			"name":    faker.Name(),
			"content": faker.Sentence(),
		}).WithHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", access_token),
		}).Send().ExpectStatus(201)

		json_resp, err := JSONHttpResponse(client.Response)

		if err != nil {
			t.Fatalf("Failed to parse JSON response: %v", err)
		}

		gist_id := json_resp["id"]

		Client(t, app).Get(fmt.Sprintf("/gists/%s", gist_id)).Send().ExpectStatus(200)

		err = factory.UserWithAuthFactory().Clean()

		if err != nil {
			t.Fatalf("Failed to clean up user: %v", err)
		}
	})

	t.Run("Verify only the owner can access private gist", func(t *testing.T) {
		app := InitServerGists()
		user_factory := factory.UserWithAuthFactory()
		alice := user_factory.Create()
		alice_access_token, err := alice.GetAccessToken()
		if err != nil {
			t.Fatalf("Failed to get access token: %v", err)
		}
		client := Client(t, app).Post("/gists").WithPayload(map[string]string{
			"name":       faker.Name(),
			"content":    faker.Sentence(),
			"visibility": "private",
		}).WithHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", alice_access_token),
		}).Send().ExpectStatus(201)
		json_resp, err := JSONHttpResponse(client.Response)
		if err != nil {
			t.Fatalf("Failed to parse JSON response: %v", err)
		}
		gist_id := json_resp["id"]
		bob := user_factory.Create()
		bob_access_token, err := bob.GetAccessToken()
		if err != nil {
			t.Fatalf("Failed to get access token: %v", err)
		}
		Client(t, app).Get(fmt.Sprintf("/gists/%s", gist_id)).WithHeaders(map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", bob_access_token),
		}).Send().ExpectStatus(403)
		err = factory.UserWithAuthFactory().Clean()
		if err != nil {
			t.Fatalf("Failed to clean up user: %v", err)
		}
	})

	t.Run("Verify only the owner can access private raw gist", func(t *testing.T) {
		app := InitServerGists()
		user_factory := factory.UserWithAuthFactory()
		alice := user_factory.Create()
		alice_access_token, err := alice.GetAccessToken()
		if err != nil {
			t.Fatalf("Failed to get access token: %v", err)
		}
		client := Client(t, app).Post("/gists").WithPayload(map[string]string{
			"name":       faker.Name(),
			"content":    faker.Sentence(),
			"visibility": "private",
		}).WithHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": fmt.Sprintf("Bearer %s", alice_access_token),
		}).Send().ExpectStatus(201)
		json_resp, err := JSONHttpResponse(client.Response)
		if err != nil {
			t.Fatalf("Failed to parse JSON response: %v", err)
		}
		gist_id := json_resp["id"]
		bob := user_factory.Create()
		bob_access_token, err := bob.GetAccessToken()
		if err != nil {
			t.Fatalf("Failed to get access token: %v", err)
		}
		Client(t, app).Get(fmt.Sprintf("/gists/raw/%s", gist_id)).WithHeaders(map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", bob_access_token),
		}).Send().ExpectStatus(403)
		err = factory.UserWithAuthFactory().Clean()
		if err != nil {
			t.Fatalf("Failed to clean up user: %v", err)
		}
	})

	t.Run("Update a public gist", func(t *testing.T) {
		app := InitServerGists()
		gist_payload := map[string]string{
			"name":        faker.Name(),
			"content":     faker.Sentence(),
			"language":    "go",
			"description": faker.Sentence(),
			"visibility":  "public",
		}

		user_factory := factory.UserWithAuthFactory()
		alice := user_factory.Create()
		bob := user_factory.Create()
		response, err := Client(t, app).Post("/gists").WithJson().LoginAs(alice).WithPayload(gist_payload).Send().JSON()
		if err != nil {
			t.Fatalf("Failed to create gist: %v", err)
		}

		gist_id := response["id"]

		update_gist_payload := map[string]string{
			"name":        faker.Name(),
			"content":     faker.Sentence(),
			"language":    "js",
			"description": faker.Sentence(),
			"visibility":  "private",
		}
		Client(t, app).Put("/gists/" + gist_id).WithJson().LoginAs(bob).WithPayload(update_gist_payload).Send().ExpectStatus(403).JSON()

		response, err = Client(t, app).Put("/gists/" + gist_id).WithJson().LoginAs(alice).WithPayload(update_gist_payload).Send().ExpectStatus(200).JSON()

		fmt.Println(response)

		for key, value := range update_gist_payload {
			if response[key] != update_gist_payload[key] {
				t.Fatalf("Expected %s to be %s, got %s", key, value, response[key])
			}
		}
		//retry writing it
		Client(t, app).Put("/gists/" + gist_id).WithJson().LoginAs(alice).WithPayload(update_gist_payload).Send().ExpectStatus(200).JSON()
	})

	t.Run("Update a private gist to private", func(t *testing.T) {
		app := InitServerGists()
		gist_payload := map[string]string{
			"name":        faker.Name(),
			"content":     faker.Sentence(),
			"language":    "go",
			"description": faker.Sentence(),
			"visibility":  "public",
		}

		user_factory := factory.UserWithAuthFactory()
		alice := user_factory.Create()
		bob := user_factory.Create()

		response, err := Client(t, app).
			Post("/gists").
			WithJson().
			LoginAs(alice).
			WithPayload(gist_payload).
			Send().
			JSON()
		if err != nil {
			t.Fatalf("Failed to create gist: %v", err)
		}

		gist_id := response["id"]

		update_gist_payload := map[string]string{
			"name":        faker.Name(),
			"content":     faker.Sentence(),
			"language":    "js",
			"description": faker.Sentence(),
			"visibility":  "private",
		}

		Client(t, app).Put("/gists/" + gist_id).WithJson().LoginAs(alice).WithPayload(update_gist_payload).Send().ExpectStatus(200).JSON()
		Client(t, app).Put("/gists/" + gist_id).WithJson().LoginAs(bob).WithPayload(update_gist_payload).Send().ExpectStatus(403).JSON()

	})
}

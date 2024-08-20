package tests

import (
	"database/sql"
	"testing"

	"github.com/gistapp/api/user"
	"github.com/gistapp/api/utils"
	"github.com/gofiber/fiber/v2"
)

func GetAuthToken(t *testing.T, app *fiber.App) string {
	// Begin the sign-up process
	beginPayload := map[string]string{
		"email": "test@test.com",
	}
	respBody, _ := utils.MakeRequest("POST", t, app, "/auth/local/begin", beginPayload, nil)
	token := respBody["token"]

	// Verify the sign-up process
	verifyPayload := map[string]string{
		"email": "test@test.com",
		"token": token,
	}
	_, resp := utils.MakeRequest("POST", t, app, "/auth/local/verify", verifyPayload, nil)

	auth_token := resp.Cookies()[0].Value
	return auth_token
}

func DeleteAuthUser(t *testing.T, auth_token string) {
	claims, _ := user.AuthService.IsAuthenticated(auth_token)

	user := user.UserSQL{
		ID: sql.NullString{
			Valid:  true,
			String: claims.Pub,
		},
	}

	err := user.Delete()

	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

}

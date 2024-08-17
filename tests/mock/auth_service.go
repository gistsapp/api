package mock

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"github.com/gistapp/api/auth"
	"github.com/gistapp/api/user"
	"github.com/gistapp/api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
)

type MockAuthService struct{}

func (m *MockAuthService) Authenticate(c *fiber.Ctx) error {
	return nil
}

func (m *MockAuthService) LocalAuth(email string) (auth.TokenSQL, error) {
	token_val := utils.GenToken(6)
	token_model := auth.TokenSQL{
		Keyword: sql.NullString{String: email, Valid: true},
		Value:   sql.NullString{String: token_val, Valid: true},
		Type:    sql.NullString{String: string(auth.LocalAuth), Valid: true},
	}

	_, err := token_model.Save()

	return token_model, err

}

func (m *MockAuthService) VerifyLocalAuthToken(token string, email string) (string, error) {
	token_model := auth.TokenSQL{
		Value:   sql.NullString{String: token, Valid: true},
		Keyword: sql.NullString{String: email, Valid: true},
		Type:    sql.NullString{String: string(auth.LocalAuth), Valid: true},
	}
	token_data, err := token_model.Get()
	if err != nil {
		return "", err
	}
	err = token_data.Delete()
	if err != nil {
		return "", errors.New("couldn't invalidate token")
	}

	//now we finish users registration
	goth_user := goth.User{
		UserID:    email,
		Name:      strings.Split(email, "@")[0],
		Email:     email,
		AvatarURL: "https://vercel.com/api/www/avatar/?u=" + email + "&s=80",
	}

	if user, _, err := m.GetUser(goth_user); err == nil {
		jwt_token, err := utils.CreateToken(user.Email, user.ID)
		if err != nil {
			return "", err
		}
		return jwt_token, nil
	}

	user, err := m.Register(goth_user)

	if err != nil {
		return "", err
	}

	jwt_token, err := utils.CreateToken(user.Email, user.ID)

	return jwt_token, err
}

func (m *MockAuthService) Callback(c *fiber.Ctx) (string, error) {
	return "", nil
}

func (m *MockAuthService) GetUser(goth_user goth.User) (*user.User, *auth.AuthIdentity, error) {
	auth_and_user, err := new(auth.AuthIdentitySQL).GetWithUser(goth_user.UserID)
	if err != nil {
		return nil, nil, err
	}

	return &auth_and_user.User, &auth_and_user.AuthIdentity, nil
}

func (m *MockAuthService) Register(auth_user goth.User) (*user.User, error) {
	data, err := json.Marshal(auth_user)
	if err != nil {
		return nil, errors.New("couldn't marshal user")
	}

	user_model := user.UserSQL{
		ID:      sql.NullString{String: auth_user.UserID, Valid: true},
		Email:   sql.NullString{String: auth_user.Email, Valid: true},
		Name:    sql.NullString{String: auth_user.Name, Valid: true},
		Picture: sql.NullString{String: auth_user.AvatarURL, Valid: true},
	}

	user_data, err := user_model.Save()

	if err != nil {
		return nil, err
	}

	auth_identity_model := auth.AuthIdentitySQL{
		Data:       sql.NullString{String: string(data), Valid: true},
		Type:       sql.NullString{String: auth_user.Provider, Valid: true},
		OwnerID:    sql.NullString{String: user_data.ID, Valid: true},
		ProviderID: sql.NullString{String: auth_user.UserID, Valid: true},
	}

	_, err = auth_identity_model.Save()
	return user_data, err
}

func (a *MockAuthService) IsAuthenticated(token string) (*auth.JWTClaim, error) {
	claims, err := utils.VerifyJWT(token)

	if err != nil {
		return nil, err
	}

	jwtClaim := new(auth.JWTClaim)
	jwtClaim.Pub = claims["pub"].(string)
	jwtClaim.Email = claims["email"].(string)

	return jwtClaim, nil
}

func (a *MockAuthService) RegisterProviders() {
}
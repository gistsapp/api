package mock

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"github.com/gistapp/api/user"
	"github.com/gistapp/api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
)

type MockAuthService struct {
}

func (m *MockAuthService) Authenticate(c *fiber.Ctx) error {
	return nil
}

func (m *MockAuthService) LocalAuth(email string) (user.TokenSQL, error) {
	token_val := utils.GenToken(6)
	token_model := user.TokenSQL{
		Keyword: sql.NullString{String: email, Valid: true},
		Value:   sql.NullString{String: token_val, Valid: true},
		Type:    sql.NullString{String: string(user.LocalAuth), Valid: true},
	}

	_, err := token_model.Save()

	return token_model, err

}

func (m *MockAuthService) VerifyLocalAuthToken(token string, email string) (*user.Tokens, error) {
	token_model := user.TokenSQL{
		Value:   sql.NullString{String: token, Valid: true},
		Keyword: sql.NullString{String: email, Valid: true},
		Type:    sql.NullString{String: string(user.LocalAuth), Valid: true},
	}
	token_data, err := token_model.Get()
	if err != nil {
		return nil, err
	}
	err = token_data.Delete()
	if err != nil {
		return nil, errors.New("couldn't invalidate token")
	}

	//now we finish users registration
	goth_user := goth.User{
		UserID:    email,
		Name:      strings.Split(email, "@")[0],
		Email:     email,
		AvatarURL: "https://vercel.com/api/www/avatar/?u=" + email + "&s=80",
	}

	if auth_user, _, err := m.GetUser(goth_user); err == nil {
		access_token, err := utils.CreateAccessToken(auth_user.Email, auth_user.ID)
		refresh_token, err := utils.CreateRefreshToken(auth_user.ID)
		if err != nil {
			return nil, err
		}
		return &user.Tokens{
			AccessToken:  access_token,
			RefreshToken: refresh_token,
		}, nil
	}

	auth_user, err := m.Register(withEmailPrefix(goth_user))

	if err != nil {
		return nil, err
	}

	jwt_token, err := utils.CreateAccessToken(auth_user.Email, auth_user.ID)
	refresh_token, err := utils.CreateRefreshToken(auth_user.ID)

	return &user.Tokens{
		AccessToken:  jwt_token,
		RefreshToken: refresh_token,
	}, err
}

func (m *MockAuthService) Callback(c *fiber.Ctx) (*user.Tokens, error) {
	return nil, nil
}

func (a *MockAuthService) GetUser(auth_user goth.User) (*user.User, *user.AuthIdentity, error) {
	return user.AuthService.GetUser(auth_user)
}

func (a *MockAuthService) IsAuthenticated(token string) (*user.JWTClaim, error) {
	claims, err := utils.VerifyJWT(token)

	if err != nil {
		return nil, err
	}

	jwtClaim := new(user.JWTClaim)
	jwtClaim.Pub = claims["pub"].(string)
	jwtClaim.Email = claims["email"].(string)

	return jwtClaim, nil
}

func withEmailPrefix(auth_user goth.User) *user.RegistrationOptions {
	return &user.RegistrationOptions{
		SqlUser: &user.UserSQL{
			ID:      sql.NullString{String: auth_user.UserID, Valid: true},
			Email:   sql.NullString{String: auth_user.Email, Valid: true},
			Name:    sql.NullString{String: strings.Split(auth_user.Email, "@")[0], Valid: true},
			Picture: sql.NullString{String: auth_user.AvatarURL, Valid: true},
		},
		AuthUser: auth_user,
	}
}

func withGithubUsername(auth_user goth.User) *user.RegistrationOptions {
	return &user.RegistrationOptions{SqlUser: &user.UserSQL{
		ID:    sql.NullString{String: auth_user.UserID, Valid: true},
		Email: sql.NullString{String: auth_user.Email, Valid: true},
		Name:  sql.NullString{String: auth_user.RawData["login"].(string), Valid: true},
	}, AuthUser: auth_user}
}

func (a *MockAuthService) Register(options *user.RegistrationOptions) (*user.User, error) {
	auth_user := options.AuthUser
	data, err := json.Marshal(auth_user)
	if err != nil {
		return nil, errors.New("couldn't marshal user")
	}

	user_model := options.SqlUser
	user_data, err := user_model.Save()

	if err != nil {
		return nil, err
	}

	auth_identity_model := user.AuthIdentitySQL{
		Data:       sql.NullString{String: string(data), Valid: true},
		Type:       sql.NullString{String: auth_user.Provider, Valid: true},
		OwnerID:    sql.NullString{String: user_data.ID, Valid: true},
		ProviderID: sql.NullString{String: auth_user.UserID, Valid: true},
	}

	_, err = auth_identity_model.Save()
	return user_data, err
}

func (a *MockAuthService) RegisterProviders() {
}

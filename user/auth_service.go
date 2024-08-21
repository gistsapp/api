package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"github.com/gistapp/api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/shareed2k/goth_fiber"
)

type IAuthService interface {
	Authenticate(c *fiber.Ctx) error
	LocalAuth(email string) (TokenSQL, error)
	VerifyLocalAuthToken(token string, email string) (string, error)
	Callback(c *fiber.Ctx) (string, error)
	GetUser(auth_user goth.User) (*User, *AuthIdentity, error)
	Register(options *RegistrationOptions) (*User, error)
	RegisterProviders()
	IsAuthenticated(token string) (*JWTClaim, error)
}

type AuthServiceImpl struct{}

func (a *AuthServiceImpl) Authenticate(c *fiber.Ctx) error {
	if user, err := goth_fiber.CompleteUserAuth(c); err == nil {
		log.Info(user)
		return nil
	} else {
		return goth_fiber.BeginAuthHandler(c)
	}
}

// generates a token and sends it to the user by email
func (a *AuthServiceImpl) LocalAuth(email string) (TokenSQL, error) {
	token_val := utils.GenToken(6)
	token_model := TokenSQL{
		Keyword: sql.NullString{String: email, Valid: true},
		Value:   sql.NullString{String: token_val, Valid: true},
		Type:    sql.NullString{String: string(LocalAuth), Valid: true},
	}

	_, err := token_model.Save()

	if err != nil {
		return token_model, err
	}

	err = utils.SendEmail("Gistapp: Local Auth", "Your token is: "+token_val, email)

	return token_model, err
}

// verifies the token and finishes the registration
func (a *AuthServiceImpl) VerifyLocalAuthToken(token string, email string) (string, error) {
	token_model := TokenSQL{
		Value:   sql.NullString{String: token, Valid: true},
		Keyword: sql.NullString{String: email, Valid: true},
		Type:    sql.NullString{String: string(LocalAuth), Valid: true},
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

	if user, _, err := a.GetUser(goth_user); err == nil {
		jwt_token, err := utils.CreateToken(user.Email, user.ID)
		if err != nil {
			return "", err
		}
		return jwt_token, nil
	}

	user, err := a.Register(withEmailPrefix(goth_user))

	if err != nil {
		return "", err
	}

	jwt_token, err := utils.CreateToken(user.Email, user.ID)

	return jwt_token, err
}

func (a *AuthServiceImpl) Callback(c *fiber.Ctx) (string, error) {
	provider := c.Params("provider")
	auth_user, err := goth_fiber.CompleteUserAuth(c)
	if err != nil {
		log.Error(err)
		return "", ErrCantCompleteAuth
	}

	user_md, _, err := a.GetUser(auth_user)

	if err == nil {
		token, err := utils.CreateToken(user_md.Email, user_md.ID)
		if err != nil {
			return "", err
		}
		return token, nil
	}

	log.Info(auth_user.NickName)
	if provider == "github" {
		user_md, err = a.Register(withGithubUsername(auth_user))
	} else {
		user_md, err = a.Register(withEmailPrefix(auth_user))
	}

	if err != nil {
		return "", err
	}

	jwt, err := utils.CreateToken(user_md.Email, user_md.ID)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func (a *AuthServiceImpl) GetUser(auth_user goth.User) (*User, *AuthIdentity, error) {
	auth_and_user, err := new(AuthIdentitySQL).GetWithUser(auth_user.UserID)
	if err != nil {
		return nil, nil, err
	}

	return &auth_and_user.User, &auth_and_user.AuthIdentity, nil
}

type RegistrationOptions struct {
	AuthUser goth.User
	SqlUser  *UserSQL
}

func withEmailPrefix(user goth.User) *RegistrationOptions {
	return &RegistrationOptions{
		SqlUser: &UserSQL{
			ID:      sql.NullString{String: user.UserID, Valid: true},
			Email:   sql.NullString{String: user.Email, Valid: true},
			Name:    sql.NullString{String: strings.Split(user.Email, "@")[0], Valid: true},
			Picture: sql.NullString{String: user.AvatarURL, Valid: true},
		},
		AuthUser: user,
	}
}

func withGithubUsername(user goth.User) *RegistrationOptions {
	return &RegistrationOptions{SqlUser: &UserSQL{
		ID:      sql.NullString{String: user.UserID, Valid: true},
		Email:   sql.NullString{String: user.Email, Valid: true},
		Name:    sql.NullString{String: user.NickName, Valid: true},
		Picture: sql.NullString{String: user.AvatarURL, Valid: true},
	}, AuthUser: user}
}

func (a *AuthServiceImpl) Register(options *RegistrationOptions) (*User, error) {
	log.Info(options)
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

	auth_identity_model := AuthIdentitySQL{
		Data:       sql.NullString{String: string(data), Valid: true},
		Type:       sql.NullString{String: auth_user.Provider, Valid: true},
		OwnerID:    sql.NullString{String: user_data.ID, Valid: true},
		ProviderID: sql.NullString{String: auth_user.UserID, Valid: true},
	}

	auth_identity, err := auth_identity_model.Save()
	log.Info(auth_identity)
	return user_data, err
}

func (a *AuthServiceImpl) RegisterProviders() {
	goth.UseProviders(
		google.New(utils.Get("GOOGLE_KEY"), utils.Get("GOOGLE_SECRET"), utils.Get("PUBLIC_URL")+"/auth/callback/google"),
		github.New(utils.Get("GITHUB_KEY"), utils.Get("GITHUB_SECRET"), utils.Get("PUBLIC_URL")+"/auth/callback/github"),
	)
}

func (a *AuthServiceImpl) IsAuthenticated(token string) (*JWTClaim, error) {
	claims, err := utils.VerifyJWT(token)

	if err != nil {
		return nil, err
	}

	jwtClaim := new(JWTClaim)
	jwtClaim.Pub = claims["pub"].(string)
	jwtClaim.Email = claims["email"].(string)

	return jwtClaim, nil
}

var AuthService AuthServiceImpl = AuthServiceImpl{}
var ErrCantCompleteAuth = errors.New("can't complete auth")

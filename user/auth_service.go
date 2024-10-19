package user

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gistapp/api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/shareed2k/goth_fiber"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type IAuthService interface {
	Authenticate(c *fiber.Ctx) error
	LocalAuth(email string) (TokenSQL, error)
	VerifyLocalAuthToken(token string, email string) (*Tokens, error)
	Callback(c *fiber.Ctx) (*Tokens, error)
	GetUser(auth_user goth.User) (*User, *AuthIdentity, error)
	Register(options *RegistrationOptions) (*User, error)
	RegisterProviders()
	Renew(user_id string) (*Tokens, error)
	IsAuthenticated(token string) (*JWTClaim, error)
	CanRefresh(token string) (*JWTClaim, error)
}

type AuthServiceImpl struct {
	user_service UserServiceImpl
}

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
		log.Info(err.Error())
		if strings.Contains(err.Error(), "token_keyword_key") {
			token, err := token_model.GetByType(LocalAuth)
			log.Info(token)

			if err != nil {
				return token_model, err
			}

			err = token.Delete()

			if err != nil {
				return token_model, err
			}

			_, err = token_model.Save()

			if err != nil {
				return token_model, err
			}
		}
	}

	err = utils.SendEmail("Gistapp: Local Auth", "Your token is: "+token_val, email)

	return token_model, err
}

// verifies the token and finishes the registration
func (a *AuthServiceImpl) VerifyLocalAuthToken(token string, email string) (*Tokens, error) {
	token_model := TokenSQL{
		Value:   sql.NullString{String: token, Valid: true},
		Keyword: sql.NullString{String: email, Valid: true},
		Type:    sql.NullString{String: string(LocalAuth), Valid: true},
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

	if user, _, err := a.GetUser(goth_user); err == nil {
		access_token, err := utils.CreateAccessToken(user.Email, user.ID)
		refresh_token, err := utils.CreateRefreshToken(user.ID)

		if err != nil {
			return nil, err
		}
		return &Tokens{
			AccessToken:  access_token,
			RefreshToken: refresh_token,
		}, nil
	}

	user, err := a.Register(withEmailPrefix(goth_user))

	if err != nil {
		return nil, err
	}

	access_token, err := utils.CreateAccessToken(user.Email, user.ID)
	refresh_token, err := utils.CreateRefreshToken(user.ID)

	return &Tokens{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}, err
}

func (a *AuthServiceImpl) Callback(c *fiber.Ctx) (*Tokens, error) {
	provider := c.Params("provider")
	auth_user, err := goth_fiber.CompleteUserAuth(c)
	if err != nil {
		log.Error(err)
		return nil, ErrCantCompleteAuth
	}

	user_md, _, err := a.GetUser(auth_user)

	if err == nil {
		access_token, err := utils.CreateAccessToken(user_md.Email, user_md.ID)
		refresh_token, err := utils.CreateRefreshToken(user_md.ID)
		if err != nil {
			return nil, err
		}
		return &Tokens{
			AccessToken:  access_token,
			RefreshToken: refresh_token,
		}, nil
	}

	log.Info(auth_user.NickName)
	if provider == "github" {
		user_md, err = a.Register(withGithubUsername(auth_user))
	} else {
		user_md, err = a.Register(withEmailPrefix(auth_user))
	}

	if err != nil {
		return nil, err
	}

	access_token, err := utils.CreateAccessToken(user_md.Email, user_md.ID)
	if err != nil {
		return nil, err
	}
	refresh_token, err := utils.CreateRefreshToken(user_md.ID)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}, nil
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

func (a *AuthServiceImpl) CanRefresh(token string) (*JWTClaim, error) {
	claims, err := utils.VerifyJWT(token)

	if err != nil {
		return nil, err
	}

	jwtClaim := new(JWTClaim)
	jwtClaim.Pub = claims["pub"].(string)
	jwtClaim.Email = "" //we don't care about email in a refresh token since it's not set

	return jwtClaim, nil
}

// renew the token, its asserted that the access token is correct since it has been verified previously inside the auth middleware
func (a *AuthServiceImpl) Renew(user_id string) (*Tokens, error) {

	user, err := a.user_service.GetUserByID(user_id)

	if err != nil {
		return nil, err
	}

	access_token, err := utils.CreateAccessToken(user.Email, user_id)
	refresh_token, err := utils.CreateRefreshToken(user_id)

	if err != nil {
		return nil, err
	}

	fmt.Println(access_token)
	return &Tokens{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}, nil
}

var AuthService AuthServiceImpl = AuthServiceImpl{
	user_service: UserService,
}
var ErrCantCompleteAuth = errors.New("can't complete auth")

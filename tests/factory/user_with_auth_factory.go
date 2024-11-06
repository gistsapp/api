package factory

import (
	"database/sql"

	"github.com/gistapp/api/user"
	"github.com/go-faker/faker/v4"
)

type TUserWithAuthFactory struct {
	UserPool []user.AuthIdentityAndUser
}

func UserWithAuthFactory() *TUserWithAuthFactory {
	return &TUserWithAuthFactory{}
}

func (f *TUserWithAuthFactory) Create() user.AuthIdentityAndUser {
	u := user.UserSQL{
		ID: sql.NullString{
			String: "",
			Valid:  false,
		},
		Email: sql.NullString{
			String: faker.Email(),
			Valid:  true,
		},
		Name: sql.NullString{
			String: faker.Name(),
			Valid:  true,
		},
		Picture: sql.NullString{
			String: faker.URL(),
			Valid:  true,
		},
	}

	user_res, _ := u.Save()

	ai := user.AuthIdentitySQL{
		ID: sql.NullString{
			String: "",
			Valid:  false,
		},
		Data: sql.NullString{
			String: "{}",
			Valid:  true,
		},
		Type: sql.NullString{
			String: "local",
			Valid:  true,
		},
		ProviderID: sql.NullString{
			String: user_res.Email,
			Valid:  true,
		},
		OwnerID: sql.NullString{
			String: user_res.ID,
			Valid:  true,
		},
	}

	ai_res, err := ai.Save()
	if err != nil {
		panic(err)
	}
	ai_user := user.AuthIdentityAndUser{
		AuthIdentity: *ai_res,
		User:         *user_res,
	}

	f.UserPool = append(f.UserPool, ai_user)

	return ai_user
}

func (f *TUserWithAuthFactory) CreateMany(count int) []user.AuthIdentityAndUser {
	var users []user.AuthIdentityAndUser
	for i := 0; i < count; i++ {
		users = append(users, f.Create())
	}
	return users
}

func (f *TUserWithAuthFactory) Clean() error {

	for _, u := range f.UserPool {
		u_sql := user.UserSQL{
			ID: sql.NullString{
				String: u.User.ID,
				Valid:  true,
			},
		}
		err := u_sql.Delete()
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *TUserWithAuthFactory) Get() user.AuthIdentityAndUser {
	if len(f.UserPool) == 0 {
		return f.Create()
	}
	return f.UserPool[0]
}

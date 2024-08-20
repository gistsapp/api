package user

import "database/sql"

type UserServiceImpl struct{}

func (u *UserServiceImpl) GetUserByID(id string) (*User, error) {
	user := UserSQL{
		ID: sql.NullString{
			String: id,
			Valid:  true,
		},
	}
	return user.GetByID()
}

var UserService = UserServiceImpl{}

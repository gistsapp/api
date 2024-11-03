package organizations

import (
	"database/sql"

	"github.com/gistapp/api/storage"
	"github.com/gofiber/fiber/v2/log"
)

type MemberSQL struct {
	MemberID sql.NullString
	OrgID    sql.NullString
	UserID   sql.NullString
	Role     sql.NullString
}

type Member struct {
	MemberID string `json:"member_id"`
	OrgID    string `json:"org_id"`
	UserID   string `json:"user_id"`
	Role     Role   `json:"role"`
}

func (m *MemberSQL) Get() (*Member, error) {
	row, err := storage.Database.Query("SELECT member_id, org_id, user_id, role FROM member WHERE org_id = $1 AND user_id = $2", m.OrgID.String, m.UserID.String)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	row.Next()

	var member Member

	err = row.Scan(&member.MemberID, &member.OrgID, &member.UserID, &member.Role)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &member, nil
}

package utils

import "database/sql"

// go support for nil/zero values unmarshalling in golang

type ZeroString struct {
	content string
	valid   bool
}

func (z *ZeroString) UnmarshalJSON(data []byte) error {
	if string(data) == "" {
		z.content = ""
		z.valid = false
		return nil
	}
	z.content = string(data)
	return nil
}

func (z *ZeroString) MarshalJSON() ([]byte, error) {
	if z.valid {
		return []byte(z.content), nil
	}
	return []byte("null"), nil
}

func (z *ZeroString) String() string {
	return z.content
}

func (z *ZeroString) SqlString() sql.NullString {
	return sql.NullString{
		String: z.content,
		Valid:  z.valid,
	}
}

func FromSQL(sqlString sql.NullString) ZeroString {
	return ZeroString{
		content: sqlString.String,
		valid:   sqlString.Valid,
	}
}

func ToNullString(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}

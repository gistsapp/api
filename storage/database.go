package storage

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/gistapp/api/utils"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type DatabaseV1 struct {
	user     string
	password string
	host     string
	port     string
	database string
}

type IDatabase interface {
	Connect() (*sql.DB, error)
	Query(query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
}

func newDatabase() *DatabaseV1 {
	return &DatabaseV1{
		utils.Get("PG_USER"),
		utils.Get("PG_PASSWORD"),
		utils.Get("PG_HOST"),
		utils.Get("PG_PORT"),
		utils.Get("PG_DATABASE"),
	}
}

func (db *DatabaseV1) Connect() (*sql.DB, error) {
	connStr := "user=" + db.user + " password=" + db.password + " host=" + db.host + " port=" + db.port + " dbname=" + db.database + " sslmode=disable"
	return sql.Open("postgres", connStr)
}

func (db *DatabaseV1) Query(query string, args ...any) (*sql.Rows, error) {
	conn, err := db.Connect()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer conn.Close()

	if len(args) == 0 {
		return conn.Query(query)
	}
	return conn.Query(query, args...)
}

var Database IDatabase = newDatabase()

func (db *DatabaseV1) Exec(query string, args ...any) (sql.Result, error) {
	conn, err := db.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return conn.Exec(query, args...)
}

func Migrate() error {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	m, err := migrate.New(
		"file://"+exPath+"/migrations",
		"postgres://"+utils.Get("PG_USER")+":"+utils.Get("PG_PASSWORD")+"@"+utils.Get("PG_HOST")+":"+utils.Get("PG_PORT")+"/"+utils.Get("PG_DATABASE")+"?sslmode=disable",
	)
	if err != nil {
		return err
	}
	m.Up()
	return nil
}

package sqlite

import (
	"database/sql"

	"github.com/MadhavKrishanGoswami/students-api/internal/config"
	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
)

type Sqlite struct {
	DB *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		age INTEGER NOT NULL CHECK(age >= 0 AND age <= 120)
	);`)

	return &Sqlite{DB: db}, nil
}

func (s *Sqlite) CreateStudent(name, email string, age int) (int64, error) {
	stmt, err := s.DB.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

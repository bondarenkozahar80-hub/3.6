package postgres

import (
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	DB *sqlx.DB
}

// New  - конструктор БД
func New(databaseUri string) (*Postgres, error) {
	db, err := sqlx.Connect("pgx", databaseUri)
	if err != nil {
		return nil, fmt.Errorf("[postgres] failed to connect to DB: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("[postgres] ping failed: %w", err)
	}
	log.Println("[postgres] connect to DB successfully")
	return &Postgres{
		DB: db,
	}, nil
}

// Close закрывает соединение с БД
func (p *Postgres) Close() error {
	if p.DB != nil {
		log.Println("[postgres] closing connection to DB")
		return p.DB.Close()
	}
	return nil
}

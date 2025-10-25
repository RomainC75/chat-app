package db

import (
	"chat/config"
	"chat/utils"
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Store interface {
	Querier
	ExecTx(ctx context.Context, fn func(*Queries) error) error
}

var DbStore *Store

type SqlStore struct {
	*Queries
	db *sql.DB
}

func (store *SqlStore) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

func NewStore(db *sql.DB) Store {
	return &SqlStore{
		db:      db,
		Queries: New(db),
	}
}

func Connect() {
	cfg := config.Get()
	utils.PrettyDisplay("cfg ", cfg)
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)

	conn, err := sql.Open("postgres", dsn)
	store := NewStore(conn)
	if err != nil {
		log.Fatal("error trying to connect to the database : ", err)
	}
	if err := conn.Ping(); err != nil {
		log.Fatal("error trying to ping the db ", err)
	}
	fmt.Println("DB : connected")

	DbStore = &store
}

func GetConnection() *Store {
	return DbStore
}

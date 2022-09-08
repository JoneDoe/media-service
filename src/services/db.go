package services

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"

	"istorage/logger"
	"istorage/models"
)

type DbEngine struct {
	pool *pgxpool.Pool
}

func InitDb() *DbEngine {
	dbpool, err := pgxpool.Connect(context.Background(), GetDataSourceName("connect_timeout=5"))

	if err != nil {
		logger.Fatal(err)
	}

	return &DbEngine{
		pool: dbpool,
	}
}

func GetDataSourceName(params interface{}) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
		params,
	)
}

func (e *DbEngine) Close() {
	e.pool.Close()
}

func (e *DbEngine) CreateRecord(mediaFile *models.MediaFile) error {
	_, err := e.pool.Exec(context.Background(), "insert into media values($1, $2)", mediaFile.Uuid, mediaFile)

	return err
}

func (e *DbEngine) DeleteRecord(uuid string) error {
	_, err := e.pool.Exec(context.Background(), "DELETE FROM media WHERE uuid = $1", uuid)

	return err
}

func (e *DbEngine) GetRecord(uuid string) (map[string]interface{}, error) {
	data := make(map[string]interface{})

	err := e.pool.QueryRow(context.Background(), "SELECT data FROM media WHERE uuid = $1", uuid).Scan(&data)

	return data, err
}

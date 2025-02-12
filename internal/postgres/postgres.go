package postgres

import (
	"fmt"

	"github.com/MikeMalyshev/SkillRocks/internal/service"
	"github.com/jackc/pgx"
)

type db struct {
	pool *pgx.ConnPool
}

func New() *db {
	return &db{}
}

func (d *db) Connect() error {
	// TODO - через пул соединений
	return nil
}

func (d *db) Create() error {
	connection, err := pgx.Connect(pgx.ConnConfig{
		Host:     "localhost",
		User:     "postgres",
		Password: "postgres",
		Database: "skillrocks",
	})
	if err != nil {
		return fmt.Errorf("Unable to connect to database: %v\n", err)
	}
	defer connection.Close()

	_, err = connection.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title TEXT,
		description TEXT,
		status TEXT CHECK(status IN ('new', 'in_progress', 'done')) DEFAULT 'new',
		created_at TIMESTAMP DEFAULT now(),
		updated_at TIMESTAMP DEFAULT now())`)

	if err != nil {
		return fmt.Errorf("Unable to create table: %v\n", err)
	}
	return nil
}

func (d *db) AddTask(service.Task) error {

	return nil
}

func (d *db) GetTasks() ([]service.Task, error) {
	return []service.Task{}, nil
}

func (d *db) UpdateTask(service.Task) error {
	return nil
}

func (d *db) DeleteTask(id int) error {
	return nil
}

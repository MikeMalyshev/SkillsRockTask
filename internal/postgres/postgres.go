package postgres

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/MikeMalyshev/SkillRocks/internal/service"
	"github.com/jackc/pgx"
)

const (
	ErrorNoConnection = "No connection to database, %v"
)

// db is a struct that contains connection to database
type db struct {
	connection *pgx.Conn
}

// New creates new db struct. Connection to database is not established. Use Connect() method to establish connection and Close() to close it
func New() *db {
	d := &db{}
	if !d.TableExists() {
		err := d.Create()
		if err != nil {
			fmt.Println("error")
			log.Fatal(err)
			return nil
		}
		fmt.Println("db created")
	}
	return d
}

// Connect to database. If connection is already established, do nothing
// Connection must be closed manually with Close() method
func (d *db) Connect() error {
	if d.CheckConnection() {
		return nil
	}

	var err error
	d.connection, err = pgx.Connect(pgx.ConnConfig{
		Host:     "localhost",
		Port:     5433,
		User:     "postgres",
		Password: "postgres",
		Database: "postgres",
	})
	if err != nil {
		return fmt.Errorf("Unable to connect to database: %v", err)
	}
	return nil
}

// Close connection to database
func (d *db) Close() error {
	err := d.connection.Close()
	d.connection = nil
	return err
}

// CheckConnection returns true if connection to database is established, false otherwise
func (d *db) CheckConnection() bool {
	if d.connection == nil {
		return false
	}

	err := d.connection.Ping(context.Background())
	if err != nil {
		return false
	}
	return true
}

// Create creates table tasks in database if it does not exist
func (d *db) Create() error {
	err := d.Connect()
	defer d.Close()

	_, err = d.connection.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title TEXT,
		description TEXT,
		status TEXT CHECK(status IN ('new', 'in_progress', 'done')) DEFAULT 'new',
		created_at TIMESTAMP DEFAULT now(),
		updated_at TIMESTAMP DEFAULT now())`)

	if err != nil {
		return fmt.Errorf("Unable to create table: %v", err)
	}

	return nil
}

// CheckTable checks if table tasks exists in database
func (d *db) TableExists() bool {
	if err := d.Connect(); err != nil {
		return false
	}

	req := `SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name = 'tasks')`
	var exists bool
	err := d.connection.QueryRow(req).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (d *db) PrepareRequestData(task service.Task) (columns, placeholders []string, data []interface{}) {
	if task.Title != nil {
		columns = append(columns, "title")
		data = append(data, *task.Title)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(columns)))
	}
	if task.Description != nil {
		columns = append(columns, "description")
		data = append(data, *task.Description)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(columns)))
	}
	if task.Status != nil {
		columns = append(columns, "status")
		data = append(data, *task.Status)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(columns)))
	}
	if task.CreatedAt != nil {
		columns = append(columns, "created_at")
		data = append(data, *task.CreatedAt)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(columns)))
	}
	if task.UpdatedAt != nil {
		columns = append(columns, "updated_at")
		data = append(data, *task.UpdatedAt)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(columns)))
	}
	return
}

// AddTask adds task to database
func (d *db) AddTask(task service.Task) error {
	if err := d.Connect(); err != nil {
		return fmt.Errorf(ErrorNoConnection, err)
	}

	columns, placeholders, data := d.PrepareRequestData(task)
	req := fmt.Sprintf(`INSERT INTO tasks (%s) VALUES (%s) RETURNING id`, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	var insertedID int
	rows := d.connection.QueryRow(req, data...)
	err := rows.Scan(&insertedID)
	if err != nil {
		return fmt.Errorf("Unable to insert task: %v", err)
	}

	return nil
}

// GetTasks returns all tasks from database
func (d *db) GetTasks() ([]service.Task, error) {
	if err := d.Connect(); err != nil {
		return []service.Task{}, fmt.Errorf(ErrorNoConnection, err)
	}

	req := `SELECT id, title, description, status, created_at, updated_at FROM tasks`
	rows, err := d.connection.Query(req)
	if err != nil {
		return []service.Task{}, fmt.Errorf("Unable to get tasks: %v", err)
	}
	defer rows.Close()

	var tasks []service.Task
	for rows.Next() {
		task := service.Task{}
		err = rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return []service.Task{}, fmt.Errorf("Unable to scan task: %v", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// UpdateTask updates task in database
func (d *db) UpdateTask(task service.Task) error {
	if err := d.Connect(); err != nil {
		return fmt.Errorf(ErrorNoConnection, err)
	}
	columns, placeholders, data := d.PrepareRequestData(task)

	req := `UPDATE tasks SET `

	for i, column := range columns {
		if i > 0 {
			req += ", "
		}
		req += fmt.Sprintf("%s = %s", column, placeholders[i])
	}

	req += fmt.Sprintf(" WHERE id = %d", task.ID)

	ct, err := d.connection.Exec(req, data...)
	if err != nil {
		return fmt.Errorf("Unable to update task: %v", err)
	}
	if ct.RowsAffected() != 1 {
		return fmt.Errorf("Unable to update task with id %d", task.ID)
	}
	return nil
}

// DeleteTask deletes task from database
func (d *db) DeleteTask(id int) error {
	if err := d.Connect(); err != nil {
		return fmt.Errorf(ErrorNoConnection, err)
	}

	req := `DELETE FROM tasks WHERE id = $1`
	ct, err := d.connection.Exec(req, id)
	if err != nil {
		return fmt.Errorf("Unable to delete task: %v", err)
	}
	if ct.RowsAffected() != 1 {
		return fmt.Errorf("Unable to delete task with id %d", id)
	}

	return nil
}

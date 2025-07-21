package dbwork

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"subscriptions/pkg/models"

	"github.com/golang-migrate/migrate/v4"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var DB DataBase

type DataBase interface {
	AddSubscriptions(sub models.Subscriptions) error
	DeleteSubscriptions(id int) error
	GetSubscriptions(id int) (models.Subscriptions, error)
	UpdateSubscriptions(id int, sub models.Subscriptions) error
	GetAllSubscriptions(sub models.Subscriptions) ([]models.Subscriptions, error)
}

type PostgresDBParams struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
	SSLMode  string
}

type PostgresDataBase struct {
	db *sql.DB
}

func (postgre *PostgresDataBase) AddSubscriptions(sub models.Subscriptions) error {
	createQuery := `INSERT INTO subscriptions
	                (user_id, service_name, price, start_date, end_date)
	                VALUES($1, $2, $3, $4, $5);`

	_, err := postgre.db.Query(
		createQuery,
		sub.UserID,
		sub.ServiceName,
		sub.Price,
		sub.StartDate,
		sub.EndDate,
	)
	if err != nil {
		return err
	}

	return nil
}

func (postgre *PostgresDataBase) DeleteSubscriptions(id int) error {
	deleteQuery := `DELETE FROM subscriptions WHERE id = $1`
	_, err := postgre.db.Exec(deleteQuery, id)
	if err != nil {
		return err
	}

	return nil
}

func (postgre *PostgresDataBase) GetSubscriptions(id int) (models.Subscriptions, error) {
	selectQuery := `SELECT id, user_id, price, service_name, start_date, end_date FROM subscriptions WHERE id=$1`
	result := models.Subscriptions{}

	rows, err := postgre.db.Query(selectQuery, id)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&result.ID,
			&result.UserID,
			&result.Price,
			&result.ServiceName,
			&result.StartDate,
			&result.EndDate,
		)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

func (postgre *PostgresDataBase) UpdateSubscriptions(id int, sub models.Subscriptions) error {
	updateQuery := `UPDATE subscriptions
                  SET `
	query, args, index := postgre.selectParameters(sub, ", ")
	updateQuery += query
	if index == 1 {
		return fmt.Errorf(
			"Для выполнения UPDATE запроса необходимо указать хотя бы 1 изменяемый параметр",
		)
	}

	updateQuery += fmt.Sprintf(` WHERE id=$%d`, index)
	args = append(args, id)
	_, err := postgre.db.Exec(updateQuery, args...)
	if err != nil {
		log.Println(updateQuery, args)
		return err
	}
	return nil
}

func (postgre *PostgresDataBase) GetAllSubscriptions(
	sub models.Subscriptions,
) ([]models.Subscriptions, error) {
	var rows *sql.Rows
	var err error
	result := []models.Subscriptions{}
	if sub.All == false {
		selectQuery := `SELECT id, user_id, price, service_name, start_date, end_date FROM subscriptions WHERE `
		query, args, _ := postgre.selectParameters(sub, " AND ")

		selectQuery += query

		rows, err = postgre.db.Query(selectQuery, args...)
		if err != nil {
			return result, err
		}
	} else {
		selectQuery := `SELECT id, user_id, price, service_name, start_date, end_date FROM subscriptions`
		rows, err = postgre.db.Query(selectQuery)
		if err != nil {
			return result, err
		}
	}
	defer rows.Close()
	for rows.Next() {
		temp := models.Subscriptions{}
		err = rows.Scan(
			&temp.ID,
			&temp.UserID,
			&temp.Price,
			&temp.ServiceName,
			&temp.StartDate,
			&temp.EndDate,
		)
		if err != nil {
			return result, err
		}
		result = append(result, temp)
	}

	return result, nil
}

func (postgre *PostgresDataBase) selectParameters(
	sub models.Subscriptions,
	sep string,
) (string, []any, int) {
	query := ""
	args := []any{}
	index := 1
	if !sub.StartDate.IsZero() {
		query += fmt.Sprintf(`start_date=$%d`, index)
		args = append(args, sub.StartDate)
		index++
	}

	if sub.ServiceName != "" {
		if index != 1 {
			query += sep
		}
		query += fmt.Sprintf("service_name=$%d", index)
		args = append(args, sub.ServiceName)
		index++
	}

	if sub.Price != 0 {
		if index != 1 {
			query += sep
		}
		query += fmt.Sprintf("price=$%d", index)
		args = append(args, sub.Price)
		index++
	}

	if !sub.EndDate.IsZero() {
		if index != 1 {
			query += sep
		}
		query += fmt.Sprintf(`end_date=$%d`, index)
		args = append(args, sub.EndDate)
		index++
	}

	if sub.UserID != uuid.Nil {
		if index != 1 {
			query += sep
		}
		query += fmt.Sprintf(`user_id=$%d`, index)
		args = append(args, sub.UserID)
		index++
	}

	return query, args, index
}

func InitializationPostgresDB(config PostgresDBParams) error {
	connStr := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
		config.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	postgre := &PostgresDataBase{db: db}

	if err := postgre.runMigrations(); err != nil {
		return err
	}

	DB = postgre
	return nil
}

func (postgre PostgresDataBase) runMigrations() error {
	driver, err := postgres.WithInstance(postgre.db, &postgres.Config{})
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	path := "file://" + filepath.Join(wd, "pkg/dbwork/migrations")

	m, err := migrate.NewWithDatabaseInstance(
		path,
		"postgres", driver,
	)
	if err != nil {
		return err
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return err
	}

	if dirty {
		if err := m.Force(int(version)); err != nil {
			return err
		}
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		if dirtyErr, ok := err.(migrate.ErrDirty); ok {
			if err := m.Force(int(dirtyErr.Version)); err != nil {
				return err
			}
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				return err
			}
		}
		return err
	}
	return nil
}

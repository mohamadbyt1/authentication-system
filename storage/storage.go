package storage
import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"github.com/mohamadbyt1/authentication-system/models"
	"os"
	_ "github.com/lib/pq"
)
type Storage interface{
	CreaateUser(ctx context.Context,user *models.User) error 
	GetUser(ctx context.Context,username string) (*models.User , error)
	DeleteUser(ctx context.Context,username string) error
	UpdateUser(ctx context.Context,user *models.UserSignup)error
	UserCheck(ctx context.Context,username string) (bool,error)
	IdCheck(ctx context.Context,id string) (bool,error)
}
	type PostgreDb struct{
		Db *sql.DB
		CTX    context.Context
	}
	func NewDb() (*PostgreDb, error) {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbSSLMode := os.Getenv("POSTGRES_SSLMODE")
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)
		db, err := sql.Open("postgres", dbInfo)
		if err != nil {
			return nil, err
		}
		if err := db.Ping(); err != nil {
			return nil, err
		}
		createTableSQL := `
		CREATE TABLE IF NOT EXISTS users (
			username VARCHAR(60) UNIQUE NOT NULL,
			firstname VARCHAR(60) NOT NULL,
			lastname VARCHAR(60) NOT NULL,
			password VARCHAR(60) NOT NULL,
			id SERIAL PRIMARY KEY
		);
		  
    `
		_ , err = db.Exec(createTableSQL)
		if err != nil {
			return nil,err
		}
		return &PostgreDb{
			Db:  db,
		}, nil
}

	func (d *PostgreDb) CreaateUser(ctx context.Context, user *models.User) error {
		query := `
			INSERT INTO users (username, firstname, lastname, password)
			VALUES ($1, $2, $3, $4)
		`
		_, err := d.Db.ExecContext(ctx,query, user.Username, user.FirstName, user.LastName, user.Password)
	if err != nil {
		return err
	}
	return nil
	}
	func (d *PostgreDb) GetUser(ctx context.Context,username string) (*models.User, error) {
		var user models.User
		query := "SELECT username, firstname, lastname, password ,id FROM users WHERE username = $1"
		row := d.Db.QueryRowContext(ctx,query, username)
		err := row.Scan(&user.Username, &user.FirstName, &user.LastName, &user.Password, &user.Id)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("user not found")
			}
			return nil, err
		}
		return &user, nil
	}
	
	
	func (d *PostgreDb) DeleteUser(ctx context.Context,id string) error {
		query := "DELETE FROM users WHERE id = $1"
		_ , err := d.Db.ExecContext(ctx,query,id)
		if err != nil {
			return  err
		}
		return nil
	}

	func (d *PostgreDb) UpdateUser(ctx context.Context,user *models.UserSignup)error{
		query := `
		UPDATE users
		SET firstname = $1, lastname = $2, password = $3
		WHERE username = $4
	`	
	_, err:=  d.Db.ExecContext(ctx,query,user.FirstName,user.LastName,user.Password,user.Username)
	if err != nil {
		return err
	}
	return nil
	}
	func (d *PostgreDb) UserCheck(ctx context.Context,username string) (bool,error) {
		log.Println("user in func:",username)
		query := `SELECT COUNT(*) FROM users WHERE username = $1`
		var c int

		err := d.Db.QueryRowContext(ctx,query, username).Scan(&c)
		if err != nil {
			log.Println("++Failed to query", err)
			return false,err
		}
		return c > 0,nil
	}
	func (d *PostgreDb) IdCheck(ctx context.Context,id string) (bool,error) {
		log.Println("user in func:",id)
		query := `SELECT COUNT(*) FROM users WHERE id = $1`
		var c int

		err := d.Db.QueryRowContext(ctx,query, id).Scan(&c)
		if err != nil {
			log.Println("++Failed to query", err)
			return false,err
		}
		return c > 0,nil
	}
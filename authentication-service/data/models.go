package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 3

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		Users: User{},
	}
}

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"-"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := "SELECT id, email, first_name, last_name, user_active, created_at, updated_at FROM users ORDER BY last_name"

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Active, &user.CreatedAt, &user.UpdatedAt)

		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := "SELECT id, email, first_name, last_name, user_active, created_at, updated_at FROM users WHERE email = $1"

	err := db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Active, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) GetOne(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := "SELECT id, email, first_name, last_name, user_active, created_at, updated_at FROM users WHERE id = $1"

	err := db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Active, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := "UPDATE users SET email = $1, first_name = $2, last_name = $3, user_active = $4, updated_at = $5 WHERE id = $6"

	_, err := db.ExecContext(ctx, query, u.Email, u.FirstName, u.LastName, u.Active, time.Now(), u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := "DELETE FROM users WHERE id = $1"

	_, err := db.ExecContext(ctx, query, u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) DeleteByID(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := "DELETE FROM users WHERE id = $1"

	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) Insert() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := "INSERT INTO users (email, first_name, last_name, password, user_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"

	err := db.QueryRowContext(ctx, query, u.Email, u.FirstName, u.LastName, u.Password, u.Active, time.Now(), time.Now()).Scan(&u.ID)
	if err != nil {
		return 0, err
	}

	return u.ID, nil
}

func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt := "UPDATE users SET password = $1 WHERE id = $2"

	_, err = db.ExecContext(ctx, stmt, hashedPassword, hashedPassword, u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

type Models struct {
	Users User
}

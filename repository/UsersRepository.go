package repository

import (
	"database/sql"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

type User struct {
	ID       int
	Username string
	Email    string
}

func (r *UserRepository) GetAllUsers() ([]User, error) {
	rows, err := r.DB.Query("select id, username, email from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Username, &u.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepository) GetUserByID(id int) (*User, error) {
	row := r.DB.QueryRow("select id, username, email from users where id = ?", id)
	var u User
	err := row.Scan(&u.ID, &u.Username, &u.Email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) CreateUser(username, email string) error {
	_, err := r.DB.Exec("insert into users (username, email) values (?, ?)", username, email)
	return err
}

func (r *UserRepository) UpdateUser(id int, username, email string) error {
	_, err := r.DB.Exec("update users set username = ?, email = ? where id = ?", username, email, id)
	return err
}

func (r *UserRepository) DeleteUserByID(id int) error {
	_, err := r.DB.Exec("delete from users where id = ?", id)
	return err
}
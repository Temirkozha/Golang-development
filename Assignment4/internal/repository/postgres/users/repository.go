package users

import (
	"database/sql"
	"errors"

	"assignment3/pkg/modules"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(u modules.User) (int, error) {
	query := `INSERT INTO users (name, email, age, password) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err := r.db.QueryRow(query, u.Name, u.Email, u.Age, u.Password).Scan(&id)
	return id, err
}

func (r *UserRepository) GetUserByID(id int) (*modules.User, error) {
	query := `SELECT id, name, email, age, password FROM users WHERE id = $1`
	var u modules.User

	err := r.db.QueryRow(query, id).Scan(&u.ID, &u.Name, &u.Email, &u.Age, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {

			return nil, errors.New("пользователь с таким ID не найден")
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetAllUsers() ([]modules.User, error) {
	query := `SELECT id, name, email, age, password FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []modules.User
	for rows.Next() {
		var u modules.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Age, &u.Password); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) UpdateUser(id int, u modules.User) error {
	query := `UPDATE users SET name=$1, email=$2, age=$3, password=$4 WHERE id=$5`
	result, err := r.db.Exec(query, u.Name, u.Email, u.Age, u.Password, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("user for update not found")
	}
	return nil
}

func (r *UserRepository) DeleteUser(id int) (int64, error) {
	query := `DELETE FROM users WHERE id=$1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	if rowsAffected == 0 {
		return 0, errors.New("user to delete not found")
	}
	return rowsAffected, nil
}

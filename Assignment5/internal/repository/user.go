package repository

import (
	"assignment5/internal/models"
	"database/sql"
	"strconv"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// 1. Пагинация + Фильтрация + Сортировка (Требование для макс балла)
func (r *UserRepository) GetPaginatedUsers(page, pageSize int, name, orderBy string) (models.PaginatedResponse, error) {
	var users []models.User
	offset := (page - 1) * pageSize

	// Базовый запрос
	query := `SELECT id, name, email, gender, birth_date FROM users WHERE 1=1`
	args := []interface{}{}
	argCounter := 1

	// Динамическая фильтрация по имени
	if name != "" {
		query += ` AND name ILIKE $` + strconv.Itoa(argCounter)
		args = append(args, "%"+name+"%")
		argCounter++
	}

	// Динамическая сортировка (защита от SQL инъекций)
	if orderBy == "name" || orderBy == "email" || orderBy == "id" {
		query += ` ORDER BY ` + orderBy + ` ASC`
	} else {
		query += ` ORDER BY id ASC` // По умолчанию
	}

	// Добавляем лимит и отступ
	query += ` LIMIT $` + strconv.Itoa(argCounter) + ` OFFSET $` + strconv.Itoa(argCounter+1)
	args = append(args, pageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return models.PaginatedResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return models.PaginatedResponse{}, err
		}
		users = append(users, u)
	}

	return models.PaginatedResponse{
		Data:       users,
		TotalCount: len(users), // Упрощенный подсчет для практики
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// 2. Получение общих друзей одним JOIN запросом (Требование для макс балла)
func (r *UserRepository) GetCommonFriends(userID1, userID2 int) ([]models.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.gender, u.birth_date 
		FROM users u
		JOIN user_friends uf1 ON u.id = uf1.friend_id
		JOIN user_friends uf2 ON u.id = uf2.friend_id
		WHERE uf1.user_id = $1 AND uf2.user_id = $2
	`
	rows, err := r.db.Query(query, userID1, userID2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate); err != nil {
			return nil, err
		}
		friends = append(friends, u)
	}
	return friends, nil
}

// 3. Добавление в друзья (Транзакции)
func (r *UserRepository) AddFriend(userID, friendID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO user_friends (user_id, friend_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", userID, friendID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("INSERT INTO user_friends (user_id, friend_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", friendID, userID)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

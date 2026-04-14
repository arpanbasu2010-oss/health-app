package impl

import (
	"database/sql"
	"errors"
	"go-api/model"
)

type UserRepoImpl struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepoImpl {
	return &UserRepoImpl{DB: db}
}

// FindAll fetches all users
func (r *UserRepoImpl) FindAll() ([]model.User, error) {
	rows, err := r.DB.Query(
		`SELECT id, name, email, age, created_at, updated_at FROM users ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(
			&u.ID, &u.Name, &u.Email, &u.Age, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// FindByID fetches a user by ID
func (r *UserRepoImpl) FindByID(id int) (*model.User, error) {
	var u model.User
	err := r.DB.QueryRow(
		`SELECT id, name, email, age, created_at, updated_at FROM users WHERE id = $1`, id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Age, &u.CreatedAt, &u.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return &u, err
}

// FindByEmail checks if email already exists
func (r *UserRepoImpl) FindByEmail(email string) (*model.User, error) {
	var u model.User
	err := r.DB.QueryRow(
		`SELECT id, name, email, age, created_at, updated_at FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Age, &u.CreatedAt, &u.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil // no user found = email is available
	}
	return &u, err
}

// Create inserts a new user
func (r *UserRepoImpl) Create(user *model.User) (*model.User, error) {
	err := r.DB.QueryRow(
		`INSERT INTO users (name, email, age) 
		 VALUES ($1, $2, $3) 
		 RETURNING id, name, email, age, created_at, updated_at`,
		user.Name, user.Email, user.Age,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt, &user.UpdatedAt)

	return user, err
}

// Update modifies an existing user
func (r *UserRepoImpl) Update(user *model.User) (*model.User, error) {
	err := r.DB.QueryRow(
		`UPDATE users 
		 SET name=$1, email=$2, age=$3, updated_at=NOW() 
		 WHERE id=$4 
		 RETURNING id, name, email, age, created_at, updated_at`,
		user.Name, user.Email, user.Age, user.ID,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt, &user.UpdatedAt)

	return user, err
}

// Delete removes a user by ID
func (r *UserRepoImpl) Delete(id int) error {
	result, err := r.DB.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

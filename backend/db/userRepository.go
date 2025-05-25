package db

import (
	"backend/models"
	"database/sql"
	"errors"
)

// UserRepository handles all database operations related to users.
// Fields:
// - DB: a pointer to the SQL database connection
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new instance of UserRepository.
// Parameters:
// - db: a pointer to the SQL database connection
// Returns: a pointer to the newly created UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

// CreateUser adds a new user to the database.
// Parameters:
// - user: a pointer to the User object to be added
// Returns: the ID of the newly created user, or an error if the insertion fails
func (r *UserRepository) CreateUser(user *models.User) (uint32, error) {
	query := `INSERT INTO users (name, email, pub_key, login_salt, encryption_salt, hmac_salt, hmac_type, encryption_type) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.DB.Exec(query, user.Name, user.Email, user.PubKey, user.LoginSalt, user.EncryptionSalt,
		user.HMACSalt, user.HMACType, user.EncryptionType)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint32(id), nil
}

// GetUserByEmail finds a user by their email address.
// Parameters:
// - email: the email address of the user to find
// Returns: a pointer to the retrieved User object, or an error if no user is found or a query error occurs
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, name, email, pub_key, login_salt, encryption_salt, hmac_salt, hmac_type, encryption_type 
              FROM users WHERE email = ?`

	var user models.User
	err := r.DB.QueryRow(query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.PubKey, &user.LoginSalt,
		&user.EncryptionSalt, &user.HMACSalt, &user.HMACType, &user.EncryptionType)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByID finds a user by their ID.
// Parameters:
// - id: the ID of the user to find
// Returns: a pointer to the retrieved User object, or an error if no user is found or a query error occurs
func (r *UserRepository) GetUserByID(id uint32) (*models.User, error) {
	query := `SELECT id, name, email, pub_key, login_salt, encryption_salt, hmac_salt, hmac_type, encryption_type, login_salt
              FROM users WHERE id = ?`

	var user models.User
	err := r.DB.QueryRow(query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.PubKey, &user.LoginSalt,
		&user.EncryptionSalt, &user.HMACSalt, &user.HMACType, &user.EncryptionType, &user.LoginSalt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// UpdateUser updates an existing user in the database.
// Parameters:
// - user: a pointer to the User object containing updated information
// Returns: an error if the update operation fails
func (r *UserRepository) UpdateUser(user *models.User) error {
	query := `
		UPDATE users 
		SET name = ?, email = ?
		WHERE id = ?
	`
	_, err := r.DB.Exec(query,
		user.Name,
		user.Email,
		user.ID,
	)
	return err
}

// DeleteUserByID deletes a user from the database by their ID.
// Parameters:
// - id: the ID of the user to delete
// Returns: an error if the deletion operation fails
func (r *UserRepository) DeleteUserByID(id uint32) error {
	query := `DELETE FROM users WHERE id = ?`

	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

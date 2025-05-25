package db

import (
	"backend/models"
	"database/sql"
	"errors"
)

// ChallengeRepository handles all database operations related to login challenges.
// Fields:
// - DB: a pointer to the SQL database connection
type ChallengeRepository struct {
	DB *sql.DB
}

// NewChallengeRepository creates a new instance of ChallengeRepository.
// Parameters:
// - db: a pointer to the SQL database connection
// Returns: a pointer to the newly created ChallengeRepository
func NewChallengeRepository(db *sql.DB) *ChallengeRepository {
	return &ChallengeRepository{
		DB: db,
	}
}

// CreateChallenge adds a new challenge to the database.
// Parameters:
// - challenge: a pointer to the Challenge object to be added
// Returns: the ID of the newly created challenge, or an error if the insertion fails
func (r *ChallengeRepository) CreateChallenge(challenge *models.Challenge) (uint32, error) {
	query := `INSERT INTO challenges (user_id, challenge_value, expires_at, used) 
              VALUES (?, ?, ?, ?)`

	result, err := r.DB.Exec(query, challenge.UserID, challenge.ChallengeValue,
		challenge.ExpiresAt, challenge.Used)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint32(id), nil
}

// GetChallenge finds a challenge by its value and email.
// Parameters:
// - challengeValue: the value of the challenge to find
// - Email: the email associated with the challenge
// Returns: a pointer to the retrieved Challenge object, or an error if no challenge is found or a query error occurs
func (r *ChallengeRepository) GetChallenge(challengeValue string, Email string) (*models.Challenge, error) {
	query := `SELECT c.id, c.user_id, c.challenge_value, c.expires_at, c.used, c.created_at 
			  FROM challenges c 
			  JOIN users u ON c.user_id = u.id 
			  WHERE c.challenge_value = ? AND u.email = ?`

	var challenge models.Challenge
	err := r.DB.QueryRow(query, challengeValue, Email).Scan(
		&challenge.ID, &challenge.UserID, &challenge.ChallengeValue,
		&challenge.ExpiresAt, &challenge.Used, &challenge.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("challenge not found")
		}
		return nil, err
	}

	return &challenge, nil
}

// MarkChallengeAsUsed marks a challenge as used.
// Parameters:
// - challengeID: the ID of the challenge to mark as used
// Returns: an error if the update operation fails
func (r *ChallengeRepository) MarkChallengeAsUsed(challengeID uint32) error {
	query := `UPDATE challenges SET used = true WHERE id = ?`

	_, err := r.DB.Exec(query, challengeID)
	return err
}

// DeleteChallenge removes a specific challenge from the database.
// Parameters:
// - challengeID: the ID of the challenge to delete
// Returns: an error if the deletion operation fails
func (r *ChallengeRepository) DeleteChallenge(challengeID uint32) error {
	query := `DELETE FROM challenges WHERE id = ?`

	_, err := r.DB.Exec(query, challengeID)
	return err
}

// ClearExpiredChallenges removes all expired challenges from the database.
// Returns: an error if the deletion operation fails
func (r *ChallengeRepository) ClearExpiredChallenges() error {
	query := `DELETE FROM challenges WHERE expires_at < NOW()`

	_, err := r.DB.Exec(query)
	return err
}

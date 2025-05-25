package db

import (
	"backend/models"
	"database/sql"
	"fmt"
)

// BlockRepository provides methods to interact with the blocks database table.
// Fields:
// - DB: a pointer to the SQL database connection
type BlockRepository struct {
	DB *sql.DB
}

// NewBlockRepository creates a new instance of BlockRepository.
// Parameters:
// - db: a pointer to the SQL database connection
// Returns: a pointer to the newly created BlockRepository
func NewBlockRepository(db *sql.DB) *BlockRepository {
	return &BlockRepository{
		DB: db,
	}
}

// GetNoteBlock retrieves the latest block for a specific note and user.
// Parameters:
// - userID: the ID of the user
// - noteID: the ID of the note
// Returns: a pointer to the retrieved block, or an error if no block is found or a query error occurs
func (r *BlockRepository) GetNoteBlock(userID uint32, noteID uint) (*models.Block, error) {
	const query = `
        SELECT note_id, user_id, prev_hash, timestamp, iv, iv_title, cipher_title, ciphertext, mac, signature
        FROM blocks
        WHERE note_id = ? AND user_id = ?
        ORDER BY timestamp DESC
        LIMIT 1
    `
	row := r.DB.QueryRow(query, noteID, userID)

	block := &models.Block{}
	if err := row.Scan(
		&noteID,
		&userID,
		&block.PrevHash,
		&block.Timestamp,
		&block.IV,
		&block.IVTitle,
		&block.CipherTitle,
		&block.Ciphertext,
		&block.MAC,
		&block.Signature,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no block found for noteID %d and userID %d", noteID, userID)
		}
		return nil, fmt.Errorf("error scanning block: %v", err)
	}

	return block, nil
}

// GetNoteBlockChain retrieves the entire blockchain for a specific note and user.
// Parameters:
// - userID: the ID of the user
// - noteID: the ID of the note
// Returns: a pointer to the NoteBlockChain containing all blocks, or an error if a query error occurs
func (r *BlockRepository) GetNoteBlockChain(userID uint32, noteID uint) (*models.NoteBlockChain, error) {
	const query = `
        SELECT prev_hash, timestamp, iv, iv_title, cipher_title, ciphertext, mac, signature
        FROM blocks
        WHERE note_id = ? AND user_id = ?
        ORDER BY timestamp ASC
    `

	rows, err := r.DB.Query(query, noteID, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying blocks: %v", err)
	}
	defer rows.Close()

	var blocks []models.Block
	for rows.Next() {
		var block models.Block
		if err := rows.Scan(
			&block.PrevHash,
			&block.Timestamp,
			&block.IV,
			&block.IVTitle,
			&block.CipherTitle,
			&block.Ciphertext,
			&block.MAC,
			&block.Signature,
		); err != nil {
			return nil, fmt.Errorf("error scanning block: %v", err)
		}
		blocks = append(blocks, block)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return &models.NoteBlockChain{
		NoteID: noteID,
		Blocks: blocks,
	}, nil
}

// CreateBlock inserts a new block into the database for a specific note and user.
// Parameters:
// - userID: the ID of the user
// - noteID: the ID of the note
// - block: a pointer to the block to be inserted
// Returns: the timestamp of the inserted block, or an error if the insertion fails
func (r *BlockRepository) CreateBlock(userID uint32, noteID uint, block *models.Block) error {
	const query = `
		INSERT INTO blocks (note_id, user_id, prev_hash, timestamp, iv, iv_title, cipher_title, ciphertext, mac, signature)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Execute the insert query
	_, err := r.DB.Exec(query,
		noteID,
		userID,
		block.PrevHash,
		block.Timestamp,
		block.IV,
		block.IVTitle,
		block.CipherTitle,
		block.Ciphertext,
		block.MAC,
		block.Signature,
	)
	if err != nil {
		return err
	}

	return nil
}

// CreateNewNote creates a new note by incrementing the note ID for the user and inserting the first block.
// Parameters:
// - userID: the ID of the user
// - block: a pointer to the block to be inserted
// Returns: the new note ID, the timestamp of the inserted block, or an error if the operation fails
func (r *BlockRepository) CreateNewNote(userID uint32, block *models.Block) (uint, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// First get the new note_id
	const getMaxNoteID = `SELECT COALESCE(MAX(note_id), 0) + 1 FROM blocks WHERE user_id = ?`
	var noteID uint
	err = tx.QueryRow(getMaxNoteID, userID).Scan(&noteID)
	if err != nil {
		return 0, err
	}

	// Then insert the new block
	const insertQuery = `
		INSERT INTO blocks (note_id, user_id, prev_hash, timestamp, iv, iv_title, cipher_title, ciphertext, mac, signature)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = tx.Exec(insertQuery,
		noteID,
		userID,
		block.PrevHash,
		block.Timestamp,
		block.IV,
		block.IVTitle,
		block.CipherTitle,
		block.Ciphertext,
		block.MAC,
		block.Signature,
	)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return noteID, nil
}

// GetTitles retrieves the latest cipher title and timestamp for each note of a user.
// Parameters:
// - userID: the ID of the user
// Returns: a slice of Title objects containing the note ID, cipher title, IV, and timestamp, or an error if a query error occurs
func (r *BlockRepository) GetTitles(userID uint32) ([]*models.Title, error) {
	// Updated query to return the latest block for each note_id
	const query = `
		SELECT b.note_id, b.cipher_title, b.iv_title, b.timestamp
		FROM blocks b
		INNER JOIN (
			SELECT note_id, MAX(timestamp) AS max_timestamp
			FROM blocks
			WHERE user_id = ?
			GROUP BY note_id
		) latest_blocks
		ON b.note_id = latest_blocks.note_id AND b.timestamp = latest_blocks.max_timestamp
		ORDER BY b.timestamp DESC
	`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var titles []*models.Title
	for rows.Next() {
		title := &models.Title{}
		if err := rows.Scan(
			&title.NoteID,
			&title.CipherTitle,
			&title.IV,
			&title.Timestamp,
		); err != nil {
			return nil, err
		}
		titles = append(titles, title)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return titles, nil
}

// DeleteNoteBlocks deletes all blocks associated with a specific note ID and user ID.
// Parameters:
// - userID: the ID of the user
// - noteID: the ID of the note
// Returns: an error if the deletion fails or no blocks are found
func (r *BlockRepository) DeleteNoteBlocks(userID uint32, noteID uint) error {
	const query = `
		DELETE FROM blocks
		WHERE note_id = ? AND user_id = ?
	`

	result, err := r.DB.Exec(query, noteID, userID)
	if err != nil {
		return fmt.Errorf("error deleting blocks: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no blocks found for noteID %d and userID %d", noteID, userID)
	}

	return nil
}

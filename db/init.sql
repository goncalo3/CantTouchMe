-- The MySQL Docker image automatically handles:
-- 1. Database creation (based on MYSQL_DATABASE env var)
-- 2. User creation (based on MYSQL_USER and MYSQL_PASSWORD env vars)
-- 3. Granting privileges to that user on the specified database

-- Users table for storing registered users
CREATE TABLE users (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    pub_key TEXT NOT NULL,
    login_salt VARCHAR(255) NOT NULL,
    encryption_salt VARCHAR(255) NOT NULL,
    hmac_salt VARCHAR(255) NOT NULL,
    hmac_type VARCHAR(30) NOT NULL,
    encryption_type VARCHAR(30) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX (email)
);

-- Challenges table for login authentication
CREATE TABLE challenges (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id INT UNSIGNED NOT NULL,
    challenge_value VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX (challenge_value),
    INDEX (expires_at)
);

-- Notes table for storing encrypted notes
CREATE TABLE blocks (
    note_id INT UNSIGNED NOT NULL,
    user_id INT UNSIGNED NOT NULL,
    prev_hash VARCHAR(255) NOT NULL, -- Cant't be UNIQUE as all blockchains get initial block with prev_hash = 0
    timestamp TIMESTAMP NOT NULL,
    iv VARCHAR(255) NOT NULL,
    iv_title VARCHAR(255) NOT NULL,
    cipher_title TEXT NOT NULL,
    ciphertext LONGTEXT NOT NULL,
    mac VARCHAR(255) NOT NULL,
    signature TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE, -- if a user is deleted, their notes are also deleted
    PRIMARY KEY (note_id, user_id, prev_hash), -- Composite primary key
    INDEX (user_id),
    INDEX (prev_hash)
);

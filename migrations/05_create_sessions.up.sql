CREATE TABLE IF NOT EXISTS sessions (
                                        id SERIAL PRIMARY KEY,
                                        user_id INTEGER NOT NULL REFERENCES users(id),
                                        session_id VARCHAR(255) NOT NULL,
                                        expires_at TIMESTAMP NOT NULL,
                                        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                        updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
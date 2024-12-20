
CREATE TABLE IF NOT EXISTS addresses (
                                         id SERIAL PRIMARY KEY,
                                         user_id INT REFERENCES users(id) ON DELETE CASCADE,
                                         latitude DOUBLE PRECISION NOT NULL,
                                         longitude DOUBLE PRECISION NOT NULL,
                                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
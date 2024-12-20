CREATE TABLE IF NOT EXISTS restaurants (
                                           id SERIAL PRIMARY KEY,
                                           name VARCHAR(255) NOT NULL,
                                           owner_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                           latitude DOUBLE PRECISION NOT NULL,
                                           longitude DOUBLE PRECISION NOT NULL,
                                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
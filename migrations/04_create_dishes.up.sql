CREATE TABLE IF NOT EXISTS dishes (
                                      id SERIAL PRIMARY KEY,
                                      name VARCHAR(255) NOT NULL,
                                      price DECIMAL(10, 2) NOT NULL,
                                      restaurant_id INT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
                                      tags TEXT[] DEFAULT '{}',
                                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
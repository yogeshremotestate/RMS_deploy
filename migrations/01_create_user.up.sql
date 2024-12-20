DO $$
    BEGIN
        -- Create the roles ENUM type if it does not exist
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'roles') THEN
            CREATE TYPE roles AS ENUM ('admin', 'sub-admin', 'user');
        END IF;
    END $$;

-- Create the users table
CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     name VARCHAR(255) NOT NULL,
                                     email VARCHAR(255) UNIQUE NOT NULL,
                                     password_hash TEXT NOT NULL,
                                     role roles NOT NULL DEFAULT 'user',
                                     created_by INT REFERENCES users(id) ON DELETE SET NULL,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

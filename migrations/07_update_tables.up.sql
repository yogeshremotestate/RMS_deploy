-- Add `deleted_at` column to the `users` table
ALTER TABLE users
    ADD COLUMN deleted_at TIMESTAMP;

-- Add `deleted_at` column to the `restaurants` table
ALTER TABLE restaurants
    ADD COLUMN deleted_at TIMESTAMP;

-- Add `deleted_at` column to the `dishes` table
ALTER TABLE dishes
    ADD COLUMN deleted_at TIMESTAMP;

-- Add `deleted_at` column to the `addresses` table
ALTER TABLE addresses
    ADD COLUMN deleted_at TIMESTAMP;

-- Add `name` column to the `addresses` table
ALTER TABLE addresses
    ADD COLUMN name VARCHAR(255);



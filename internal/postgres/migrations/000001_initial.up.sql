CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    first_name TEXT,
    last_name TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_no BIGINT NOT NULL,
    card_type INTEGER NOT NULL,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    amount REAL NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Insert sample users
INSERT INTO users (email, password, first_name, last_name) VALUES
('john.doe@example.com', 'hashed_password_1', 'John', 'Doe'),
('jane.smith@example.com', 'hashed_password_2', 'Jane', 'Smith'),
('robert.johnson@example.com', 'hashed_password_3', 'Robert', 'Johnson'),
('sarah.williams@example.com', 'hashed_password_4', 'Sarah', 'Williams'),
('michael.brown@example.com', 'hashed_password_5', 'Michael', 'Brown');
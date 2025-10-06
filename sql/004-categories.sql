CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    code VARCHAR(32) UNIQUE NOT NULL,
    name VARCHAR(256) NOT NULL
);

-- Insert the 3 categories with CATGORY codes
INSERT INTO categories (code, name) VALUES
('CATGORY001', 'Clothing'),
('CATGORY002', 'Shoes'),
('CATGORY003', 'Accessories');

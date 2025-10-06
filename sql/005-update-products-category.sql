-- Add category_id column to products table
ALTER TABLE products ADD COLUMN IF NOT EXISTS category_id INTEGER REFERENCES categories(id);

-- Update products with their categories based on the assignment requirements
-- Clothing: PROD001, PROD004, PROD007
UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'CATGORY001') WHERE code IN ('PROD001', 'PROD004', 'PROD007');

-- Shoes: PROD002, PROD006
UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'CATGORY002') WHERE code IN ('PROD002', 'PROD006');

-- Accessories: PROD003, PROD005, PROD008
UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'CATGORY003') WHERE code IN ('PROD003', 'PROD005', 'PROD008');

-- Tabel "users"
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    deposit_amount DECIMAL(10, 2) DEFAULT 0
);

-- Tabel "book_inventory"
CREATE TABLE book_inventory (
    book_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    stock_availability INT NOT NULL,
    rental_costs DECIMAL(10, 2) NOT NULL,
    category VARCHAR(50)
);

-- Tabel "rental_history"
CREATE TABLE rental_history (
    rental_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    book_id INT REFERENCES book_inventory(book_id),
    rental_date DATE NOT NULL,
    return_date DATE,
    rental_cost DECIMAL(10, 2) NOT NULL
);

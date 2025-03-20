CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    user_username VARCHAR(50) UNIQUE NOT NULL,
    user_password VARCHAR(255) NOT NULL,
    user_role VARCHAR(20) CHECK (user_role IN ('admin', 'user')) NOT NULL
);

CREATE TABLE customers (
    customer_id SERIAL PRIMARY KEY,
    customer_nik VARCHAR(16) UNIQUE NOT NULL,
    customer_full_name VARCHAR(100) NOT NULL,
    customer_legal_name VARCHAR(100) NOT NULL,
    customer_birth_place VARCHAR(50) NOT NULL,
    customer_birth_date DATE NOT NULL,
    customer_salary DECIMAL(15,2) NOT NULL,
    customer_ktp_photo TEXT NOT NULL,
    customer_selfie_photo TEXT NOT NULL,
    customer_created_by INT NOT NULL REFERENCES users(user_id),
    customer_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE limits (
    limit_id SERIAL PRIMARY KEY,
    limit_nik VARCHAR(16) NOT NULL REFERENCES customers(customer_nik) ON DELETE CASCADE,
    limit_tenor INT CHECK (limit_tenor IN (1, 2, 3, 6)) NOT NULL,
    limit_amount DECIMAL(15,2) NOT NULL,
    limit_used_amount DECIMAL(15,2) DEFAULT 0,
    limit_remaining_amount DECIMAL(15,2) NOT NULL,
    limit_created_by INT NOT NULL REFERENCES users(user_id),
    limit_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transactions (
    transaction_id SERIAL PRIMARY KEY,
    transaction_contract_number VARCHAR(50) UNIQUE NOT NULL,
    transaction_nik VARCHAR(16) NOT NULL REFERENCES customers(customer_nik) ON DELETE CASCADE,
    transaction_limit INT NOT NULL REFERENCES limits(limit_id) ON DELETE CASCADE,
    transaction_otr DECIMAL(15,2) NOT NULL,
    transaction_admin_fee DECIMAL(15,2) NOT NULL,
    transaction_installment INT NOT NULL,
    transaction_interest DECIMAL(5,2) NOT NULL,
    transaction_asset_name VARCHAR(100) NOT NULL,
    transaction_created_by INT NOT NULL REFERENCES users(user_id),
    transaction_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
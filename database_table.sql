CREATE TABLE customers ( 
    customers_id SERIAL PRIMARY KEY,
    customers_nik VARCHAR(16) UNIQUE NOT NULL,
    customers_full_name VARCHAR(255) NOT NULL,
    customers_legal_name VARCHAR(255),
    customers_birth_place VARCHAR(255),
    customers_birth_date DATE,
    customers_salary DECIMAL,
    customers_ktp_photo text,
    customers_selfie_photo text,
    customers_created_at timestamp DEFAULT NOW()
);

CREATE TABLE limits ( 
    limits_id SERIAL PRIMARY KEY,
    limits_customers_id INT REFERENCES customers(customers_id) ON DELETE CASCADE,
    limits_tenor INT CHECK (limits_tenor IN (1,2,3,6)),
    limits_amount DECIMAL NOT NULL
);

CREATE TABLE transactions ( 
    transactions_id SERIAL PRIMARY KEY,
    transactions_contract_no varchar(255) UNIQUE NOT NULL,
    transactions_customers_id INT REFERENCES customers(customers_id) ON DELETE CASCADE,
    transactions_otr DECIMAL NOT NULL,
    transactions_admin_fee DECIMAL NOT NULL,
    transactions_instalment DECIMAL NOT NULL,
    transactions_interest DECIMAL NOT NULL,
    transactions_asset_name varchar(255),
    transactions_date TIMESTAMP DEFAULT NOW()
);
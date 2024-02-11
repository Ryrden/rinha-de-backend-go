CREATE TABLE IF NOT EXISTS clients (
    id SERIAL PRIMARY KEY, 
    client_name text NOT NULL,
    balance_limit integer NOT NULL,
    balance integer NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS transactions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id integer NOT NULL,
    amount integer NOT NULL,
    kind char NOT NULL CHECK (kind IN ('c', 'd')), -- 'c' para crédito, 'd' para débito
    description text NOT NULL CHECK (char_length(description) >= 1 AND char_length(description) <= 10),
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (client_id) REFERENCES clients(id)
);

INSERT INTO clients (client_name, balance_limit)
VALUES
('o barato sai caro', 1000 * 100),
('zan corp ltda', 800 * 100),
('les cruders', 10000 * 100),
('padaria joia de cocaia', 100000 * 100),
('kid mais', 5000 * 100);
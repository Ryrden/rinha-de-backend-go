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
    kind char NOT NULL,
    description text NOT NULL,
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


-- Index client_id and created_at columns

CREATE INDEX IF NOT EXISTS transactions_client_id_idx ON transactions(client_id);
CREATE INDEX IF NOT EXISTS transactions_created_at_idx ON transactions(created_at);
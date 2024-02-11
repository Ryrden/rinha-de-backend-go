CREATE TABLE IF NOT EXISTS public.client (
    id uuid PRIMARY KEY NOT NULL,
    balance_limit integer NOT NULL,
    balance integer NOT NULL
);

CREATE TABLE IF NOT EXISTS public.transactions (
    id uuid PRIMARY KEY NOT NULL,
    client_id uuid NOT NULL,
    amount integer NOT NULL,
    kind char NOT NULL CHECK (kind IN ('c', 'd')), -- 'c' para crédito, 'd' para débito
    description text NOT NULL CHECK (char_length(description) >= 1 AND char_length(description) <= 10),
    created_at timestamp NOT NULL,
    FOREIGN KEY (client_id) REFERENCES public.client(id)
);

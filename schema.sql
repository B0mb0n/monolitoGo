-- Se ejecuta en automatico cuando se arranca el
-- contenedor de Postgres por primera vez.

CREATE TABLE IF NOT EXISTS users (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    email      VARCHAR(150) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Datos de ejemplo para probar.
INSERT INTO users (name, email) VALUES
    ('Aarón Salto', 'aaron@example.com'),
    ('Alejandra', 'ale@example.com')

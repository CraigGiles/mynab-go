CREATE TABLE IF NOT EXISTS public.users(
   id SERIAL PRIMARY KEY,
   name TEXT,
   age INT
);

INSERT INTO users (id, name, age) VALUES (1, 'craig', 42);

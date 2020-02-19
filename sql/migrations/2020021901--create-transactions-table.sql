CREATE TABLE IF NOT EXISTS public.transactions(
   id TEXT PRIMARY KEY,
   date TEXT,
   payee TEXT,
   category TEXT,
   amount NUMERIC
);

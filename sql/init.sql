CREATE DATABASE mynab OWNER postgres;
\c mynab;

CREATE TABLE public.table_migrations(
    version TEXT,
    script TEXT,
    file_name TEXT,
    file_hash TEXT,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users(
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    description TEXT NOT NULL
);
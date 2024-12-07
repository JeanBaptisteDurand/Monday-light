CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    discord_id VARCHAR(255),
    discord_pseudo VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    categories TEXT[] NOT NULL DEFAULT '{}'
);

CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(255),
    project_id INT REFERENCES projects(id),
    status VARCHAR(50) DEFAULT 'to_assign',
    estimated_time INT DEFAULT 0,
    real_time INT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    taken_from TIMESTAMP
);

-- Table de liaison pour relation many-to-many entre users et tasks
CREATE TABLE IF NOT EXISTS user_tasks (
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    task_id INT REFERENCES tasks(id) ON DELETE CASCADE,
    PRIMARY KEY(user_id, task_id)
);

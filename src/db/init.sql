-- src/db/init.sql
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    categories TEXT[] NOT NULL DEFAULT '{}'
);

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(255),
    project_id INT REFERENCES projects(id),
    status VARCHAR(50) DEFAULT 'to_assign',
    assigned_users INT[],
    estimated_time INT DEFAULT 0,
    real_time INT DEFAULT 0
);

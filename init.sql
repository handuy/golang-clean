CREATE TABLE users (
    id varchar(50) PRIMARY KEY,
    email varchar(50) UNIQUE,
    password text,
    created_at datetime DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE notes (
    id int AUTO_INCREMENT PRIMARY KEY,
    title text,
    status bool,
    author text,
    created_at datetime DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime DEFAULT CURRENT_TIMESTAMP
);
DROP DATABASE IF EXISTS DATABASE_NAME;

CREATE DATABASE DATABASE_NAME;

USE DATABASE_NAME;

CREATE TABLE users(
  id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL
);

CREATE INDEX user_email
ON users(email);

CREATE TABLE sessions(
  token VARCHAR(255) NOT NULL PRIMARY KEY,
  user_id INT UNSIGNED NOT NULL,
  expires_at DATETIME NOT NULL,
  FOREIGN KEY(user_id) REFERENCES users(id)
);

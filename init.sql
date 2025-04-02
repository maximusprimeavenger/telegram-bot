CREATE DATABASE IF NOT EXISTS notifier;
USE notifier;
CREATE TABLE users (
    id          INT AUTO_INCREMENT PRIMARY KEY,
    email       VARCHAR(191) NOT NULL UNIQUE,
    name        VARCHAR(30) NOT NULL,
    username    VARCHAR(191) NOT NULL UNIQUE,
    notifier_id VARCHAR(191) NOT NULL UNIQUE,
    notify_mode BOOL NOT NULL
);
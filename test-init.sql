CREATE DATABASE IF NOT EXISTS notifier_test;
USE notifier_test;

CREATE USER IF NOT EXISTS  'testuser'@'%' IDENTIFIED BY 'testpass';

GRANT ALL PRIVILEGES ON notifier_test.* TO 'testuser'@'%';

FLUSH PRIVILEGES;

CREATE TABLE IF NOT EXISTS test_users(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    email       VARCHAR(191) NOT NULL UNIQUE,
    name        VARCHAR(30) NOT NULL
)
-- init/01-init.sql
-- Initial database schema for HTTP API

-- CREATE DATABASE IF NOT EXISTS SCHOOL_DB;
-- USE SCHOOL_DB;
-- Create tables for your models
CREATE TABLE IF NOT EXISTS teachers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    class VARCHAR(255) NOT NULL,
    subject VARCHAR(255) NOT NULL,
    INDEX idx_email (email),
    INDEX idx_class (class)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS students (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    class VARCHAR(255) NOT NULL,
    FOREIGN KEY (class) REFERENCES teachers(class) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES
('Sarah', 'Johnson', 'sarah.johnson@school.com', '10A', 'Mathematics'),
('Michael', 'Chen', 'michael.chen@school.com', '10B', 'Science'),
('Emily', 'Rodriguez', 'emily.rodriguez@school.com', '10C', 'English'),
('David', 'Patel', 'david.patel@school.com', '11A', 'Physics'),
('Lisa', 'Anderson', 'lisa.anderson@school.com', '11B', 'Chemistry'),
('James', 'Williams', 'james.williams@school.com', '11C', 'Biology'),
('Maria', 'Garcia', 'maria.garcia@school.com', '9A', 'History'),
('Robert', 'Martinez', 'robert.martinez@school.com', '9B', 'Geography'),
('Jennifer', 'Taylor', 'jennifer.taylor@school.com', '9C', 'Art'),
('William', 'Brown', 'william.brown@school.com', '12A', 'Computer Science'),
('Patricia', 'Davis', 'patricia.davis@school.com', '12B', 'Economics'),
('Richard', 'Miller', 'richard.miller@school.com', '12C', 'Literature'),
('Jessica', 'Wilson', 'jessica.wilson@school.com', '10A', 'Physical Education'),
('Joseph', 'Moore', 'joseph.moore@school.com', '10B', 'Music'),
('Susan', 'Thomas', 'susan.thomas@school.com', '11A', 'Spanish'),
('Charles', 'Jackson', 'charles.jackson@school.com', '11B', 'French'),
('Karen', 'White', 'karen.white@school.com', '11C', 'German'),
('Christopher', 'Harris', 'christopher.harris@school.com', '9A', 'Philosophy'),
('Nancy', 'Martin', 'nancy.martin@school.com', '9B', 'Psychology'),
('Daniel', 'Thompson', 'daniel.thompson@school.com', '9C', 'Sociology'),
('Betty', 'Garcia', 'betty.garcia@school.com', '12A', 'Statistics'),
('Matthew', 'Martinez', 'matthew.martinez@school.com', '12B', 'Calculus'),
('Margaret', 'Robinson', 'margaret.robinson@school.com', '12C', 'Algebra'),
('Anthony', 'Clark', 'anthony.clark@school.com', '10A', 'Geometry'),
('Amanda', 'Lewis', 'amanda.lewis@school.com', '10B', 'Trigonometry')
ON DUPLICATE KEY UPDATE email=email;  -- Skip if email already exists

-- INSERT INTO students (first_name, last_name, email, grade, teacher_id) VALUES
--     ('Alice', 'Johnson', 'alice@student.com', 10, 1),
--     ('Bob', 'Williams', 'bob@student.com', 11, 2);

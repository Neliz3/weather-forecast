CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    frequency VARCHAR(10) NOT NULL CHECK (frequency IN ('hourly', 'daily')),
    confirmed BOOLEAN DEFAULT FALSE
);

CREATE UNIQUE INDEX idx_unique_email_city ON subscriptions(email, city);

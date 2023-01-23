CREATE TABLE person (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(100) DEFAULT NULL,
    date_of_birth VARCHAR(100) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS relationships (
    id SERIAL PRIMARY KEY,
    person_id INTEGER REFERENCES person(id),
    related_person_id INTEGER REFERENCES person(id),
    relationship VARCHAR(255) NOT NULL
);
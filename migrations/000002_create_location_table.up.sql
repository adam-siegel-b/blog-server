CREATE TABLE IF NOT EXISTS location(
    id VARCHAR(256) PRIMARY KEY NOT NULL UNIQUE,
    name VARCHAR(256),
    latlon point NOT NULL 
);

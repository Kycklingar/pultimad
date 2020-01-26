CREATE TABLE creators (
	id INTEGER NOT NULL PRIMARY KEY,
	name TEXT NOT NULL,
	added TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	downloaded TIMESTAMP,
	download BOOL NOT NULL DEFAULT false
);
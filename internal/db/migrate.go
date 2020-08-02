package db

const (
	// InitialMigration is the script to create all the tables in the database
	InitialMigration = `
CREATE TABLE IF NOT EXISTS bank (
	id		INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	alias	VARCHAR(256) UNIQUE NOT NULL,
	name	VARCHAR(256) NOT NULL
);

CREATE TABLE IF NOT EXISTS deposit (
	id					INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	alias				VARCHAR(256) NOT NULL,
	name				VARCHAR(256) NOT NULL,
	bank_id 			INTEGER NOT NULL,
	minimal_amount		INTEGER DEFAULT 0,
	rate				REAL NOT NULL,
	has_replenishment	BOOLEAN DEFAULT FALSE,
	is_updated			BOOLEAN DEFAULT TRUE,
	detail				VARCHAR(1024),
	previous_rate		REAL,
	off					BOOLEAN DEFAULT FALSE,
	is_exist            BOOLEAN DEFAULT TRUE,
	updated_at			DATE DEFAULT CURRENT_DATE,
	FOREIGN KEY (bank_id) REFERENCES bank(id),
	UNIQUE (bank_id, alias)
);

CREATE TABLE IF NOT EXISTS deposit_details (
	id					INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	deposit_id			INTEGER NOT NULL,
	full_description	TEXT,
	FOREIGN KEY (deposit_id) REFERENCES deposit(id),
	UNIQUE (deposit_id)
)`
)

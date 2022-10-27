package storage

const (
	createUrls = `CREATE TABLE IF NOT EXISTS urls (
            	short varchar(255) PRIMARY KEY,
                original varchar(255),
    			userid varchar(64),
    			UNIQUE(original)
                )`
	addQuery = `
	INSERT INTO urls (short, original, userid)
	VALUES ($1, $2, $3)
	ON CONFLICT DO NOTHING
	RETURNING short`
	countURL       = `SELECT COUNT(*) FROM urls`
	getQuery       = `SELECT original FROM urls WHERE short = $1`
	getByUserQuery = `SELECT * FROM urls WHERE userid = $1`
	getShort       = `SELECT short FROM urls WHERE original = $1`
	drop           = `DROP TABLE urls`
)

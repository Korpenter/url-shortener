package storage

const (
	createUrls = `CREATE TABLE IF NOT EXISTS urls (
            	short varchar(255) PRIMARY KEY,
                original varchar(255),
    			userid varchar(64),
    			deleted boolean DEFAULT false,
    			UNIQUE(original)
                )`
	addQuery = `
	INSERT INTO urls (short, original, userid)
	VALUES ($1, $2, $3)
	ON CONFLICT DO NOTHING
	RETURNING short`
	updateDeleteQuery = `DELETE FROM urls WHERE short IN (SELECT unnest($1::text[])) AND userid = $2`
	getQuery          = `SELECT original, deleted FROM urls WHERE short = $1`
	getByUserQuery    = `SELECT * FROM urls WHERE userid = $1`
	getShort          = `SELECT short FROM urls WHERE original = $1`
	countUserURLs     = "SELECT count(*) FROM urls WHERE userid = $1"
	drop              = `DROP TABLE urls`
)

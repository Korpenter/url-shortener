package storage

const (
	// createUrls creates tables if they don't exist.
	createUrls = `CREATE TABLE IF NOT EXISTS urls (
            	short varchar(255) PRIMARY KEY,
                original varchar(255),
    			userid varchar(64),
    			deleted boolean DEFAULT false,
    			UNIQUE(original)
                )`
	// addQuery inserts a new URL into the 'urls' table, returning existing short ID if it already exists.
	addQuery = `
	INSERT INTO urls (short, original, userid)
	VALUES ($1, $2, $3)
	ON CONFLICT DO NOTHING
	RETURNING short`
	// updateDeleteQuery marks the urls from the list and created by a specific user as deleted.
	updateDeleteQuery = `UPDATE urls SET DELETED=TRUE WHERE short IN (SELECT unnest($1::text[])) AND userid = $2`
	// getQuery retrieves a single URL from the 'urls' table.
	getQuery = `SELECT original, deleted FROM urls WHERE short = $1`
	// getByUserQuery retrieves all URLs belonging to a specific user from the 'urls' table.
	getByUserQuery = `SELECT * FROM urls WHERE userid = $1`
	// getShort retrieves the short URL for a given original URL from the 'urls' table.
	getShort = `SELECT short FROM urls WHERE original = $1`
	// countUserURLs counts the number of URLs belonging to a specific user in the 'urls' table.
	countUserURLs = "SELECT count(*) FROM urls WHERE userid = $1"
	// get count of registered users and urls
	getStats = "SELECT COUNT(*), COUNT(DISTINCT(userid)) FROM urls;"
	// drop drops the 'urls' table.
	drop = `DROP TABLE urls`
)

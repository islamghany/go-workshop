package main

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

/*
to connect to the database we will need a data source name (DNS), which is basically a
string that contains neccessary connection parameters, the exact format of dns will depend
on which database driver we are using, put when using pg the dsn looks like

potgres://username:password@localhost/dbname.

in our case the dsn will be like

postgres://test:islamghany@localhost/test.
*/

const (
	dsn = "postgres://test:islamghany@localhost/test"
)

/*
Establishing a connection pool
*/

// The openDB() function returns a sql.DB connection pool.
func OpenDB() (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN.
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Set the maximum number of open (in-use + idle) connections in the pool. Note that
	// passing a value less than or equal to 0 will mean there is no limit.
	db.SetMaxOpenConns(25)

	// Set the maximum number of idle connections in the pool. Again, passing a value
	// less than or equal to 0 will mean there is no limit.
	db.SetMaxIdleConns(25)

	// Use the time.ParseDuration() function to convert the idle timeout duration string
	// to a time.Duration type.
	duration, err := time.ParseDuration("15m")
	if err != nil {
		return nil, err
	}

	// Set the maximum idle timeout.
	db.SetConnMaxIdleTime(duration)

	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	// Return the sql.DB connection pool.
	return db, nil

}

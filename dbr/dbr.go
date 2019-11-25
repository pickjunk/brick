package dbr

import (
	"os"
	"strconv"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	dbr "github.com/gocraft/dbr"
)

// DB session instance, base on dbr.Session
type DB struct {
	*dbr.Session
}

// Tx transaction instance, base on dbr.Tx
type Tx struct {
	*dbr.Tx
}

// Begin creates a transaction for the given session.
func (db *DB) Begin() (*Tx, error) {
	t, err := db.Session.Begin()
	return &Tx{t}, err
}

// New dbr.Session
// if optionalDSN is omited, config from ENV will be used
// Supported ENV:
// MYSQL_DSN - dsn for connection
// MYSQL_MAX_IDLE - max idle connections for pool
// MYSQL_MAX_OPEN - max open connections for pool
func New(optionalDSN ...string) *DB {
	var dsn string

	if len(optionalDSN) > 0 {
		dsn = optionalDSN[0]
	}

	if dsn == "" {
		dsn = os.Getenv("MYSQL_DSN")
	}

	if dsn == "" {
		log.Panic().Str("field", "MYSQL_DSN").Msg("env required")
	}

	// open connection
	conn, err := dbr.Open("mysql", dsn, log)
	if err != nil {
		log.Panic().Err(err).Send()
	}

	maxIdleConns, _ := strconv.Atoi(os.Getenv("MYSQL_MAX_IDLE"))
	if maxIdleConns == 0 {
		maxIdleConns = 1
	}
	conn.DB.SetMaxIdleConns(maxIdleConns)

	maxOpenConns, _ := strconv.Atoi(os.Getenv("MYSQL_MAX_OPEN"))
	if maxOpenConns == 0 {
		maxOpenConns = 1
	}
	conn.DB.SetMaxOpenConns(maxOpenConns)

	// ping
	err = conn.DB.Ping()
	if err != nil {
		log.Panic().Err(err).Send()
	}

	log.Info().
		Str("dsn", dsn).
		Int("maxIdleConns", maxIdleConns).
		Int("maxOpenConns", maxOpenConns).
		Msg("dbr open")

	return &DB{
		conn.NewSession(nil),
	}
}

// export dbr expression function for convenience
var (
	// And creates AND from a list of conditions.
	And = dbr.And
	// Or creates OR from a list of conditions.
	Or = dbr.Or
	// Eq is `=`.
	// When value is nil, it will be translated to `IS NULL`.
	// When value is a slice, it will be translated to `IN`.
	// Otherwise it will be translated to `=`.
	Eq = dbr.Eq
	// Neq is `!=`.
	// When value is nil, it will be translated to `IS NOT NULL`.
	// When value is a slice, it will be translated to `NOT IN`.
	// Otherwise it will be translated to `!=`.
	Neq = dbr.Neq
	// Gt is `>`.
	Gt = dbr.Gt
	// Gte is '>='.
	Gte = dbr.Gte
	// Lt is '<'.
	Lt = dbr.Lt
	// Lte is `<=`.
	Lte = dbr.Lte
	// Like is `LIKE`, with an optional `ESCAPE` clause
	Like = dbr.Like
	// NotLike is `NOT LIKE`, with an optional `ESCAPE` clause
	NotLike = dbr.NotLike
	// Expr allows raw expression to be used when current SQL syntax is
	// not supported by gocraft/dbr.
	Expr = dbr.Expr
	// Union builds
	Union = dbr.Union
	// UnionAll builds
	UnionAll = dbr.UnionAll
)

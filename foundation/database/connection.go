package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	sqldblogger "github.com/simukti/sqldb-logger"
)

type Config struct {
	Host       string
	Port       int
	Database   string
	Username   string
	Password   string
	DisableTLS bool
	CertPath   string
	Logger     sqldblogger.Logger
}

// URL returns database config in URL presentation
func (c Config) URL() *url.URL {
	sslMode := "verify-full"
	if c.DisableTLS {
		sslMode = "disable"
	}
	q := make(url.Values)
	q.Set("sslmode", sslMode)
	if !c.DisableTLS {
		q.Set("sslrootcert", c.CertPath)
	}
	q.Set("timezone", "utc")

	return &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.Username, c.Password),
		Host:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Path:     c.Username,
		RawQuery: q.Encode(),
	}
}

// NewConnection creates a new database connection
func NewConnection(conf Config) (*sqlx.DB, error) {
	dsn := conf.URL().String()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if conf.Logger != nil {
		db = sqldblogger.OpenDriver(
			dsn,
			db.Driver(),
			conf.Logger,
			sqldblogger.WithExecerLevel(sqldblogger.LevelDebug),
			sqldblogger.WithQueryerLevel(sqldblogger.LevelDebug),
			sqldblogger.WithPreparerLevel(sqldblogger.LevelDebug),
		)
	}
	sqx := sqlx.NewDb(db, "postgres")
	if err := StatusCheck(context.Background(), sqx); err != nil {
		return nil, err
	}
	return sqx, nil
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *sqlx.DB) error {
	// Run a simple query to determine connectivity. The db has a "Ping" method
	// but it can false-positive when it was previously able to talk to the
	// database but the database has since gone away. Running this query forces a
	// round trip to the database.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}

package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	// _ kullanmamızın nedeni bu dosyada bu paketlerin kullanılmamasından dolayı error göstermemesidir.
)

// DB holds the database connection pool.
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{} // initializes empty DB, burada dbConn pointer (& kullandık çünkü aynı DB objemizi işaret etmesini istiyoruz)
// Normalde tip kullanarak pointer oluştururken, var dbConn *DB (ilk değerini vermedik)
// Ama ilk değerini vererek bir pointer oluştururken, var dbConn = &DB{}

// some parameters about db connection pool
const maxOpenDBConn = 10
const maxIdleDBConn = 5
const maxDBLifetime = 5 * time.Minute

// Creates a new database for the application
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	// pings the databse
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil // everything is fine
}

// Creates database pool for postgres, set some parameters
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}
	// sets some parameters about db connection pool
	d.SetMaxOpenConns(maxOpenDBConn)
	d.SetMaxIdleConns(maxIdleDBConn)
	d.SetConnMaxLifetime(maxDBLifetime)

	// (*dbConn).SQL şeklinde açıkça yazabiliriz çünkü dereference ediyoruz, ama structlarda GO bunu anlar.
	dbConn.SQL = d  // parametreleri ayarladıktan sonra package level variable olan dbConn'a erişip SQL field değerini atıyoruz.
	err = testDB(d) // pings the database
	if err != nil {
		return nil, err
	}
	return dbConn, nil // everything is fine
}

// Tries to ping the database
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

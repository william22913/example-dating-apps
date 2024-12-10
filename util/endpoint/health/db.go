package health

import "database/sql"

type dbConn struct {
	db *sql.DB
}

func (d *dbConn) Ping() string {
	if err := d.db.Ping(); err != nil {
		return "DOWN"
	}

	return "UP"
}

func NewDBConnChecker(
	db *sql.DB,
) *dbConn {
	return &dbConn{
		db: db,
	}
}

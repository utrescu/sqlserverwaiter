package mssql

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/utrescu/sqlserverwait/cmd"
	database "github.com/utrescu/sqlserverwait/db"
)

// Connection Defines a connection with Sql Server
type Connection struct {
	connection *sql.DB
}

// New creates a connection with Sql Server
func New(connectionString string) (database.RepositoryReady, error) {
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return nil, err
	}

	sqlconnection := &Connection{
		connection: db,
	}

	return sqlconnection, nil
}

// IsAlive checks if Sql Server is up and accepts connections
func (m *Connection) IsAlive() error {

	if err := m.connection.Ping(); err != nil {
		if err.Error() == "EOF" {
			return errors.New("database not ready")
		}
		return err
	}

	// When Collation is not null database is ready
	row := m.connection.QueryRow(fmt.Sprintf("SELECT DATABASEPROPERTYEX('%s', 'Collation') AS Collation", cmd.Database))
	var collation string
	err := row.Scan(&collation)
	if err != nil {
		return errors.New("database not ready")
	}

	// Database ready!!
	return nil
}

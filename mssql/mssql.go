package mssql

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb" // mssql driver
	ready "github.com/utrescu/sqlserverwaiter/ready"
)

// Connection Defines a connection with Sql Server
type Connection struct {
	connection *sql.DB
	Name       string
}

// New creates a connection with Sql Server
func New(connectionString string, name string) (ready.RepositoryReady, error) {
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return nil, err
	}

	sqlconnection := &Connection{
		connection: db,
		Name:       name,
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
	row := m.connection.QueryRow(fmt.Sprintf("SELECT DATABASEPROPERTYEX('%s', 'Collation') AS Collation", m.Name))
	var collation string
	err := row.Scan(&collation)
	if err != nil {
		return errors.New("database not ready")
	}

	// Database ready!!
	return nil
}

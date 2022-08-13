package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/go-sql-driver/mysql"
)

var (
	ErrInvalidUser  = errors.New("DBInfo struct has invalid user")
	ErrInvalidPass  = errors.New("DBInfo struct has invalid password")
	ErrInvalidName  = errors.New("DBInfo struct has invalid database name")
	ErrInvalidTable = errors.New("DBInfo struct has invalid table")
	ErrInvalidHost  = errors.New("DBInfo struct has invalid host")
	ErrInvalidPort  = errors.New("DBInfo struct has invalid port")
)

// All required database information.
type DBInfo struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"database"`
	Table    string `json:"table"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

func (db *DBInfo) Validate(ctx context.Context) error {
	switch {
	case db.User == "":
		return ErrInvalidUser
	case db.Password == "":
		return ErrInvalidPass
	case db.Name == "":
		return ErrInvalidName
	case db.Table == "":
		return ErrInvalidTable
	case db.Host == "":
		return ErrInvalidHost
	case db.Port == "":
		return ErrInvalidPort
	}

	return nil
}

// Print DBInfo object information.
func (db *DBInfo) Print() {
	fmt.Printf("%s: %s\n", "User", db.User)
	fmt.Printf("%s: %s\n", "Password", db.Password)
	fmt.Printf("%s: %s\n", "Database", db.Name)
	fmt.Printf("%s: %s\n", "Table", db.Table)
	fmt.Printf("%s: %s\n", "Host", db.Host)
	fmt.Printf("%s: %s\n", "Port", db.Port)
}

// Return DB pointer.
func DB(ctx context.Context, cfgFile string) (*sql.DB, string) {
	log := logr.FromContextOrDiscard(ctx)

	// Get config filename from context.
	config, err := LoadConfig(ctx, cfgFile)
	if err != nil {
		return nil, ""
	}

	// Capture connection properties.
	cfg := mysql.Config{
		User:   config.User,
		Passwd: config.Password,
		Net:    "tcp",
		Addr:   config.Host + ":" + config.Port,
		DBName: config.Name,
	}

	// Get a database handle.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Error(err, "Could not open database.")

		return nil, ""
	}

	if err := db.Ping(); err != nil {
		log.Error(err, "Could not ping database.")

		return nil, ""
	}

	log.Info("Connected to database.")

	return db, config.Table
}

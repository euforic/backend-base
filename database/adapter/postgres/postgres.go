package postgres

import (
	"context"
	"time"

	"github.com/euforic/backend-base/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const todoSchema = "todos"

var _ database.Adapter = &Database{} // or &myType{} or [&]myType if scalar

type Database struct {
	db      *gorm.DB
	connStr string
}

func New(connStr string) *Database {
	d := Database{
		connStr: connStr,
	}
	return &d
}

// Connect opens a connection to to buntdb
func (d *Database) Connect() error {
	db, err := gorm.Open(postgres.Open(d.connStr), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: false,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "todos.", // schema name
			SingularTable: false,
		},
	})
	if err != nil {
		return err
	}

	sqldb, err := db.DB()
	if err != nil {
		return err
	}

	sqldb.SetConnMaxLifetime(time.Minute)
	if err != nil {
		return err
	}

	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS " + todoSchema).Error; err != nil {
		return err
	}

	if err := db.AutoMigrate(&database.Todo{}); err != nil {
		panic(err)
	}

	d.db = db
	return nil
}

// session ...
func (d Database) session(ctx context.Context) *gorm.DB {
	return d.db.Session(&gorm.Session{FullSaveAssociations: false, Context: ctx})
}

// Close closes out BuntDB connection
func (d Database) Close() error {
	//Deprecated in gorm because of connection pooling
	return nil
}

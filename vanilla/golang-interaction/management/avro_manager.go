package management

import (
	"database/sql"
	"fmt"
	"kafkaing/utils"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

type AvroManager interface {
	CreateSchema(name, version, schema string) bool
	UpdateSchema(name, version, schema string) bool
	DeleteSchema(name, version string) bool
	GetSchema(name, version string) string
}

type SqliteAvroManager struct {
	conn *sql.DB
	sema semaphore.Weighted
}

func Init(driverName, dbName string) *SqliteAvroManager {
	//sql.Open("sqlite3_custom", ":memory:")
	db, err := sql.Open("sqlite3", "./avro.db")
	if err != nil {
		zap.L().Sugar().Fatal(err)
	}

	_, err = db.Exec(TableCreationQuery)
	if err != nil {
		zap.L().Sugar().Error("%q: %s\n", err, TableCreationQuery)
	}

	go utils.RunCleanup(utils.PrepareSigtermChannel(), func() {
		zap.L().Sugar().Info("Shutting down SqliteAvroManager")
		db.Close()
	})

	return &SqliteAvroManager{conn: db, sema: *semaphore.NewWeighted(5)}
}

func (avm *SqliteAvroManager) CreateSchema(name, version, schema string) bool {
	_, err := avm.conn.Exec(SchemaInsertQuery, name, schema, version)
	if err != nil {
		zap.L().Sugar().Error(err)
		return false
	}

	return true
}

func (avm *SqliteAvroManager) UpdateSchema(name, version, schema string) bool {
	return false
}

func (avm *SqliteAvroManager) DeleteSchema(name, version string) bool {
	return false
}

func (avm *SqliteAvroManager) GetSchema(name, version string) string {

	return ""
}

func TestConnection() {
	db, err := sql.Open("sqlite3", "./avro.db")
	if err != nil {
		zap.L().Sugar().Fatal(err)
	}
	_, err = db.Exec(TableCreationQuery)
	if err != nil {
		zap.L().Sugar().Error("%q: %s\n", err, TableCreationQuery)
		return
	}

	defer db.Close()
}

func RemoveDbFile(dbFile string) error {
	if !strings.HasSuffix(dbFile, ".db") {
		return fmt.Errorf("File not of type db")
	}
	return os.Remove(dbFile)
}

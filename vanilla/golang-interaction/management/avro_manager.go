package management

import (
	"context"
	"database/sql"
	"fmt"
	"kafkaing/utils"
	"os"
	"strings"

	"github.com/hamba/avro"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

type AvroManager interface {
	CreateSchema(name, version, schema string) bool
	UpdateSchema(name, version, schema string) bool
	DeleteSchema(name, version string) bool
	UpsertSchema(name, version, schema string) bool
	GetSchema(name, version string) avro.Schema
}

type SqliteAvroManager struct {
	conn        *sql.DB
	sema        semaphore.Weighted
	context     context.Context
	schemaCache avro.SchemaCache
}

func InitAvroManager(driverName, dbName string) *SqliteAvroManager {
	//sql.Open("sqlite3_custom", ":memory:")
	//sql.Open("sqlite3", "./avro.db")
	db, err := sql.Open(driverName, dbName)
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

	return &SqliteAvroManager{conn: db, sema: *semaphore.NewWeighted(5), context: context.Background(), schemaCache: avro.SchemaCache{}}
}

func (avm *SqliteAvroManager) CreateSchema(name, version, schema string) bool {
	avm.sema.Acquire(avm.context, 1)
	defer avm.sema.Release(1)

	{
		avroSchema := avm.decodeSchema(schema)
		if avroSchema == nil {
			return false
		}
		avm.updateSchemaCache(name, version, avroSchema)
	}

	{
		res, err := avm.conn.Exec(SchemaInsertQuery, name, version, schema)
		if err != nil {
			zap.L().Sugar().Error(err)
			return false
		}

		rowId, _ := res.LastInsertId()
		zap.L().Sugar().Debugf("Scehma created with identifier (%s,%s) and rowId - %s", name, version, rowId)
	}

	return true
}

func (avm *SqliteAvroManager) UpdateSchema(name, version, schema string) bool {
	avm.sema.Acquire(avm.context, 1)
	defer avm.sema.Release(1)

	{
		avroSchema := avm.decodeSchema(schema)
		if avroSchema == nil {
			return false
		}
		avm.updateSchemaCache(name, version, avroSchema)
	}

	_, err := avm.conn.Exec(SchemaUpdateQuery, schema, name, version)

	if err != nil {
		zap.L().Sugar().Error(err)
		return false
	}

	return false
}

func (avm *SqliteAvroManager) UpsertSchema(name, version, schema string) bool {
	avm.sema.Acquire(avm.context, 1)
	defer avm.sema.Release(1)
	{
		_, err := avm.conn.Exec(SchemaUpsertQuery, name, version, schema)
		if err != nil {
			zap.L().Sugar().Error(err)
			return false
		}
	}

	{
		avroSchema := avm.decodeSchema(schema)
		if avroSchema == nil {
			return false
		}
		avm.updateSchemaCache(name, version, avroSchema)
	}

	return true
}

func (avm *SqliteAvroManager) DeleteSchema(name, version string) bool {
	avm.sema.Acquire(avm.context, 1)
	defer avm.sema.Release(1)

	avm.removeSchemaFromCache(name, version)
	_, err := avm.conn.Exec(SchemaDeleteQuery, name, version)

	if err != nil {
		zap.L().Sugar().Error(err)
		return false
	}

	return true
}

func (avm *SqliteAvroManager) GetSchema(name, version string) avro.Schema {
	avm.sema.Acquire(avm.context, 1)
	defer avm.sema.Release(1)

	{
		res := avm.checkSchemaCache(name, version)
		if res != nil {
			return res
		}
	}

	resRow := avm.conn.QueryRow(SchemaFetchQuery, name, version)
	var schemaStringified string
	resRow.Scan(schemaStringified)

	return avm.decodeSchema(schemaStringified)
}

func (avm *SqliteAvroManager) decodeSchema(schema string) avro.Schema {
	res, err := avro.Parse(schema)
	if err != nil {
		zap.L().Sugar().Error(err)
	}
	return res
}

func (avm *SqliteAvroManager) updateSchemaCache(name, version string, schema avro.Schema) {
	avm.schemaCache.Add(fmt.Sprintf("%s::%s", name, version), schema)
}

func (avm *SqliteAvroManager) checkSchemaCache(name, version string) avro.Schema {
	return avm.schemaCache.Get(fmt.Sprintf("%s::%s", name, version))
}

func (avm *SqliteAvroManager) removeSchemaFromCache(name, version string) {
	avm.schemaCache.Add(fmt.Sprintf("%s::%s", name, version), nil)
}

func RemoveDbFile(dbFile string) error {
	if !strings.HasSuffix(dbFile, ".db") {
		return fmt.Errorf("file not of type db")
	}
	return os.Remove(dbFile)
}

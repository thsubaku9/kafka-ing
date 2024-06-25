package management

const TableCreationQuery string = `
CREATE TABLE IF NOT EXISTS avro_store(id INTEGER PRIMARY KEY AUTOINCREMENT, 
		name TEXT NOT NULL, 
		schema BLOB NOT NULL, 
		version int NOT NULL,
		created_at DATETIME DEFAULT (now()),
	CONSTRAINT unique_schema_version UNIQUE (name, version)
	);
CREATE UNIQUE INDEX IF NOT EXISTS name_lookup_index ON avro_store (name);
`

const SchemaBulkFetchQuery string = `SELECT (schema,version) FROM avro_store WHERE name = ? ORDERBY(version) DESC`
const SchemaFetchQuery string = `SELECT schema FROM avro_store WHERE name = ? AND version = ?`
const SchemaDeleteQuery string = `DELETE * FROM avro_store WHERE name = ? AND version = ?`
const SchemaBulkDeleteQuery string = `DELETE * FROM avro_store WHERE name = ?`
const SchemaInsertQuery string = `INSERT INTO avro_store(name, version, schema) VALUES(?,?,?)`
const SchemaUpdateQuery string = `UPDATE avro_store SET schema = ? WHERE name = ?`

package management

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

const tableCreationQuery string = `
CREATE TABLE IF NOT EXISTS avro_store(id INTEGER NOT NULL PRIMARY KEY, name TEXT UNIQUE, schema BLOB);
CREATE UNIQUE INDEX IF NOT EXISTS name_lookup_index ON avro_store (name);
`

func TestConnection() {
	for _, driver := range sql.Drivers() {
		println(driver)
	}

	db, err := sql.Open("sqlite3", "./avro.db")
	if err != nil {
		zap.L().Sugar().Fatal(err)
	}
	_, err = db.Exec(tableCreationQuery)
	if err != nil {
		log.Printf("%q: %s\n", err, tableCreationQuery)
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

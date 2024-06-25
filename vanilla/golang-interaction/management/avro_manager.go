package management

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func TestConnection() {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		zap.L().Sugar().Fatal(err)
	}
	defer db.Close()
}

func RemoveDbFile(dbFile string) error {
	if !strings.HasSuffix(dbFile, ".db") {
		return fmt.Errorf("File not of type db")
	}
	return os.Remove(dbFile)
}

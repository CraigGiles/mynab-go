package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Migration struct {
	version    string
	script     string
	file_name  string
	file_hash  string
	created_at time.Time
}

func hash_file_md5(file_path string) string {
	// NOTE: since we've already read the contents of the file before
	//   we assume we can open it. Therefore no error needed
	file, _ := os.Open(file_path)
	defer file.Close()

	hash := md5.New()
	io.Copy(hash, file)

	// get the 16 bytes hash
	hash_in_bytes := hash.Sum(nil)[:16]
	result := hex.EncodeToString(hash_in_bytes)

	return result
}

func create_database_connection(username string, password string, database string) *sql.DB {
	connection_string := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", username, password, database)

	var db, err = sql.Open("postgres", connection_string)

	if err != nil {
		log.Fatal("Error opening connection to database: %s", err)
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal("Error with database connection: %s", err)
		panic(err)
	}

	return db
}

func migration_files() []string {
	var result []string

	root := "./sql/migrations/" // TODO do i want this to be configurable?
	// Note(Kyle): Directory walked in lexical order.
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, "sql") {
			result = append(result, path)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	sort.Strings(result)

	return result
}

func parse_version_and_script(path string) (string, string) {
	var result_version string = ""
	var result_script string = ""

	index := strings.LastIndex(path, "/")

	if index == -1 {
		index = 0
	}

	filename := path[index+1:]

	results := strings.Split(filename, "--")
	if len(results) > 0 {
		result_version = results[0]
		result_script = results[1]
	}

	return result_version, result_script
}

func migration_exists(migration Migration, version string, script string, hash string) bool {
	result := false

	if migration.version == version && migration.script == script {
		if migration.file_hash != hash {
			log.Printf("Migration %s--%s has been run with a different file hash\n", version, script)
			result = true
		} else {
			log.Printf("%s--%s exists. Skipping...\n", version, script)
			result = true
		}
	}

	return result
}

func main() {
	db := create_database_connection("postgres", "root", "mynab")

	for _, path := range migration_files() {
		var migration = Migration{}
		var bytes, io_error = ioutil.ReadFile(path)

		if io_error != nil {
			log.Printf("Unable to open %s\n%s", path, io_error.Error())
			continue
		}

		var version, script = parse_version_and_script(path)
		var hash = hash_file_md5(path)

		var row = db.QueryRow("SELECT * FROM table_migrations WHERE version=$1 AND script=$2",
			version, script)
		_ = row.Scan(&migration.version, &migration.script,
			&migration.file_name, &migration.file_hash, &migration.created_at)

		if !migration_exists(migration, version, script, hash) {
			var sql = string(bytes)

			var tx, _ = db.Begin()
			var _, db_error = db.Exec(sql)

			if db_error != nil {
				tx.Rollback()
				log.Printf("Migration Error: %s\n", db_error.Error())
			} else {
				migration_insert_sql := `
                    INSERT INTO table_migrations (version, script, file_name, file_hash)
                    VALUES ($1, $2, $3, $4)`
				var _, migration_insert_error = db.Exec(migration_insert_sql, version, script, path, hash)

				if migration_insert_error != nil {
					log.Printf("Error inserting the migration log. Rolling back transaction for %s--%s: %s\n", version, script, migration_insert_error.Error())
					tx.Rollback()
				} else {
					tx.Commit()
				}
			}

		}
	}
}

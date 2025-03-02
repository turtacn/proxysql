// main.go
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/auth"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/information_schema"
)

func main() {
	engine := sqle.NewDefault(
		sql.NewDatabaseProvider(
			createTpccDatabase(),
			information_schema.NewInformationSchemaDatabase(),
		))

	config := server.Config{
		Protocol: "tcp",
		Address:  "localhost:3306",
		Auth:     auth.NewNativeSingle("root", "123", auth.AllPermissions),
	}

	s, err := server.NewDefaultServer(config, engine)
	if err != nil {
		panic(err)
	}

	fmt.Println("MySQL server listening on localhost:3306")
	if err := s.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func createTpccDatabase() *memory.Database {
	dbName := "tpcc"
	db := memory.NewDatabase(dbName)

	sqlFiles := []string{"create_table.sql", "add_fkey_idx.sql"} // Assuming these files are in same directory
	ctx := sql.NewEmptyContext()
	e := sqle.NewDefault(sql.NewDatabaseProvider(db, information_schema.NewInformationSchemaDatabase()))

	for _, file := range sqlFiles {
		sqlContent, err := os.ReadFile(file)
		if err != nil {
			panic(fmt.Sprintf("Error reading SQL file %s: %v. Please make sure %s and add_fkey_idx.sql are in the same directory as main.go", file, err, "create_table.sql"))
		}
		queries := strings.Split(string(sqlContent), ";")

		for _, query := range queries {
			query = strings.TrimSpace(query)
			if query == "" {
				continue
			}
			_, _, err = e.Query(ctx, fmt.Sprintf("USE %s;", dbName))
			if err != nil {
				panic(fmt.Sprintf("Error using database %s: %v", dbName, err))
			}

			_, _, err = e.Query(ctx, query+";")
			if err != nil {
				fmt.Printf("Error executing SQL: %s\n", query) // Print the failing query
				panic(fmt.Sprintf("Error executing SQL from %s: %v", file, err))
			}
		}
		fmt.Printf("Executed SQL from %s\n", file)
	}
	fmt.Println("TPCC database and tables created.")
	return db
}

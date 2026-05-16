package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "usage: dump_schema <dsn> <output.sql>\n")
		os.Exit(1)
	}
	dsn := os.Args[1]
	outPath := os.Args[2]

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "ping: %v\n", err)
		os.Exit(1)
	}

	var dbName string
	if err := db.QueryRow("SELECT DATABASE()").Scan(&dbName); err != nil {
		fmt.Fprintf(os.Stderr, "database: %v\n", err)
		os.Exit(1)
	}

	tables, err := getTables(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tables: %v\n", err)
		os.Exit(1)
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("-- Database: %s\n", dbName))
	b.WriteString(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;\n", dbName))
	b.WriteString(fmt.Sprintf("USE `%s`;\n\n", dbName))

	for _, table := range tables {
		createSQL := getCreateTableSQL(db, table)
		if createSQL == "" {
			continue
		}
		createSQL = ensureCreateTableIfNotExists(createSQL)
		b.WriteString(createSQL)
		b.WriteString(";\n\n")
	}

	if err := os.WriteFile(outPath, []byte(b.String()), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("exported %d tables to %s\n", len(tables), outPath)
}

func getTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tables = append(tables, name)
	}
	return tables, rows.Err()
}

func getCreateTableSQL(db *sql.DB, table string) string {
	rows, err := db.Query("SHOW CREATE TABLE `" + table + "`")
	if err != nil {
		fmt.Fprintf(os.Stderr, "show create %s: %v\n", table, err)
		return ""
	}
	defer rows.Close()
	var tableName, createSQL string
	for rows.Next() {
		if err := rows.Scan(&tableName, &createSQL); err != nil {
			fmt.Fprintf(os.Stderr, "scan %s: %v\n", table, err)
			return ""
		}
	}
	return createSQL
}

// ensureCreateTableIfNotExists makes re-running the script safe when tables already exist.
func ensureCreateTableIfNotExists(createSQL string) string {
	const prefix = "CREATE TABLE "
	if strings.HasPrefix(createSQL, prefix) && !strings.HasPrefix(createSQL, "CREATE TABLE IF NOT EXISTS ") {
		return "CREATE TABLE IF NOT EXISTS " + strings.TrimPrefix(createSQL, prefix)
	}
	return createSQL
}

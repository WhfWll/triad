package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"os"
)

// 数据库 字段
type Column struct {
	Field   string      `json:"field"`
	Type    string      `json:"type"`
	Null    string      `json:"null"`
	Key     string      `json:"key"`
	Default interface{} `json:"default"`
	Extra   string      `json:"extra"`
}

type Index struct {
	Table        string      `db:"Table"`
	NonUnique    int         `db:"Non_unique"`
	KeyName      string      `db:"Key_name"`
	SeqInIndex   int         `db:"Seq_in_index"`
	ColumnName   string      `db:"Column_name"`
	Collation    string      `db:"Collation"`
	Cardinality  int         `db:"Cardinality"`
	SubPart      *int        `db:"Sub_part"`
	Packed       *string     `db:"Packed"`
	Null         string      `db:"Null"`
	IndexType    string      `db:"Index_type"`
	Comment      string      `db:"Comment"`
	IndexComment string      `db:"Index_comment"`
	Visible      string      `db:"Visible"`
	Expression   interface{} `db:"Expression"`
}

// 数据库 表 结构
type Table struct {
	Name    string    `json:"name"`
	Create  string    `json:"create"`
	Columns []*Column `json:"columns"`
}

// 数据库 结构
type DbStruct struct {
	Name   string  `json:"name"`
	Tables []Table `json:"tables"`
}

func GetNewestDbStruct(sourceDecisionDBUrl, sourceSmartDBUrl string) {
	// 连接数据库
	var err error
	decisionDB, err := sql.Open("mysql", sourceDecisionDBUrl)
	if err != nil {
		logrus.WithError(err).Fatal("the url of sourceDB is wrong")
	}
	err = decisionDB.Ping()
	if err != nil {
		logrus.WithError(err).Fatal("cannot connect to sourceDB")
	}
	decisionTables, err := getTables(decisionDB)
	if err != nil {
		logrus.WithError(err).Fatalf("fail to get all tables from %s", sourceDecisionDBUrl)
	}
	dbStructList := make([]DbStruct, 0)
	decisionDbStruct := DbStruct{Name: "decision"}
	for _, stb := range decisionTables {
		columns, err := getColumns(decisionDB, stb)
		if err != nil {
			fmt.Println(err)
			continue
		}
		decisionDbStruct.Tables = append(decisionDbStruct.Tables, Table{
			Name:    stb,
			Columns: columns,
			Create:  getCreateTableSql(decisionDB, stb),
		})
	}
	dbStructList = append(dbStructList, decisionDbStruct)

	// 连接 smart 数据库
	smartDB, err := sql.Open("mysql", sourceSmartDBUrl)
	if err != nil {
		logrus.WithError(err).Fatal("the url of sourceDB is wrong")
	}
	err = smartDB.Ping()
	if err != nil {
		logrus.WithError(err).Fatal("cannot connect to sourceDB")
	}
	smartTables, err := getTables(smartDB)
	smartDbStruct := DbStruct{Name: "smart"}
	for _, stb := range smartTables {
		columns, err := getColumns(smartDB, stb)
		if err != nil {
			fmt.Println(err)
			continue
		}
		smartDbStruct.Tables = append(smartDbStruct.Tables, Table{
			Name:    stb,
			Columns: columns,
			Create:  getCreateTableSql(smartDB, stb),
		})
	}
	dbStructList = append(dbStructList, smartDbStruct)

	dbStructByte, err := json.Marshal(dbStructList)
	// 定义要写入的文件名
	filename := "dbStruct.json"
	// 打开或创建文件用于写入
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// 将JSON数据写入文件
	_, err = file.Write(dbStructByte)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
	}
	fmt.Println("JSON data written to file successfully:", filename)
}

// 还原数据库结构
func RestoreDbStruct(ctx context.Context, filename string, targetDecisionDBUrl, targetSmartDBUrl string) error {
	// 连接数据库
	var err error
	jsonBytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	var dbStructList []DbStruct
	err = json.Unmarshal(jsonBytes, &dbStructList)
	if err != nil {
		logrus.Fatalf("Error unmarshalling JSON: %v", err)
	}
	for _, dbStruct := range dbStructList {
		fmt.Println("check db name:", dbStruct.Name, "===========================================")
		if dbStruct.Name == "decision" {
			printUpdateSql(ctx, targetDecisionDBUrl, dbStruct.Tables)
		}
		if dbStruct.Name == "smart" {
			printUpdateSql(ctx, targetSmartDBUrl, dbStruct.Tables)
		}
	}
	return nil
}

func printUpdateSql(ctx context.Context, dbUrl string, sourceTables []Table) {
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		logrus.WithError(err).Fatal("the url of sourceDB is wrong")
	}
	err = db.Ping()
	if err != nil {
		logrus.WithError(err).Fatal("cannot connect to sourceDB")
	}
	targetTables, err := getTables(db)
	tableNames := make([]string, 0)
	for _, src := range sourceTables {
		tableNames = append(tableNames, src.Name)
	}
	printDbDiff(tableNames, targetTables, sourceTables)

	fmt.Println("table need to modify:", ".....................")
	for _, ttb := range targetTables {
		for _, stb := range sourceTables {
			if ttb == stb.Name {
				printSqlDiff(db, ttb, stb.Columns)
			}
		}
	}
}

func printDbDiff(tableNames []string, targetTables []string, sourceTables []Table) {
	var dbDropDiff []string
	var dbAddDiff []string
	for _, table := range targetTables {
		if !tableExists(tableNames, table) {
			dbDropDiff = append(dbDropDiff, fmt.Sprintf("DROP TABLE %s", table))
		}
	}
	if len(dbDropDiff) != 0 {
		fmt.Println("table need to drop: ...............")
		for _, sqlStr := range dbDropDiff {
			fmt.Println(sqlStr + ";")
		}
		fmt.Println()
	}
	createTableSql := getCreateTables(targetTables, sourceTables)
	dbAddDiff = append(dbAddDiff, createTableSql...)
	if len(dbAddDiff) != 0 {
		fmt.Println("table need to add: .......................")
		for _, sqlStr := range dbAddDiff {
			fmt.Println(sqlStr + ";")
		}
		fmt.Println()
	}
}

func printSqlDiff(targetDB *sql.DB, tableName string, sourceColumns []*Column) {
	// 获取目标表的列信息
	targetColumns, err := getColumns(targetDB, tableName)
	if err != nil {
		logrus.Fatal(err)
	}

	//sComments, err := getColumnsComments(sourceDB, tableName)
	//if err != nil {
	//	logrus.WithError(err).Fatal("fail to get sComments")
	//}

	//tComments, err := getColumnsComments(targetDB, tableName)
	//if err != nil {
	//	logrus.WithError(err).Fatal("fail to get sComments")
	//}

	//tIndexes, err := getIndexes(targetDB, tableName)
	if err != nil {
		logrus.WithError(err).Fatal("fail to get targetIndexes")
	}
	//sIndexes, err := getIndexes(sourceDB, tableName)
	//if err != nil {
	//	logrus.WithError(err).Fatal("fail to get sourceIndexes")
	//}

	// 比较源表和目标表的列差异
	columnDiff := make([]string, 0, 100)
	indexDiff := make([]string, 0, 10)
	sComments := make(map[string]string, 0)
	tComments := make(map[string]string, 0)
	addColumn(sourceColumns, targetColumns, sComments, &columnDiff)
	dropColumn(targetColumns, sourceColumns, &columnDiff)
	modifyColumn(targetColumns, sourceColumns, sComments, tComments, &columnDiff)

	//比较索引差异
	//addIndex(sIndexes, tIndexes, &indexDiff)
	//dropIndex(sIndexes, tIndexes, &indexDiff)

	// 生成表变更的 DDL 语句
	printResult(columnDiff, indexDiff, tableName)
}

func printResult(columnDiff []string, indexDiff []string, tableName string) {
	if len(columnDiff) > 0 {
		ddlStatement := fmt.Sprintf("ALTER TABLE %s\n%s", tableName, joinStrings(columnDiff, ",\n"))
		fmt.Printf("%s", ddlStatement)
		fmt.Println()
		fmt.Println()
	}

	if len(indexDiff) > 0 {
		fmt.Printf("#Index update statement for table %s:\n", tableName)
		for _, idxStatement := range indexDiff {
			fmt.Println(idxStatement)
		}
		fmt.Println()
		fmt.Println()
	}
}

func addIndex(sIndexes []Index, tIndexes []Index, indexDiff *[]string) {
	for _, sIdx := range sIndexes {
		if !indexExists(&tIndexes, sIdx.KeyName) {
			uniqueStr := " "
			if sIdx.NonUnique == 0 {
				uniqueStr = " UNIQUE "
			}
			sqlStr := fmt.Sprintf("CREATE%sINDEX %s ON %s(%s);", uniqueStr, sIdx.KeyName, sIdx.Table, sIdx.ColumnName)
			*indexDiff = append(*indexDiff, sqlStr)
		}
	}
}

func dropIndex(sIndexes []Index, tIndexes []Index, indexDiff *[]string) {
	for _, tIdx := range tIndexes {
		if !indexExists(&sIndexes, tIdx.KeyName) {
			sqlStr := fmt.Sprintf("DROP INDEX %s ON %s;", tIdx.KeyName, tIdx.Table)
			*indexDiff = append(*indexDiff, sqlStr)
		}
	}
}

func modifyColumn(targetColumns []*Column, sourceColumns []*Column, sComments map[string]string, tComments map[string]string, columnDiff *[]string) {
	for _, tclm := range targetColumns {
		for _, sclm := range sourceColumns {
			if tclm.Field == sclm.Field {
				if tclm.Type != sclm.Type || tclm.Null != sclm.Null || tclm.Key != sclm.Key || !defaultEqual(tclm.Default, sclm.Default) || sclm.Extra != tclm.Extra || !commentEqual(sComments, tComments, sclm.Field, tclm.Field) {
					sqlStr := fmt.Sprintf("MODIFY COLUMN %s %s", sclm.Field, sclm.Type)
					if sclm.Null == "NO" {
						sqlStr += " NOT NULL"
					}
					if len(sclm.Key) != 0 {
						switch sclm.Key {
						case "PRI":
							sqlStr += " PRIMARY KEY"
						}
					}
					if sclm.Default != nil {
						bts := sclm.Default.(string)
						if len(bts) > 0 {
							sqlStr += " DEFAULT " + string(bts)
						}
					}
					if len(sclm.Extra) != 0 && sclm.Extra == "AUTO_INCREMENT" {
						sqlStr += " " + sclm.Extra
					}
					if cmt, ok := sComments[sclm.Field]; ok {
						sqlStr += " COMMENT " + "'" + cmt + "'"
					}
					*columnDiff = append(*columnDiff, sqlStr)
				}

			}
		}
	}

}

func dropColumn(targetColumns []*Column, sourceColumns []*Column, columnDiff *[]string) {
	for _, targetColumn := range targetColumns {
		columnName := targetColumn.Field

		if !columnExists(sourceColumns, columnName) {
			*columnDiff = append(*columnDiff, fmt.Sprintf("DROP COLUMN %s", columnName))
		}
	}
}

func addColumn(sourceColumns []*Column, targetColumns []*Column, sComments map[string]string, columnDiff *[]string) {
	for _, clm := range sourceColumns {
		columnName := clm.Field
		columnType := clm.Type

		if !columnExists(targetColumns, columnName) {
			addColumn := fmt.Sprintf("ADD COLUMN %s %s", columnName, columnType)
			if clm.Null == "NO" {
				addColumn += " NOT NULL"
			}
			if len(clm.Key) != 0 {
				switch clm.Key {
				case "PRI":
					addColumn += " PRIMARY KEY"
				}
			}
			if clm.Default != nil {
				bts := clm.Default.(string)
				if len(bts) > 0 {
					addColumn += " DEFAULT " + string(bts)
				}
			}
			if len(clm.Extra) != 0 && clm.Extra == "AUTO_INCREMENT" {
				addColumn += " " + clm.Extra
			}
			if cmt, ok := sComments[columnName]; ok {
				addColumn += " COMMENT " + "'" + cmt + "'"
			}
			*columnDiff = append(*columnDiff, addColumn)
		}
	}
}

func getColumnsComments(db *sql.DB, tableName string) (ans map[string]string, err error) {
	row, err := db.Query("select database()")
	if err != nil {
		return
	}
	var schemaName string
	for row.Next() {
		row.Scan(&schemaName)
	}
	rows, err := db.Query("SELECT COLUMN_NAME, COLUMN_COMMENT FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = ? AND TABLE_SCHEMA= ?", tableName, schemaName)
	if err != nil {
		return nil, err
	}
	ans = make(map[string]string)
	defer rows.Close()
	for rows.Next() {
		var field string
		var comment string
		err = rows.Scan(&field, &comment)
		if err != nil {
			logrus.WithError(err).Errorln("fail in scan result from mysql")
			continue
		}
		if len(comment) != 0 {
			ans[field] = comment
		}
	}
	return
}

// 获取表的列信息
func getColumns(db *sql.DB, tableName string) (columns []*Column, err error) {
	rows, err := db.Query("SHOW COLUMNS FROM " + tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		column := new(Column)
		err := rows.Scan(
			&column.Field,
			&column.Type,
			&column.Null,
			&column.Key,
			&column.Default,
			&column.Extra,
		)
		if err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}

	return columns, nil
}

func getIndexes(db *sql.DB, tbName string) (ans []Index, err error) {
	rows, err := db.Query("SHOW INDEX FROM " + tbName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var index Index
		err := rows.Scan(
			&index.Table,
			&index.NonUnique,
			&index.KeyName,
			&index.SeqInIndex,
			&index.ColumnName,
			&index.Collation,
			&index.Cardinality,
			&index.SubPart,
			&index.Packed,
			&index.Null,
			&index.IndexType,
			&index.Comment,
			&index.IndexComment,
		)
		if err != nil {
			logrus.WithError(err).Errorln("fail to get index from db")
			continue
		}
		flag := false
		for i, idx := range ans {
			if index.KeyName == idx.KeyName {
				idx.ColumnName += "," + index.ColumnName
				ans[i] = idx
				flag = true
				break
			}
		}
		if flag {
			continue
		}

		// 将当前行的数据添加到切片中
		ans = append(ans, index)
	}

	return ans, nil
}

// 获取表的建表语句
func getTables(db *sql.DB) (tableNames []string, err error) {
	rows, err := db.Query("SHOW TABLES")
	for rows.Next() {
		var tbName string
		err = rows.Scan(&tbName)
		if err != nil {
			logrus.WithError(err).Errorln("fail to get tables from db")
			continue
		}
		tableNames = append(tableNames, tbName)
	}
	return
}

func getCreateTableSql(sourceDB *sql.DB, table string) string {
	rows, err := sourceDB.Query("SHOW CREATE TABLE " + table)
	if err != nil {
		logrus.WithError(err).Errorln("fail to show create table " + table)
	}
	var tableName string
	var createTableSql string
	for rows.Next() {
		err := rows.Scan(&tableName, &createTableSql)
		if err != nil {
			logrus.WithError(err).Errorln("fail to get create table sql")
			continue
		}
	}
	return createTableSql
}

func getCreateTables(targetTables []string, sourceTables []Table) (ans []string) {
	for _, table := range sourceTables {
		if !tableExists(targetTables, table.Name) {
			ans = append(ans, table.Create)
		}
	}
	return
}
func defaultEqual(d1, d2 interface{}) bool {
	if d1 == nil && d2 == nil {
		return true
	}
	var bts1 []byte
	var bts2 []byte
	if d1 == nil {
		bts2 = d2.([]byte)
		if len(bts2) == 0 {
			return true
		}
		return false
	}
	if d2 == nil {
		bts1 = d1.([]byte)
		if len(bts1) == 0 {
			return true
		}
		return false
	}
	return string(bts1) == string(bts2)
}

func commentEqual(sComments, tComments map[string]string, name1, name2 string) bool {
	c1, ok1 := sComments[name1]
	c2, ok2 := tComments[name2]
	if !ok1 && !ok2 {
		return true
	}
	if !ok2 || !ok1 {
		return false
	}
	return c1 == c2
}

// 检查列是否存在
func columnExists(columns []*Column, columnName string) bool {
	for _, column := range columns {
		if column.Field == columnName {
			return true
		}
	}
	return false
}

func indexExists(indexes *[]Index, indexName string) bool {
	for _, idx := range *indexes {
		if idx.KeyName == indexName {
			return true
		}
	}
	return false
}

func tableExists(tables []string, table string) bool {
	for _, s := range tables {
		if s == table {
			return true
		}
	}
	return false
}

// 将字符串切片连接为一个字符串
func joinStrings(strings []string, separator string) string {
	if len(strings) == 0 {
		return ""
	}
	if len(strings) == 1 {
		return strings[0]
	}
	return strings[0] + separator + joinStrings(strings[1:], separator)
}

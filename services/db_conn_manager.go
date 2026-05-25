package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DBConnConfig struct {
	DBType   int
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
	Timeout  time.Duration
}

type DBConnection struct {
	Config   *DBConnConfig
	MySQLDB  *sql.DB
	MongoDB  *mongo.Client
	RedisDB  *redis.Client
	LastUsed time.Time
}

type DBConnManager struct {
	mu       sync.RWMutex
	connPool map[string]*DBConnection
}

var (
	globalDBConnManager *DBConnManager
	dbConnOnce          sync.Once
)

func GetDBConnManager() *DBConnManager {
	dbConnOnce.Do(func() {
		globalDBConnManager = &DBConnManager{
			connPool: make(map[string]*DBConnection),
		}
	})
	return globalDBConnManager
}

func dbConnPoolKey(dbType int, host string, port int, dbName string) string {
	return fmt.Sprintf("%d-%s:%d/%s", dbType, host, port, dbName)
}

func (m *DBConnManager) GetConnection(ctx context.Context, config *DBConnConfig) (*DBConnection, error) {
	key := dbConnPoolKey(config.DBType, config.Host, config.Port, config.DBName)

	m.mu.RLock()
	conn, exists := m.connPool[key]
	m.mu.RUnlock()

	if exists {
		conn.LastUsed = time.Now()
		return conn, nil
	}

	timeout := config.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	conn = &DBConnection{
		Config:   config,
		LastUsed: time.Now(),
	}

	var err error
	switch config.DBType {
	case 1: // MySQL
		err = m.connectMySQL(ctx, conn, timeout)
	case 2: // PostgreSQL
		err = m.connectPostgreSQL(ctx, conn, timeout)
	case 3: // MongoDB
		err = m.connectMongoDB(ctx, conn, timeout)
	case 4: // Redis
		err = m.connectRedis(ctx, conn, config, timeout)
	case 5: // CouchDB
		err = m.connectCouchDB(ctx, config, timeout)
	default:
		return nil, fmt.Errorf("unsupported db type: %d", config.DBType)
	}

	if err != nil {
		return nil, fmt.Errorf("connect to %s failed: %v", key, err)
	}

	m.mu.Lock()
	m.connPool[key] = conn
	m.mu.Unlock()

	return conn, nil
}

func (m *DBConnManager) connectMySQL(ctx context.Context, conn *DBConnection, timeout time.Duration) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&timeout=%s",
		conn.Config.Username, conn.Config.Password, conn.Config.Host, conn.Config.Port, conn.Config.DBName, timeout)
	var err error
	conn.MySQLDB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	conn.MySQLDB.SetMaxOpenConns(5)
	conn.MySQLDB.SetMaxIdleConns(2)
	conn.MySQLDB.SetConnMaxLifetime(5 * time.Minute)
	return conn.MySQLDB.PingContext(ctx)
}

func (m *DBConnManager) connectPostgreSQL(ctx context.Context, conn *DBConnection, timeout time.Duration) error {
	addr := net.JoinHostPort(conn.Config.Host, fmt.Sprintf("%d", conn.Config.Port))
	dialer := net.Dialer{Timeout: timeout}
	tcpConn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("postgresql dial failed: %v", err)
	}
	tcpConn.Close()
	return nil
}

func (m *DBConnManager) connectMongoDB(ctx context.Context, conn *DBConnection, timeout time.Duration) error {
	authSource := "admin"
	if strings.TrimSpace(conn.Config.DBName) != "" {
		authSource = strings.TrimSpace(conn.Config.DBName)
	}
	uri := fmt.Sprintf("mongodb://%s:%d/?directConnection=true&serverSelectionTimeoutMS=%d&connectTimeoutMS=%d",
		conn.Config.Host, conn.Config.Port, timeout.Milliseconds(), timeout.Milliseconds())
	clientOpts := options.Client().ApplyURI(uri)
	if conn.Config.Username != "" {
		clientOpts.SetAuth(options.Credential{
			Username:   conn.Config.Username,
			Password:   conn.Config.Password,
			AuthSource: authSource,
		})
	}
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return fmt.Errorf("mongodb connect failed: %v", err)
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		_ = client.Disconnect(context.Background())
		return fmt.Errorf("mongodb ping failed: %v", err)
	}
	conn.MongoDB = client
	return nil
}

func (m *DBConnManager) connectRedis(ctx context.Context, conn *DBConnection, config *DBConnConfig, timeout time.Duration) error {
	redisAddr := net.JoinHostPort(config.Host, fmt.Sprintf("%d", config.Port))
	conn.RedisDB = redis.NewClient(&redis.Options{
		Addr:        redisAddr,
		Password:    config.Password,
		DB:          0,
		DialTimeout: timeout,
	})
	_, err := conn.RedisDB.Ping(ctx).Result()
	return err
}

func (m *DBConnManager) connectCouchDB(ctx context.Context, config *DBConnConfig, timeout time.Duration) error {
	couchURL := fmt.Sprintf("http://%s:%d", config.Host, config.Port)
	client := &http.Client{Timeout: timeout}
	req, _ := http.NewRequestWithContext(ctx, "GET", couchURL+"/_up", nil)
	if config.Username != "" || config.Password != "" {
		req.SetBasicAuth(config.Username, config.Password)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("couchdb ping failed: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("couchdb returned status: %d", resp.StatusCode)
	}
	return nil
}

func (m *DBConnManager) ExecuteQuery(ctx context.Context, conn *DBConnection, query string) ([]map[string]string, error) {
	switch conn.Config.DBType {
	case 1:
		if conn.MySQLDB == nil {
			return nil, fmt.Errorf("mysql connection not established")
		}
		return m.executeSQLQuery(ctx, conn.MySQLDB, query)
	case 2:
		return m.executePostgresQuery(ctx, conn, query)
	case 3:
		return m.executeMongoCommand(ctx, conn, query)
	case 4:
		return m.executeRedisCommand(ctx, conn.RedisDB, query)
	case 5:
		return m.executeCouchHTTP(ctx, conn, query)
	}
	return nil, fmt.Errorf("unsupported db type: %d", conn.Config.DBType)
}

func (m *DBConnManager) executeSQLQuery(ctx context.Context, db *sql.DB, query string) ([]map[string]string, error) {
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]string
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]string)
		for i, col := range columns {
			row[col] = sqlCellToString(values[i])
		}
		results = append(results, row)
	}

	if results == nil {
		results = make([]map[string]string, 0)
	}
	return results, nil
}

func (m *DBConnManager) executeMongoCommand(ctx context.Context, conn *DBConnection, query string) ([]map[string]string, error) {
	if conn.MongoDB == nil {
		return nil, fmt.Errorf("mongodb connection not established")
	}

	var cmd bson.D
	switch query {
	case "/", "/_isMaster", "/hello":
		cmd = bson.D{{Key: "hello", Value: 1}}
	case "/buildInfo":
		cmd = bson.D{{Key: "buildInfo", Value: 1}}
	case "/_cmdLineOpts":
		cmd = bson.D{{Key: "getCmdLineOpts", Value: 1}}
	default:
		return nil, fmt.Errorf("unsupported mongodb query: %s", query)
	}

	result := bson.M{}
	err := conn.MongoDB.Database("admin").RunCommand(ctx, cmd).Decode(&result)
	if err != nil && query == "/_isMaster" {
		result = bson.M{}
		err = conn.MongoDB.Database("admin").RunCommand(ctx, bson.D{{Key: "isMaster", Value: 1}}).Decode(&result)
	}

	payload := mongoProbePayload{
		Command: query,
		Host:    conn.Config.Host,
		Port:    conn.Config.Port,
	}
	if err != nil {
		payload.Error = err.Error()
	} else {
		payload.Result = result
	}
	raw, _ := json.Marshal(payload)
	return []map[string]string{{"response": string(raw)}}, nil
}

func (m *DBConnManager) executePostgresQuery(ctx context.Context, conn *DBConnection, query string) ([]map[string]string, error) {
	return []map[string]string{{"note": "postgres query requires lib/pq driver"}}, nil
}

func (m *DBConnManager) executeRedisCommand(ctx context.Context, client *redis.Client, cmd string) ([]map[string]string, error) {
	args := strings.Fields(cmd)
	if len(args) == 0 {
		return nil, fmt.Errorf("empty redis command")
	}
	interfaceArgs := make([]interface{}, len(args))
	for i, arg := range args {
		interfaceArgs[i] = arg
	}

	val, err := client.Do(ctx, interfaceArgs...).Result()
	if err != nil {
		return nil, err
	}
	row := map[string]string{"result": fmt.Sprintf("%v", val)}
	return []map[string]string{row}, nil
}

func (m *DBConnManager) executeCouchHTTP(ctx context.Context, conn *DBConnection, path string) ([]map[string]string, error) {
	useAuth := conn.Config.Username != ""
	switch {
	case strings.HasPrefix(path, "public:"):
		useAuth = false
		path = strings.TrimPrefix(path, "public:")
	case strings.HasPrefix(path, "admin:"):
		useAuth = true
		path = strings.TrimPrefix(path, "admin:")
	}

	baseURL := fmt.Sprintf("http://%s:%d", conn.Config.Host, conn.Config.Port)
	client := &http.Client{Timeout: 30 * time.Second}
	req, _ := http.NewRequestWithContext(ctx, "GET", baseURL+path, nil)
	if useAuth && conn.Config.Username != "" {
		req.SetBasicAuth(conn.Config.Username, conn.Config.Password)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	headers := map[string]string{}
	for _, key := range []string{"Access-Control-Allow-Origin", "Access-Control-Allow-Credentials", "Server", "Content-Type"} {
		if val := resp.Header.Get(key); val != "" {
			headers[key] = val
		}
	}
	payload := httpProbePayload{
		StatusCode: resp.StatusCode,
		Body:       string(body),
		Headers:    headers,
		Host:       conn.Config.Host,
		Port:       conn.Config.Port,
	}
	raw, _ := json.Marshal(payload)
	row := map[string]string{"response": string(raw)}
	return []map[string]string{row}, nil
}

func sqlCellToString(v interface{}) string {
	if v == nil {
		return "NULL"
	}
	switch t := v.(type) {
	case []byte:
		return string(t)
	case string:
		return t
	case time.Time:
		return t.Format("2006-01-02 15:04:05")
	default:
		return fmt.Sprintf("%v", t)
	}
}

func (m *DBConnManager) GetDatabases(ctx context.Context, conn *DBConnection) ([]string, error) {
	switch conn.Config.DBType {
	case 1:
		return m.executeSQLList(ctx, conn.MySQLDB, "SHOW DATABASES")
	case 2:
		return nil, fmt.Errorf("postgresql requires lib/pq driver")
	case 3:
		if conn.MongoDB == nil {
			return nil, fmt.Errorf("mongodb connection not established")
		}
		return conn.MongoDB.ListDatabaseNames(ctx, bson.D{})
	case 4:
		return []string{"0"}, nil
	case 5:
		return m.executeCouchDBList(ctx, conn)
	}
	return nil, fmt.Errorf("unsupported db type: %d", conn.Config.DBType)
}

func (m *DBConnManager) executeSQLList(ctx context.Context, db *sql.DB, query string) ([]string, error) {
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			continue
		}
		list = append(list, name)
	}
	return list, nil
}

func (m *DBConnManager) executeCouchDBList(ctx context.Context, conn *DBConnection) ([]string, error) {
	results, err := m.executeCouchHTTP(ctx, conn, "/_all_dbs")
	if err != nil {
		return nil, err
	}
	if len(results) > 0 {
		var list []string
		if payload, ok := parseHTTPProbePayload(results[0]["response"]); ok {
			_ = json.Unmarshal([]byte(payload.Body), &list)
		}
		return list, nil
	}
	return nil, nil
}

func (m *DBConnManager) GetTables(ctx context.Context, conn *DBConnection, dbName string) ([]string, error) {
	switch conn.Config.DBType {
	case 1:
		return m.executeSQLList(ctx, conn.MySQLDB, fmt.Sprintf("SHOW TABLES FROM `%s`", dbName))
	case 2:
		return nil, fmt.Errorf("postgresql requires lib/pq driver")
	case 3:
		if conn.MongoDB == nil {
			return nil, fmt.Errorf("mongodb connection not established")
		}
		return conn.MongoDB.Database(dbName).ListCollectionNames(ctx, bson.D{})
	case 4:
		return []string{}, nil
	case 5:
		if strings.TrimSpace(dbName) == "" {
			return nil, nil
		}
		return []string{dbName}, nil
	}
	return nil, fmt.Errorf("unsupported db type: %d", conn.Config.DBType)
}

func (m *DBConnManager) GetColumns(ctx context.Context, conn *DBConnection, dbName, tableName string) ([]map[string]string, error) {
	switch conn.Config.DBType {
	case 1:
		return m.executeSQLQuery(ctx, conn.MySQLDB,
			fmt.Sprintf("SELECT COLUMN_NAME, DATA_TYPE, COLUMN_COMMENT FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA='%s' AND TABLE_NAME='%s'", dbName, tableName))
	case 2:
		return nil, fmt.Errorf("postgresql requires lib/pq driver")
	case 3:
		return m.inspectMongoColumns(ctx, conn, dbName, tableName)
	case 5:
		return m.inspectCouchColumns(ctx, conn, dbName)
	default:
		return nil, nil
	}
}

func (m *DBConnManager) GetFieldSamples(ctx context.Context, conn *DBConnection, dbName, tableName, fieldName string, limit int64) ([]string, error) {
	switch conn.Config.DBType {
	case 3:
		return m.mongoFieldSamples(ctx, conn, dbName, tableName, fieldName, limit)
	case 5:
		return m.couchFieldSamples(ctx, conn, dbName, fieldName, limit)
	default:
		return nil, fmt.Errorf("field sample fetch unsupported for db type: %d", conn.Config.DBType)
	}
}

func (m *DBConnManager) inspectMongoColumns(ctx context.Context, conn *DBConnection, dbName, tableName string) ([]map[string]string, error) {
	if conn.MongoDB == nil {
		return nil, fmt.Errorf("mongodb connection not established")
	}
	coll := conn.MongoDB.Database(dbName).Collection(tableName)
	cursor, err := coll.Find(ctx, bson.D{}, options.Find().SetLimit(20))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	fields := map[string]string{}
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		for key, val := range doc {
			if key == "_id" {
				continue
			}
			if _, exists := fields[key]; exists {
				continue
			}
			fields[key] = nosqlValueType(val)
		}
	}
	return formatNoSQLColumns(fields), nil
}

func (m *DBConnManager) inspectCouchColumns(ctx context.Context, conn *DBConnection, dbName string) ([]map[string]string, error) {
	rows, err := m.fetchCouchDocs(ctx, conn, dbName, 20)
	if err != nil {
		return nil, err
	}
	fields := map[string]string{}
	for _, row := range rows {
		for key, val := range row.Doc {
			if strings.HasPrefix(key, "_") {
				continue
			}
			if _, exists := fields[key]; exists {
				continue
			}
			fields[key] = nosqlValueType(val)
		}
	}
	return formatNoSQLColumns(fields), nil
}

func (m *DBConnManager) mongoFieldSamples(ctx context.Context, conn *DBConnection, dbName, tableName, fieldName string, limit int64) ([]string, error) {
	if conn.MongoDB == nil {
		return nil, fmt.Errorf("mongodb connection not established")
	}
	coll := conn.MongoDB.Database(dbName).Collection(tableName)
	cursor, err := coll.Find(ctx, bson.D{}, options.Find().SetLimit(limit).SetProjection(bson.D{{Key: fieldName, Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var out []string
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		if val, ok := doc[fieldName]; ok {
			if text := nosqlValueString(val); text != "" {
				out = append(out, text)
			}
		}
	}
	return out, nil
}

func (m *DBConnManager) couchFieldSamples(ctx context.Context, conn *DBConnection, dbName, fieldName string, limit int64) ([]string, error) {
	rows, err := m.fetchCouchDocs(ctx, conn, dbName, limit)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, row := range rows {
		if val, ok := row.Doc[fieldName]; ok {
			if text := nosqlValueString(val); text != "" {
				out = append(out, text)
			}
		}
	}
	return out, nil
}

func (m *DBConnManager) fetchCouchDocs(ctx context.Context, conn *DBConnection, dbName string, limit int64) ([]struct {
	Doc map[string]interface{} `json:"doc"`
}, error) {
	path := fmt.Sprintf("/%s/_all_docs?include_docs=true&limit=%d", dbName, limit)
	results, err := m.executeCouchHTTP(ctx, conn, path)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	payload, ok := parseHTTPProbePayload(results[0]["response"])
	if !ok {
		return nil, fmt.Errorf("parse couchdb docs response failed")
	}
	var docList struct {
		Rows []struct {
			Doc map[string]interface{} `json:"doc"`
		} `json:"rows"`
	}
	if err := json.Unmarshal([]byte(payload.Body), &docList); err != nil {
		return nil, err
	}
	return docList.Rows, nil
}

func formatNoSQLColumns(fields map[string]string) []map[string]string {
	out := make([]map[string]string, 0, len(fields))
	for name, dataType := range fields {
		out = append(out, map[string]string{
			"COLUMN_NAME": name,
			"DATA_TYPE":   dataType,
		})
	}
	return out
}

func nosqlValueType(v interface{}) string {
	switch v.(type) {
	case string:
		return "text"
	case map[string]interface{}, []interface{}, bson.M, bson.A:
		return "json"
	case bool:
		return "bool"
	case int32, int64, int, float32, float64:
		return "number"
	default:
		return "text"
	}
}

func nosqlValueString(v interface{}) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case []byte:
		return string(t)
	case map[string]interface{}, []interface{}, bson.M, bson.A:
		raw, _ := json.Marshal(t)
		return string(raw)
	default:
		return fmt.Sprintf("%v", t)
	}
}

func (m *DBConnManager) Close(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if conn, exists := m.connPool[key]; exists {
		if conn.MySQLDB != nil {
			conn.MySQLDB.Close()
		}
		if conn.MongoDB != nil {
			_ = conn.MongoDB.Disconnect(context.Background())
		}
		if conn.RedisDB != nil {
			conn.RedisDB.Close()
		}
		delete(m.connPool, key)
	}
}

func (m *DBConnManager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for key, conn := range m.connPool {
		if conn.MySQLDB != nil {
			conn.MySQLDB.Close()
		}
		if conn.MongoDB != nil {
			_ = conn.MongoDB.Disconnect(context.Background())
		}
		if conn.RedisDB != nil {
			conn.RedisDB.Close()
		}
		delete(m.connPool, key)
	}
}

package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"smart/tools/enums"

	log "github.com/sirupsen/logrus"
	_ "github.com/mattn/go-sqlite3"
)

type CveDB struct {
	db *sql.DB
}

type CveRecord struct {
	Cve               string
	TitleZh           string
	Severity          string
	SeverityLevel     int
	Product           string
	Vendor            string
	CpeConfigurations []byte
}

type CpeConfiguration struct {
	Nodes []CpeNode `json:"nodes"`
}

type CpeNode struct {
	Operator string     `json:"operator"`
	CpeMatch []CpeMatch `json:"cpe_match"`
	Children []CpeNode  `json:"children"`
}

type CpeMatch struct {
	Vulnerable            bool   `json:"vulnerable"`
	Cpe23Uri              string `json:"cpe23Uri"`
	VersionStartExcluding string `json:"versionStartExcluding"`
	VersionEndExcluding   string `json:"versionEndExcluding"`
	VersionStartIncluding string `json:"versionStartIncluding"`
	VersionEndIncluding   string `json:"versionEndIncluding"`
}

type CveMatchResult struct {
	Cve       string `json:"cve"`
	Title     string `json:"title"`
	Severity  string `json:"severity"`
	Product   string `json:"product"`
	Version   string `json:"version"`
	RiskLevel int    `json:"riskLevel"`
}

var (
	globalCveDB   *CveDB
	cveDBOnce     sync.Once
	resolvedCveDB string
)

// ResolveCveDBPath 解析 default-cve.db 路径（开发/部署多目录兼容）
func ResolveCveDBPath() string {
	if resolvedCveDB != "" {
		if _, err := os.Stat(resolvedCveDB); err == nil {
			return resolvedCveDB
		}
	}
	candidates := []string{
		"data/default-cve.db",
		"default-cve.db",
	}
	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates,
			filepath.Join(wd, "data", "default-cve.db"),
			filepath.Join(wd, "default-cve.db"),
		)
	}
	if exe, err := os.Executable(); err == nil {
		root := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(root, "data", "default-cve.db"),
			filepath.Join(root, "default-cve.db"),
		)
	}
	if enums.SystemUpgradeProjectDir != "" && enums.SystemUpgradeProjectDir != "/opt/laozhi/" {
		base := strings.TrimRight(enums.SystemUpgradeProjectDir, `/\`)
		candidates = append(candidates,
			filepath.Join(base, "data", "default-cve.db"),
			filepath.Join(base, "default-cve.db"),
		)
	}
	seen := make(map[string]struct{})
	for _, p := range candidates {
		p = filepath.Clean(p)
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
			resolvedCveDB = p
			return p
		}
	}
	return ""
}

// InitCveDB 启动时初始化 CVE 库（打日志便于排查）
func InitCveDB() {
	_ = GetCveDB()
}

func GetCveDB() *CveDB {
	cveDBOnce.Do(func() {
		dbPath := ResolveCveDBPath()
		if dbPath == "" {
			log.Warn("default-cve.db not found; tried data/default-cve.db and project paths")
			return
		}
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Errorf("open cve db %s: %v", dbPath, err)
			return
		}
		if err := db.Ping(); err != nil {
			log.Errorf("ping cve db %s: %v", dbPath, err)
			_ = db.Close()
			return
		}
		globalCveDB = &CveDB{db: db}
		log.Infof("CVE database loaded: %s", dbPath)
	})
	return globalCveDB
}

// CveDBResolvedPath 当前使用的 CVE 库路径
func CveDBResolvedPath() string {
	if resolvedCveDB != "" {
		return resolvedCveDB
	}
	return ResolveCveDBPath()
}

func (d *CveDB) IsAvailable() bool {
	return d != nil && d.db != nil
}

func (d *CveDB) Close() {
	if d.db != nil {
		d.db.Close()
	}
}

func (d *CveDB) QueryByProduct(productName string) ([]CveRecord, error) {
	keywords := getCveKeywords(productName)
	if len(keywords) == 0 {
		return nil, nil
	}

	whereClause := buildWhereClause(keywords)
	query := fmt.Sprintf(`SELECT cve, COALESCE(title_zh, ''), COALESCE(severity, ''), COALESCE(product, ''), COALESCE(vendor, ''), COALESCE(cpe_configurations, '{}') FROM cves WHERE %s LIMIT 200`, whereClause)

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query cve by product %s failed: %v", productName, err)
	}
	defer rows.Close()

	var results []CveRecord
	for rows.Next() {
		var r CveRecord
		var cpeStr string
		if err := rows.Scan(&r.Cve, &r.TitleZh, &r.Severity, &r.Product, &r.Vendor, &cpeStr); err != nil {
			continue
		}
		r.CpeConfigurations = []byte(cpeStr)
		r.SeverityLevel = severityToLevel(r.Severity)
		results = append(results, r)
	}
	return results, nil
}

func (d *CveDB) MatchCpe(productVersion string, cpeConfig []byte) (bool, []CpeMatch) {
	var config CpeConfiguration
	if err := json.Unmarshal(cpeConfig, &config); err != nil {
		return false, nil
	}
	if len(config.Nodes) == 0 {
		return false, nil
	}

	var matched []CpeMatch
	overallMatch := false

	for _, node := range config.Nodes {
		nodeMatch, nodeCpes := matchCpeNode(productVersion, node, node.Operator)
		if nodeMatch {
			overallMatch = true
			matched = append(matched, nodeCpes...)
		}
	}
	return overallMatch, matched
}

type cveKeyword struct {
	Vendor  string
	Product []string
}

func getCveKeywords(product string) []cveKeyword {
	if m, ok := cveProductMap[product]; ok {
		return m
	}
	return nil
}

func buildWhereClause(keywords []cveKeyword) string {
	var parts []string
	for _, kw := range keywords {
		var conds []string
		if kw.Vendor != "" {
			conds = append(conds, fmt.Sprintf("',' || vendor || ',' LIKE '%%,%s,%%'", kw.Vendor))
		}
		if len(kw.Product) > 0 {
			var prodConds []string
			for _, p := range kw.Product {
				prodConds = append(prodConds, fmt.Sprintf("',' || product || ',' LIKE '%%,%s,%%'", p))
			}
			if len(conds) > 0 {
				conds = append(conds, "("+strings.Join(prodConds, " OR ")+")")
			} else {
				conds = append(conds, strings.Join(prodConds, " OR "))
			}
		}
		if len(conds) > 0 {
			parts = append(parts, "("+strings.Join(conds, " AND ")+")")
		}
	}
	if len(parts) == 0 {
		return "1=0"
	}
	return strings.Join(parts, " OR ")
}

func matchCpeNode(version string, node CpeNode, operator string) (bool, []CpeMatch) {
	if operator == "" {
		operator = "AND"
	}

	var matched []CpeMatch
	var result bool
	if operator == "OR" {
		result = false
	} else {
		result = true
	}

	opFunc := func(current bool, next bool) bool {
		if operator == "OR" {
			return current || next
		}
		return current && next
	}

	for _, match := range node.CpeMatch {
		matchResult := matchCpeEntry(version, match)
		if matchResult {
			matched = append(matched, match)
		}
		result = opFunc(result, matchResult)
	}

	for _, child := range node.Children {
		childResult, childCpes := matchCpeNode(version, child, child.Operator)
		if childResult {
			matched = append(matched, childCpes...)
		}
		result = opFunc(result, childResult)
	}

	return result, matched
}

func matchCpeEntry(version string, match CpeMatch) bool {
	startEx := match.VersionStartExcluding
	endEx := match.VersionEndExcluding
	startIn := match.VersionStartIncluding
	endIn := match.VersionEndIncluding

	hasRange := startEx != "" || endEx != "" || startIn != "" || endIn != ""

	if !hasRange {
		parts := strings.Split(match.Cpe23Uri, ":")
		if len(parts) >= 6 && parts[5] != "*" {
			return compareVersions(version, parts[5]) == 0
		}
		return len(parts) >= 6 && parts[5] == "*"
	}

	cmp := compareVersions(version, startIn)
	if startIn != "" && cmp < 0 {
		return false
	}
	cmp = compareVersions(version, startEx)
	if startEx != "" && cmp <= 0 {
		return false
	}
	cmp = compareVersions(version, endIn)
	if endIn != "" && cmp > 0 {
		return false
	}
	cmp = compareVersions(version, endEx)
	if endEx != "" && cmp >= 0 {
		return false
	}
	return true
}

func compareVersions(v1, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var p1, p2 string
		if i < len(parts1) {
			p1 = parts1[i]
		} else {
			p1 = "0"
		}
		if i < len(parts2) {
			p2 = parts2[i]
		} else {
			p2 = "0"
		}

		n1, err1 := strconv.Atoi(p1)
		n2, err2 := strconv.Atoi(p2)

		if err1 == nil && err2 == nil {
			if n1 < n2 {
				return -1
			} else if n1 > n2 {
				return 1
			}
		} else {
			if p1 < p2 {
				return -1
			} else if p1 > p2 {
				return 1
			}
		}
	}
	return 0
}

func severityToLevel(severity string) int {
	switch strings.ToUpper(severity) {
	case "CRITICAL":
		return 4
	case "HIGH":
		return 3
	case "MEDIUM":
		return 2
	case "LOW":
		return 1
	default:
		return 0
	}
}

func (d *CveDB) Count() (int64, error) {
	var total int64
	err := d.db.QueryRow("SELECT COUNT(*) FROM cves").Scan(&total)
	return total, err
}

type CveSearchRecord struct {
	CveID       string
	Title       string
	Severity    string
	Product     string
	Vendor      string
	Description string
}

func (d *CveDB) Search(keyword string, page, size int) ([]CveSearchRecord, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	offset := (page - 1) * size

	var total int64
	var rows *sql.Rows
	var err error

	if keyword != "" {
		kw := "%" + keyword + "%"
		err = d.db.QueryRow("SELECT COUNT(*) FROM cves WHERE cve LIKE ? OR product LIKE ? OR title_zh LIKE ?", kw, kw, kw).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
		rows, err = d.db.Query(
			"SELECT cve, COALESCE(title_zh, ''), COALESCE(severity, ''), COALESCE(product, ''), COALESCE(vendor, ''), COALESCE(description_main_zh, '') FROM cves WHERE cve LIKE ? OR product LIKE ? OR title_zh LIKE ? ORDER BY cve LIMIT ? OFFSET ?",
			kw, kw, kw, size, offset,
		)
	} else {
		err = d.db.QueryRow("SELECT COUNT(*) FROM cves").Scan(&total)
		if err != nil {
			return nil, 0, err
		}
		rows, err = d.db.Query(
			"SELECT cve, COALESCE(title_zh, ''), COALESCE(severity, ''), COALESCE(product, ''), COALESCE(vendor, ''), COALESCE(description_main_zh, '') FROM cves ORDER BY cve LIMIT ? OFFSET ?",
			size, offset,
		)
	}

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var records []CveSearchRecord
	for rows.Next() {
		var r CveSearchRecord
		if err := rows.Scan(&r.CveID, &r.Title, &r.Severity, &r.Product, &r.Vendor, &r.Description); err != nil {
			continue
		}
		records = append(records, r)
	}
	return records, total, nil
}

var cveProductMap = map[string][]cveKeyword{
	"openssl":       {{"openssl", []string{"openssl"}}},
	"libssl":        {{"openssl", []string{"openssl"}}},
	"python":        {{"python", []string{"python"}}},
	"python3":       {{"python", []string{"python"}}},
	"php":           {{"php", []string{"php"}}},
	"mysql":         {{"mysql", []string{"mysql"}}, {"oracle", []string{"mysql"}}},
	"mariadb":       {{"mysql", []string{"mysql"}}},
	"nginx":         {{"nginx", []string{"nginx"}}},
	"apache2":       {{"apache", []string{"http_server"}}},
	"httpd":         {{"apache", []string{"http_server"}}},
	"openssh":       {{"openbsd", []string{"openssh"}}},
	"ssh":           {{"openbsd", []string{"openssh"}}},
	"libssh":        {{"libssh", []string{"libssh"}}},
	"bash":          {{"gnu", []string{"bash"}}},
	"glibc":         {{"gnu", []string{"glibc"}}},
	"libc":          {{"gnu", []string{"glibc"}}},
	"libc6":         {{"gnu", []string{"glibc"}}},
	"systemd":       {{"systemd", []string{"systemd"}}},
	"curl":          {{"curl", []string{"curl"}}},
	"libcurl":       {{"curl", []string{"curl"}}},
	"wget":          {{"gnu", []string{"wget"}}},
	"git":           {{"git", []string{"git"}}},
	"vim":           {{"vim", []string{"vim"}}},
	"sudo":          {{"sudo", []string{"sudo"}}},
	"libxml2":       {{"xmlsoft", []string{"libxml2"}}},
	"libxslt":       {{"xmlsoft", []string{"libxslt"}}},
	"zlib":          {{"zlib", []string{"zlib"}}},
	"libz":          {{"zlib", []string{"zlib"}}},
	"sqlite":        {{"sqlite", []string{"sqlite"}}},
	"sqlite3":       {{"sqlite", []string{"sqlite"}}},
	"postgresql":    {{"postgresql", []string{"postgresql"}}},
	"redis":         {{"redis", []string{"redis"}}},
	"memcached":     {{"memcached", []string{"memcached"}}},
	"docker":        {{"docker", []string{"docker"}}},
	"containerd":    {{"containerd", []string{"containerd"}}},
	"kubernetes":    {{"kubernetes", []string{"kubernetes"}}},
	"kubelet":       {{"kubernetes", []string{"kubernetes"}}},
	"tomcat":        {{"apache", []string{"tomcat"}}},
	"java":          {{"oracle", []string{"jdk"}}, {"oracle", []string{"jre"}}},
	"openjdk":       {{"oracle", []string{"jdk"}}, {"oracle", []string{"jre"}}},
	"nodejs":        {{"nodejs", []string{"node.js"}}},
	"npm":           {{"npm", []string{"npm"}}},
	"ruby":          {{"ruby", []string{"ruby"}}},
	"perl":          {{"perl", []string{"perl"}}},
	"go":            {{"golang", []string{"go"}}},
	"golang":        {{"golang", []string{"go"}}},
	"rsync":         {{"rsync", []string{"rsync"}}},
	"ntp":           {{"ntp", []string{"ntp"}}},
	"dnsmasq":       {{"dnsmasq", []string{"dnsmasq"}}},
	"bind9":         {{"isc", []string{"bind"}}},
	"named":         {{"isc", []string{"bind"}}},
	"snmpd":         {{"net-snmp", []string{"net-snmp"}}},
	"libpcap":       {{"tcpdump", []string{"libpcap"}}},
	"tcpdump":       {{"tcpdump", []string{"tcpdump"}}},
	"gcc":           {{"gnu", []string{"gcc"}}},
	"libstdc++":     {{"gnu", []string{"gcc"}}},
	"pcre":          {{"pcre", []string{"pcre"}}},
	"pcre2":         {{"pcre", []string{"pcre"}}},
	"libpcre":       {{"pcre", []string{"pcre"}}},
	"libpcre2":      {{"pcre", []string{"pcre"}}},
	"readline":      {{"gnu", []string{"readline"}}},
	"libreadline":   {{"gnu", []string{"readline"}}},
	"ncurses":       {{"gnu", []string{"ncurses"}}},
	"libncurses":    {{"gnu", []string{"ncurses"}}},
	"libpng":        {{"libpng", []string{"libpng"}}},
	"libjpeg":       {{"ijg", []string{"libjpeg"}}},
	"libtiff":       {{"libtiff", []string{"libtiff"}}},
	"libwebp":       {{"webmproject", []string{"libwebp"}}},
	"freetype":      {{"freetype", []string{"freetype"}}},
	"fontconfig":    {{"fontconfig", []string{"fontconfig"}}},
	"cairo":         {{"cairo", []string{"cairo"}}},
	"pango":         {{"pango", []string{"pango"}}},
	"gtk":           {{"gtk", []string{"gtk+"}}},
	"gtk3":          {{"gtk", []string{"gtk+"}}},
	"qt":            {{"qt", []string{"qt"}}},
	"qt5":           {{"qt", []string{"qt"}}},
	"qt6":           {{"qt", []string{"qt"}}},
	"libvirt":       {{"redhat", []string{"libvirt"}}},
	"qemu":          {{"qemu", []string{"qemu"}}},
	"xen":           {{"xen", []string{"xen"}}},
	"kernel":        {{"linux", []string{"linux_kernel"}}},
	"linux":         {{"linux", []string{"linux_kernel"}}},
	"grub":          {{"gnu", []string{"grub"}}},
	"grub2":         {{"gnu", []string{"grub"}}},
	"selinux":       {{"nsa", []string{"selinux"}}},
	"apparmor":      {{"canonical", []string{"apparmor"}}},
	"policykit":     {{"freedesktop", []string{"policykit"}}},
	"pkexec":        {{"freedesktop", []string{"policykit"}}},
	"polkit":        {{"freedesktop", []string{"policykit"}}},
	"dbus":          {{"freedesktop", []string{"dbus"}}},
	"avahi":         {{"avahi", []string{"avahi"}}},
	"cups":          {{"cups", []string{"cups"}}},
	"samba":         {{"samba", []string{"samba"}}},
	"smbd":          {{"samba", []string{"samba"}}},
	"nmbd":          {{"samba", []string{"samba"}}},
	"vsftpd":        {{"vsftpd", []string{"vsftpd"}}},
	"proftpd":       {{"proftpd", []string{"proftpd"}}},
	"pure-ftpd":     {{"pureftpd", []string{"pure-ftpd"}}},
	"ntpd":          {{"ntp", []string{"ntp"}}},
	"chrony":        {{"chrony", []string{"chrony"}}},
	"corosync":      {{"corosync", []string{"corosync"}}},
	"pacemaker":     {{"pacemaker", []string{"pacemaker"}}},
	"keepalived":    {{"keepalived", []string{"keepalived"}}},
	"haproxy":       {{"haproxy", []string{"haproxy"}}},
	"varnish":       {{"varnish", []string{"varnish"}}},
	"squid":         {{"squid", []string{"squid"}}},
	"mongodb":       {{"mongodb", []string{"mongodb"}}},
	"mongod":        {{"mongodb", []string{"mongodb"}}},
	"elasticsearch": {{"elastic", []string{"elasticsearch"}}},
	"logstash":      {{"elastic", []string{"logstash"}}},
	"kibana":        {{"elastic", []string{"kibana"}}},
	"rabbitmq":      {{"rabbitmq", []string{"rabbitmq"}}},
	"erlang":        {{"erlang", []string{"erlang"}}},
	"activemq":      {{"apache", []string{"activemq"}}},
	"kafka":         {{"apache", []string{"kafka"}}},
	"zookeeper":     {{"apache", []string{"zookeeper"}}},
	"consul":        {{"consul", []string{"consul"}}},
	"etcd":          {{"etcd", []string{"etcd"}}},
	"envoy":         {{"envoy", []string{"envoy"}}},
	"istio":         {{"istio", []string{"istio"}}},
	"prometheus":    {{"prometheus", []string{"prometheus"}}},
	"grafana":       {{"grafana", []string{"grafana"}}},
	"jenkins":       {{"jenkins", []string{"jenkins"}}},
	"gitlab":        {{"gitlab", []string{"gitlab"}}},
	"jira":          {{"atlassian", []string{"jira"}}},
	"confluence":    {{"atlassian", []string{"confluence"}}},
	"nexus":         {{"sonatype", []string{"nexus"}}},
	"sonarqube":     {{"sonarsource", []string{"sonarqube"}}},
	"hadoop":        {{"apache", []string{"hadoop"}}},
	"spark":         {{"apache", []string{"spark"}}},
	"flink":         {{"apache", []string{"flink"}}},
	"storm":         {{"apache", []string{"storm"}}},
	"hbase":         {{"apache", []string{"hbase"}}},
	"hive":          {{"apache", []string{"hive"}}},
	"presto":        {{"presto", []string{"presto"}}},
	"trino":         {{"trino", []string{"trino"}}},
	"druid":         {{"druid", []string{"druid"}}},
	"clickhouse":    {{"clickhouse", []string{"clickhouse"}}},
	"influxdb":      {{"influxdb", []string{"influxdb"}}},
	"timescaledb":   {{"timescale", []string{"timescaledb"}}},
	"couchdb":       {{"apache", []string{"couchdb"}}},
	"neo4j":         {{"neo4j", []string{"neo4j"}}},
	"cassandra":     {{"apache", []string{"cassandra"}}},
	"cockroachdb":   {{"cockroach", []string{"cockroachdb"}}},
	"tidb":          {{"pingcap", []string{"tidb"}}},
	"minio":         {{"minio", []string{"minio"}}},
	"ceph":          {{"ceph", []string{"ceph"}}},
	"glusterfs":     {{"gluster", []string{"glusterfs"}}},
	"nfs":           {{"linux", []string{"linux_kernel"}}},
	"nfs-utils":     {{"linux", []string{"linux_kernel"}}},
	"openldap":      {{"openldap", []string{"openldap"}}},
	"slapd":         {{"openldap", []string{"openldap"}}},
	"389-ds":        {{"fedoraproject", []string{"389_directory_server"}}},
	"sssd":          {{"freedesktop", []string{"sssd"}}},
	"krb5":          {{"mit", []string{"kerberos"}}},
	"kerberos":      {{"mit", []string{"kerberos"}}},
	"libkrb5":       {{"mit", []string{"kerberos"}}},
	"heimdal":       {{"heimdal", []string{"heimdal"}}},
	"pam":           {{"linux", []string{"pam"}}},
	"libpam":        {{"linux", []string{"pam"}}},
	"shadow":        {{"shadow", []string{"shadow"}}},
	"login":         {{"shadow", []string{"shadow"}}},
	"util-linux":    {{"kernel", []string{"util-linux"}}},
	"coreutils":     {{"gnu", []string{"coreutils"}}},
	"findutils":     {{"gnu", []string{"findutils"}}},
	"diffutils":     {{"gnu", []string{"diffutils"}}},
	"patch":         {{"gnu", []string{"patch"}}},
	"make":          {{"gnu", []string{"make"}}},
	"autoconf":      {{"gnu", []string{"autoconf"}}},
	"automake":      {{"gnu", []string{"automake"}}},
	"libtool":       {{"gnu", []string{"libtool"}}},
	"m4":            {{"gnu", []string{"m4"}}},
	"bison":         {{"gnu", []string{"bison"}}},
	"flex":          {{"flex", []string{"flex"}}},
	"cmake":         {{"cmake", []string{"cmake"}}},
	"meson":         {{"meson", []string{"meson"}}},
	"ninja":         {{"ninja", []string{"ninja"}}},
	"bazel":         {{"bazel", []string{"bazel"}}},
	"maven":         {{"apache", []string{"maven"}}},
	"ant":           {{"apache", []string{"ant"}}},
	"gradle":        {{"gradle", []string{"gradle"}}},
	"pip":           {{"python", []string{"python"}}},
	"pip3":          {{"python", []string{"python"}}},
	"setuptools":    {{"python", []string{"python"}}},
	"wheel":         {{"python", []string{"python"}}},
	"cryptography":  {{"python", []string{"python"}}},
	"requests":      {{"python", []string{"python"}}},
	"urllib3":       {{"python", []string{"python"}}},
	"pyyaml":        {{"python", []string{"python"}}},
	"pillow":        {{"python", []string{"python"}}},
	"numpy":         {{"numpy", []string{"numpy"}}},
	"scipy":         {{"scipy", []string{"scipy"}}},
	"scikit-learn":  {{"scikit-learn", []string{"scikit-learn"}}},
	"tensorflow":    {{"google", []string{"tensorflow"}}},
	"pytorch":       {{"pytorch", []string{"pytorch"}}},
	"ansible":       {{"ansible", []string{"ansible"}}},
	"puppet":        {{"puppet", []string{"puppet"}}},
	"chef":          {{"chef", []string{"chef"}}},
	"salt":          {{"saltstack", []string{"salt"}}},
	"terraform":     {{"hashicorp", []string{"terraform"}}},
	"packer":        {{"hashicorp", []string{"packer"}}},
	"vagrant":       {{"hashicorp", []string{"vagrant"}}},
	"vault":         {{"hashicorp", []string{"vault"}}},
	"nomad":         {{"hashicorp", []string{"nomad"}}},
	"boundary":      {{"hashicorp", []string{"boundary"}}},
	"waypoint":      {{"hashicorp", []string{"waypoint"}}},
}

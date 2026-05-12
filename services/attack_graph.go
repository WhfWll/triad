package services

import (
	"context"
	"smart/api/typespec"
	"smart/models/mysqls"
	"smart/tools/enums"
	"strings"
)

type AttackGraph struct{}

func (a *AttackGraph) BuildGraphByVulIds(ctx context.Context, vulIds []int) ([]typespec.GraphDataNodes, []typespec.GraphDataLinks, error) {
	var vulLib mysqls.VulLibraries
	libs, err := vulLib.AllVulLibrariesForIds(ctx, vulIds)
	if err != nil {
		return nil, nil, err
	}
	if len(libs) > 1000 {
		libs = libs[:1000]
	}

	base := []string{
		enums.ScriptNamePortScan,
		enums.ScriptNameFingerPrint,
		enums.ScriptNameSecondDirBrute,
		enums.ScriptNameCrawlerx,
		//enums.ScriptNameCdnDetect,
		//enums.ScriptNameWafDetect,
		//enums.ScriptNameWhois,
	}

	names := make([]string, 0)
	typeMap := make(map[string]string)
	for _, b := range base {
		names = append(names, b)
		if _, ok := typeMap[b]; !ok {
			typeMap[b] = "base"
		}
	}

	seen := make(map[string]struct{})
	uniq := make([]string, 0)
	for _, n := range names {
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		uniq = append(uniq, n)
	}

	nodes := make([]typespec.GraphDataNodes, 0, len(uniq))
	index := make(map[string]int)
	for i, n := range uniq {
		nodes = append(nodes, typespec.GraphDataNodes{
			Id:         int64(i),
			Name:       n,
			Category:   typeMap[n],
			SymbolSize: "40",
		})
		index[n] = i
	}

	links := make([]typespec.GraphDataLinks, 0)
	if i, ok := index[enums.ScriptNamePortScan]; ok {
		if j, ok := index[enums.ScriptNameFingerPrint]; ok {
			links = append(links, typespec.GraphDataLinks{Source: i, Target: j})
		}
	}
	serviceKeywords := []string{"mysql", "ssh", "ftp", "redis", "mongodb", "postgres", "mssql", "rdp", "smb", "telnet", "tomcat", "vnc", "oracle", "ldap", "smtp", "pop3", "imap", "memcached", "http"}
	serviceSet := make(map[string]struct{})
	for _, lib := range libs {
		lower := strings.ToLower(lib.Component + " " + lib.Name)
		for _, kw := range serviceKeywords {
			if strings.Contains(lower, kw) {
				serviceSet[kw] = struct{}{}
				break
			}
		}
	}
	serviceNodes := make([]string, 0)
	for kw := range serviceSet {
		serviceNodes = append(serviceNodes, kw)
	}
	for _, sn := range serviceNodes {
		if _, ok := index[sn]; ok {
			continue
		}
		typeMap[sn] = "service"
		nodes = append(nodes, typespec.GraphDataNodes{
			Id:         int64(len(nodes)),
			Name:       sn,
			Category:   typeMap[sn],
			SymbolSize: "40",
		})
		index[sn] = len(nodes) - 1
	}
	if i, ok := index[enums.ScriptNameFingerPrint]; ok {
		for _, sn := range serviceNodes {
			if j, ok := index[sn]; ok {
				links = append(links, typespec.GraphDataLinks{Source: i, Target: j})
			}
		}
		if j, ok := index[enums.ScriptNameSecondDirBrute]; ok {
			links = append(links, typespec.GraphDataLinks{Source: i, Target: j})
		}
	}
	if i, ok := index[enums.ScriptNameSecondDirBrute]; ok {
		if k, ok := index[enums.ScriptNameCrawlerx]; ok {
			links = append(links, typespec.GraphDataLinks{Source: i, Target: k})
		}
		typeEnum := enums.AllTypeEnumMap()
		specialList := []string{"命令注入", "表单爆破检测", "目录索引", "跨站脚本攻击", "未加密的链接", "readme信息泄露", "XPath注入", "基于浏览器的表单爆破检测", "MinIO信息泄露漏洞", "http到https的不安全重定向", "明文密码传输", "SQL注入", "git信息泄露", "HttpOnly属性缺失", "任意文件读取", "备份文件泄露", "PHP哈希碰撞拒绝服务漏洞", "PHPInfo信息泄露", "Apache Shiro默认密钥漏洞", "目录遍历", "开放重定向", "任意文件上传", "开发文件泄露", ".htaccess敏感文件泄露", "文件包含", "CRLF注入", "xml实体注入", "敏感信息泄露", "不安全的js库", "内容安全策略未实施", "未授权访问", "Access数据库文件泄露", "服务器端请求伪造", "ldap注入", "不安全的TLS配置", "主机头攻击", "敏感配置文件泄露", "代码执行"}
		specialSet := make(map[string]struct{})
		for _, v := range specialList {
			specialSet[strings.ToLower(v)] = struct{}{}
		}
		crawlerIdx, hasCrawler := index[enums.ScriptNameCrawlerx]
		for _, lib := range libs {
			baseSkip := false
			for _, b := range base {
				if strings.EqualFold(lib.Name, b) {
					baseSkip = true
					break
				}
			}
			if baseSkip {
				continue
			}
			lower := strings.ToLower(lib.Component + " " + lib.Name)
			isService := false
			for _, sn := range serviceNodes {
				if strings.Contains(lower, sn) {
					isService = true
					break
				}
			}
			if isService {
				continue
			}
			nodeName := lib.Name
			category := typeEnum[lib.Type]
			if category == "" {
				category = "vul"
			}
			src := i
			if hasCrawler {
				if _, ok := specialSet[strings.ToLower(lib.Name)]; ok {
					src = crawlerIdx
				}
			}
			if j, ok := index[nodeName]; ok {
				links = append(links, typespec.GraphDataLinks{Source: src, Target: j})
			} else {
				typeMap[nodeName] = category
				nodes = append(nodes, typespec.GraphDataNodes{
					Id:         int64(len(nodes)),
					Name:       nodeName,
					Category:   category,
					SymbolSize: "40",
				})
				index[nodeName] = len(nodes) - 1
				links = append(links, typespec.GraphDataLinks{Source: src, Target: len(nodes) - 1})
			}
		}
	}

	return nodes, links, nil
}

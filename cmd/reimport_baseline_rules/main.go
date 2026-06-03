package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"gitlabee.4dogs.cn/common/config"
	"gitlabee.4dogs.cn/common/mysql"

	"smart/models/mysqls"
	"smart/services"
)

type fileRule struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Category        int      `json:"category"`
	Risk            int      `json:"risk"`
	OSType          int      `json:"osType"`
	Commands        []string `json:"commands"`
	ExpectedValue   string   `json:"expectedValue"`
	MatchType       string   `json:"matchType"`
	FixSuggestion   string   `json:"fixSuggestion"`
	RiskDescription string   `json:"riskDescription"`
}

func main() {
	var (
		rulesDir  = flag.String("dir", "data/baseline", "baseline JSON directory")
		sources   = flag.String("sources", "compliance_rules_zh.json,custom_rules.json", "comma-separated JSON files to import")
		dryRun    = flag.Bool("dry-run", false, "only print counts, do not write DB")
		confirm   = flag.Bool("yes", false, "truncate host_baseline_rule and reimport")
		configArg = flag.String("config", "", "config.json path")
	)
	flag.Parse()

	cfgPath := resolveConfig(*configArg)
	if cfgPath == "" {
		fmt.Println("找不到 config.json，请用 -config 指定")
		os.Exit(1)
	}
	if err := config.NewConfig(cfgPath); err != nil {
		fmt.Printf("load config: %v\n", err)
		os.Exit(1)
	}
	mysql.Setup()

	rules, skippedPlaceholder, skippedEmpty, err := loadRules(*rulesDir, strings.Split(*sources, ","))
	if err != nil {
		fmt.Printf("load rules: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("loaded %d rules (skipped placeholder=%d empty=%d)\n", len(rules), skippedPlaceholder, skippedEmpty)
	if len(rules) < 1200 {
		log.Warnf("rule count %d is below product requirement 1200+", len(rules))
	}

	if *dryRun {
		printByOS(rules)
		return
	}
	if !*confirm {
		fmt.Println("add -yes to truncate host_baseline_rule and import (use -dry-run first)")
		os.Exit(1)
	}

	ctx := mysql.NewContext(context.Background(), mysql.GetDB())
	var model mysqls.HostBaselineRule
	if err := model.DeleteAll(ctx); err != nil {
		fmt.Printf("delete all rules: %v\n", err)
		os.Exit(1)
	}

	code := 0
	imported := 0
	seen := make(map[string]bool)
	for _, r := range rules {
		key := fmt.Sprintf("%s|%d|%d", r.Name, r.Category, r.OSType)
		if seen[key] {
			continue
		}
		seen[key] = true
		code++
		cmdsJSON, _ := json.Marshal(r.Commands)
		mt := r.MatchType
		if mt == "" {
			mt = "contains"
		}
		row := &mysqls.HostBaselineRule{
			RuleCode:        code,
			Name:            r.Name,
			Description:     r.Description,
			Category:        r.Category,
			Risk:            r.Risk,
			OSType:          r.OSType,
			CommandsJSON:    string(cmdsJSON),
			ExpectedValue:   r.ExpectedValue,
			MatchType:       mt,
			FixSuggestion:   r.FixSuggestion,
			RiskDescription: r.RiskDescription,
			Enabled:         1,
		}
		if err := model.Create(ctx, row); err != nil {
			fmt.Printf("insert failed name=%q err=%v\n", r.Name, err)
			continue
		}
		imported++
	}

	if err := services.ReloadBaselineRulesFromDB(ctx); err != nil {
		fmt.Printf("reload engine warning: %v\n", err)
	}
	fmt.Printf("done: imported %d rules into host_baseline_rule\n", imported)
}

func resolveConfig(explicit string) string {
	if explicit != "" {
		return explicit
	}
	for _, p := range []string{"config.json", "../config.json", `D:\goproject\triad\config.json`} {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func loadRules(dir string, files []string) ([]fileRule, int, int, error) {
	var all []fileRule
	skippedPH, skippedEmpty := 0, 0
	for _, name := range files {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("read %s: %w", path, err)
		}
		var chunk []fileRule
		if err := json.Unmarshal(data, &chunk); err != nil {
			return nil, 0, 0, fmt.Errorf("parse %s: %w", path, err)
		}
		for _, r := range chunk {
			if r.Name == "" || len(r.Commands) == 0 {
				skippedEmpty++
				continue
			}
			if services.RuleHasUnresolvedPlaceholder(r.Commands, r.ExpectedValue) {
				skippedPH++
				continue
			}
			if r.OSType < 1 || r.OSType > 4 {
				r.OSType = 1
			}
			all = append(all, r)
		}
		fmt.Printf("  %s: %d rows\n", name, len(chunk))
	}
	return all, skippedPH, skippedEmpty, nil
}

func printByOS(rules []fileRule) {
	counts := map[int]int{}
	for _, r := range rules {
		counts[r.OSType]++
	}
	fmt.Println("by osType:")
	for osType, n := range counts {
		fmt.Printf("  osType=%d: %d\n", osType, n)
	}
}

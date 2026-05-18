package services

import (
	"context"
	"fmt"
	"smart/tools/enums"
	"strings"
	"sync"
	"time"
)

type HostVulnScanner struct {
	connManager *HostConnManager
	cveDB       *CveDB
}

type VulnScanTask struct {
	TaskID        int
	TargetID      int
	Host          string
	Port          int
	Username      string
	Password      string
	Key           string
	OSType        int
	Transport     int
	WinRMUseHttps bool
}

type VulnScanReport struct {
	TaskID      int
	TargetID    int
	TargetIP    string
	OSType      int
	Packages    int
	MatchedVulns int
	CriticalCount int
	HighCount   int
	MediumCount int
	LowCount    int
	Results     []VulnScanResult
	StartTime   time.Time
	EndTime     time.Time
}

type VulnScanResult struct {
	PackageName    string `json:"packageName"`
	PackageVersion string `json:"packageVersion"`
	Cve            string `json:"cve"`
	Title          string `json:"title"`
	Severity       string `json:"severity"`
	RiskLevel      int    `json:"riskLevel"`
}

type PackageInfo struct {
	Name    string
	Version string
}

var (
	globalHostVulnScanner *HostVulnScanner
	hostVulnOnce          sync.Once
)

func GetHostVulnScanner() *HostVulnScanner {
	hostVulnOnce.Do(func() {
		globalHostVulnScanner = &HostVulnScanner{
			connManager: GetHostConnManager(),
			cveDB:       GetCveDB(),
		}
	})
	return globalHostVulnScanner
}

func (s *HostVulnScanner) RunVulnScan(ctx context.Context, task *VulnScanTask) (*VulnScanReport, error) {
	report := &VulnScanReport{
		TaskID:    task.TaskID,
		TargetID:  task.TargetID,
		TargetIP:  task.Host,
		OSType:    task.OSType,
		StartTime: time.Now(),
	}

	if s.cveDB == nil || !s.cveDB.IsAvailable() {
		return nil, fmt.Errorf("CVE database not available")
	}

	connConfig := &HostConnConfig{
		Host:       task.Host,
		Port:       task.Port,
		Username:   task.Username,
		Password:   task.Password,
		PrivateKey: task.Key,
		OSType:     task.OSType,
		Transport:  task.Transport,
		UseHTTPS:   task.WinRMUseHttps,
		Timeout:    30 * time.Second,
	}

	conn, err := s.connManager.GetConnection(ctx, connConfig)
	if err != nil {
		return nil, fmt.Errorf("connect to host %s failed: %v", task.Host, err)
	}

	packages, err := s.collectPackages(ctx, conn, task.OSType)
	if err != nil {
		return nil, fmt.Errorf("collect packages from %s failed: %v", task.Host, err)
	}

	report.Packages = len(packages)

	for _, pkg := range packages {
		select {
		case <-ctx.Done():
			return report, ctx.Err()
		default:
		}

		cves, err := s.cveDB.QueryByProduct(pkg.Name)
		if err != nil || len(cves) == 0 {
			continue
		}

		for _, cve := range cves {
			matched, _ := s.cveDB.MatchCpe(pkg.Version, cve.CpeConfigurations)
			if !matched {
				continue
			}

			result := VulnScanResult{
				PackageName:    pkg.Name,
				PackageVersion: pkg.Version,
				Cve:            cve.Cve,
				Title:          cve.TitleZh,
				Severity:       cve.Severity,
				RiskLevel:      cve.SeverityLevel,
			}
			report.Results = append(report.Results, result)
			report.MatchedVulns++

			switch cve.SeverityLevel {
			case 4:
				report.CriticalCount++
			case 3:
				report.HighCount++
			case 2:
				report.MediumCount++
			case 1:
				report.LowCount++
			}
		}
	}

	report.EndTime = time.Now()
	return report, nil
}

func (s *HostVulnScanner) collectPackages(ctx context.Context, conn *HostConnection, osType int) ([]PackageInfo, error) {
	switch osType {
	case enums.BaselineOSTypeLinux, enums.BaselineOSTypeDomestic:
		return s.collectLinuxPackages(ctx, conn)
	case enums.BaselineOSTypeWindows:
		return s.collectWindowsPackages(ctx, conn)
	default:
		return s.collectLinuxPackages(ctx, conn)
	}
}

func (s *HostVulnScanner) collectLinuxPackages(ctx context.Context, conn *HostConnection) ([]PackageInfo, error) {
	output, err := s.connManager.ExecuteCommand(ctx, conn, "dpkg -l 2>/dev/null || rpm -qa --queryformat '%{NAME} %{VERSION}\\n' 2>/dev/null")
	if err != nil {
		return nil, err
	}

	var packages []PackageInfo
	seen := make(map[string]bool)

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		if fields[0] == "dpkg-query:" || fields[0] == "Desired=" || fields[0] == "|" {
			continue
		}

		if strings.HasPrefix(fields[0], "ii") && len(fields) >= 3 {
			name := fields[1]
			version := fields[2]
			version = cleanVersion(version)
			key := name + "@" + version
			if !seen[key] {
				packages = append(packages, PackageInfo{Name: name, Version: version})
				seen[key] = true
			}
		} else if !strings.HasPrefix(fields[0], "ii") && !strings.HasPrefix(fields[0], "un") {
			name := fields[0]
			version := fields[1]
			version = cleanVersion(version)
			key := name + "@" + version
			if !seen[key] {
				packages = append(packages, PackageInfo{Name: name, Version: version})
				seen[key] = true
			}
		}
	}

	return packages, nil
}

func (s *HostVulnScanner) collectWindowsPackages(ctx context.Context, conn *HostConnection) ([]PackageInfo, error) {
	output, err := s.connManager.ExecuteCommand(ctx, conn, "Get-WmiObject -Class Win32_Product | Select-Object Name,Version")
	if err != nil {
		return nil, err
	}

	var packages []PackageInfo
	seen := make(map[string]bool)

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "Name") || strings.Contains(line, "---") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 2 {
			version := fields[len(fields)-1]
			name := strings.Join(fields[:len(fields)-1], " ")
			version = cleanVersion(version)
			key := name + "@" + version
			if !seen[key] {
				packages = append(packages, PackageInfo{Name: name, Version: version})
				seen[key] = true
			}
		}
	}

	return packages, nil
}

func cleanVersion(version string) string {
	version = strings.TrimLeft(version, "=:>~ ")
	idx := strings.IndexAny(version, "+-~")
	if idx > 0 {
		version = version[:idx]
	}
	idx = strings.Index(version, "ubuntu")
	if idx > 0 {
		version = version[:idx]
	}
	idx = strings.Index(version, "debian")
	if idx > 0 {
		version = version[:idx]
	}
	idx = strings.Index(version, "el")
	if idx > 0 {
		version = version[:idx]
	}
	return strings.TrimRight(version, ". ")
}
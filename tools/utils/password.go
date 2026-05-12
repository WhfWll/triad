package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

type PasswordConfig struct {
	MinLength int
	MaxLength int
	UseLower  bool
	UseUpper  bool
	UseDigits bool
	UseSymbol bool
}

// 默认配置：8-13位随机长度，包含所有字符类型
var DefaultConfig = PasswordConfig{
	MinLength: 8,
	MaxLength: 13,
	UseLower:  true,
	UseUpper:  true,
	UseDigits: true,
	UseSymbol: false,
}

// 生成随机密码
func GeneratePassword(config PasswordConfig) (string, error) {
	// 验证配置
	if config.MinLength < 1 {
		return "", fmt.Errorf("最小长度不能小于1")
	}
	if config.MaxLength < config.MinLength {
		return "", fmt.Errorf("最大长度不能小于最小长度")
	}
	if config.MinLength < 8 || config.MaxLength > 13 {
		return "", fmt.Errorf("密码长度范围应在8-13位之间")
	}

	// 随机选择密码长度
	length, err := randomInt(config.MinLength, config.MaxLength)
	if err != nil {
		return "", err
	}

	// 构建字符池
	charPool := ""
	if config.UseLower {
		charPool += "abcdefghijklmnopqrstuvwxyz"
	}
	if config.UseUpper {
		charPool += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	if config.UseDigits {
		charPool += "0123456789"
	}
	if config.UseSymbol {
		charPool += "!@#$%^&*()_+-=[]{}|;:,.<>?"
	}

	if charPool == "" {
		return "", fmt.Errorf("必须至少选择一种字符类型")
	}

	// 生成密码
	password := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charPool))))
		if err != nil {
			return "", err
		}
		password[i] = charPool[n.Int64()]
	}

	// 确保密码包含所有选择的字符类型
	return ensureCharacterVariety(string(password), config), nil
}

// 确保密码包含所有选择的字符类型
func ensureCharacterVariety(password string, config PasswordConfig) string {
	// 检查是否已包含所有类型的字符
	hasLower := !config.UseLower
	hasUpper := !config.UseUpper
	hasDigits := !config.UseDigits
	hasSymbol := !config.UseSymbol

	// 检查当前密码包含的字符类型
	for _, char := range password {
		if config.UseLower && !hasLower && strings.ContainsRune("abcdefghijklmnopqrstuvwxyz", char) {
			hasLower = true
		}
		if config.UseUpper && !hasUpper && strings.ContainsRune("ABCDEFGHIJKLMNOPQRSTUVWXYZ", char) {
			hasUpper = true
		}
		if config.UseDigits && !hasDigits && strings.ContainsRune("0123456789", char) {
			hasDigits = true
		}
		if config.UseSymbol && !hasSymbol && strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char) {
			hasSymbol = true
		}
	}

	// 如果已包含所有类型，直接返回
	if hasLower && hasUpper && hasDigits && hasSymbol {
		return password
	}

	// 构建需要添加的字符类型列表
	neededChars := ""
	if config.UseLower && !hasLower {
		neededChars += "a" // 代表小写字母
	}
	if config.UseUpper && !hasUpper {
		neededChars += "A" // 代表大写字母
	}
	if config.UseDigits && !hasDigits {
		neededChars += "0" // 代表数字
	}
	if config.UseSymbol && !hasSymbol {
		neededChars += "!" // 代表符号
	}

	// 将密码转换为字节切片以便修改
	pwdBytes := []byte(password)

	// 随机替换位置以添加缺失的字符类型
	for _, charType := range neededChars {
		// 随机选择替换位置
		pos, _ := randomInt(0, len(pwdBytes)-1)

		// 根据字符类型选择具体的字符
		var replacement byte
		switch charType {
		case 'a': // 小写字母
			replacement = "abcdefghijklmnopqrstuvwxyz"[randomIntSimple(26)]
		case 'A': // 大写字母
			replacement = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"[randomIntSimple(26)]
		case '0': // 数字
			replacement = "0123456789"[randomIntSimple(10)]
		case '!': // 符号
			replacement = "!@#$%^&*()_+-=[]{}|;:,.<>?"[randomIntSimple(24)]
		}

		pwdBytes[pos] = replacement
	}

	return string(pwdBytes)
}

// 生成指定范围内的随机整数
func randomInt(min, max int) (int, error) {
	if min == max {
		return min, nil
	}

	diff := big.NewInt(int64(max - min + 1))
	n, err := rand.Int(rand.Reader, diff)
	if err != nil {
		return 0, err
	}
	return min + int(n.Int64()), nil
}

// 简单的随机整数生成（用于非加密安全场景）
func randomIntSimple(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}

// 批量生成密码
func GenerateMultiplePasswords(count int, config PasswordConfig) ([]string, error) {
	passwords := make([]string, 0, count)

	for i := 0; i < count; i++ {
		pwd, err := GeneratePassword(config)
		if err != nil {
			return nil, err
		}
		passwords = append(passwords, pwd)
	}

	return passwords, nil
}

/*
func main() {
	// 命令行参数解析
	minLength := flag.Int("min", 8, "最小密码长度")
	maxLength := flag.Int("max", 13, "最大密码长度")
	count := flag.Int("n", 1, "生成密码数量")
	noLower := flag.Bool("no-lower", false, "不使用小写字母")
	noUpper := flag.Bool("no-upper", false, "不使用大写字母")
	noDigits := flag.Bool("no-digits", false, "不使用数字")
	noSymbols := flag.Bool("no-symbols", false, "不使用特殊符号")

	flag.Parse()

	// 创建配置
	config := PasswordConfig{
		MinLength: *minLength,
		MaxLength: *maxLength,
		UseLower:  !*noLower,
		UseUpper:  !*noUpper,
		UseDigits: !*noDigits,
		UseSymbol: !*noSymbols,
	}

	// 生成密码
	if *count == 1 {
		password, err := GeneratePassword(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(password)
	} else {
		passwords, err := GenerateMultiplePasswords(*count, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}
		for i, pwd := range passwords {
			fmt.Printf("%d. %s (长度: %d)\n", i+1, pwd, len(pwd))
		}
	}
}*/

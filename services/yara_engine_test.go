package services

import "testing"

func TestMooseRuleDoesNotMatchProcessText(t *testing.T) {
	ruleText := `rule moose
{
    strings:
        $s24 = "sh"
        $s25 = "ps"
        $s27 = "chmod"
    condition:
        is_elf and all of them
}`
	rules := parseYaraContent(ruleText)
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}

	engine := &YaraEngine{rules: rules}
	psOutput := []byte("USER PID %CPU COMMAND\nroot 1 sh /bin/sh\nadmin ps aux | grep moose\nchmod +x script.sh")

	matches := engine.Match(psOutput)
	if len(matches) != 0 {
		t.Fatalf("moose should not match plain text output, got %d matches: %+v", len(matches), matches)
	}
}

func TestMooseRuleRequiresAllStringsOnELF(t *testing.T) {
	ruleText := `rule moose
{
    strings:
        $s24 = "sh"
        $s25 = "ps"
        $s27 = "chmod"
    condition:
        is_elf and all of them
}`
	rules := parseYaraContent(ruleText)
	engine := &YaraEngine{rules: rules}

	elfHeader := []byte{0x7f, 'E', 'L', 'F'}
	partial := append(append([]byte{}, elfHeader...), []byte("... sh only ...")...)
	if len(engine.Match(partial)) != 0 {
		t.Fatal("partial ELF payload should not satisfy all of them")
	}

	full := append(append([]byte{}, elfHeader...), []byte(" sh ps chmod ")...)
	if len(engine.Match(full)) != 1 {
		t.Fatalf("full ELF payload should match once, got %d", len(engine.Match(full)))
	}
}

func TestSimpleStringRuleMatches(t *testing.T) {
	ruleText := `rule demo
{
    strings:
        $a = "hello"
    condition:
        $a
}`
	rules := parseYaraContent(ruleText)
	engine := &YaraEngine{rules: rules}
	matches := engine.Match([]byte("say hello world"))
	if len(matches) != 1 || matches[0].Matched != "hello" {
		t.Fatalf("unexpected matches: %+v", matches)
	}
}

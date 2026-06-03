package services

import "testing"

func TestNonEmptyResultCheck_requirepass(t *testing.T) {
	check := nonEmptyResultCheck()
	cases := []struct {
		actual string
		want   bool
	}{
		{"result=[requirepass redis.4dogs.cn]", true},
		{"result=redis.4dogs.cn", true},
		{"result=[]", false},
		{"result=", false},
		{"", false},
	}
	for _, c := range cases {
		if got := check(c.actual); got != c.want {
			t.Errorf("nonEmptyResultCheck(%q) = %v, want %v", c.actual, got, c.want)
		}
	}
}

func TestRedisBindNotWildcardCheck(t *testing.T) {
	check := redisBindNotWildcardCheck()
	cases := []struct {
		actual string
		want   bool
	}{
		{"result=[bind 0.0.0.0 ::1]", false},
		{"result=[bind 127.0.0.1]", true},
		{"result=[bind 127.0.0.1 ::1]", true},
		{"result=[bind *]", false},
	}
	for _, c := range cases {
		if got := check(c.actual); got != c.want {
			t.Errorf("redisBindNotWildcardCheck(%q) = %v, want %v", c.actual, got, c.want)
		}
	}
}

func TestBuildDBRuleCheckFunc_nonemptyExpected(t *testing.T) {
	check := buildDBRuleCheckFunc("empty", "非空")
	actual := "result=[requirepass redis.4dogs.cn]"
	if !check(actual) {
		t.Fatalf("empty+非空 should pass when password set, got fail for %q", actual)
	}
}

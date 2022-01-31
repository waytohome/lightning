package ginx

import (
	"testing"
)

func TestGetHttpCode(t *testing.T) {
	code := GetHttpCode(CodeAuthFailed)
	if code != 401 {
		t.Fatalf("unexpected result, want = %d, actually = %d;", 401, code)
	}
	t.Logf("get http code = %d", code)
}

package integrations

import (
	"fmt"
	"net/http"
	"testing"
)

func TestFreeSpinsGetSuccess(t *testing.T) {
	expectCode := http.StatusOK

	code, content := SendRequest("GET",
		fmt.Sprintf("core/free_spins?session_token=%v", manager.GetSessionToken()), nil)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestFreeSpinsGetNoSession(t *testing.T) {
	expectCode := http.StatusUnprocessableEntity

	code, content := SendRequest("GET", "core/free_spins", nil)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestFreeSpinsGetWrongSession(t *testing.T) {
	expectCode := http.StatusUnauthorized

	code, content := SendRequest("GET", "core/free_spins?session_token=8d1a3d5c-17db-11ed-861d-0242ac120002", nil)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

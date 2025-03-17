package integrations

import (
	"net/http"
	"testing"
)

func TestWagerSuccess(t *testing.T) {
	req := manager.DefaultWagerRequest()
	expectCode := http.StatusOK

	code, content := SendRequest("POST", WagerPath, req)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestWagerWrongWager(t *testing.T) {
	req := manager.DefaultWagerRequest()
	expectCode := http.StatusUnprocessableEntity

	req.Wager++

	code, content := SendRequest("POST", WagerPath, req)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestWagerNegativeWager(t *testing.T) {
	req := manager.DefaultWagerRequest()
	expectCode := http.StatusUnprocessableEntity

	req.Wager += -1

	code, content := SendRequest("POST", WagerPath, req)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestWagerZeroWager(t *testing.T) {
	req := manager.DefaultWagerRequest()
	expectCode := http.StatusUnprocessableEntity

	req.Wager = 0

	code, content := SendRequest("POST", WagerPath, req)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestWagerNoSession(t *testing.T) {
	req := manager.DefaultWagerRequest()
	expectCode := http.StatusUnprocessableEntity

	req.SessionToken = ""

	code, content := SendRequest("POST", WagerPath, req)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestWagerWrongSession(t *testing.T) {
	req := manager.DefaultWagerRequest()
	expectCode := http.StatusUnauthorized

	req.SessionToken = manager.GetWrongSessionToken()

	code, content := SendRequest("POST", WagerPath, req)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestWagerWrongFreeSpin(t *testing.T) {
	req := manager.DefaultWagerRequest()
	expectCode := http.StatusConflict

	req.FreeSpinID = manager.GetWrongSessionToken()

	code, content := SendRequest("POST", WagerPath, req)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func SendSuccessWagerRequest() {
	code, content := SendRequest("POST", WagerPath, manager.DefaultWagerRequest())
	if code != http.StatusOK {
		smartPanic(string(content))
	}
}

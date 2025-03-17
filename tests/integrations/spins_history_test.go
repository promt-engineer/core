package integrations

import (
	"net/http"
	"testing"
)

func TestSpinsHistorySuccess(t *testing.T) {
	req := manager.DefaultHistoryRequest()
	expectCode := http.StatusOK

	code, content := SendRequest("GET", HistoryPath+req.String(), nil)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestSpinsHistoryNoSession(t *testing.T) {
	req := manager.DefaultHistoryRequest()
	expectCode := http.StatusUnprocessableEntity

	req.SessionToken = ""

	code, content := SendRequest("GET", HistoryPath+req.String(), nil)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestSpinsHistoryWrongSession(t *testing.T) {
	req := manager.DefaultHistoryRequest()
	expectCode := http.StatusUnauthorized

	req.SessionToken = manager.GetWrongSessionToken()

	code, content := SendRequest("GET", HistoryPath+req.String(), nil)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestSpinsHistoryNoPage(t *testing.T) {
	req := manager.DefaultHistoryRequest()
	expectCode := http.StatusUnprocessableEntity

	req.Page = nil

	code, content := SendRequest("GET", HistoryPath+req.String(), nil)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestSpinsHistoryZeroPage(t *testing.T) {
	req := manager.DefaultHistoryRequest()
	expectCode := http.StatusUnprocessableEntity

	zero := 0
	req.Page = &zero

	code, content := SendRequest("GET", HistoryPath+req.String(), nil)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestSpinsHistoryNegativePage(t *testing.T) {
	req := manager.DefaultHistoryRequest()
	expectCode := http.StatusUnprocessableEntity

	minusOne := -1
	req.Page = &minusOne

	code, content := SendRequest("GET", HistoryPath+req.String(), nil)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestSpinsHistoryNoCount(t *testing.T) {
	req := manager.DefaultHistoryRequest()
	expectCode := http.StatusUnprocessableEntity

	req.Count = nil

	code, content := SendRequest("GET", HistoryPath+req.String(), nil)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestSpinsHistoryZeroCount(t *testing.T) {
	req := manager.DefaultHistoryRequest()
	expectCode := http.StatusUnprocessableEntity

	zero := 0
	req.Count = &zero

	code, content := SendRequest("GET", HistoryPath+req.String(), nil)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestSpinsHistoryNegativeCount(t *testing.T) {
	req := manager.DefaultHistoryRequest()
	expectCode := http.StatusUnprocessableEntity

	minusOne := -1
	req.Count = &minusOne

	code, content := SendRequest("GET", HistoryPath+req.String(), nil)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

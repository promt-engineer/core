package integrations

import (
	"net/http"
	"testing"
)

func TestStateSuccess(t *testing.T) {
	req := manager.DefaultStateRequest()
	expectCode := http.StatusOK

	code, content := manager.SendStateRequest(req)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestStateNoGameField(t *testing.T) {
	req := manager.DefaultStateRequest()
	expectCode := http.StatusUnprocessableEntity

	req.Game = ""

	code, content := manager.SendStateRequest(req)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestStateWrongGameField(t *testing.T) {
	req := manager.DefaultStateRequest()
	expectCode := http.StatusUnprocessableEntity

	req.Game = "wrong field"

	code, content := manager.SendStateRequest(req)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestStateNoIntegratorField(t *testing.T) {
	req := manager.DefaultStateRequest()
	expectCode := http.StatusUnprocessableEntity

	req.Integrator = ""

	code, content := manager.SendStateRequest(req)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

func TestStateWrongIntegratorField(t *testing.T) {
	req := manager.DefaultStateRequest()
	expectCode := http.StatusUnprocessableEntity

	req.Integrator = "wrong field"

	code, content := manager.SendStateRequest(req)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

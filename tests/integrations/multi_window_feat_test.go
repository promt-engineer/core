package integrations

import "testing"

func TestMultipleWindowsError(t *testing.T) {
	oldSession := manager.GetSessionToken()
	manager.GenerateNewSessionToken()

	req := manager.DefaultWagerRequest()
	req.SessionToken = oldSession
	expectCode := StatusSessionExpired

	code, content := SendRequest("POST", "core/wager", req)
	if code != expectCode {
		t.Errorf("received status %v, %v expected\n content: %v", code, expectCode, string(content))
	}
}

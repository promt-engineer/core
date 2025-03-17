package integrations

import (
	"net/http"
	"testing"
)

func TestMasterSlotHistory(t *testing.T) {
	expectTotal := 1

	manager.GenerateNewSessionToken()
	masterSession := manager.GetSessionToken()
	masterReq := manager.DefaultWagerRequest()
	masterReq.SessionToken = masterSession

	if code, content := SendRequest("POST", WagerPath, masterReq); code != http.StatusOK {
		smartPanic(string(content))
	}

	masterReqHistory := manager.DefaultHistoryRequest()
	masterReqHistory.SessionToken = masterSession

	masterCode, masterContent := SendRequest("GET", HistoryPath+masterReqHistory.String(), nil)
	if masterCode != http.StatusOK {
		smartPanic(string(masterContent))
	}

	masterResp := manager.ParseHistoryResponse(masterContent)
	if masterResp.Total != expectTotal {
		t.Errorf("history has %v records, %v expected", masterResp.Total, expectTotal)
	}

	reqTMP := manager.DefaultStateRequest()
	reqTMP.Game = config.ReskinGame
	reqTMP.Params.Game = config.ReskinGame
	reskinSession := manager.GetSessionTokenPure(reqTMP)

	reskinReq := manager.DefaultWagerRequest()
	reskinReq.SessionToken = reskinSession

	if code, content := SendRequest("POST", WagerPath, reskinReq); code != http.StatusOK {
		smartPanic(string(content))
	}

	reskinReqHistory := manager.DefaultHistoryRequest()
	reskinReqHistory.SessionToken = reskinSession
	reskinCode, reskinContent := SendRequest("GET", HistoryPath+reskinReqHistory.String(), nil)

	if reskinCode != http.StatusOK {
		smartPanic(string(reskinContent))
	}

	reskinResp := manager.ParseHistoryResponse(reskinContent)
	if reskinResp.Total != expectTotal {
		t.Errorf("history has %v records, %v expected", reskinResp.Total, expectTotal)
	}
}

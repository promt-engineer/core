package engine

import (
	"encoding/json"
	"go.uber.org/zap"
	"reflect"
)

type Features struct {
	//RTP        *int64  `json:"rtp"`
	Volatility string `json:"volatility"`
}

func UnmarshalTo[T any](payload interface{}) (T, error) {
	var b T

	zap.S().Infof("UnmarshalTo called with payload %s of type: %s", payload, reflect.TypeOf(payload))

	bytes, err := json.Marshal(payload)
	if err != nil {
		zap.S().Errorf("json.Marshal error: %v", err)
		return b, err
	}

	if err := json.Unmarshal(bytes, &b); err != nil {
		zap.S().Errorf("json.Unmarshal error: %v", err)
		return b, err
	}

	return b, nil
}

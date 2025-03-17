package cryptolut_rgs

import (
	"bytes"
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"net/http"
)

var httpClient = &http.Client{}

type HTTPErr struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func askServer[Req, Res any](ctx context.Context,
	url string, body Req, parseTo *Res) func(withLog bool) error {
	return func(withLog bool) error {
		bufBody, err := json.Marshal(body)
		if err != nil {
			return err
		}

		if withLog {
			zap.S().Info(string(bufBody))
		}

		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bufBody))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		res, err := io.ReadAll(resp.Body)
		if err != nil {
			if withLog {
				zap.S().Info(err)
			}

			return err
		}

		if withLog {
			zap.S().Info(err, string(res))
		}

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			if err = json.Unmarshal(res, &parseTo); err != nil {
				return err
			}

			return nil
		}

		httpErr := HTTPErr{}
		if err = json.Unmarshal(res, &httpErr); err != nil {
			return err
		}

		var ok bool
		err, ok = errorMap[httpErr.Code]
		if !ok {
			return ErrCLInternalServerError
		}

		return err
	}
}

func buildURL(base, path string) string {
	return base + path
}

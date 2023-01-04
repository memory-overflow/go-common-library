package httpcall

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// JsonPut http put 请求
func JsonPut(ctx context.Context, uri string, header map[string]string, req, rsp interface{}) (err error) {
	body, ok := req.([]byte)
	if !ok {
		body, _ = json.Marshal(req)
	}
	reqCtx, err := http.NewRequestWithContext(ctx, http.MethodPut, uri, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("JsonPut make NewRequest error: %v", err)
	}
	if header == nil {
		reqCtx.Header.Set("Content-Type", "application/octet-stream")
	} else {
		for head, value := range header {
			reqCtx.Header.Set(head, value)
		}
	}
	response, err := http.DefaultClient.Do(reqCtx)
	if err != nil {
		return fmt.Errorf("JsonPut do request error: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("JsonPut response status error, statuscode: %d, reason: %s",
			response.StatusCode, response.Status)
	}
	data, _ := ioutil.ReadAll(response.Body)
	// log.Printf("JsonPut req: %v, rsp: %v \n", string(body), string(data))
	if rsp != nil {
		err = json.Unmarshal(data, rsp)
		if err != nil {
			return fmt.Errorf("JsonPut response protocol error: %v, data: %s", err, string(data))
		}
	}
	return nil
}

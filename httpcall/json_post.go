package httpcall

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// JsonPost http json post 请求
func JsonPost(ctx context.Context, uri string, header map[string]string, req, rsp interface{}) (err error) {
	body, _ := json.Marshal(req)
	reqCtx, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("HttpPost make NewRequest error: %v", err)
	}
	if header == nil {
		reqCtx.Header.Set("Content-Type", "application/json")
	} else {
		for head, value := range header {
			reqCtx.Header.Set(head, value)
		}
	}
	response, err := http.DefaultClient.Do(reqCtx)
	if err != nil {
		return fmt.Errorf("HttpPost do request error: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HttpPost response status error, statuscode: %d, reason: %s",
			response.StatusCode, response.Status)
	}
	data, _ := ioutil.ReadAll(response.Body)
	// log.Printf("JsonPost req: %v, rsp: %v \n", string(body), string(data))
	if rsp != nil {
		err = json.Unmarshal(data, rsp)
		if err != nil {
			return fmt.Errorf("HttpPost response protocol error: %v, data: %s", err, string(data))
		}
	}
	return nil
}

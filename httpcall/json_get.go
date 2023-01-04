package httpcall

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// JsonGet http json get 请求
func JsonGet(ctx context.Context, uri string, header map[string]string, rsp interface{}) (err error) {
	reqCtx, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return fmt.Errorf("JsonGet make NewRequest error: %v", err)
	}

	for head, value := range header {
		reqCtx.Header.Set(head, value)
	}

	response, err := http.DefaultClient.Do(reqCtx)
	if err != nil {
		return fmt.Errorf("JsonGet do request error: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("JsonGet response status error, statuscode: %d, reason: %s",
			response.StatusCode, response.Status)
	}
	data, _ := ioutil.ReadAll(response.Body)
	if rsp != nil {
		err = json.Unmarshal(data, rsp)
		if err != nil {
			return fmt.Errorf("JsonGet response protocol error: %v, data: %s", err, string(data))
		}
	}
	return nil
}

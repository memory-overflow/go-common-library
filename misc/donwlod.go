package misc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

// DownloadFile 下载网络文件到服务器本地
func DownloadFile(ctx context.Context, url string, dirs, fileName string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return errors.New("DownloadFile " + url + " failed! " + err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("DownloadFile " + url + " failed! " + err.Error())
	}
	if res.StatusCode != 200 {
		return errors.New("DownloadFile " + url + " failed! " + res.Status)
	}
	if err := os.MkdirAll(dirs, 0766); err != nil {
		return errors.New("DownloadFile " + url + " failed! " + err.Error())
	}
	filePath := dirs + "/" + fileName
	file, err := os.Create(filePath)
	if err != nil {
		return errors.New("DownloadFile " + url + " create file failed! " + err.Error())
	}
	if length, err := io.Copy(file, res.Body); err != nil {
		os.Remove(filePath)
		return errors.New("DownloadFile " + url + " copy failed! " + err.Error() +
			" copied length: " + strconv.Itoa(int(length)))
	}
	file.Close()
	_, err = os.Stat(filePath)
	if err != nil {
		os.Remove(filePath)
		return errors.New("DownloadFile " + url + " failed! " + err.Error())
	}
	return nil
}

// DownloadFileWithLimit 下载网络文件到服务器本地，限制文件大小
func DownloadFileWithLimit(ctx context.Context, url string, dirs, fileName string, maxBytes int64) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return errors.New("DownloadFile " + url + " failed! " + err.Error())
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("DownloadFile " + url + " failed! " + err.Error())
	}
	if res.StatusCode != 200 {
		return errors.New("DownloadFile " + url + " failed! " + res.Status)
	}
	if err := os.MkdirAll(dirs, 0766); err != nil {
		return errors.New("DownloadFile " + url + " failed! " + err.Error())
	}
	filePath := dirs + "/" + fileName
	file, err := os.Create(filePath)
	if err != nil {
		return errors.New("DownloadFile " + url + " create file failed! " + err.Error())
	}
	if length, err := io.CopyN(file, res.Body, int64(maxBytes)); err != nil && err != io.EOF {
		os.Remove(filePath)
		return errors.New("DownloadFile " + url + " copy failed! " + err.Error() +
			" copied length: " + strconv.Itoa(int(length)))
	}
	file.Close()
	nextByte := make([]byte, 1)
	nRead, _ := io.ReadFull(res.Body, nextByte)
	if nRead > 0 {
		// Yep, there's too much data
		os.Remove(filePath)
		return fmt.Errorf("DownloadFile %s exceed the file size limit %d", url, maxBytes)
	}
	_, err = os.Stat(filePath)
	if err != nil {
		os.Remove(filePath)
		return errors.New("DownloadFile " + url + " failed! " + err.Error())
	}
	return nil
}

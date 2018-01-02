package backend

import (
	"github.com/armson/bingo/utils"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"errors"
	"fmt"
)

func (b *HttpServer) buildUrl() (bool, error) {
	if b.request.Method == "GET" {
		if query := b.Query(); query != "" {
			if strings.Contains(b.url, "?") {
				b.url += "&" + query
			} else {
				b.url = b.url + "?" + query
			}
		}
	}
	url, err := url.Parse(b.url)
	if err != nil {
		return false, err
	}
	b.request.URL = url
	return true, nil
}

func (b *HttpServer) buildBody() (bool, error) {
	if !utils.Slice.In(b.request.Method,[]string{"POST","PUT","PATCH","DELETE"}) {
		return true, nil
	}
	if b.raw == nil && len(b.params) == 0  && len(b.files) == 0  && len(b.readers) == 0 {
		return  true, nil
	}

	if b.raw != nil {
		bs, err := json.Marshal(b.raw)
		if err != nil {
			return false, err
		}
		b.request.Body = ioutil.NopCloser(bytes.NewReader(bs))
		b.request.ContentLength = int64(len(bs))
		b.request.Header.Set("Content-Type", "application/json")
		return true, nil
	}

	if len(b.files) > 0 {
		for formField, fileName := range b.files {
			rc, err := os.Open(fileName)
			if err != nil {
				return false, fmt.Errorf("HttpServer:%v", err)
			}
			b.readers = append(b.readers, Reader{
				formField:formField,
				fileName:fileName,
				rc:rc,
			})
		}
		b.files = map[string]string{}
	}

	if len(b.readers) > 0 {
		buf := new(bytes.Buffer)
		bodyWriter := multipart.NewWriter(buf)

		for _, rd := range b.readers {
			fileWriter, err := bodyWriter.CreateFormFile(rd.formField, rd.fileName)
			if err != nil {
				fmt.Errorf("HttpServer:%v", err)
			}
			//iocopy
			_, err = io.Copy(fileWriter, rd.rc)
			rd.rc.Close()
			if err != nil {
				fmt.Errorf("HttpServer:%v", err)
			}
		}
		for k, v := range b.params {
			for _, vv := range v {
				bodyWriter.WriteField(k, vv)
			}
		}
		bodyWriter.Close()
		b.Header("Content-Type", bodyWriter.FormDataContentType())
		b.request.ContentLength = int64(buf.Len())
		b.request.Body = ioutil.NopCloser(buf)
		return true, nil
	}

	if len(b.params) > 0 {
		b.Header("Content-Type", "application/x-www-form-urlencoded")
		query := b.Query();
		bf := bytes.NewBufferString(query)
		b.request.Body = ioutil.NopCloser(bf)
		b.request.ContentLength = int64(len(query))
		return true, nil
	}
	return false, errors.New("HttpServer is break!")
}


func (b *HttpServer) ci() (response *http.Response, err error) {
	if _, err = b.buildUrl(); err != nil {
		return nil , err
	}

	if _, err = b.buildBody(); err != nil {
		return nil , err
	}

	trans := b.setting.Transport
	if trans == nil {
		// create default transport
		trans = &http.Transport{
			TLSClientConfig:     b.setting.TLSClientConfig,
			Proxy:               b.setting.Proxy,
			Dial:                TimeoutDialer(b.setting.ConnectTimeout, b.setting.ReadWriteTimeout),
			MaxIdleConnsPerHost: -1,
		}
	} else {
		// if b.transport is *http.Transport then set the settings.
		if t, ok := trans.(*http.Transport); ok {
			if t.TLSClientConfig == nil {
				t.TLSClientConfig = b.setting.TLSClientConfig
			}
			if t.Proxy == nil {
				t.Proxy = b.setting.Proxy
			}
			if t.Dial == nil {
				t.Dial = TimeoutDialer(b.setting.ConnectTimeout, b.setting.ReadWriteTimeout)
			}
		}
	}

	var jar http.CookieJar
	if b.setting.EnableCookie {
		if defaultCookieJar == nil {
			createDefaultCookie()
		}
		jar = defaultCookieJar
	}
	client := &http.Client{
		Transport: trans,
		Jar:       jar,
	}

	if b.setting.UserAgent != "" && b.request.Header.Get("User-Agent") == "" {
		b.request.Header.Set("User-Agent", b.setting.UserAgent)
	}

	if b.setting.CheckRedirect != nil {
		client.CheckRedirect = b.setting.CheckRedirect
	}

	start := time.Now()
	for i := 0; b.setting.Retries == -1 || i <= b.setting.Retries; i++ {
		response, err = client.Do(b.request)
		if err == nil {
			break
		}
	}
	b.cost = time.Since(start)
	return response, err
}
package anchor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

var ExteriorUri = `http://dockerpre-innerapi-yhcpcv.cupid.autohome.com.cn/v1/exteriorAnchors/recogFromImage`

func Send(ctx context.Context, uri string, bs []byte, filename string) (r Resp, err error) {
	body := &bytes.Buffer{}
	multi := multipart.NewWriter(body)
	fieldWriter, err := multi.CreateFormField("req_id")
	if err != nil {
		fmt.Printf("CreateFormField req_id error:%s", err.Error())
		return
	}

	if _, err = fieldWriter.Write([]byte(ctx.Value("traceId").(string))); err != nil {
		fmt.Printf("fieldWriter.Write traceId error:%s", err.Error())
		return
	}

	fieldWriter, _ = multi.CreateFormField("service_type")
	if _, err = fieldWriter.Write([]byte(`panoadmin`)); err != nil {
		fmt.Printf("fieldWriter.Write panoadmin error:%s", err.Error())
		return
	}

	fieldWriter, _ = multi.CreateFormField("service_key")
	if _, err = fieldWriter.Write([]byte(`LVHGSNV2LKA4RS72PTJA5J5A7ERZ5E6O`)); err != nil {
		fmt.Printf("CreateFormField service_key error:%s", err.Error())
		return
	}

	fieldWriter, _ = multi.CreateFormField("timestamp")
	if _, err = fieldWriter.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10))); err != nil {
		fmt.Printf("fieldWriter.Write timestamp error:%s", err.Error())
		return
	}

	fileWriter, _ := multi.CreateFormFile("image_file", filename)
	if _, err = fileWriter.Write(bs); err != nil {
		fmt.Printf("multi.CreateFormFile image_file error:%s", err.Error())
		return
	}

	if err = multi.Close(); err != nil {
		fmt.Printf("multi.Close() error:%s", err.Error())
		return
	}

	resp, err := http.Post(uri, multi.FormDataContentType(), body)
	if err != nil {
		fmt.Printf("http.Post error:%s", err.Error())
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll error:%s", err.Error())
		return
	}

	err = json.Unmarshal(bs, &r)
	if err != nil {
		fmt.Printf("json.Unmarsha error:%s", err.Error())
		return
	}

	return r, err
}

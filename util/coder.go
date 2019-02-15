package util

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
)

func GbkToUtf8(s string) string {
	var str string
	reader := transform.NewReader(bytes.NewReader([]byte(s)), simplifiedchinese.GB18030.NewDecoder())
	b, e := ioutil.ReadAll(reader)
	if e != nil {
		str = "{}"
	} else {
		str = string(b)
	}
	return str
}

func Utf8ToGbk(s string) string {
	var str string
	reader := transform.NewReader(bytes.NewReader([]byte(s)), simplifiedchinese.GB18030.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		str = "{}"
	} else {
		str = string(d)
	}
	return str
}

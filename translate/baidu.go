package translate

import (
	"NovelWeb/net"
	"bytes"
	"encoding/json"
	"github.com/jinzhongmin/gtra"
	"github.com/robertkrimen/otto"
	"log"
	"net/url"
	"os"
)

type BaiDu struct {
	url string
}

type BaiDuResult struct {
	TransResult struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Domain string `json:"domain"`
		Type   int    `json:"type"`
		Status int    `json:"status"`
		Data   []struct {
			Dst        string          `json:"dst"`
			PrefixWrap int             `json:"prefixWrap"`
			Src        string          `json:"src"`
			Relation   []interface{}   `json:"relation"`
			Result     [][]interface{} `json:"result"`
		} `json:"data"`
		Keywords []struct {
			Means []string `json:"means"`
			Word  string   `json:"word"`
		} `json:"keywords"`
	} `json:"trans_result"`
	DictResult []interface{} `json:"dict_result"`
	LijuResult struct {
		Double string `json:"double"`
		Single string `json:"single"`
	} `json:"liju_result"`
	Logid int64 `json:"logid"`
}

func NewBaiDu() BaiDu {
	return BaiDu{
		url: "https://fanyi.baidu.com/v2transapi",
	}
}

func (baidu BaiDu) TranslateGoogle(source string) string {

	t := gtra.NewTranslater()
	_, content := t.Translate(source)

	log.Print("[翻译]", source, "  ", content)
	return content
}

// 翻译，支持超过翻译字数限制
func (baidu BaiDu) Translate(source string) string {
	var content bytes.Buffer

	for i := 0; i < len(source); i += 1000 {
		var endPoint int
		if i+1000 < len(source) {
			endPoint = i + 1000
		} else {
			endPoint = len(source)
		}

		var part string
		if endPoint == len(source)-1 {
			part = baidu.TranslateLimit(source[i:])
		} else {
			part = baidu.TranslateLimit(source[i:endPoint])
		}
		content.WriteString(part)
	}
	return content.String()
}

func (bd BaiDu) TranslateLimit(source string) string {
	var api = "https://fanyi.baidu.com/v2transapi"

	sign := bd.Sign(source)
	values := url.Values{
		"query":             {source},
		"from":              {"zh"},
		"to":                {"en"},
		"transtype":         {"translang"},
		"sign":              {sign},
		"token":             {"7b35624ba7fe34e692ea909140d9582d"},
		"simple_means_flag": {"3"},
		"version":           {"2.1"},
		"keyfrom":           {"fanyi.web"},
		"action":            {"FY_BY_CLICKBUTTION"},
		"typoResult":        {"typoResult"}}

	var headers = make(map[string]string)
	headers["Cookie"] = "BAIDUID=16FFA1EAAF1A387C647A22DB9FC81DAE:FG=1; BIDUPSID=16FFA1EAAF1A387C647A22DB9FC81DAE; PSTM=1514024118; __cfduid=d2f7fd3a024d1ee90b8a817ddd866d9bc1514812370; REALTIME_TRANS_SWITCH=1; FANYI_WORD_SWITCH=1; HISTORY_SWITCH=1; SOUND_SPD_SWITCH=1; SOUND_PREFER_SWITCH=1; BDUSS=E83bFplYnBRen5hV3FKNDZCM3ZSZldWMVBjV1ZGbW9uTTVwcjFYVDQzTGVFTlJhQVFBQUFBJCQAAAAAAAAAAAEAAADnIjgkcG9obG9ndTQxNTA3AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAN6DrFreg6xaS; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; H_PS_PSSID=1429_21103_18559_22075; PSINO=3; locale=zh; Hm_lvt_64ecd82404c51e03dc91cb9e8c025574=1525654866,1525657996,1525658015,1525671031; Hm_lpvt_64ecd82404c51e03dc91cb9e8c025574=1525671031; to_lang_often=%5B%7B%22value%22%3A%22en%22%2C%22text%22%3A%22%u82F1%u8BED%22%7D%2C%7B%22value%22%3A%22zh%22%2C%22text%22%3A%22%u4E2D%u6587%22%7D%5D; from_lang_often=%5B%7B%22value%22%3A%22zh%22%2C%22text%22%3A%22%u4E2D%u6587%22%7D%2C%7B%22value%22%3A%22en%22%2C%22text%22%3A%22%u82F1%u8BED%22%7D%5D"

	var back = net.Get(api+"?"+values.Encode(), headers)
	var result BaiDuResult
	err := json.Unmarshal([]byte(back), &result)
	if err != nil {
		log.Print(err,"  ",back)
	}

	return result.TransResult.Data[0].Dst
}

func (bd BaiDu) Sign(r string) string {
	var gtk = "320305.131321201"

	f, err := os.Open("./tk/Baidu.js")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(f); err != nil {
		panic(err)
	}

	runtime := otto.New()
	if _, err := runtime.Run(buf.String()); err != nil {
		panic(err)
	}

	result, err := runtime.Call("token", nil, r, gtk)
	if err != nil {
		panic(err)
	}
	value, _ := result.ToString()

	return value
}

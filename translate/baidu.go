package translate

import (
	pb "NovelWeb/baidu"
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-resty/resty"
	"google.golang.org/grpc"
	"log"
	"regexp"
	"time"
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
}

func NewBaiDu() BaiDu {
	return BaiDu{
		url: "https://fanyi.baidu.com/v2transapi",
	}
}

// 翻译，支持超过翻译字数限制
func (baidu BaiDu) Translate(source string) string {
	var content bytes.Buffer

	re := regexp.MustCompile("[^\\s]+")
	paragraphs := re.FindAllStringSubmatch(source, -1)

	for _, value := range paragraphs {
		part := baidu.TranslateLimit(value[0], 0)
		content.WriteString(part)
		content.WriteString("\n\n")
	}
	return content.String()
}

func (bd BaiDu) TranslateLimit(source string, t int32) string {
	var api = "https://fanyi.baidu.com/v2transapi"

	sign := bd.sigh(source)
	var headers = make(map[string]string)
	headers["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.139 Safari/537.36"
	headers["Cookie"] = "BAIDUID=16FFA1EAAF1A387C647A22DB9FC81DAE:FG=1; BIDUPSID=16FFA1EAAF1A387C647A22DB9FC81DAE; PSTM=1514024118; __cfduid=d2f7fd3a024d1ee90b8a817ddd866d9bc1514812370; REALTIME_TRANS_SWITCH=1; FANYI_WORD_SWITCH=1; HISTORY_SWITCH=1; SOUND_SPD_SWITCH=1; SOUND_PREFER_SWITCH=1; BDUSS=E83bFplYnBRen5hV3FKNDZCM3ZSZldWMVBjV1ZGbW9uTTVwcjFYVDQzTGVFTlJhQVFBQUFBJCQAAAAAAAAAAAEAAADnIjgkcG9obG9ndTQxNTA3AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAN6DrFreg6xaS; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; H_PS_PSSID=1429_21103_18559_22075; PSINO=3; locale=zh; Hm_lvt_64ecd82404c51e03dc91cb9e8c025574=1525654866,1525657996,1525658015,1525671031; Hm_lpvt_64ecd82404c51e03dc91cb9e8c025574=1525671031; to_lang_often=%5B%7B%22value%22%3A%22en%22%2C%22text%22%3A%22%u82F1%u8BED%22%7D%2C%7B%22value%22%3A%22zh%22%2C%22text%22%3A%22%u4E2D%u6587%22%7D%5D; from_lang_often=%5B%7B%22value%22%3A%22zh%22%2C%22text%22%3A%22%u4E2D%u6587%22%7D%2C%7B%22value%22%3A%22en%22%2C%22text%22%3A%22%u82F1%u8BED%22%7D%5D"

	resp, _ := resty.R().
		SetFormData(map[string]string{
			"query":             source,
			"from":              "zh",
			"to":                "en",
			"transtype":         "translang",
			"sign":              sign,
			"token":             "7b35624ba7fe34e692ea909140d9582d",
			"simple_means_flag": "3",
		}).SetHeaders(headers).
		Post(api)

	var result BaiDuResult
	err := json.Unmarshal(resp.Body(), &result)

	var resultTrans string
	if err != nil {
		log.Print(err)
		if t > 2 {
			resultTrans = ""
		} else {
			return bd.TranslateLimit(source, t+1)
		}
	} else {
		resultTrans = result.TransResult.Data[0].Dst
	}
	log.Print("百度翻译", source, " == ", resultTrans, string(resp.Body()))
	return resultTrans
}

func (bd BaiDu) sigh(text string) string {
	gtk := "320305.131321201"

	// Set up a connection to the server.
	conn, err := grpc.Dial("192.168.40.116:8088", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewTokenServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	r, err := c.GenerateToken(ctx, &pb.InputParams{Text: text, Gtk: gtk})

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	return r.Token
}

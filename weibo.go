package weibo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Handler struct {
	ua     string
	cookie string
}

func New(cookie string) *Handler {
	return &Handler{
		ua:     "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
		cookie: cookie,
	}

}

type SuperTopicsResult struct {
	Data struct {
		Name        string        `json:"name"`
		List        []*SuperTopic `json:"list"`
		TotalNumber int           `json:"total_number"`
		Pagetype    string        `json:"pageType"`
		Fixeddata   []struct {
			Name     string        `json:"name"`
			List     []interface{} `json:"list"`
			Pagetype string        `json:"pageType"`
		} `json:"fixedData"`
	} `json:"data"`
	Ok int `json:"ok"`
}

type SuperTopic struct {
	Title        string `json:"title"`
	Content1     string `json:"content1"`
	Content2     string `json:"content2"`
	Link         string `json:"link"`
	Linktype     string `json:"linkType"`
	Pic          string `json:"pic"`
	Oid          string `json:"oid"`
	TopicName    string `json:"topic_name"`
	LatestStatus string `json:"latest_status"`
	Scheme       string `json:"scheme"`
	StatusCount  int    `json:"status_count"`
	FollowCount  int    `json:"follow_count"`
	Intro        string `json:"intro"`
}

func (h *Handler) SuperTopics(ctx context.Context) (*SuperTopicsResult, error) {
	urlStr := "https://weibo.com/ajax/profile/topicContent?tabid=231093_-_chaohua&page={}"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cookie", h.cookie)
	req.Header.Set("User-Agent", h.ua)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := &SuperTopicsResult{}
	if err := json.Unmarshal(data, result); err != nil {
		return nil, err
	}
	if result.Ok != 1 {
		return nil, fmt.Errorf("weibo: unexpected super topics result, %s", string(data))
	}

	return result, nil
}

func (h *Handler) SuperTopicSignIn(ctx context.Context, id string) error {
	urlStr := "https://weibo.com/p/aj/general/button?"
	param := url.Values{
		"ajwvr":    {"6"},
		"api":      {"http://i.huati.weibo.com/aj/super/checkin"},
		"texta":    {"签到"},
		"textb":    {"已签到"},
		"status":   {"0"},
		"id":       {id},
		"location": {"page_100808_super_index"},
		"timezone": {"GMT 0800"},
		"lang":     {"zh-cn"},
		"plat":     {"MacIntel"},
		"ua":       {h.ua},
		"screen":   {"2560*1440"},
		"__rnd":    {strconv.FormatInt(time.Now().Unix()*1000, 10)},
	}

	urlStr += param.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Cookie", h.cookie)
	req.Header.Set("User-Agent", h.ua)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	switch result.Code {
	case 100000:
	case 382004: // signed
	default:
		return fmt.Errorf("weibo: super topic sign in error, %s", string(data))
	}
	return nil
}

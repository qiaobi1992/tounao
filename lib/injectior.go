package lib

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"io/ioutil"
	"strings"
	"strconv"
	"github.com/elazarl/goproxy"
	"fmt"
)

var cache = make(map[string]Question)

func Injection(bytes []byte, ctx *goproxy.ProxyCtx) (data []byte) {

	data = bytes
	content := string(bytes)

	//log.Printf("path:%s\ncontent:%s", ctx.Req.URL.Path, content)

	if strings.Contains(content, "roomID") && strings.Contains(content, "quizNum") {
		values, _ := url.ParseQuery(content)
		roomId, _ := strconv.Atoi(values.Get("roomID"))
		cacheKey := fmt.Sprintf("roomID=%s", strconv.Itoa(roomId))
		ctx.UserData = cacheKey
	} else if strings.Contains(content, "quiz") && strings.Contains(content, "options") {
		data = injectQuestionResponse(bytes, ctx)
	} else if strings.Contains(content, "score") && strings.Contains(content, "totalScore") {
		go injectChooseResponse(bytes)
	}

	return
}

func injectChooseResponse(bytes []byte) {

	var resp ChooseResp

	json.Unmarshal(bytes, &resp)

	cacheKey := fmt.Sprintf("roomID=%s", strconv.Itoa(resp.Data.RoomID))

	question := cache[cacheKey]

	if question.Quiz != "" {

		question.Answer = question.Options[resp.Data.Answer-1]
		pushAnswerToCache(question)

		delete(cache, cacheKey)
	}

	return
}

func injectQuestionResponse(bytes []byte, ctx *goproxy.ProxyCtx) (data []byte) {
	var resp QuestionResp
	var origin QuestionResp

	json.Unmarshal(bytes, &resp)
	json.Unmarshal(bytes, &origin)

	cacheKey := ctx.UserData.(string)

	cache[cacheKey] = NewQuestion(origin)

	answer := fetchAnswerFromCache(resp.Data.Quiz)

	if answer != "" {
		for index, option := range resp.Data.Options {
			if option == answer {
				resp.Data.Options[index] = option + "[标答]"
			}
		}
	} else {

		page := search(resp.Data.Quiz)

		for index, option := range resp.Data.Options {
			count := strings.Count(page, option)
			resp.Data.Options[index] = option + "[" + strconv.Itoa(count) + "]"
		}

	}

	log.Println(resp.Data.Quiz)
	for _, item := range resp.Data.Options {
		log.Println(item)
	}

	data, _ = json.Marshal(resp)

	return
}

func search(question string) string {
	req, _ := http.NewRequest("GET", "http://www.baidu.com/s?wd="+url.QueryEscape(question), nil)

	resp, _ := http.DefaultClient.Do(req)

	content, _ := ioutil.ReadAll(resp.Body)

	return string(content)
}

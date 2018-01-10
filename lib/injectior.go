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
	"wangzhe/util"
	"time"
)

var cache = make(map[string]Question)

var tapSwitch = false

func Injection(bytes []byte, ctx *goproxy.ProxyCtx) (data []byte) {

	data = bytes
	content := string(bytes)

	//log.Printf("path:%s\ncontent:%s", ctx.Req.URL.Path, content)

	if strings.Contains(content, "roomID") && strings.Contains(content, "quizNum") {
		//请求题目和发送答案的时候停止点击
		tapSwitch = false
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

//收到结果
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

	if resp.Data.Num == 5 {
		log.Println("答题完毕！！！")
		gameRestart()
	}

	return
}

//收到题目开始点
func injectQuestionResponse(bytes []byte, ctx *goproxy.ProxyCtx) (data []byte) {
	var resp QuestionResp
	var origin QuestionResp

	json.Unmarshal(bytes, &resp)
	json.Unmarshal(bytes, &origin)

	cacheKey := ctx.UserData.(string)

	cache[cacheKey] = NewQuestion(origin)

	start := time.Now()

	answer := fetchAnswerFromCache(resp.Data.Quiz)

	//收到题目开始点答案

	tapSwitch = true

	guss := 0

	if answer != "" {
		for index, option := range resp.Data.Options {
			if option == answer {
				resp.Data.Options[index] = option + "[标答]"
				guss = index
			}
		}
	} else {

		page := search(resp.Data.Quiz)

		//如果题干中包含 '不' 字结果反向取
		var max, min, reverse = 0, 65535,
			strings.Contains(resp.Data.Quiz, "不是") ||
				strings.Contains(resp.Data.Quiz, "不属于") ||
				strings.Contains(resp.Data.Quiz, "不包括")

		for index, option := range resp.Data.Options {
			words := util.Split(option)

			grade := strings.Count(page, option)

			//log.Println(option + "加了" + strconv.Itoa(grade) + "权重")

			if len(words) > 1 {
				for _, word := range words {
					if len(word) > 1 {
						//分词的权重计算
						g := int(float32(strings.Count(page, option)) * (1 / float32(len(words))))
						grade += g
						//log.Println(word + "加了" + strconv.Itoa(g) + "权重")
					}
				}
			}

			resp.Data.Options[index] = option + "[" + strconv.Itoa(grade) + "]"

			if reverse {
				if grade < min {
					min = grade
					guss = index
				}
			} else {
				if grade > max {
					max = grade
					guss = index
				}
			}

		}

	}

	end := time.Now()
	delta := end.Sub(start)
	//log.Printf("查找答案耗时: %s\n", delta)

	tap(guss, (3333*time.Millisecond)-delta)

	log.Println(resp.Data.Quiz)

	for _, item := range resp.Data.Options {
		log.Println(item)
	}

	data, _ = json.Marshal(resp)

	return
}

// 循环点按直到返回结果，不同分辨率按钮位置不同
func tap(i int, delay time.Duration) {
	log.Println("延时点按", string(97+i), delay)
	go func() {
		time.Sleep(delay)
		times := 1
		for tapSwitch {
			switch i {
			case 0:
				util.RunWithAdb("shell", "input tap 540 1040")
				break
			case 1:
				util.RunWithAdb("shell", "input tap 540 1240")
				break
			case 2:
				util.RunWithAdb("shell", "input tap 540 1440")
				break
			case 3:
				util.RunWithAdb("shell", "input tap 540 1640")
				break
			}
			times++

			//遇到特殊情况，对方退出，需要重新进入游戏
			if times > 10 {
				gameRestart()
			}
		}
	}()

}

//答题完毕后点击 继续游戏 ，但是这里可能会遇到弹出升级框的情况，有待优化
func gameRestart() {
	go func() {
		time.Sleep(10 * time.Second)
		util.RunWithAdb("shell", "input tap 540 1440")
		//time.Sleep(2 * time.Second)
		util.RunWithAdb("shell", "input tap 540 1740")
	}()
}

func search(question string) string {
	req, _ := http.NewRequest("GET", "http://www.baidu.com/s?wd="+url.QueryEscape(question), nil)
	resp, _ := http.DefaultClient.Do(req)
	content, _ := ioutil.ReadAll(resp.Body)
	return string(content)
}

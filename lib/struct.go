package lib

//{"data":{"quiz":"微信中界面、按钮的主要颜色是？", "options":["橙色", "黄色", "蓝色", "绿色"], "num":1,"school":"生活", "type":"日常", "contributor":"", "endTime":1515413835, "curTime":1515413820}, "errcode":0}
type QuestionResp struct {
	Data struct {
		Quiz        string   `json:"quiz"`
		Options     []string `json:"options"`
		Num         int      `json:"num"`
		School      string   `json:"school"`
		Type        string   `json:"type"`
		Contributor string   `json:"contributor"`
		EndTime     int      `json:"endTime"`
		CurTime     int      `json:"curTime"`
	} `json:"data"`
	Errcode int `json:"errcode"`
}

type ChooseResp struct {
	Data struct {
		UID         int  `json:"uid"`
		Num         int  `json:"num"`
		Answer      int  `json:"answer"`
		Option      int  `json:"option"`
		Yes         bool `json:"yes"`
		Score       int  `json:"score"`
		TotalScore  int  `json:"totalScore"`
		RowNum      int  `json:"rowNum"`
		RowMult     int  `json:"rowMult"`
		CostTime    int  `json:"costTime"`
		RoomID      int  `json:"roomId"`
		EnemyScore  int  `json:"enemyScore"`
		EnemyAnswer int  `json:"enemyAnswer"`
	} `json:"data"`
	Errcode int `json:"errcode"`
}

type Question struct {
	Quiz    string   `json:"quiz"`
	Options []string `json:"options"`
	School  string   `json:"school"`
	Type    string   `json:"type"`
	Answer  string   `json:"answer"`
}

func NewQuestion(resp QuestionResp) Question {
	question := new(Question)
	question.Quiz = resp.Data.Quiz
	question.Options = resp.Data.Options
	question.School = resp.Data.School
	question.Type = resp.Data.Type
	return *question
}

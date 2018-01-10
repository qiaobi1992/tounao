package lib

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"log"
	"encoding/json"
	"fmt"
)

const (
	SQL_INIT = `
PRAGMA foreign_keys = FALSE;
CREATE TABLE IF NOT EXISTS questions (
	quiz VARCHAR(2000) NOT NULL,
	school VARCHAR(20) NOT NULL,
	type VARCHAR(20) NOT NULL,
	options VARCHAR(2000) NOT NULL,
	answer VARCHAR(2000) NOT NULL,
	PRIMARY KEY (quiz)
);
PRAGMA foreign_keys = TRUE;
`
)

var (
	db *sql.DB
)

func init() {
	db, _ = sql.Open("sqlite3", "data/data.db")

	db.Exec(SQL_INIT)

	log.Println("database init success")
}

func main() {

}

func fetchAnswerFromCache(quiz string) (answer string) {

	s := "SELECT answer FROM questions WHERE quiz = '" + quiz + "' LIMIT 1"

	//log.Println(s)

	rows, _ := db.Query(s)

	if rows != nil {
		defer rows.Close()

		if rows.Next() {
			rows.Scan(&answer)
		}
	}

	return answer
}

func pushAnswerToCache(question Question) {

	s := fmt.Sprintf("INSERT INTO questions(quiz, school, type, options, answer) VALUES ('%s','%s','%s','%s','%s')",
		question.Quiz, question.School, question.Type, question.Options, question.Answer)
	db.Exec(s)

}

func loadAll() []Question {

	questions := make([]Question, 0)

	if db != nil {
		rows, _ := db.Query("SELECT * FROM questions")

		if rows != nil {
			defer rows.Close()
			for rows.Next() {
				question := Question{}
				var q, s, t, o, a string
				rows.Scan(&q, &s, &t, &o, &a)

				question.Quiz = q
				question.School = s
				question.Type = t
				json.Unmarshal([]byte(o), &question.Options)
				question.Answer = a
				questions = append(questions, question)
				//log.Printf("q:%s,s:%s,t:%s,o:%s,a:%s\n", q, s, t, o, a)
			}
		}

	} else {
		log.Panicln("数据库未连接")
	}

	return questions
}

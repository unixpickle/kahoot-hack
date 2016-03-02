package kahoot

import (
	"encoding/json"
	"errors"
)

type QuizActionType int

const (
	QuestionIntro QuizActionType = iota
	QuestionAnswers
)

type QuizAction struct {
	Type       QuizActionType
	NumAnswers int
	Index      int
}

type Quiz struct {
	conn *Conn
}

func NewQuiz(c *Conn) *Quiz {
	return &Quiz{c}
}

// Receive receives the next QuizAction.
// This may be a QuestionIntro, indicating a new question is starting,
// or QuestionAnswers, indicating that the user may now submit an answer.
func (q *Quiz) Receive() (*QuizAction, error) {
	for {
		packet, err := q.conn.Receive("/service/player")
		if err != nil {
			return nil, err
		}
		var content Message
		if data, ok := packet["data"].(map[string]interface{}); !ok {
			continue
		} else if id, ok := data["id"].(float64); !ok {
			continue
		} else if contentStr, ok := data["content"].(string); !ok {
			continue
		} else if json.Unmarshal([]byte(contentStr), &content) != nil {
			continue
		} else if numArray, ok := content["quizQuestionAnswers"].([]interface{}); !ok {
			continue
		} else if questionIndex, ok := content["questionIndex"].(float64); !ok {
			continue
		} else if int(questionIndex) >= len(numArray) || int(questionIndex) < 0 {
			continue
		} else if numAnswers, ok := numArray[int(questionIndex)].(float64); !ok {
			continue
		} else {
			var t QuizActionType
			if id == 1 {
				t = QuestionIntro
			} else if id == 2 {
				t = QuestionAnswers
			} else {
				continue
			}
			return &QuizAction{
				Type:       t,
				NumAnswers: int(numAnswers),
				Index:      int(questionIndex),
			}, nil
		}
	}
}

// Send responds to a server's QuestionAnswers action with an answer index.
func (q *Quiz) Send(index int) error {
	content := Message{
		"choice": index,
		"meta": Message{
			"lag": 22,
			"device": Message{
				"userAgent": "hack",
				"screen": Message{
					"width":  1337,
					"height": 1337,
				},
			},
		},
	}
	encodedContent, _ := json.Marshal(content)
	message := Message{
		"data": Message{
			"id":      6,
			"type":    "message",
			"gameid":  q.conn.gameId,
			"host":    "kahoot.it",
			"content": string(encodedContent),
		},
	}
	if err := q.conn.Send("/service/controller", message); err != nil {
		return err
	}
	if controllerMsg, err := q.conn.Receive("/service/controller"); err != nil {
		return err
	} else if success, ok := controllerMsg["successful"].(bool); !ok || !success {
		return errors.New("did not receive successful response")
	}
	return nil
}

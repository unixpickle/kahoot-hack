package kahoot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// QuizChoice represents a possible answer for a QuizQuestion.
type QuizChoice struct {
	Answer  string `json:"answer"`
	Correct bool   `json:"correct"`
}

// QuizVideo is an optional video for a QuizQuestion.
type QuizVideo struct {
	FullUrl   string  `json:"fullUrl"`
	Id        string  `json:"id"`
	StartTime float64 `json:"startTime"`
	EndTime   float64 `json:"endTime"`
	Service   string  `json:"service"`
}

// QuizQuestion is a question in a quiz.
type QuizQuestion struct {
	NumberOfAnswers int          `json:"numberOfAnswers"`
	Image           string       `json:"image"`
	Video           QuizVideo    `json:"video"`
	Question        string       `json:"question"`
	QuestionFormat  int          `json:"questionFormat"`
	Time            int          `json:"time"`
	Points          bool         `json:"points"`
	Choices         []QuizChoice `json:"choices"`
	Resources       string       `json:"resources"`
	Type            string       `json:"type"`
}

// QuizMetadata stores metadata about a quiz.
type QuizMetadata struct {
	Resolution string         `json:"resolution"`
	Moderation QuizModeration `json:"moderation"`
}

// QuizModeration stores moderator information for a quiz.
type QuizModeration struct {
	FlaggedTimestamp    float64 `json:"flaggedTimestamp"`
	TimestampResolution float64 `json:"timestampResolution"`
	Resolution          string  `json:"resolution"`
}

type userToken struct {
	Email        string            `json:"email"`
	PublicAccess bool              `json:"public_access"`
	PrimaryUsage string            `json:"primary_usage"`
	BannersShown map[string]int    `json:"banners_shown"`
	Metadata     map[string]string `json:"metadata"`
	Picture      string            `json:"picture"`
	Uuid         string            `json:"uuid"`
	Activated    bool              `json:"activated"`
	Created      int64             `json:"created"`
	Modified     int64             `json:"modified"`
	Type         string            `json:"type"`
	Username     string            `json:"username"`
	Birthday     []int             `json:"birthday"`
}

type token struct {
	AccessToken        string            `json:"access_token"`
	Expires            int64             `json:"expires"`
	User               userToken         `json:"user"`
	Roles              []string          `json:"roles"`
	CountryCode        string            `json:"countryCode"`
	CampaignAttributes map[string]string `json:"campaignAttributes"`
}

// QuizInfo stores information about a quiz, including
// the correct answers.
type QuizInfo struct {
	Uuid                string         `json:"uuid"`
	QuizType            string         `json:"quizType"`
	Cover               string         `json:"cover"`
	Modified            int64          `json:"modified"`
	Creator             string         `json:"creator"`
	Audience            string         `json:"audience"`
	Title               string         `json:"title"`
	Description         string         `json:"description"`
	Type                string         `json:"type"`
	Created             int64          `json:"created"`
	Language            string         `json:"language"`
	CreatorPrimaryUsage string         `json:"creator_primary_usage"`
	Questions           []QuizQuestion `json:"questions"`
	Image               string         `json:"image"`
	Video               QuizVideo      `json:"video"`
	Metadata            QuizMetadata   `json:"metadata"`
	Resources           string         `json:"resources"`
	CreatorUsername     string         `json:"creator_username"`
	Visibility          int64          `json:"visibility"`
}

// AccessToken returns an access token from the
// kahoot rest api.
func AccessToken(email, password string) (string, error) {
	client := &http.Client{}
	rawauth := map[string]string{"username": email, "password": password, "grant_type": "password"}
	authentication, err := json.Marshal(rawauth)
	if err != nil {
		return "", err
	}
	request, err := http.NewRequest("POST", "https://create.kahoot.it/rest/authenticate", bytes.NewReader(authentication))
	request.Header.Add("content-type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	receivedtoken := &token{}
	err = json.NewDecoder(response.Body).Decode(receivedtoken)
	if err != nil {
		return "", err
	}
	if receivedtoken.User.Activated == false {
		return "", errors.New("401 unauthorized error:email or password is incorrect")
	}
	return receivedtoken.AccessToken, nil
}

// QuizInformation returns all quiz information for a
// specific kahoot id.
func QuizInformation(token, quizid string) (*QuizInfo, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", fmt.Sprintf("https://create.kahoot.it/rest/kahoots/%s", quizid), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("content-type", "application/json")
	request.Header.Add("authorization", token)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	kahootquiz := &QuizInfo{}
	err = json.NewDecoder(response.Body).Decode(kahootquiz)
	if err != nil {
		return nil, err
	}
	return kahootquiz, nil
}

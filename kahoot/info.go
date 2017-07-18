package kahoot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ChoiceStruct struct {
	Answer  string `json:"answer"`
	Correct bool   `json:"correct"`
}
type VideoStruct struct {
	FullUrl   string  `json:"fullUrl"`
	Id        string  `json:"id"`
	StartTime float64 `json:"startTime"`
	EndTime   float64 `json:"endTime"`
	Service   string  `json:"service"`
}
type QuestionStruct struct {
	NumberOfAnswers int            `json:"numberOfAnswers"`
	Image           string         `json:"image"`
	Video           VideoStruct    `json:"video"`
	Question        string         `json:"question"`
	QuestionFormat  int            `json:"questionFormat"`
	Time            int            `json:"time"`
	Points          bool           `json:"points"`
	Choices         []ChoiceStruct `json:"choices"`
	Resources       string         `json:"resources"`
	Type            string         `json:"type"`
}
type MetadataStruct struct {
	Resolution string           `json:"resolution"`
	Moderation ModerationStruct `json:"moderation"`
}
type ModerationStruct struct {
	FlaggedTimestamp    float64 `json:"flaggedTimestamp"`
	TimestampResolution float64 `json:"timestampResolution"`
	Resolution          string  `json:"resolution"`
}

//############################//
type UserStruct struct {
	Email        string            `json:"email"`
	PublicAccess bool              `json:"public_access"`
	PrimaryUsage string            `json:"primary_usage"`
	BannersShown map[string]int    `json:"banners_shown"`
	Metadata     map[string]string `json:"metadata"`
	Picture      string            `json:"picture"`
	Uuid         string            `json:"uuid"`
	Activated    bool              `json:"activated"`
	Created      int               `json:"created"`
	Modified     int               `json:"modified"`
	Type         string            `json:"type"`
	Username     string            `json:"username"`
	Birthday     []int             `json:"birthday"`
}
type Token struct {
	AccessToken        string            `json:"access_token"`
	Expires            int               `json:"expires"`
	User               UserStruct        `json:"user"`
	Roles              []string          `json:"roles"`
	CountryCode        string            `json:"countryCode"`
	CampaignAttributes map[string]string `json:"campaignAttributes"`
}
type KahootQuiz struct {
	Uuid                string           `json:"uuid"`
	QuizType            string           `json:"quizType"`
	Cover               string           `json:"cover"`
	Modified            int              `json:"modified"`
	Creator             string           `json:"creator"`
	Audience            string           `json:"audience"`
	Title               string           `json:"title"`
	Description         string           `json:"description"`
	Type                string           `json:"type"`
	Created             int              `json:"created"`
	Language            string           `json:"language"`
	CreatorPrimaryUsage string           `json:"creator_primary_usage"`
	Questions           []QuestionStruct `json:"questions"`
	//Type                string           `json:"type"`
	Image           string         `json:"image"`
	Video           VideoStruct    `json:"video"`
	Metadata        MetadataStruct `json:"metadata"`
	Resources       string         `json:"resources"`
	CreatorUsername string         `json:"creator_username"`
	Visibility      int            `json:"visibility"`
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
	//decode response
	token := &Token{}
	err = json.NewDecoder(response.Body).Decode(token)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

// QuizInformation returns all quiz information for a
// specific kahoot id.
func QuizInformation(token, quizid string) *KahootQuiz {
	client := &http.Client{}
	request, err := http.NewRequest("GET", fmt.Sprintf("https://create.kahoot.it/rest/kahoots/%s", quizid), nil)
	if err != nil {
		panic(err)
	}
	request.Header.Add("content-type", "application/json")
	request.Header.Add("authorization", token)
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	kahootquiz := &KahootQuiz{}
	err = json.NewDecoder(response.Body).Decode(kahootquiz)
	if err != nil {
		panic(err)
	}
	return kahootquiz
}

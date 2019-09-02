package function

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/nlopes/slack"
	"golang.org/x/net/context"
)

type Submission struct {
	ID            int     `json:"id"`
	EpochSecond   int64   `json:"epoch_second"`
	ProblemID     string  `json:"problem_id"`
	ContestID     string  `json:"contest_id"`
	UserID        string  `json:"user_id"`
	Language      string  `json:"language"`
	Point         float64 `json:"point"`
	Length        int     `json:"length"`
	Result        string  `json:"result"`
	ExecutionTime int     `json:"execution_time"`
}

type Problem struct {
	ID                   string      `json:"id"`
	ContestID            string      `json:"contest_id"`
	Title                string      `json:"title"`
	ShortestSubmissionID int         `json:"shortest_submission_id"`
	ShortestProblemID    string      `json:"shortest_problem_id"`
	ShortestContestID    string      `json:"shortest_contest_id"`
	ShortestUserID       string      `json:"shortest_user_id"`
	FastestSubmissionID  int         `json:"fastest_submission_id"`
	FastestProblemID     string      `json:"fastest_problem_id"`
	FastestContestID     string      `json:"fastest_contest_id"`
	FastestUserID        string      `json:"fastest_user_id"`
	FirstSubmissionID    int         `json:"first_submission_id"`
	FirstProblemID       string      `json:"first_problem_id"`
	FirstContestID       string      `json:"first_contest_id"`
	FirstUserID          string      `json:"first_user_id"`
	SourceCodeLength     int         `json:"source_code_length"`
	ExecutionTime        int         `json:"execution_time"`
	Point                interface{} `json:"point"`
	Predict              float64     `json:"predict"`
	SolverCount          int         `json:"solver_count"`
}

func NotifyReview(ctx context.Context, msg *pubsub.Message) error {

	atcoderUser := os.Getenv("ATCODER_USER")
	slackAPIToken := os.Getenv("SLACK_API_TOKEN")
	slackChannel := os.Getenv("SLACK_CHANNEL")

	if atcoderUser == "" || slackAPIToken == "" || slackChannel == "" {
		return errors.New("missing environments variable")
	}

	// get user_id: args[0] submission data
	filtered, err := GetSubmissionData(atcoderUser)
	if err != nil {
		return err
	}

	// get all problems data
	mapProblem, err := GetProblemData()
	if err != nil {
		return err
	}

	// post message (thread parent)
	ts, err := PostMessage("今週の復習問題を通知します", "", slackAPIToken, slackChannel, "")
	if err != nil {
		return err
	}

	// post message (thread child)
	for _, data := range filtered {
		msgUrl := fmt.Sprintf("https://atcoder.jp/contests/%s/tasks/%s", data.ContestID, data.ProblemID)
		ts, err = PostMessage(mapProblem[data.ProblemID], msgUrl, slackAPIToken, slackChannel, ts)
	}

	return nil
}

// get user's submission data from atcoder-api
func GetSubmissionData(user_id string) ([]Submission, error) {
	values := url.Values{}
	resp, err := http.Get("https://kenkoooo.com/atcoder/atcoder-api/results?user=" + user_id + values.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var submission []Submission

	err = json.Unmarshal(body, &submission)
	if err != nil {
		return nil, err
	}

	// specify the range of past data to review
	now := time.Now()
	startUnix := now.AddDate(0, 0, -21).Unix()
	endUnix := now.AddDate(0, 0, -14).Unix()

	var filtered []Submission
	for _, data := range submission {
		f1 := startUnix <= data.EpochSecond && data.EpochSecond <= endUnix
		f2 := data.Result == "AC"
		f3 := data.Point >= 300
		if !f1 || !f2 || !f3 {
			continue
		}
		filtered = append(filtered, data)
	}

	// remove duplicates
	var noDuplicatedFiltered []Submission
	for i := 0; i < len(filtered); i++ {
		dup := false
		for j := i + 1; j < len(filtered); j++ {
			if filtered[i].ProblemID == filtered[j].ProblemID {
				dup = true
				break
			}
		}
		if dup == true {
			continue
		}
		noDuplicatedFiltered = append(noDuplicatedFiltered, filtered[i])
	}

	return noDuplicatedFiltered, nil
}

func GetProblemData() (map[string]string, error) {
	values := url.Values{}
	resp, err := http.Get("https://kenkoooo.com/atcoder/resources/merged-problems.json" + values.Encode())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var problems []Problem

	err = json.Unmarshal(body, &problems)
	if err != nil {
		return nil, err
	}

	mapProblem := make(map[string]string)

	for _, data := range problems {
		mapProblem[data.ID] = fmt.Sprintf("%s: %s (%#v)", data.ContestID, data.Title, data.Point)
	}

	return mapProblem, nil
}

func PostMessage(msg, mainMsg, slackAPIToken, slackChannel, parentTimestamp string) (string, error) {

	api := slack.New(slackAPIToken)
	attachment := slack.Attachment{
		Color: "#0054a3",
		Text:  mainMsg,
	}

	var channelID, timestamp string
	var err error
	if mainMsg == "" {
		channelID, timestamp, err = api.PostMessage(slackChannel, slack.MsgOptionText(msg, false))
	} else {
		channelID, timestamp, err = api.PostMessage(slackChannel, slack.MsgOptionText(msg, false), slack.MsgOptionAttachments(attachment), slack.MsgOptionPostMessageParameters(slack.PostMessageParameters{ThreadTimestamp: parentTimestamp}))
	}
	if err != nil {
		return "", err
	}
	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	return timestamp, nil
}

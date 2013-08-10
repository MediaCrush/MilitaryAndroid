package main

import (
	"encoding/json"
	"fmt"
	"github.com/jdiez17/go-irc"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

var issueRegexp *regexp.Regexp = regexp.MustCompile("\\#([0-9]+)")

const (
	githubIssueUrl string = "https://api.github.com/repos/%s/%s/issues/%d"
	issueResponse         = "#%d: %s (%s)"
)

type githubIssue struct {
	Html_url string
	Title    string
}

func getGithubIssue(owner, repo string, issue int) (*githubIssue, error) {
	rqurl := fmt.Sprintf(githubIssueUrl, owner, repo, issue)
	res, err := http.Get(rqurl)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	ghissue := new(githubIssue)
	err = json.Unmarshal(bytes, ghissue)
	if err != nil {
		return nil, err
	}

	return ghissue, nil
}

func expandGithubIssue(c *irc.Connection, e *irc.Event) {
	match := issueRegexp.FindString(e.Payload["message"])
	issue := -1
	if match != "" {
		issue, _ = strconv.Atoi(match[1:])
	}

	if issue != -1 {
		ghissue, err := getGithubIssue("MediaCrush", "MediaCrush", issue)
		if err != nil {
			message := fmt.Sprintf("You're not going to space today. GitHub is broken. (%s)", err.Error())
			e.ReactToChannel(c, message)
			return
		}
		message := fmt.Sprintf(issueResponse, issue, ghissue.Title, ghissue.Html_url)
		e.ReactToChannel(c, message)
	}
}

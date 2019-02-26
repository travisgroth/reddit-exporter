package handlers

import (
	"os"
	"testing"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/turnage/graw/reddit"
)

func TestNewRegex(t *testing.T) {
	_, err := NewRegex("test")
	assert.Nil(t, err)
}

func TestAddMatch(t *testing.T) {
	r, _ := NewRegex("test")
	name, regexString := "testMatch", "^testMessage$"

	err := r.AddMatch(name, regexString)
	assert.Nil(t, err)
}

func TestAddMatchFail(t *testing.T) {
	r, _ := NewRegex("test")
	name, regexString := "testMatchFail", "((1)"

	err := r.AddMatch(name, regexString)
	assert.NotNil(t, err)
}

func TestRegex(t *testing.T) {
	comments := []reddit.Comment{
		{
			Subreddit: "test",
			Body:      "FindMe",
		},
		{
			Subreddit: "test",
			Body:      "DontFindMe",
		},
	}

	posts := []reddit.Post{
		{
			Subreddit: "test",
			SelfText:  "FindMe",
		},
		{
			Subreddit: "test",
			SelfText:  "DontFindMe",
		},
		{
			Subreddit: "testTitle",
			Title:     "FindMe",
		},
	}

	r, _ := NewRegex("testHandler")
	r.AddMatch("testGroupNoMatch", "^NoMatch")
	r.AddMatch("testGroup", "^FindMe")

	for _, comment := range comments {
		r.Comment(&comment)
	}

	for _, post := range posts {
		r.Post(&post)
	}

	testRegexData, _ := os.Open("regex_test_data.txt")
	assert.Nil(t, testutil.GatherAndCompare(prom.DefaultGatherer, testRegexData, "subreddit_comment_matches_total", "subreddit_post_matches_total"))

}

package handlers

import (
	"regexp"

	prom "github.com/prometheus/client_golang/prometheus"

	"github.com/turnage/graw/reddit"
)

type regexMatch struct {
	Name  string
	Regex *regexp.Regexp
}

type Regex struct {
	Name    string
	Regexes []regexMatch
}

var (
	Comments = prom.NewCounterVec(
		prom.CounterOpts{Name: "subreddit_comment_regex", Help: "Comment regex counters by matchgroup name, match name, and sub"},
		[]string{"subreddit", "matchgroup", "match"},
	)
	Posts = prom.NewCounterVec(
		prom.CounterOpts{Name: "subreddit_post_regex", Help: "Post regex counters by matchgroup name, match name and sub"},
		[]string{"subreddit", "matchgroup", "match"},
	)
)

func init() {
	prom.MustRegister(Comments)
	prom.MustRegister(Posts)
}

func NewRegex(name string) (*Regex, error) {

	r := &Regex{
		Name:    name,
		Regexes: make([]regexMatch, 0),
	}
	return r, nil
}

func (r *Regex) AddMatch(name string, regexString string) error {

	regex, err := regexp.Compile(regexString)

	if err != nil {
		return err
	}

	match := regexMatch{
		Name:  name,
		Regex: regex,
	}

	r.Regexes = append(r.Regexes, match)
	return nil
}

func (r *Regex) Comment(comment *reddit.Comment) error {
	for _, regex := range r.Regexes {
		if regex.Regex.MatchString(comment.Body) {
			labels := prom.Labels{
				"subreddit":  comment.Subreddit,
				"matchgroup": r.Name,
				"match":      regex.Name,
			}
			Comments.With(labels).Inc()
		}
	}
	return nil
}

func (r *Regex) Post(post *reddit.Post) error {
	for _, regex := range r.Regexes {
		if regex.Regex.MatchString(post.SelfText) || regex.Regex.MatchString(post.Title) {
			labels := prom.Labels{
				"subreddit":  post.Subreddit,
				"matchgroup": r.Name,
				"match":      regex.Name,
			}
			Posts.With(labels).Inc()
		}
	}
	return nil
}

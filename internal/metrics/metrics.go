package metrics

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

var (
	Comments = prom.NewCounterVec(
		prom.CounterOpts{Name: "subreddit_comments_total", Help: "Comment counters by type and sub"},
		[]string{"subreddit", "type"},
	)
	Posts = prom.NewCounterVec(
		prom.CounterOpts{Name: "subreddit_posts_total", Help: "Post count by sub, flair and type"},
		[]string{"subreddit", "flair", "type"},
	)
)

func init() {
	prom.MustRegister(Comments)
	prom.MustRegister(Posts)
}

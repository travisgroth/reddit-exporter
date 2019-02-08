package handlers

import (
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/travisgroth/reddit-exporter/internal/metrics"
	"github.com/turnage/graw/reddit"
)

type Submissions struct{}

func (b *Submissions) Comment(post *reddit.Comment) error {
	commentType := "response"
	if post.IsTopLevel() {
		commentType = "discussion"
	}
	metrics.Comments.With(
		prom.Labels{"subreddit": post.Subreddit, "type": commentType}).Inc()
	return nil
}

func (b *Submissions) Post(post *reddit.Post) error {
	postType := "link"
	if post.IsSelf {
		postType = "self"
	}
	metrics.Posts.With(prom.Labels{"subreddit": post.Subreddit, "flair": post.LinkFlairText, "type": postType}).Inc()
	return nil
}

package handlers

import (
	log "github.com/sirupsen/logrus"
	"github.com/turnage/graw/reddit"
)

type DebugPrinter struct{}

func (d *DebugPrinter) Comment(post *reddit.Comment) error {
	log.WithFields(log.Fields{"subreddit": post.Subreddit, "author": post.Author, "parent": post.ParentID}).Debug("Comment from ", post.Author)
	return nil
}

func (d *DebugPrinter) Post(post *reddit.Post) error {
	log.WithFields(log.Fields{"subreddit": post.Subreddit, "author": post.Author}).Debug("Post from ", post.Author)
	return nil
}

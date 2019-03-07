package scanner

import (
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/travisgroth/reddit-exporter/internal/handlers"

	"github.com/turnage/graw"
	"github.com/turnage/graw/botfaces"
	"github.com/turnage/graw/reddit"
)

type dispatcher struct {
	commentHandlers []botfaces.CommentHandler
	postHandlers    []botfaces.PostHandler
}

func (d *dispatcher) Comment(comment *reddit.Comment) error {

	fixedComment := *comment
	fixedComment.Subreddit = strings.ToLower(fixedComment.Subreddit)
	for _, handler := range d.commentHandlers {
		err := handler.Comment(&fixedComment)
		if err != nil {
			log.Error("Failed to handle comment: ", err)
		}
	}
	return nil
}

func (d *dispatcher) Post(post *reddit.Post) error {
	fixedPost := *post
	fixedPost.Subreddit = strings.ToLower(fixedPost.Subreddit)
	for _, handler := range d.postHandlers {
		err := handler.Post(&fixedPost)
		if err != nil {
			log.Error("Failed to handle post: ", err)
		}
	}
	return nil
}

type Scanner struct {
	Cfg             graw.Config
	Script          reddit.Script
	CommentHandlers []botfaces.CommentHandler
	PostHandlers    []botfaces.PostHandler
	GrawScan        func(interface{}, reddit.Script, graw.Config) (func(), func() error, error)
}

func NewScanner(cfg graw.Config, script reddit.Script, grawScan func(interface{}, reddit.Script, graw.Config) (func(), func() error, error)) *Scanner {
	commentHandlers := make([]botfaces.CommentHandler, 0)
	postHandlers := make([]botfaces.PostHandler, 0)

	commentHandlers = append(commentHandlers, new(handlers.DebugPrinter))
	postHandlers = append(postHandlers, new(handlers.DebugPrinter))

	commentHandlers = append(commentHandlers, new(handlers.Submissions))
	postHandlers = append(postHandlers, new(handlers.Submissions))

	s := &Scanner{
		Cfg:             cfg,
		Script:          script,
		GrawScan:        grawScan,
		CommentHandlers: commentHandlers,
		PostHandlers:    postHandlers,
	}

	return s
}

func (s *Scanner) AddCommentHandler(h botfaces.CommentHandler) {

	s.CommentHandlers = append(s.CommentHandlers, h)
}

func (s *Scanner) AddPostHandler(h botfaces.PostHandler) {
	s.PostHandlers = append(s.PostHandlers, h)
}

func (s Scanner) Run() {

	for {
		log.Info("Starting scan for subreddits: ", s.Cfg.Subreddits)
		d := &dispatcher{
			commentHandlers: s.CommentHandlers,
			postHandlers:    s.PostHandlers,
		}

		_, wait, err := s.GrawScan(d, s.Script, s.Cfg)
		if err != nil {
			log.Error("Failed to start scan: ", err)
			continue
		}
		log.Info("Scan started...")

		err = wait()
		log.Info("Scan Terminated", err)
		if err.Error() == "test ended" {
			break
		}
		time.Sleep(time.Second * 10)
	}

}

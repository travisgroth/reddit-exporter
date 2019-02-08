package scanner

import (
	"errors"
	"testing"
	"time"

	"github.com/turnage/graw/botfaces"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

func TestScannerRun(t *testing.T) {

	testSubs := []string{"testsub"}
	testScript, _ := reddit.NewScript("graw:reddit-exporter:0.1.0", time.Second*1)
	testCfg := graw.Config{SubredditComments: testSubs, Subreddits: testSubs}
	scanCalled := false

	mockGrawScan := func(h interface{}, script reddit.Script, cfg graw.Config) (func(), func() error, error) {

		assert.IsType(t, &dispatcher{}, h)
		assert.Equal(t, testCfg, cfg)
		assert.Equal(t, testScript, script)

		scanCalled = true
		return func() {},
			func() error { return errors.New("test ended") },
			nil
	}

	scanner := NewScanner(
		testCfg,
		testScript,
		mockGrawScan,
	)

	scanner.Run()
	assert.True(t, scanCalled)
}

func TestScannerAddHandler(t *testing.T) {

	h := new(mockHandler)

	scanner := &Scanner{
		PostHandlers:    make([]botfaces.PostHandler, 0),
		CommentHandlers: make([]botfaces.CommentHandler, 0),
	}

	scanner.AddPostHandler(h)
	scanner.AddCommentHandler(h)

	assert.Equal(t, scanner.PostHandlers[0], h)
	assert.Equal(t, scanner.CommentHandlers[0], h)
}

type mockHandler struct {
	Fail bool
	mock.Mock
}

func (m *mockHandler) Comment(comment *reddit.Comment) error {
	m.Called(comment)
	if m.Fail {
		return errors.New("Comment Failed")
	}

	return nil
}

func (m *mockHandler) Post(post *reddit.Post) error {
	m.Called(post)

	if m.Fail {
		return errors.New("Post Failed")
	}

	return nil
}

func TestCommentDispatch(t *testing.T) {
	m := new(mockHandler)

	mockComment := new(reddit.Comment)
	mockPost := new(reddit.Post)

	m.On("Comment", mockComment).Return(nil)
	m.On("Post", mockPost).Return(nil)

	d := new(dispatcher)
	d.commentHandlers = []botfaces.CommentHandler{m}
	d.postHandlers = []botfaces.PostHandler{m}

	d.Post(mockPost)
	d.Comment(mockComment)

	m.AssertExpectations(t)

	m.Fail = true
	assert.Nil(t, d.Post(mockPost))
	assert.Nil(t, d.Comment(mockComment))

}

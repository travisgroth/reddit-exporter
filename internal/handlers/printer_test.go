package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turnage/graw/reddit"
)

func TestDebugPrinter(t *testing.T) {

	testComment := &reddit.Comment{
		Subreddit: "test",
		Author:    "testAuthor",
		ParentID:  "testParent",
	}

	testPost := &reddit.Post{
		Subreddit: "test",
		Author:    "testAuthor",
	}
	d := new(DebugPrinter)
	assert.Nil(t, d.Comment(testComment), "Comment handling failed")
	assert.Nil(t, d.Post(testPost), "Post handling failed")

}

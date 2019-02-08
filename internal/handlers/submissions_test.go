package handlers

import (
	"os"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"

	prom "github.com/prometheus/client_golang/prometheus"

	"github.com/stretchr/testify/assert"
	"github.com/turnage/graw/reddit"
)

func TestSubmissions(t *testing.T) {

	testComments := []reddit.Comment{
		{Subreddit: "test", Author: "testAuthor", ParentID: "t1_"},
		{Subreddit: "test", Author: "testAuthor", ParentID: "t3_"},
	}

	testPosts := []reddit.Post{
		{Subreddit: "test", Author: "testAuthor"},
		{Subreddit: "test", Author: "testAuthor", IsSelf: true},
		{Subreddit: "test", Author: "testAuthor", LinkFlairText: "testFlair"},
	}

	d := new(Submissions)
	for _, comment := range testComments {
		assert.Nil(t, d.Comment(&comment), "Comment handling failed")
	}

	for _, post := range testPosts {
		assert.Nil(t, d.Post(&post), "Post handling failed")
	}

	testCommentData, _ := os.Open("submissions_test_data_comments.txt")
	assert.Nil(t, testutil.GatherAndCompare(prom.DefaultGatherer, testCommentData, "subreddit_comment"))

	testPostData, _ := os.Open("submissions_test_data_posts.txt")
	assert.Nil(t, testutil.GatherAndCompare(prom.DefaultGatherer, testPostData, "subreddit_post"))

}

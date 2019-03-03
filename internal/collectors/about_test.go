package collectors

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jarcoal/httpmock"
)

var (
	mockAboutData = `{ "data": { "accounts_active": 99, "subscribers": 1 } }`
)

func TestGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://api.reddit.com/r/pass/about.json",
		httpmock.NewStringResponder(200, mockAboutData),
	)
	httpmock.RegisterResponder("GET", "https://api.reddit.com/r/fail/about.json",
		httpmock.NewStringResponder(469, "Error"),
	)

	// pass
	c := NewAboutSubredditCollector("pass", new(http.Client))
	info, _ := c.get()
	assert.Equal(t, float64(99), info.AccountsActive)
	assert.Equal(t, float64(1), info.Subscribers)

	//fail
	c = NewAboutSubredditCollector("fail", new(http.Client))
	info, err := c.get()
	assert.NotNil(t, err)
	assert.Nil(t, info)

	//assert.Equal
	//prom.MustRegister(c)

	//fmt.Println(testutil.GatherAndCompare(prom.DefaultGatherer, strings.NewReader("test"), "subreddit_active_users"))
}

func TestCollect(t *testing.T) {

}

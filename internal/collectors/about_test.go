package collectors

import (
	"net/http"
	"os"
	"testing"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"

	"github.com/jarcoal/httpmock"
)

var (
	mockAboutData = `{ "data": { "accounts_active": 99, "subscribers": 1 } }`
)

func setup() {
	httpmock.Activate()

	httpmock.RegisterResponder("GET", "https://api.reddit.com/r/pass/about.json",
		httpmock.NewStringResponder(200, mockAboutData),
	)
	httpmock.RegisterResponder("GET", "https://api.reddit.com/r/fail/about.json",
		httpmock.NewStringResponder(469, "Error"),
	)
}

func TestGet(t *testing.T) {
	setup()

	defer httpmock.DeactivateAndReset()

	c := NewAboutSubredditCollector([]string{"pass", "fail"}, new(http.Client))

	// pass
	info, _ := c.getSubredditInfo("pass")
	assert.Equal(t, float64(99), info.AccountsActive)
	assert.Equal(t, float64(1), info.Subscribers)

	//fail
	info, err := c.getSubredditInfo("fail")
	assert.NotNil(t, err)
	assert.Nil(t, info)
}

func TestCollect(t *testing.T) {
	setup()

	defer httpmock.DeactivateAndReset()
	c := NewAboutSubredditCollector([]string{"pass"}, new(http.Client))
	prom.MustRegister(c)

	testAboutData, _ := os.Open("about_test_data.txt")
	assert.Nil(t, testutil.GatherAndCompare(prom.DefaultGatherer, testAboutData, "subreddit_active_users", "subreddit_subscriber_users"))
}

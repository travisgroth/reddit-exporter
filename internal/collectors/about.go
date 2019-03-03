package collectors

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	prom "github.com/prometheus/client_golang/prometheus"

	log "github.com/sirupsen/logrus"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type AboutSubredditCollector struct {
	Client    HTTPClient
	Subreddit string
}

type subInfo struct {
	AccountsActive float64 `json:"accounts_active"`
	Subscribers    float64 `json:"subscribers"`
}

type aboutResponse struct {
	Data subInfo `json:"data"`
}

var (
	activeUsersDesc = prom.NewDesc("subreddit_active_users",
		"Current active users by sub",
		[]string{"subreddit"}, nil,
	)

	subscribersDesc = prom.NewDesc("subreddit_subscriber_users",
		"Current subscriber users by sub",
		[]string{"subreddit"}, nil,
	)
)

func NewAboutSubredditCollector(subreddit string, client HTTPClient) *AboutSubredditCollector {

	collector := &AboutSubredditCollector{
		Client:    client,
		Subreddit: subreddit,
	}

	return collector
}

func (collector *AboutSubredditCollector) Collect(c chan<- prom.Metric) {

	info, err := collector.get()
	if err != nil {
		return
	}

	c <- prometheus.MustNewConstMetric(activeUsersDesc, prom.GaugeValue, info.AccountsActive, collector.Subreddit)
	c <- prometheus.MustNewConstMetric(subscribersDesc, prom.GaugeValue, info.Subscribers, collector.Subreddit)

}

func (collector *AboutSubredditCollector) Describe(c chan<- *prom.Desc) {
	c <- activeUsersDesc
	c <- subscribersDesc
}

func (c *AboutSubredditCollector) get() (*subInfo, error) {
	client := &http.Client{}
	aboutURL := fmt.Sprintf("https://api.reddit.com/r/%s/about.json", c.Subreddit)
	req, _ := http.NewRequest("GET", aboutURL, nil)
	req.Header.Set("User-Agent", "Golang_Reddit_Exporter")
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Could not get info for subreddit '%s': %s", c.Subreddit, err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		err = fmt.Errorf("Bad return code: %d.  %s", resp.StatusCode, body)

		log.Errorf("Could not get info for subreddit '%s': %s", c.Subreddit, err)
		return nil, err
	}

	defer resp.Body.Close()

	var info aboutResponse
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		log.Errorf("Could not decode info for subreddit '%s': %s", c.Subreddit, err)
		return nil, err
	}

	return &info.Data, err
}

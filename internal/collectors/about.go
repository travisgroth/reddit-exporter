package collectors

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	prom "github.com/prometheus/client_golang/prometheus"

	log "github.com/sirupsen/logrus"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type AboutSubredditCollector struct {
	Client     HTTPClient
	Subreddits []string
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

func NewAboutSubredditCollector(subreddits []string, client HTTPClient) *AboutSubredditCollector {

	collector := &AboutSubredditCollector{
		Client:     client,
		Subreddits: subreddits,
	}

	return collector
}

func (collector *AboutSubredditCollector) Collect(c chan<- prom.Metric) {
	wg := new(sync.WaitGroup)
	for _, subreddit := range collector.Subreddits {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			info, err := collector.getSubredditInfo(s)
			if err == nil {
				c <- prometheus.MustNewConstMetric(activeUsersDesc, prom.GaugeValue, info.AccountsActive, s)
				c <- prometheus.MustNewConstMetric(subscribersDesc, prom.GaugeValue, info.Subscribers, s)
			}
		}(subreddit)
	}

	wg.Wait()

}

func (collector *AboutSubredditCollector) Describe(c chan<- *prom.Desc) {
	c <- activeUsersDesc
	c <- subscribersDesc
}

func (c *AboutSubredditCollector) getSubredditInfo(subreddit string) (*subInfo, error) {
	client := c.Client
	aboutURL := fmt.Sprintf("https://api.reddit.com/r/%s/about.json", subreddit)
	req, _ := http.NewRequest("GET", aboutURL, nil)
	req.Header.Set("User-Agent", "Golang_Reddit_Exporter")
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Could not get info for subreddit '%s': %s", subreddit, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		err = fmt.Errorf("Bad return code: %d.  %s", resp.StatusCode, body)

		log.Errorf("Could not get info for subreddit '%s': %s", subreddit, err)
		return nil, err
	}

	var info aboutResponse
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		log.Errorf("Could not decode info for subreddit '%s': %s", subreddit, err)
		return nil, err
	}

	return &info.Data, err
}

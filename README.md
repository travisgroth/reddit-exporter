# reddit-exporter
![Docker Pulls](https://img.shields.io/docker/pulls/travisgroth/reddit-exporter.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/travisgroth/reddit-exporter)](https://goreportcard.com/report/github.com/travisgroth/reddit-exporter)
![Docker Automated build](https://img.shields.io/docker/automated/travisgroth/reddit-exporter.svg)
![CircleCI (all branches)](https://img.shields.io/circleci/project/github/travisgroth/reddit-exporter.svg)

A prometheus compatible exporter for generating stats about the traffic on a given subreddit.  It is meant to help capture traffic and utilization trends over time without crawling over historic data.  The most likely use case is for monitoring a community related to a particular product or service, though more novel usage may arise.

A single exporter can handle multiple subreddits and breaks down all stats by subreddit in addition to other dimensions.

Internally, reddit-exporter runs a continous scan on _new_ content to a subreddit to accumulate metrics on posts and comments.  Any "cheap" point-in-time metrics are retrieved at scrape time.

## Current Capabilities
* Post counter by type
* Comment counter by type (top level discussion or a reply)
* Configurable regex matching counters for both post and comment content
* Subscriber gauge
* Active user gauge

## Maybe Capabilities
* Custom metrics via RPC to a downstream service

# Install

Building reddit-exporter uses go modules and requires go 1.11+

## Go Install

```
go install github.com/travisgroth/reddit-exporter/cmd/reddit-exporter
```

## Compile

```
go build cmd/reddit-exporter/*.go
```

## Docker

```
docker run -p 8000:8000 travisgroth/reddit-exporter -s askreddit
```

## Helm
See chart source [here](helm/reddit-exporter)
```
$ helm repo add reddit-exporter https://reddit-exporter-chart.storage.googleapis.com
"reddit-exporter" has been added to your repositories

helm upgrade --install reddit-exporter-test reddit-exporter/reddit-exporter --set subreddits="{wtf,askreddit}"
```

# Usage

## Command Line 

```
Usage:
  reddit-exporter [flags]

Flags:
      --address string      Metrics server bind address (default "0.0.0.0")
  -h, --help                help for reddit-exporter
  -p, --port int            Metrics server port (default 8000)
      --regexfile strings   File containing regex matches in format 'name;regex'
                            Can be specified multiple times or comma separated
  -s, --subreddit strings   Subreddit(s) to monitor (Required).
                            Can be specified multiple times or comma separated
  -v, --verbose             Turn on verbose mode
  ```

  You may list multiple subreddits either comma separated or with multiple '-s' flags.

  Verbose mode logs individual events and is good for debugging.

  See next section for how to use regexfiles.

  Typical command:

  ```
$ ./reddit-exporter -v -s askreddit,wtf
time="2019-02-23T15:52:36-05:00" level=info msg="Listening on 0.0.0.0:8000"
time="2019-02-23T15:52:36-05:00" level=info msg="Starting scan for subreddits: [askreddit wtf]"
time="2019-02-23T15:52:39-05:00" level=info msg="Scan started..."
time="2019-02-23T15:52:40-05:00" level=debug msg="Post from Nar40e03" author=Nar40e03 subreddit=AskReddit
```

## Regex Configuration

You may configure reddit-exporter with multiple files containing labeled regexes.  You may use json, yaml or properties file format for your regexfile.  Matches to these regexes will be kept track of in a counter broken down by matchgroup and matchname.  

Example yaml file:

```
MyMatchGroup:
    MyFirstMatchName: ^first$
    MySecondMatchName: ^second$
```

Equivilent properties file:
```
MyMatchGroup.MyFirstMatchName=^first$
MyMatchGroup.MySecondMatchName=^first$
```

Each regex has a matchgroup name as well as a match name.  

* Match Groups should be used to label related regexes for easier aggregation and querying later.  

* Match Names indicate the exact regex that matched

* Both the match name and match group appear as labels on the metric created for the regex

Regexes are not wildcard padded so if you are looking for mid-string matching, be sure to include '.*' on both ends of your regex.

reddit-exporter uses [viper](https://github.com/spf13/viper) under the hood, so as long as your chosen format presents key/value matches as expected, any supported file format should work.  

# Sample Metrics

Assuming a regexfile containing the following:
```
test.the=.* the .*
```

```
# HELP subreddit_active_users Current active users by sub
# TYPE subreddit_active_users gauge
subreddit_active_users{subreddit="askreddit"} 95117
# HELP subreddit_comments_total Comment counters by type and sub
# TYPE subreddit_comments_total counter
subreddit_comments_total{subreddit="askreddit",type="discussion"} 56
subreddit_comments_total{subreddit="askreddit",type="response"} 45
subreddit_comments_total{subreddit="wtf",type="discussion"} 3
# HELP subreddit_comment_matches_total Comment regex counters by matchgroup name, match name, and sub
# TYPE subreddit_comment_matches_total counter
subreddit_comment_matches_total{match="the",matchgroup="test",subreddit="askreddit"} 42
# HELP subreddit_posts_total Post count by sub, flair and type
# TYPE subreddit_posts_total counter
subreddit_posts_total{flair="",subreddit="askreddit",type="self"} 7
subreddit_posts_total{flair="",subreddit="wtf",type="link"} 1
# HELP subreddit_post_matches_total Post regex counters by matchgroup name, match name and sub
# TYPE subreddit_post_matches_total counter
subreddit_post_matches_total{match="the",matchgroup="test",subreddit="askreddit"} 3
# HELP subreddit_subscriber_users Current subscriber users by sub
# TYPE subreddit_subscriber_users gauge
subreddit_subscriber_users{subreddit="askreddit"} 2.1760036e+07
```



package main

import (
	"fmt"
	"testing"

	"github.com/travisgroth/reddit-exporter/internal/handlers"

	"github.com/turnage/graw/botfaces"

	"github.com/stretchr/testify/mock"
)

type mockScanner struct {
	mock.Mock
}

func (m *mockScanner) AddCommentHandler(h botfaces.CommentHandler) {
	m.Called(h)
}

func (m *mockScanner) AddPostHandler(h botfaces.PostHandler) {
	m.Called(h)
}

func TestAddregexHandlers(t *testing.T) {
	m := new(mockScanner)
	var testProperties = map[string]map[string]string{
		"testgroup1": map[string]string{
			"matchname1": "match1",
		},
		"testgroup2": map[string]string{
			"matchname2": "match2",
		},
	}

	for groupName, match := range testProperties {
		r, _ := handlers.NewRegex(groupName)
		for matchName, regex := range match {
			r.AddMatch(matchName, regex)
			fmt.Println(regex)
		}
		m.On("AddCommentHandler", r).Return()
		m.On("AddPostHandler", r).Return()
	}
	addRegexHandlers(m, "test.properties")
	m.AssertExpectations(t)
}

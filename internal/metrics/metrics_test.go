package metrics

import (
	"testing"

	prom "github.com/prometheus/client_golang/prometheus"

	"github.com/stretchr/testify/assert"
)

func TestRegistry(t *testing.T) {
	assert.NotEmpty(t, prom.DefaultRegisterer)
	assert.NotNil(t, Comments)
	assert.NotNil(t, Posts)
}

package accounts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Unittest-1 - Make sure the SetBaseURL() function on our custom struct URL type, works as expected.
func TestSetBaseURL(t *testing.T) {
	DefaultUrl.SetBaseURL("super.fake.com")
	assert.Equal(t, DefaultUrl.BaseURL, "super.fake.com", "setting custom BaseURL failed")
}

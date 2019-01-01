package webhook_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/yakshaving.art/alertsnitch/webhook"
)

func TestParsingPayloadWithEmptyPayloadFails(t *testing.T) {
	_, err := webhook.Parse([]byte(""))
	assert.EqualError(t, err, "failed to decode json webhook payload: unexpected end of JSON input")
}

func TestParsingPayloadWithInvalidPayloadFails(t *testing.T) {
	_, err := webhook.Parse([]byte("error"))
	assert.EqualError(t, err, "failed to decode json webhook payload: invalid character 'e' looking for beginning of value")
}

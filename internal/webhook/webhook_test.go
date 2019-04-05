package webhook_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/yakshaving.art/alertsnitch/internal/webhook"
)

func TestParsingPayloadWithEmptyPayloadFails(t *testing.T) {
	_, err := webhook.Parse([]byte(""))
	assert.EqualError(t, err, "failed to decode json webhook payload: unexpected end of JSON input")
}

func TestParsingPayloadWithInvalidPayloadFails(t *testing.T) {
	_, err := webhook.Parse([]byte("error"))
	assert.EqualError(t, err, "failed to decode json webhook payload: invalid character 'e' looking for beginning of value")
}

func TestParsingValidPayloadWorks(t *testing.T) {
	a := assert.New(t)
	b, err := ioutil.ReadFile("sample-payload.json")

	a.NoError(err)

	d, err := webhook.Parse(b)

	a.NoError(err)
	a.NotNil(d)

	a.Equal(d.Status, "resolved")
	a.Equal(d.ExternalURL, "http://alertmanager:9093")
}

func TestParsingValidPayloadWithoutEndsAtWorks(t *testing.T) {
	a := assert.New(t)
	b, err := ioutil.ReadFile("sample-payload-invalid-ends-at.json")

	a.NoError(err)

	d, err := webhook.Parse(b)

	a.NoError(err)
	a.NotNil(d)

	a.Equal(d.Status, "resolved")
	a.Equal(d.ExternalURL, "http://alertmanager:9093")
	a.True(d.Alerts[0].EndsAt.Before(d.Alerts[0].StartsAt))
}

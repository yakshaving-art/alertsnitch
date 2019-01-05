package webhook

import (
	"encoding/json"
	"fmt"

	"gitlab.com/yakshaving.art/alertsnitch/internal"
)

// Parse gets a webhook payload and parses it returning a prometheus
// template.Data object if successful
func Parse(payload []byte) (*internal.AlertGroup, error) {
	d := internal.AlertGroup{}
	err := json.Unmarshal(payload, &d)
	if err != nil {
		return nil, fmt.Errorf("failed to decode json webhook payload: %s", err)
	}
	return &d, nil
}

package webhook

import (
	"encoding/json"
	"fmt"

	"github.com/prometheus/alertmanager/template"
)

func parsePayload(payload []byte) (*template.Data, error) {
	d := template.Data{}
	err := json.Unmarshal(payload, &d)
	if err != nil {
		return nil, fmt.Errorf("failed to decode json webhook payload: %s", err)
	}
	return &d, nil
}

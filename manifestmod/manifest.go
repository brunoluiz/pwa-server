package manifestmod

import (
	"encoding/json"
	"fmt"
)

// ChangeBaseURL changes manifest.json start_url and scope using baseURL
func ChangeBaseURL(buf []byte, baseURL string) ([]byte, error) {
	manifest := make(map[string]interface{})

	if err := json.Unmarshal(buf, &manifest); err != nil {
		return nil, fmt.Errorf("error on unmarshaling manifest.json: %s", err)
	}

	for k, v := range manifest {
		if k != "start_url" && k != "scope" {
			continue
		}

		if _, ok := v.(string); ok {
			manifest[k] = baseURL
		}
	}

	js, err := json.Marshal(manifest)
	if err != nil {
		return nil, fmt.Errorf("error on marshalling manifest.json: %s", err)
	}

	return js, nil
}

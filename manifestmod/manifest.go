package manifestmod

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

// ChangeBaseURL changes manifest.json start_url and scope using baseURL
func ChangeBaseURL(f io.Reader, baseURL string) ([]byte, error) {
	manifest := make(map[string]interface{})
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error on opening manifest.json: %s", err)
	}

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

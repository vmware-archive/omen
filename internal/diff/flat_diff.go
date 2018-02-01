package diff

import (
	"bytes"
	"strings"

	"encoding/json"

	"github.com/kylelemons/godebug/diff"
)

func FlatDiff(object1 interface{}, object2 interface{}) (string, error) {
	// This marshal-unmarshal dance must be done so that the struct JSON tags
	// of manifests are applied to the output. For example, we want the top-level
	// key to be manifests, not Data. See manifest.Manifests for more information.
	aJson, err := json.Marshal(object1)
	if err != nil {
		return "", err
	}

	bJson, err := json.Marshal(object2)
	if err != nil {
		return "", err
	}

	var a interface{}
	err = json.Unmarshal([]byte(aJson), &a)
	if err != nil {
		return "", err
	}

	var b interface{}
	err = json.Unmarshal([]byte(bJson), &b)
	if err != nil {
		return "", err
	}

	d := diff.Diff(Flatten(a), Flatten(b))

	sb := bytes.Buffer{}

	for _, s := range strings.Split(d, "\n") {
		switch {
		case strings.HasPrefix(s, "+"):
			sb.Write([]byte(s + "\n"))
		case strings.HasPrefix(s, "-"):
			sb.Write([]byte(s + "\n"))
		}
	}

	return sb.String(), nil
}

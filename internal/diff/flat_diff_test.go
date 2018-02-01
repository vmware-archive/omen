package diff_test

import (
	"github.com/pivotal-cloudops/omen/internal/diff"

	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FlatDiff", func() {
	It("diffs a simple example", func() {
		simple1, err := convertStrToJson(SIMPLE1)
		Expect(err).ToNot(HaveOccurred())

		simple2, err := convertStrToJson(SIMPLE2)
		Expect(err).ToNot(HaveOccurred())

		actual, err := diff.FlatDiff(simple1, simple2)
		Expect(err).NotTo(HaveOccurred())
		Expect(actual).To(ContainSubstring(SIMPLE_EXPECTED))
	})
})

func convertStrToJson(jsonStr string) (interface{}, error) {
	var data interface{}

	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return "", err
	}

	return data, nil
}

const SIMPLE1 string = `{ "a": { "field": "value1" } }`
const SIMPLE2 string = `{ "a": { "field": "value2" } }`

const SIMPLE_EXPECTED string = `-a.field=value1
+a.field=value2
`

package diff_test

import (
	"github.com/pivotal-cloudops/omen/internal/diff"

	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Flatten", func() {
	It("flattens a json string", func() {
		var data interface{}

		err := json.Unmarshal([]byte(FLAT_SIMPLE1), &data)
		Expect(err).ToNot(HaveOccurred())

		actual := diff.Flatten(data)

		Expect(actual).To(ContainSubstring(FLAT_SIMPLE_EXPECTED))
	})

	It("flattens a struct", func() {
		data := Outer{
			Person: []Inner{
				{
					Name:  "james",
					Phone: "whatevs",
				},
				{
					Name:  "Tom",
					Phone: "yokey",
				},
			},
			Birthday: "whenevs",
		}

		actual := diff.Flatten(data)

		Expect(actual).To(ContainSubstring(FLAT_STRUCT_EXPECTED))
	})

	It("flattens a complex json string", func() {
		var data interface{}

		err := json.Unmarshal([]byte(FLAT_COMPLEX1), &data)
		Expect(err).ToNot(HaveOccurred())

		actual := diff.Flatten(data)

		Expect(actual).To(ContainSubstring(FLAT_COMPLEX1_EXPECTED))
	})

	DescribeTable("Flatten handles different data types",
		func(a string, b string) {
			var data interface{}
			err := json.Unmarshal([]byte(a), &data)
			Expect(err).ToNot(HaveOccurred())

			actual := diff.Flatten(data)
			Expect(actual).To(ContainSubstring(b))
		},
		Entry("Is a Decimal", FLAT_DECIMAL, FLAT_DECIMAL_EXPECTED),
		Entry("Is a float", FLAT_FLOAT, FLAT_FLOAT_EXPECTED),
		Entry("Is a bool", FLAT_BOOL, FLAT_BOOL_EXPECTED),
		Entry("Is a array", FLAT_ARRAY, FLAT_ARRAY_EXPECTED),
	)
})

type Inner struct {
	Name  string
	Phone string
}

type Outer struct {
	Person   []Inner
	Birthday string
}

const FLAT_STRUCT_EXPECTED = `Birthday=whenevs
Person.Name=Tom
Person.Name=james
Person.Phone=whatevs
Person.Phone=yokey
`

const FLAT_SIMPLE1 string = `{ "a": { "field": "value1" } }`
const FLAT_SIMPLE_EXPECTED string = "a.field=value1"

const FLAT_COMPLEX1 string = `
{
  "properties": {
    "forwarder": {
      "ca_cert": "-----BEGIN CERTIFICATE-----\nMIIC+zCCAeOgADAfMQswCQYDVQQGEwJVUzEQ\nMA4GA1UECgwHUGl2b3RhbDAeFw0xNzA1MDkxMTMyNDJaFw0yMTA1MTAxMTMyNDJa\nMB8xCzAJBgNVBAYTAlVTMRAwDgYDVQQKDAdQaXZvdGFsMIIBIjANBgkqhki",
      "port": 13322,
	  "semver": "1.0.14",
	  "some-decimal": 1.01,
      "server_cert": "-----BEGIN CERTIFICATE-----\nMIIDbTCCAlWgAwIBAgIUOJtHGpzb/7CYcvJ+QCeCjUIHNCcwDQYJKwHhcNMTcxMDAxMDk0\nNjI0WhcNMTkxMDAxMDk0NjI0WjA1MQswCQYDVQQGEwJVUzEQMA4GA1U",
      "server_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAv5gc27P0Phth1jQHQNBmj5BtUWzSioXVmh2KAZOKE0BLezwaPkzUgXph/15mk2mTXRMk1bvNqybMtBN\nFW4Dy+SmqrtDWjTVIc331UeZkam7mVOnqfdgUAs0X+yZdmZ/Jnqa",
      "server_name": "healthwatch"
    },
    "loggregator": {
      "tls": {
        "ca_cert": "-----BEGIN CERTIFICATE-----\nMIIC+zCCAeOgAwIBAgIBADANswCQYDVQQGEwJVUzEQ\nMA4GA1UECgwHUGl2b3RhbDAeFw0xNzA1MDkxMTMyNDJaFw0yMTA1MTAxMTMyNDJa\nMB8xCzAJBgNVBAYTAlVTMRAwDgYDVQQKDAdQaXZvdGFsMIIBIjANBgkqh",
        "metron": {
          "cert": "-----BEGIN CERTIFICATE-----\nMIIDZDCCAkygAwIBAgIVANEORp9qlJPYvVa4B+rWE552u42XMA0GCSAdQaXZvdGFsMB4XDTE3MDUxNjE1\nNTc1N1oXDTE5MDUxNjE1NTc1N1owMDELMAkGA1UEBhMCVVMxEDAOBgNVBA",
          "key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAyO1ft52OAeoa6n9SCiY4OHblXNVo5kN84wK0OqPk3jIZMHmz\na\nuNqb0cEOQfbvvw/QIj0JJSHCf7vcytDE36mLUD/49OsZVzBX2YJdZfq"
        }
      }
    }
  }
}
`
const FLAT_COMPLEX1_EXPECTED string = `properties.forwarder.ca_cert=-----BEGIN CERTIFICATE-----MIIC+zCCAeOgADAfMQswCQYDVQQGEwJVUzEQMA4GA1UECgwHUGl2b3RhbDAeFw0xNzA1MDkxMTMyNDJaFw0yMTA1MTAxMTMyNDJaMB8xCzAJBgNVBAYTAlVTMRAwDgYDVQQKDAdQaXZvdGFsMIIBIjANBgkqhki
properties.forwarder.port=13322
properties.forwarder.semver=1.0.14
properties.forwarder.server_cert=-----BEGIN CERTIFICATE-----MIIDbTCCAlWgAwIBAgIUOJtHGpzb/7CYcvJ+QCeCjUIHNCcwDQYJKwHhcNMTcxMDAxMDk0NjI0WhcNMTkxMDAxMDk0NjI0WjA1MQswCQYDVQQGEwJVUzEQMA4GA1U
properties.forwarder.server_key=-----BEGIN RSA PRIVATE KEY-----MIIEpAIBAAKCAQEAv5gc27P0Phth1jQHQNBmj5BtUWzSioXVmh2KAZOKE0BLezwaPkzUgXph/15mk2mTXRMk1bvNqybMtBNFW4Dy+SmqrtDWjTVIc331UeZkam7mVOnqfdgUAs0X+yZdmZ/Jnqa
properties.forwarder.server_name=healthwatch
properties.forwarder.some-decimal=1.01
properties.loggregator.tls.ca_cert=-----BEGIN CERTIFICATE-----MIIC+zCCAeOgAwIBAgIBADANswCQYDVQQGEwJVUzEQMA4GA1UECgwHUGl2b3RhbDAeFw0xNzA1MDkxMTMyNDJaFw0yMTA1MTAxMTMyNDJaMB8xCzAJBgNVBAYTAlVTMRAwDgYDVQQKDAdQaXZvdGFsMIIBIjANBgkqh
properties.loggregator.tls.metron.cert=-----BEGIN CERTIFICATE-----MIIDZDCCAkygAwIBAgIVANEORp9qlJPYvVa4B+rWE552u42XMA0GCSAdQaXZvdGFsMB4XDTE3MDUxNjE1NTc1N1oXDTE5MDUxNjE1NTc1N1owMDELMAkGA1UEBhMCVVMxEDAOBgNVBA
properties.loggregator.tls.metron.key=-----BEGIN RSA PRIVATE KEY-----MIIEowIBAAKCAQEAyO1ft52OAeoa6n9SCiY4OHblXNVo5kN84wK0OqPk3jIZMHmzauNqb0cEOQfbvvw/QIj0JJSHCf7vcytDE36mLUD/49OsZVzBX2YJdZfq
`
const FLAT_DECIMAL = `{"properties": {"port": 13322} }`
const FLAT_DECIMAL_EXPECTED = `properties.port=13322`

const FLAT_FLOAT = `{"properties": {"height": 133.222} }`
const FLAT_FLOAT_EXPECTED = `properties.height=133.222`

const FLAT_BOOL = `{"properties": {"fun": true} }`
const FLAT_BOOL_EXPECTED = `properties.fun=true`

const FLAT_ARRAY = `{"properties": {"fun": [
"TOM",
"JAMES",
"GARIMA"
]
} }`
const FLAT_ARRAY_EXPECTED = `properties.fun=GARIMA
properties.fun=JAMES
properties.fun=TOM`

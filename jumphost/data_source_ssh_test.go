package jumphost

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	key = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACDuzSye5Tr7FcwyuEm5oDF15DWuUlkHPQsU0VJq9um3+QAAAJBjJvpdYyb6
XQAAAAtzc2gtZWQyNTUxOQAAACDuzSye5Tr7FcwyuEm5oDF15DWuUlkHPQsU0VJq9um3+Q
AAAEBaDLAdQm00RB4AGO5/uCoTCZsQaHwBvPK+czoQabv2su7NLJ7lOvsVzDK4SbmgMXXk
Na5SWQc9CxTRUmr26bf5AAAAC3Rlc3RzQGxvY2FsAQI=
-----END OPENSSH PRIVATE KEY-----
`

	httpServerDataSourcePassword = `
	data jumphost_ssh "http_server" {
		hostname = "web"
		port = 8080
	}

	output local_port {
		value = data.jumphost_ssh.http_server.local_port
	}

	provider jumphost {
		port = %d
		username = "terraform"
		password = "1234"
	}
	`

	httpServerDataSourcePublicKey = `
	data jumphost_ssh "http_server" {
		hostname = "web"
		port = 8080
	}

	output local_port {
		value = data.jumphost_ssh.http_server.local_port
	}

	provider jumphost {
		port = %d
		username = "terraform"
		private_key = <<EOT
%s
EOT
	}
	`
)

func init() {
	composeUp()
}

func TestSingleConnectionUsingPassword(t *testing.T) {
	resourceName := "data.jumphost_ssh.http_server"
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       fmt.Sprintf(httpServerDataSourcePassword, sshPort),
				Check: func(s *terraform.State) error {
					r := s.RootModule()
					localPort := r.Outputs["local_port"].Value.(string)
					address := fmt.Sprintf("http://localhost:%s/status/418", localPort)
					response, err := http.Get(address)
					if err != nil {
						t.Fatalf("failed to call forwarded service %v", err)
					}
					assert.Equal(t, response.StatusCode, 418)
					var output bytes.Buffer
					_, err = output.ReadFrom(response.Body)
					if err != nil {
						t.Fatalf("failed to read service's response %v", err)
					}
					assert.Equal(t, output.String(), "I'm a teapot!")
					return nil
				},
			},
		},
	})
}

func TestSingleConnectionUsingPublicKey(t *testing.T) {
	resourceName := "data.jumphost_ssh.http_server"
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       fmt.Sprintf(httpServerDataSourcePublicKey, sshPort, key),
				Check: func(s *terraform.State) error {
					r := s.RootModule()
					localPort := r.Outputs["local_port"].Value.(string)
					address := fmt.Sprintf("http://localhost:%s/status/418", localPort)
					response, err := http.Get(address)
					if err != nil {
						t.Fatalf("failed to call forwarded service %v", err)
					}
					assert.Equal(t, response.StatusCode, 418)
					var output bytes.Buffer
					_, err = output.ReadFrom(response.Body)
					if err != nil {
						t.Fatalf("failed to read service's response %v", err)
					}
					assert.Equal(t, output.String(), "I'm a teapot!")
					return nil
				},
			},
		},
	})
}

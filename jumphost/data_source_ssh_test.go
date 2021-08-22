package jumphost

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	httpServerDataSource = `
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
)

func init() {
	composeUp()
}

func TestSingleConnection(t *testing.T) {
	resourceName := "data.jumphost_ssh.http_server"
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: resourceName,
				Config:       fmt.Sprintf(httpServerDataSource, sshPort),
				Check: func(s *terraform.State) error {
					r := s.RootModule()
					localPort := r.Outputs["local_port"].Value.(string)
					address := fmt.Sprintf("http://localhost:%s/status/418", localPort)
					response, err := http.Get(address)
					if err != nil {
						t.Fatalf("failed to call forwarded service %s", err)
					}
					assert.Equal(t, response.StatusCode, 418)
					var output bytes.Buffer
					_, err = output.ReadFrom(response.Body)
					if err != nil {
						t.Fatalf("failed to read service's response %s", err)
					}
					assert.Equal(t, output.String(), "I'm a teapot!")
					return nil
				},
			},
		},
	})
}

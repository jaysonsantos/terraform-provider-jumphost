package jumphost

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

var (
	providerFactories = map[string]func() (*schema.Provider, error){
		"jumphost": newTestProvider,
	}
)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func newTestProvider() (*schema.Provider, error) {
	provider := Provider()
	return provider, nil
}

func composeUp() {
	composeLock.Lock()
	defer composeLock.Unlock()

	cmd := exec.Command("docker-compose", "up", "-d")
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	cmd = exec.Command("docker-compose", "port", "ssh", "2222")
	var output bytes.Buffer
	cmd.Stdout = &output
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	port := strings.Split(strings.Trim(output.String(), "\n"), ":")[1]
	sshPort, err = strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
}

func Test_currentUser(t *testing.T) {
	currentUserCmd := exec.Command("whoami")
	currentLocalUser, err := currentUserCmd.Output()
	if err != nil {
		t.Fatalf("failed to get current user %v", err)
	}
	tests := []struct {
		name  string
		want  string
		setup func()
	}{
		{
			name:  "current user",
			want:  strings.TrimSpace(string(currentLocalUser)),
			setup: func() { os.Setenv("SSH_USER", "") },
		},
		{
			name:  "overriden user",
			want:  "ssh-user",
			setup: func() { os.Setenv("SSH_USER", "ssh-user") },
		},
	}
	for _, tt := range tests {
		tt.setup()
		got, err := currentUser()
		assert.Nil(t, err)
		assert.Equal(t, got, tt.want)
	}
}

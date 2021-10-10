package jumphost

import (
	"context"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	hostNameAttr   = "hostname"
	portAttr       = "port"
	usernameAttr   = "username"
	passwordAttr   = "password"
	privateKeyAttr = "private_key"
	useAgentAttr   = "use_agent"
)

var (
	composeLock = &sync.Mutex{}
	sshPort     = 0
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			hostNameAttr: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "localhost",
			},
			portAttr: {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  22,
			},
			usernameAttr: {
				Type:     schema.TypeString,
				Optional: true,
			},
			passwordAttr: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			privateKeyAttr: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			useAgentAttr: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"jumphost_ssh": dataSourceSsh(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	port := d.Get(portAttr).(int)
	client := NewSshClient(
		d.Get(hostNameAttr).(string),
		d.Get(usernameAttr).(string),
		d.Get(passwordAttr).(string),
		d.Get(privateKeyAttr).(string),
		d.Get(useAgentAttr).(bool),
		port,
	)
	return &client, diags
}

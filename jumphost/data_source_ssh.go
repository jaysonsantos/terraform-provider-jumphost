package jumphost

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const localPortAttr = "local_port"

func dataSourceSsh() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceSshRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			hostNameAttr: {
				Type:     schema.TypeString,
				Required: true,
			},
			portAttr: {
				Type:     schema.TypeInt,
				Required: true,
			},
			localPortAttr: {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceSshRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	client, ok := m.(*SshClient)
	if !ok {
		diags = append(diags, diag.Errorf("the provider is not a valid SshClient")...)
		return diags
	}

	tunnel, err := client.GetTunnel(ctx, d)

	if err != nil {
		diags = append(diags, diag.FromErr(err)...)
		return diags
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	err = d.Set(localPortAttr, tunnel.LocalPort)
	if err != nil {
		diags = append(diags, diag.FromErr(fmt.Errorf("failed to set local_port %w", err))...)
	}

	return diags
}

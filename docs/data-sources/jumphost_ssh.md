---
page_title: "jumphost_ssh Data Source - terraform-provider-jumphost"
subcategory: ""
description: |-
  The jumphost_ssh data source uses an existing connection to a jumphost to start a dynamic tunnel to a target server/port using a random local port.
---

# Data Source `jumphost_ssh`

The jumphost_ssh data source allows you make tunnels using the jumphost to a remote server.

## Example Usage

```terraform
	data jumphost_ssh "http_server" {
		hostname = "remote-server"
		port = 8080
	}
```

## Attributes Reference

The following attributes are exported.
- `local_port` -  The local port assigned to connect to the remote server.

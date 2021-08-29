---
page_title: "Provider: Jumphost"
subcategory: ""
description: |-
  Terraform provider for to deal with jumphosts, creating dynamic tunnels using data resources.
---

# Terraform Provider
In case you have an infrastructure where you don't have a VPN and relies on jumphosts, this provider can help by connecting to the jumphost and creating arbitray connections by using data resources.

## Example Usage

Do not keep your authentication password in HCL for production environments, use Terraform environment variables.

```terraform
	provider jumphost {
		port = 22
    hostname = "localhost"
		username = "terraform"
		password = "1234"
	}
```

## Schema

### Optional

- **hostname** (String, Optional) Jumphost's hostname (defaults to `localhost`)
- **port** (Integer, Optional) Jumpost's port (defaults to `22`)
- **username** (String, Optional) Username to authenticate to the jumphost (in future it will try to guess from ssh config)
- **password** (String, Optional) Password to authenticate to the jumphost (can fallback to ssh agent)

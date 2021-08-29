# Terraform Provider Jumphost

In case you have an infrastructure where you don't have a VPN and relies on jumphosts, this provider can help by connecting to the jumphost and creating arbitray connections by using data resources.

## Build provider

Run the following command to build the provider

```shell
$ go build -o terraform-provider-jumphost
```

## Test sample configuration

First, build and install the provider.

```shell
$ make install
```

Then, navigate to the `examples` directory.

```shell
$ cd examples
```

Run the following command to initialize the workspace and apply the sample configuration.

```shell
$ terraform init && terraform apply
```

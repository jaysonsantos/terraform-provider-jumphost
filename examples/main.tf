terraform {
  required_providers {
    jumphost = {
      version = "~> 0.0.1"
      source  = "jaysonsantos/jumphost"
    }
    http = {
      version = "~> 2.1"
      source  = "hashicorp/http"
    }
  }
}

provider "jumphost" {
  hostname = "localhost"
  username = "terraform"
  port     = 2222
}

data "jumphost_ssh" "httpbin" {
  hostname = "httpbin.org"
  port     = "443"
}

data "http" "example" {
  url = "https://localhost:${data.jumphost_ssh.httpbin.local_port}/get"

  # Optional request headers
  request_headers = {
    Accept = "application/json"
  }
}

output "host" {
  value = jsondecode(data.http.example.body).headers.Host
}

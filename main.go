package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/terraform-providers/terraform-provider-ns1/ns1"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ns1.Provider})
}

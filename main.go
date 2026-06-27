// terraform-provider-bagre lets you manage Bagre IPAM as code.
//
// Built on the Terraform Plugin Framework (MPL-2.0). The same binary works with
// both OpenTofu and Terraform — they share the plugin protocol. Bagre is
// open-source, so this provider is developed OpenTofu-first while staying fully
// compatible with Terraform.
package main

import (
	"context"
	"flag"
	"log"

	"github.com/fabgcruz/terraform-provider-bagre/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// version is set at build time via -ldflags="-X main.version=…".
var version = "dev"

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "run the provider with support for debuggers like delve")
	flag.Parse()

	err := providerserver.Serve(context.Background(), provider.New(version), providerserver.ServeOpts{
		// Source address used to match the provider in OpenTofu/Terraform.
		Address: "registry.opentofu.org/fabgcruz/bagre",
		Debug:   debug,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}

package main

import (
	"context"
	"github.com/falcosecurity/kilt/runtimes/terraform/tf"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"log"
)

func main() {
	err := tfsdk.Serve(context.Background(), tf.New, tfsdk.ServeOpts{
		Name: "kilt",
	})
	if err != nil {
		log.Fatalf("terraform-provider-kilt plugin failed: %+v", err)
	}
}

package main

import (
	"context"
	"fmt"
	"github.com/falcosecurity/kilt/runtimes/cloudformation/cfnpatcher"
	"github.com/rs/zerolog"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s KILT_DEFINITION TEMPLATE\n", os.Args[0])
		return
	}
	kiltDef, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Cannot read kilt definition %s: %s\n", os.Args[1], err)
		return
	}

	template, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Cannot read template %s: %s\n", os.Args[2], err)
		return
	}

	config := &cfnpatcher.Configuration{
		Kilt:            string(kiltDef),
		ImageAuthSecret: "",
		OptIn:           false,
	}
	ctx := context.Background()
	l := zerolog.New(os.Stderr).With().Timestamp().Logger()
	ctx = l.WithContext(ctx)
	result, err := cfnpatcher.Patch(ctx, config, template)

	if err != nil {
		panic(fmt.Errorf("could not patch template: %w", err))
	}

	fmt.Printf("%s\n", string(result))

}

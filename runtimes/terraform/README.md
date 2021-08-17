# Terraform Kilt plugin

## Installation

### Building from source

1. Clone the repository and build the terraform provider:

       cd cmd/terraform-provider-kilt
       go build

2. Set the path to the kilt provider in your .terraformrc:

       provider_installation {
         dev_overrides {
            "falcosecurity/kilt" = "/.../full/path/to/cmd/terraform-provider-kilt"
         }

         direct {}
       }
3. Set `TF_CLI_CONFIG_FILE=/path/to/.terraformrc` with the above contents

## Usage

The kilt provider modifies your Fargate task definitions according to the spec in the kilt definition file.

### Enabling and configuring the provider

1. Add the `kilt` provider to your terraform module:

       terraform {
         required_providers {
           kilt = {
             source = "falcosecurity/kilt"
           }
         }
       }

2. Create an instance of the provider

       provider kilt {
         image_auth_secret = ""
       }

`image_auth_secret` is the only supported parameter. It is used to set RepositoryCredentials.CredentialsParameter for any containers injected by `kilt`.

### Using the provider to modify task definitions

Without the provider, your Fargate task definition will probably use the `aws_ecs_task_definition` resource, like this:

    resource "aws_ecs_task_definition" "test" {
      family                = "test"
      container_definitions = CONTAINER_DEFINITIONS_JSON
      # ...
    }

To pass the task definition through `kilt`,

1. Add a `kilt_fargate_container_definitions` data source:

       data "kilt_fargate_container_definitions" "sample" {
         container_definitions = CONTAINER_DEFINITIONS_JSON
         kilt_definition = KILT_DEFINITION
       }
2. Use the output of the data source as the new container definitions:

       resource "aws_ecs_task_definition" "test" {
         family                = "test"
         container_definitions = "${kilt_fargate_container_definitions.sample.output_container_definitions}"
         # ...
       }

If your kilt template uses variable substitutions, e.g. `${config.foo}`, you can set them in the `recipe_config` field of the data source:

       data "kilt_fargate_container_definitions" "sample" {
         container_definitions = CONTAINER_DEFINITIONS_JSON
         kilt_definition = KILT_DEFINITION
         recipe_config = {
           foo = "bar"
         }
       }

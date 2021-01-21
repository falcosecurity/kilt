Most of these steps are performed automatically by [Kilt](https://github.com/falcosecurity/kilt). In case you want to roll your own solution or test it out manually use the following instructions.

## Uploading sysdig instrumentation image onto personal aws account
In this example we're going to assume the target aws region to be us-east-1.

## Instrumenting manually an existing Task Definition
* Add 1 new containers to your task definition
    - The image name doesn't matter but we'll need it afterwards so we'll use FalcoInstrumentation
    - The entrypoint/command fields can be left empty
    - As for the image itself, use falcosecurity/falco-userspace:latest
* Add another container to your task definition
    - The image name doesn't matter but we'll need it afterwards so we'll use KiltUtils
    - The entrypoint/command fields can be left empty
    - As for the image itself, use falcosecurity/kilt-utilities:latest

* Edit the containers that you want to instrument
    - Add a startup dependency on the FalcoInstrumentation and KiltUtils containers created before
    - Mount volumes from FalcoInstrumentation and KiltUtils
    - Add `SYS_PTRACE` capability to your container. See [this](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_KernelCapabilities.html)
    - Set the following as your entry point: `/kilt/launcher,/vendor/falco/bin/pdig,$YOUR_COMMAND,--`
    - Set the following as your command: `/vendor/falco/bin/falco,-u,-c,/falco/falco.yaml,--alternate-lua-dir,/vendor/falco/share/lua`
    - Add environment variable `__CW_LOG_GROUP` to set the output log group 
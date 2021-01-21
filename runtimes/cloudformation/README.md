# Cloud Formation Kilt Runtime
This kilt runtime alters cloud formation templates to apply kilt definitions.
It installs a Cloud Formation Macro that will alter the incoming template.

##Components

* installer - installer for the macro
* cmd/handler - the golang lambda functions powering the Macro
* cmd/cfn-apply-kilt - applies kilt transformation to a CFN template
* cmd/cfn-image-info - gets configuration for the image from repository


# Usage
The installer will create a CFN macro that you can use to apply automatically
instrumentation to task definitions. To use it with a macro called `MyMacro` add
`Transform: MyMacro` or `Transform: ["MyMacro"]` to the root of your CFN Template.

There are 2 modes of operation for the macro, both selected during install. *opt-in*
and *opt-out*. You can use the following tags to include or exclude pieces of your 
task definition:

* `kilt-include:<any-value>` - will apply instrumentation in opt-in mode of operation
* `kilt-ignore:<any_value>` - will not apply instrumentation in opt-out mode of operation
* `kilt-include-containers:containerA,ContainerB` - value is comma separated list of 
  container names. Will include only some contaiers in opt-in mode
* `kilt-ignore-containers:containerA,containerB` - will exclude some containers in 
  opt-out mode
  

# Metadata offered
Inside the kilt definitions `original.metadata` is a hashmap containing tags of the
task definition.
# An Antidote for Messy Viper Configurations

A simplified use-pattern for the Cobra-Viper toolbox

## Motivation, or Why Extend Viper?

[Viper](https://github.com/spf13/viper) with
[Cobra](https://github.com/spf13/cobra) provides an extremely flexible, full
featured set of argument and option utilities for Go programs, equivalent to,
e.g., [Argparse](https://docs.python.org/3/library/argparse.html) in
Python-land. However, there are a dozen or so ways to utilize the Cobra plus
Viper combination.

The most straightforward and utilitarian use pattern is via the [Cobra command
generator](https://github.com/spf13/cobra/blob/master/cobra/README.md). With
`cobra init` and subsequent `cobra add` commands a clean scaffold for command
arguments and options is generated. However, for the common use case of a
12-factor harvest of startup settings, we can do better.

### What is Mithri?

Short for "[Mithridates'
antidote](https://en.wikipedia.org/wiki/Mithridates_VI#Mithridates'_antidote)",
for the curious.

### Why Mithri?

Besides the Hellanic anti-toxin reference, `mithri` is unlikely to have many
name collisions in the Go ecosystem.

### What are We Trying to Solve?

Some stumbling blocks we've seen in ad-hoc Viper builds:

- We run different services (http, message drive jobs, serverless) with the same
  codebase and build. Some configuration values are common, some
  unique. However, we mash them all together and it's unclear what config is
  used with what service.
- A developer is often setup for many different service codebases
  locally. Generic configuration names in the environment, `DATABASE_URL` for
  example, are a problem if used in multiple codes and a different value is
  needed.
- It's hard for a developer not intimately familar with the configuration of a
  codebase to know what values are required, what's optional, or even what all
  config keys are in play. Then, to fish out config value types and find where
  they're accessed and how they're used becomes a real task.
- A `config-sample` file in the code repo is helpful, but brittle to
  maintain. Such a file is often sparsely commented and new or changed
  configuration easily forgotten.

## Try Viper Mithri Out

This repository can be built as a small usage sample, for experiment, test, and
debug. In fact, Viper has a number of idiosyncracies that are much easier to
work out in the small.

```console
git clone git@github.com:aspiration-labs/viper-mithri.git
cd viper-mithri
go run main.go --help
```

```console
A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
  viper-mithri [flags]
  viper-mithri [command]

Available Commands:
  config-tool       Read, unmarshal, then write config to file or '-' for stdout
  help              Help about any command
  serve             A brief description of your command
  serve-config-tool Read, unmarshal, then write config to file or '-' for stdout

Flags:
      --config-env string          use env with prefix
      --config-file string         config file to read
      --config-type string         format of config (default is yaml) (default "yaml")
  -h, --help                       help for viper-mithri
      --serve-config-file string   config file to read
      --serve-config-type string   format of config (default is yaml) (default "yaml")

Use "viper-mithri [command] --help" for more information about a command.
```

### Ok, Now What?

Let's work the sample app from the command line, then we'll walk through the
code behind it. Mithri automatically adds commands like `config-tool` to the
Cobra based invocations. The point of `config-tool` and `serve-config-tool` are
to easily inspect, save, and modify configurations based on what the code says
can be configurated.

Try `config-tool`

```console
go run main.go config-tool
```

Which outputs

```yaml
api_url: http://localhost/api
auth_password: 12fa
auth_username: zzyzx
```

That's the complete set of configurables for the root command. The defaults are
compiled in values, nominal settings to get started with, but that certainly
will need to be overridden, both for local development and deployed
operation. Let's modify these built in defaults.

```console
go run main.go config-tool >config.yaml
```

Then edit `config.yaml` by deleting the `auth_password` and `auth_username`
lines and changing the `api_url` to `https://server/api`. Your new `config.yaml`
should look like

```yaml
api_url: https://server/api
```

Rerun `config-tool` with

```console
go run main.go config-tool --config-file config.yaml
```

You'll get

```yaml
api_url: https://server/api
auth_password: 12fa
auth_username: zzyzx
```

The `api_url` setting came from `config.yaml`, the rest from the compiled in
defaults. Config files have higher priority than builtin defaults. This builtin
plus config file is often enough for local development.

### What About Environment Variables?

Yes, environment variables are the preferred option in deployment. But, by
default, they are turned _off_ in Mithri. Try this.

```console
AUTH_USERNAME=amboy AUTH_PASSWORD=secret go run main.go config-tool --config-file config.yaml
```

...which produces the same output as the previous example. How do we get
environment variables into play? Environment variables must be called in
explicitly with a command line option, and that option _requires a prefix_.

```console
MITHRI_AUTH_USERNAME=amboy MITHRI_AUTH_PASSWORD=secret go run main.go config-tool --config-file config.yaml --config-env mithri
```

and...voilà. In case you missed it, check the `--config-env mithri`, which both
turned on environment vars and set a prefix to prepend (with `_`) to the config
keys.

```yaml
api_url: https://server/api
auth_password: secret
auth_username: amboy
```

Finally, run the root command of the app, without `config-tool`

```console
MITHRI_AUTH_USERNAME=amboy MITHRI_AUTH_PASSWORD=secret go run main.go --config-file config.yaml --config-env mithri
```

The root command is setup to output the app config _from a go struct_.

```console
root called with cmd.RootConfig{ApiUrl:"https://server/api", Hostname:"", Auth:cmd.RootAuthConfig{Username:"amboy", Password:"secret"}}, cmd.ServeConfig{ServePort:8080, ServeHost:"127.0.0.1"}
```

The Mithri practice is to access configuration values in code from a Go struct,
and not a bunch of `GetThis()` and `GetThat()` (Viper joke there) sprinkled
about the startup code. Instead, you should see lines like

```go
http.ListenAndServe(":" + ServeAppConfig.Port, nil)
```

Before we move on to code, notice the `cmd.ServeConfig` in the output above and
recall there's a `serve` subcommand, along with `server-config-tool`. Mithri
allows for multiple parallel configurations and generates a distinct config tool subcommand for each. Try

```console
go run main.go serve-config-tool
```

and see what you get. Feel free to experiment a bit.

## Integrating Viper Mithri into Your Cobra App

Let's look at how Mithri integrates into the sample app. The recommended pattern
is based on the one in this app.

The Mithri sample app was kicked off with the Cobra generator, like so:

```console
mkdir test-mithri
go mod init test-mithri
cobra init --pkg-name test-mithri
cobra add serve
```

Then, we added the config to the root command, `cmd/root_config.go` in the
sample app. Picking that source file apart:

```go
type RootConfig struct {
	ApiUrl string `mapstructure:"api_url"`
	Hostname string `mapstructure:"hostname"`
	Auth RootAuthConfig `mapstructure:",squash"`
}

type RootAuthConfig struct {
	Username string `mapstructure:"auth_username"`
	Password string `mapstructure:"auth_password"`
}

var RootAppConfig RootConfig
```

The composite (nested) struct `RootConfig` defines a type for our
configuration. `RootAppConfig` is the instance Mithri will unmarshal the final
values into, after resolving compiled defaults, config file and environment
vars.

We use the [mapstructure](https://godoc.org/github.com/mitchellh/mapstructure)
package (per the
[Viper.Unmarshal](https://pkg.go.dev/github.com/spf13/viper#Unmarshal) function)
to define the struct mappings.

**Note** especially the use of the `squash` option, to omit the composition of
`RootAuthConfig` inside the `RootConfig` struct. In this case we manually add
`auth_` to the `mapstructure` tags in `RootAuthConfig`. While many use cases
don't require nested structs for configuration purposes, it can be
convenient. When used though it his highly recommend to squash, flatten, and
manually guarantee no collisions in the `mapstructure` tags.

On to the default values.

```go
var rootDefaults = map[string]interface{}{
	"api_url": "http://localhost/api",
	"auth_username": "zzyzx",
	"auth_password": "12fa",
}
```

...a straightforward mapping into the `mapstructure` tagged `RootConfig` type.

And finally note that we've made `RootAppConfig` externally available to other
modules in the application. The defaults, `rootDefaults` here should only be
needed once, within the `cmd` package. That part is up next.

### Hooking the Config Struct and Defaults Using Mithri

The code integration is easy, Cobra friendly, and the reason this modest package
was developed. (That's in addition to addressing the problems we're trying to
solve.)

If you've just started an app with the Cobra generator you'll have a

`cmd/root.go` with an `init()` that looks like

```go
func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.test-mithri.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
```

You can replace the contents to look like

```go
func init() {
	mithri.AddCommand(rootCmd, rootDefaults, &RootAppConfig, "")
}
```

...where `rootDefaults` and `RootAppConfig` come from the config setup you've
built. (Be very sure the 3rd arg to `mithri.AddCommand` is by reference, with
an `&`; that goes into an empty interface type which will not type check
whether a struct or a pointer is passed.)

For the `serve` subcommand it's about the same, but make sure to preserve the
Viper `AddCommand` in `serve.go`.

```go
func init() {
	rootCmd.AddCommand(serveCmd)
	mithri.AddCommand(rootCmd, serveDefaults, &ServeAppConfig,"serve")
}
```

## And That...

Should be all there is to it.

## Acknowledgements

- [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper) (obviously)
- [Managing configuration with Viper](https://scene-si.org/2017/04/20/managing-configuration-with-viper/) for the inspiration to make those ideas more better
- [Charlie McElfresh](https://github.com/charliemcelfresh) for framing the problem and great ideas

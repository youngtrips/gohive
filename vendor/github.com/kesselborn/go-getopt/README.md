getopt for go
=============

This is a getopt implementation for go with the following feature set:

 * proper short- and long-opt support
 * automatically create help and usage text
 * set options via environment variables
 * config file support

Installation
------------

    goinstall github.com/kesselborn/go-getopt

Usage example
-------------

The following examle is included in the `example` folder, in order to run it, install go-getopt, `cd example` and do a

    make

and try out different calls (`make` will give instructions)

### Source code

    package main

    import (
      "fmt"
      "os"
      getopt "github.com/kesselborn/go-getopt"
    )

    func main() {
      optionDefinition := getopt.Options{
        {"debug|d|DEBUG",            "debug mode",             getopt.Optional | getopt.Flag,                 false},
        {"config|c",                 "config file",            getopt.IsConfigFile | getopt.ExampleIsDefault, "./config_sample.conf"},
        {"ports|p|PORTS",            "ports",                  getopt.Optional | getopt.ExampleIsDefault,     []int64{3000, 3001, 3002}},
        {"sports|s|SECONDARY_PORTS", "secondary ports",        getopt.Optional | getopt.NoLongOpt,            []int{5000, 5001, 5002}},
        {"instances||INSTANCES",     "instances",              getopt.Required,                               4}
        {"keys||KEYS",               "keys",                   getopt.Required,                               []{"foo", "bar", "baz"}},
        {"logfile||LOGFILE",         "logfile",                getopt.Optional | getopt.NoEnvHelp,            "/var/log/foo.log"},
        {"file",                     "files",                  getopt.IsArg,                                  ""}
        {"directories",              "directories",            getopt.IsArg | getopt.Optional,                ""},
        {"pass through",             "pass through arguments", getopt.IsPassThrough | getopt.Optional,        ""},
      }

      options, arguments, passThrough, e := optionDefinition.ParseCommandLine()

      if e != nil {
        description := "this is a small sample application for getopt demonstration"
        exit_code := 0
        switch {
          case e.ErrorCode == getopt.WantsUsage:
            fmt.Print(optionDefinition.Usage())
          case e.ErrorCode == getopt.WantsHelp:
            fmt.Print(optionDefinition.Help(description))
          default:
            fmt.Println("**** Error: ", e.Message, "\n", optionDefinition.Help(description))
            exit_code = e.ErrorCode
        }
        os.Exit(exit_code)
      }

      fmt.Printf("options:\n")
      fmt.Printf("debug: %#v\n",          options["debug"].Bool)
      fmt.Printf("config: %#v\n",         options["config"].String)
      fmt.Printf("ports: %#v\n",          options["ports"].IntArray)
      fmt.Printf("secondaryports: %#v\n", options["sports"].IntArray)
      fmt.Printf("instances: %#v\n",      options["instances"].Int)
      fmt.Printf("keys: %#v\n",           options["keys"].StrArray)
      fmt.Printf("logfile: %#v\n",        options["logfile"].String)
      fmt.Printf("files: %#v\n",          options["files"].StrArray)

      fmt.Printf("arguments: %#v\n", arguments)
      fmt.Printf("passThrough: %#v\n", passThrough)
    }

### Help output
Calling the help of this programs generates the following output:

    $ ./getopt-sample-app --help
    Usage: getopt-sample-app [-d] -c <config> [-p <ports>] [-s <sports>] --instances=<instances> --keys=<keys> [--logfile=<logfile>] <file> [<directories>] [-- <pass through>]
    this is a small sample application for getopt demonstration

    Options:
        -d, --debug                       debug mode (e.g. false); setable via $DEBUG
        -c, --config=<config>             config file (default: ./config_sample.conf)
        -p, --ports=<ports>               ports (default: 3000,3001,3002); setable via $PORTS
        -s <sports>                       secondary ports (e.g. 5000,5001,5002); setable via $SECONDARY_PORTS
            --instances=<instances>       instances (e.g. 4); setable via $INSTANCES
            --keys=<keys>                 keys (e.g. foo,bar,baz); setable via $KEYS
            --logfile=<logfile>           logfile (e.g. /var/log/foo.log)
        -h, --help                        usage (-h) / detailed help text (--help)

    Arguments:
        <file>                            files
        <directories>                     directories

    Pass through arguments:
        <pass through>                    pass through arguments

play around with the program to see how it behaves with ENV variables and the different included configuration files --
there should not be any surprises.

Options struct explained
------------------------

The options you pass to the `Options` struct have the following structure:

    {"<longopt>|<shortopt>|<ENVVAR>", "<description for help text>", <options>         , <default or example value>}

  * **&lt;longopt&gt;**: the long option name that can passed to your program with
`--<long_opt>`. Furthermore, this value will be the key under which this
value is available in the options map. Long opt values need to be separated
by a whitespace or an '=': `--logfile /tmp/log.txt` or `--logfile=/tmp/log.txt`.
If you don't want this option to have a long-opt style, pass `getopt.NoLong` in the options.
  * **&lt;shortopt&gt;**: short option letter ... leave it out if you only want a
long opt style for this option. Short opt values can be separated by a
whitespace: `-l /tmp/log.txt' or nothing: `-l/tmp/log.txt'
  * **&lt;ENVVAR&gt;**: if you want to let users set this option via an
environment varialbe, put the name of the env variable here. If you want
long opt style + env variable but not short opt style, pass in
`"<longopt>||ENV_VAR" `
  * the string in description will be used to create the help text
  * **&lt;options&gt;**:
    * **getopt.Required**: this options is required. If it is not passed in, `ParseCommandLine`
will return an error
    * **getopt.Optional**: can be set
    * **getopt.Flag**: this option does not have a value ... it'll toggle the default
value
    * **getopt.NoLongOpt**: don't accept longopt for this option
    * **getopt.ExampleIsDefault**: needs to go along with **getopt.Optional** -- if
this option is not set, its defaut value is taken (see <default or example value>)
    * **getopt.IsArg**: this is an argument, not an option (see example above: file / directories)
    * **getopt.Usage**: this is the shortopt option that shows the usage text (if no
options has this flag, '-h' is the default option for usage)
    * **getopt.Help**: this is the longopt option that shows the help text (if no options
has this flag, '--help' is the default option for help)
    * **getopt.IsPassThrough**: pass through arguments
    * **getopt.IsConfigFile**: this option accepts the name of a config file
    * **getopt.NoEnvHelp**: don't show ENV variable help text for this option
  * **&lt;default or example value&gt;**: is a default value for optional options that have the
`getopt.ExampleIsDefault` flag set and an example vaule for required options.
Can be empty strings. For optional options, if a `nil` example is set, the
`options` map won't contain an entry if this option is not passed in by a user.
If set to a value different to nil, the `options` map will contain the default value
if the user does not pass in the option. Types of the values will be saved as stated here.

Configuration File
------------------
I you want your program to have a configuration file, you can use the included config file parser.
Configuration have to be of the form of bash env var definitions -- all other lines are excluded:

    # this line will be ignored as it is not of the form /^[A-z0-9_.,]+=.*$/
    PORT=8121
    HOST=zookeeper.apache.org

the option for setting the config file needs to get the `IsConfigFile` flag (see example)

The ParseCommandLine method
----------------------------

The ParseCommandLine method returns four values:

 * **options**: a map with the parsed options; access an option by doing `options["<longopt>"].<type>`
where type is one of the following (depending on the example value you passed in for the respective
option:
   * `Bool` of type `bool`
   * `String` of type `string`
   * `Int` of type `int64`
   * `StrArray` of type `[]string`
   * `IntArray` of type `[]int64`
 * **arguments**: arguments the command received of type `[]string`
 * **passThrough**: pass through arguments of type `[]string`that were passed
after the delimiting `--`, helpful if you want to pass through parameters for a sub command
 * **err**: error of type `*GetOptError`, is `nil` if no error occured

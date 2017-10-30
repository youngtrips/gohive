
0.4.1 / 2012-03-30 
==================

  * don't mask error messages if scope is missing

0.4.0 / 2012-03-27 
==================

  * this release introduces SubCommandOptions and SubSubCommandOptions for git
    like subcommand parsing. The naming will most probably change in the future
    as it is kind of bad. Code works but needs to be refactored

  * WantsHelp and WantsUsage is not an error anymore but is added to the options
    map


0.3.8 / 2012-01-06
==================

  * make it possible to have more than one passThrough arg in the usage text; don't include the '--' in the brackets if passThrough is optional

0.3.7 / 2012-01-06
==================

  * make sure help or usage options in passThrough section are not displaying help or usage

0.3.6 / 2012-01-04
==================

  * add readme with example calls for example program
  * explain config file support in readme

0.3.5 / 2012-01-03
==================

  * make sure that we don't have required arguments after optional arguments

0.3.4 / 2012-01-03
==================

  * proper error message when a required argument was not provided

0.3.2 / 2012-01-03
==================

  * check for error message

0.3.1 / 2012-01-03
==================

  * if option is ExampleIsDefault, handle it as an Optional option

0.3.0 / 2011-12-20
==================

  * split up UsageOrHelp errors in WantsUsage and WantsHelp
  * change ParseCommandLine(description string, flags int) -> ParseCommandLine()
  * add example conf files; add NoEnvHelp
  * fix bug: env var from NoLongOpt option wasn't recognized correctly

0.2.1 / 2011-12-19
==================

  * add readme disclaimer
  * make sure IsConfigFile is an option that has a value

0.2.0 / 2011-12-19
==================

  * Merge branch 'feature/config-file-support' into develop
  * implement config file functionality
  * implement config parsing
  * fix equalStringArray + tests that passed because equalStringArray was buggy

0.1.0 / 2011-12-18
==================

  * change Parse(args []string, defaults []string, description string, flags int) -> ParseCommandLine(description string, flags int) ... read args and env by os.Args and os.Environ
  * use path/filepath.Base function

0.0.2 / 2011-12-13
==================

  * move from VALUE to <value>
  * goinstall info

0.0.1 / 2011-12-13
==================

  * use gofmt
  * export variables
  * fixes
  * implement first version

0.0.0 / 2011-11-28
==================

  * first commit

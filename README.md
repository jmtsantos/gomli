# gomli

![gomli](assets/gomli.png)

`go + smali ~= gomli` a dynamic smali patching utility

## Usage

A simple replace operation `gomli replace -s script.js -d work/com.example_decompiled`

Or you are welcome to use `--help`

```
A go parser for smali code

Usage:
  gomli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  find        find call chains XREFs for a particular class
  help        Help about any command
  replace     replace the values of const-strings globally in the application

Flags:
  -h, --help      help for gomli
  -v, --verbose   enable verbose/debug mode

Use "gomli [command] --help" for more information about a command.
```


## Replace

In order do maintain some code consistency when using a decompiler tool like `jadx` specific rules when it comes to value replacement have to be considered to each instruction, obviously `const-string` is the most simple just replace the string argument from the instruction. Arrays `fill-array-data` are a bit trickier.

### Replacing `const-string`

### Replacing `fill-array-data`

## TODO

- [ ] remove newlines
- [ ] Allow direct ingestion of APK files
- [ ] More efficient search process when running the execute command
- [ ] Auto frida script generation based on search rules

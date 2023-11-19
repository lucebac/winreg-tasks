# winreg-tasks

This repository contains structure definitions and some tooling for the BLOBs found in the TaskCache registry key on Windows (`HKLM\Microsoft\Windows NT\CurrentVersion\Schedule\TaskCache`). [I wrote a blog post](https://cyber.wtf/?p=1923) which gives some background knowledge on what this is all about.

This tool can be run in two different modes:
* accessing live Registry data (this requires administrative privileges so that all registry keys can be accessed)
* using a raw Hive file (log files can be applied if provided)

Since v0.4.0, winreg-tasks is platform-independent with official support for both Windows and Linux. You probably can also compile it for any other platform supported by Golang, but official support is only provided for Windows and Linux.

winreg-tasks has several features. Consult the help pages for more specific information on the sub commands. The main help page lists all available commands for your platform (example shows the help page on Linux using winreg-tasks v0.5.0):

```
Usage: winreg-tasks-linux-amd64 <command>

Flags:
  -h, --help                       Show context-sensitive help.
  -f, --file=FILE                  If provided, use this Hive file instead of the System's live one.
  -x, --log-files=LOG-FILES,...    If provided, these log files will be applied to the Hive file.

Commands:
  actions <task-id>
    Dump the Actions of a given Task.

  dump
    Dump the Task list to a file

  dynamicinfo <task-id>
    Dump the DynamicInfo of a given Task.

  parseall
    Parses all existing Tasks

  triggers <task-id>
    Dump the Triggers of a given Task.

Run "winreg-tasks-linux-amd64 <command> --help" for more information on a command.
```

Here's a list of example commands:
```powershell
# display the available commands:
.\winreg-tasks.exe [-h|--help]

# iterates all tasks and prints a list of Actions, Triggers,
# and the DynamicInfo of all tasks registered to the system
.\winreg-tasks.exe parseall
# same as before but only prints errors; most useful when you
# changed something and want to see if anything broke
.\winreg-tasks.exe parseall -q

# dumps all tasks on the system and prints the result to stdout
.\winreg-tasks.exe dump
# dumps all tasks and writes the data into a JSONL file
.\winreg-tasks.exe dump -o tasks.jsonl

# get a more detailed dump of the actions of a task:
.\winreg-tasks.exe actions '{00000000-1111-2222-3333-444444444444}'
# you can pass the path alternatively (leading backslash required!):
.\winreg-tasks.exe actions '\My Task'

# get the triggers of a given task:
.\winreg-tasks.exe triggers '{00000000-1111-2222-3333-444444444444}'
# you can pass the path alternatively (leading backslash required!):
.\winreg-tasks.exe triggers '\My Task'

# get the DynamicInfo of a given task:
.\winreg-tasks.exe dynamicinfo '{00000000-1111-2222-3333-444444444444}'
# you can pass the path alternatively (leading backslash required!):
.\winreg-tasks.exe dynamicinfo '\My Task'
```

Several commands support the `-d` or `--dump` flag which prints the data read from the registry key as a 16 bytes wide hex dump. I found it much more easy to work with these dumps than exporting a value of a key with regedit and then converting the `hex:00,11,22,...` notation to something more readable.

## Generate

If you want to re-generate the source files, just use the `generate.sh` script. If you need to output the source for another language, please adapt the script to your needs. The script requires Kaitai v0.10 or later to function properly, please refer to the Kaitai's documentation for [installation instructions](https://kaitai.io/#download).


## Build

If you did not change anything and just want to use the tool, simply download and run any pre-built executable files from the releases.

If you want to build the executables by yourself or want introduced changes to the source code, you can use the provided Makefile; all output files are written to the `out` folder. Just make sure, you have a working installation of Golang 1.21 (or later) and then run:
```bash
make
```
The Makefile has several targets (e.g. `windows` or `linux`) which limit the executables that will be built. Use `make targets` to display a full list of targets.

## Install From Source
Alternatively, you can just install and run the package from source:
```bash
go install github.com/lucebac/winreg-tasks/cmd@latest
```

# Using the Generated Code
To get started, add winreg-tasks to you dependencies: `go get github.com/lucebac/winreg-tasks`.

After that, you can use the `FromBytes` methods of any of the packages (e.g. `actions`, `dynamicinfo`, `triggers`) to convert the raw byte representation of these information to a nicely accessible Go structure.

Minimum example for Windows (except for the Registry access in this example, there's no difference in the use of the package between the platforms):
```golang
package main

import (
	"log"

	"github.com/lucebac/winreg-tasks/dynamicinfo"
	"golang.org/x/sys/windows/registry"
)

func main() {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion\Schedule\TaskCache\Tasks\{B75AF762-3C5C-4C74-ADB1-B99F98FDE0E5}`, registry.QUERY_VALUE)
	if err != nil {
		log.Fatalf("cannot open task key: %v", err)
	}
	defer key.Close()

	dynamicInfoRaw, _, err := key.GetBinaryValue("DynamicInfo")
	if err != nil {
		log.Fatalf("cannot get dynamic info for task: %v", err)
	}

	dynamicInfo, err := dynamicinfo.FromBytes(dynamicInfoRaw)
	if err != nil {
		log.Fatalf("cannot parse bytes into DynamicInfo: %v", err)
	}

	lastErrorCode := dynamicInfo.LastErrorCode
	log.Printf("Last Error Code: 0x%08x", lastErrorCode)
}


```


# Licensing
winreg-tasks is released under the MIT license. See the LICENSE file for more information.

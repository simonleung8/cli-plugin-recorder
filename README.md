CF CLI Recorder [![Build Status](https://travis-ci.org/simonleung8/cli-plugin-recorder.png?branch=master)](https://travis-ci.org/simonleung8/cli-plugin-recorder)
=================
This plugin records sets of CF CLI commands, and allows you to playback a set or sets commands anytime.


##Usage
```
$ cf record <name>

>> {enter cf commands as usual} ...
>> {enter cf commands as usual} ..
>> stop
```
After recording, play back with `replay`, you can play back 1 or more recorded command sets.
```
$ cf replay <name>
```

##Installation
#####Install from CLI (v.6.10.0 and up)
  ```
  $ cf add-plugin-repo CF-Community http://plugins.cloudfoundry.org/
  $ cf install-plugin CLI-Recorder -r CF-Community
  ```
  
  
#####Install with binary
- Download the binary [`win64`](https://github.com/simonleung8/cli-plugin-recorder/raw/master/bin/win64/cli-recorder.exe) [`linux64`](https://github.com/simonleung8/cli-plugin-recorder/raw/master/bin/linux64/cli-recorder.linux64) [`osx`](https://github.com/simonleung8/cli-plugin-recorder/raw/master/bin/osx/cli-recorder.osx)
- Install plugin `$ cf install-plugin <binary_name>`



#####Install from Source (need to have [Go](http://golang.org/dl/) installed)
  ```
  $ go get github.com/cloudfoundry/cli
  $ go get github.com/simonleung8/cli-plugin-recorder
  $ cd $GOPATH/src/github.com/simonleung8/cli-plugin-recorder
  $ go build -o cli-recorder main.go
  $ cf install-plugin cli-recorder
  ```

##Full Command List

| command | usage | description|
| :--------------- |:---------------| :------------|
|`record`| `cf record Cmd_Name` |record a set of commands|
|`record -l`|`cf record -l`|list all recorded command sets|
|`record -n`|`cf record -n <Cmd_Name>`|list all commands within a set|
|`record -d`|`cf record -d <Cmd_Name>`|delete a recorded command set|
|`record -clear`|`cf record -clear`|delete all recorded command sets|
|`replay`|`cf replay <Cmd_Name...>`|replay a command set or sets|
|`rp`|`cf rp <Cmd_Name...>`|alias of `replay`|

##Help Command

| command | usage | description|
| :--------------- |:---------------| :------------|
|`record -h`| `cf record -h` |show `record` usage|


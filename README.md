CF CLI Recorder
=================
This plugin records sets of CF CLI commands, and allows you to replay a set or sets command anytime.

For example
```shell
cf record bosh-lite
Please start entering CF commands
For example: 'cf api http://api.10.244.0.34.xip.io --skip-ssl-validation'

type 'stop' to stop recording and save
type 'quit' to quit recording without saving

(recording) >> cf api http://api.10.244.0.34.xip.io --skip-ssl-validation

Setting api endpoint to http://api.10.244.0.34.xip.io...
Warning: Insecure http API endpoint detected: secure https API endpoints are recommended

OK


API endpoint:   http://api.10.244.0.34.xip.io (API version: 2.22.0)
Not logged in. Use 'cf login' to log in.

(recording) >> cf auth admin admin

API endpoint: http://api.10.244.0.34.xip.io
Authenticating...
OK
Use 'cf target' to view or set your target org and space

(recording) >> stop
```
Now you can replay by `cf replay bosh-lite`

Commands
===============


Options

| command | usage | description|
| :---------------: |:---------------:| :------------:|
|`record`| `cf record Cmd_Name` |record a set of commands|
|`record -l`|`cf record -l`|list all recorded command sets|
|`record -n`|`cf record -n [Cmd_Name]`|list all commands within a set|
|`record -d`|`cf record -d [Cmd_Name]`|delete a recorded command set|
|`record -clear`|`cf record -clear`|delete all recorded command sets|
|`replay`|`cf replay [Cmd_Name...]`|replay a command set or sets|
|`rp`|`cf rp [Cmd_Name...]`|alias of `replay`|

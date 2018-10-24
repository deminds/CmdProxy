# CmdProxy
Execute local shell commands and manage remote telnet devices via this http service

## Example usage

#### CURL

##### Connect
```
curl -v -X GET http://localhost:25505/api/v1.0/console/connect
```

##### Execute command
```
curl -v -d '{"sessionid":"219602104153538926", "command":"ls -lah /home/"}' -X POST http://localhost:25505/api/v1.0/console/command
```

##### Disconnect
```
curl -v -X GET http://localhost:25505/api/v1.0/console/disconnect?sessionid=219602104153538926
```

##### Test Handlers
You can test *CmdProxy* via tool `testHandler.py`  
Use `./testHandler.py -h` for more information

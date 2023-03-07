# Group Checker

Simple API for checking if user is part of the user-check ldap group and counting users in user-check group.

- check if user is part of the user-check group (returns true or false)
- count users in user-check group (returns int/number of users in user-check group)

# Building

After cloning the repo, you can:

## Option 1: Build locally

```shell
cd src
go get -d -v ./...
go build -o user-check
```

Or:

```shell
swag init --parseDependency && go build main.go && ./main -s -d 
```

__!!!NB!!!__ swag is required for this

```shell
go install github.com/swaggo/swag/cmd/swag@v1.8.7
```

To change the port use -p option


## Options 2: Build the docker image

```shell
cd $GOPATH/src
docker build -t user-check -f go-user-check-api/user-check.Dockerfile .
```

# Local run

Execute the binary and pass necessary command-line parameters.

```shell
./user-check [opts]
```

Examples

```shell
./user-check -d -s
```

(open browser: http://localhost:8080/swagger/index.html#/)

# Command line parameters

You may specify a number of command-line parameters which change the behavior of the application

| Short | Long | Default | Usable in prod | Description |
|-----|-----|-----|-----|-----|
| -t | --timeout | 60 | Yes | Time to wait for graceful shutdown on SIGTERM/SIGINT in seconds |
| -p | --port | 8080 | Yes | TCP port for the HTTP listener to bind to |
| -s | --swagger | | No | Activate swagger. Do not use this in Production! |
| -d | --devel | | No | Start in development mode. Implies --swagger. Do not use this in Production! |
| -l | --tls | | No | Enable TLS. Implies having cert and key files.Use this in Production! |

# Environment variables and options

Default values used by user-check

(check configuration.go file)

```go
// ldap server
appConfig.LdapServerAddress = utils.EnvOrDefault("LDAP_ADDR", "ldaps://server.com:636")
// NPA account info
// user
appConfig.NpaUser = utils.EnvOrDefault("NPA_USER", "npa@domain.com")
// password
// yes, i know how to use a vault and store password there, but this prj waay tooo simple
appConfig.NpaPassword = utils.EnvOrDefault("NPA_PASSWORD", "XXXX")
// search people
appConfig.SearchPeople = "OU=eCore Office,OU=People Accounts,DC=domain,DC=com"
// user-check group
appConfig.user-checkGroup = utils.EnvOrDefault("user-check_GROUP", "mygroup")
// certificate file
appConfig.CertFile = utils.EnvOrDefault("LDAP_CERT_FILE", "cert.crt")
// API SSL CRT file
appConfig.ApiCertCrtFile = utils.EnvOrDefault("API_CERT_CRT_FILE", "server.crt")
appConfig.ApiCertKeyFile = utils.EnvOrDefault("API_CERT_KEY_FILE", "private.key")
```

Export new env vars if different values are required

Eg:

```shell
export LDAP_ADDR=ldaps://.....:636
```


# LDAP

## LDAP ENDPOINT

go-user-check is using internal ldap to validate users and count 

Ldap server endpoint can be set via envvars

configuration.go file:

```shell
// ldap server
appConfig.LdapServerAddress = utils.EnvOrDefault("LDAP_ADDR", "ldaps://server.com:636")
```

```shell
export LDAP_ADDR=ldaps://server.com:636
```


# TLS

## Enable tls

In order to enable TLS use -l  (or --tls) param

```shell
./main -l
```

By default, the API server is using server.crt and server.key files to provide secure connection.


```shell
appConfig.ApiCertCrtFile = utils.EnvOrDefault("API_CERT_CRT_FILE", "server.crt")
appConfig.ApiCertKeyFile = utils.EnvOrDefault("API_CERT_KEY_FILE", "private.key")
```

To use different names for the files:

```shell
export API_CERT_CRT_FILE="my_server.crt"
export API_CERT_KEY_FILE="my_private.key"
```

## Check if TLS is working

```shell
 curl -X 'GET' 'https://localhost:8080/api/v1/status' -H 'accept: application/json'
```

# API request sample

## Health checkpoint

```shell
curl -X 'GET' \
  'http://localhost:8080/api/v1/status' \
  -H 'accept: application/json'
```

## Check if user is part of user-check group

```shell
curl -X 'GET' \
  'http://localhost:8080/api/v1/usercheck/bordeanu' \
  -H 'accept: application/json'
```

## Count users in user-check ldap group

```shell
curl -X 'GET' \
  'http://localhost:8080/api/v1/usercount/' \
  -H 'accept: application/json'
```

# Run functional tests

## Run all tests

```shell
cd test
go test -v 
```

## Run a specific test

```shell
cd test
go test -run TestUserExists
```

## Mock the API

### Build and start

```shell
cd test
go build mock.go
./mock -p 9999
```

### Usage and examples

```shell
Usage of ./mock:
  -c, --howmany int32         HowMany users are in the security group. Default:1000 (default 1000)
  -l, --myldapstatus string   Status of the API. Default:up (default "up")
  -p, --port int32            TCP port for the HTTP listener to bind to. Default: 8082 (default 8080)
  -u, --usercheck string      User to check if exists. Default:bordeanu (default "bordeanu")
  -r, --userresponse string   Response value checking user. Default:true (default "true")
```


```shell
./mock -p 8080 -c 100 -l up -u bordeanu -r true
```

This will return:

```shell
curl -X 'GET' http://localhost:8080/api/v1/usercount -H 'accept:application/json'
{"code":200,"id":"mock-62a1dee8-1acf-429d-91c0-eefa95b62371","message":"Success","data":100}

curl -X 'GET' http://localhost:8080/api/v1/status -H 'accept:application/json'
{"code":200,"id":"f21609a2-643a-4dc4-9c30-7e63c08d8283","message":"Success","data":{"LdapStatus":"up","ProcessId":1234}}

 curl -X 'GET' http://localhost:8080/api/v1/usercheck/bordeanu -H 'accept:application/json'
{"code":200,"id":"mock-62a1dee8-1acf-429d-91c0-eefa95b62371","message":"Success","data":"true"}
      
```


### Run curl


Query if user is part of the group

```shell
curl -X 'GET' http://localhost:8082/api/v1/usercheck/test -H 'accept:application/json'
{"Code":"200","Id":"mock-62a1dee8-1acf-429d-91c0-eefa95b62371","Data":"false","Message":"Success"}
```

Count users in aad group

```shell
curl -X 'GET' http://localhost:8082/api/v1/usercount -H 'accept:application/json'
{"Code":"200","Id":"mock-62a1dee8-1acf-429d-91c0-eefa95b62371","Data":999999,"Message":"Success"}%
```

Check ldap status

```shell
curl -X 'GET' http://localhost:8082/api/v1/status -H 'accept:application/json'
{"Code":"200","Id":"mock-62a1dee8-1acf-429d-91c0-eefa95b62371","Data":{"LdapStatus":"up","ProcessId":1234},"Message":"Success"}
```

## Goconvey UI for tests

### Install goconvey

```shell
go get github.com/smartystreets/goconvey
```

### Run tests in UI

```shell
goconvey -port 8099
```

## Test endpoint value

__!!!NB__ tested endpoint can be changed via env vars

```shell
export TEST_ENDPOINT="http://localhost:8080/"
```

# Master mind

Dan Bordeanu
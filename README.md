# Token-Generator

** NOTE **

Please note that this project is not production ready.

## Summary

This service is desgined to return Oauth2 access tokens. These access tokens can be used as a form of validation, on behalf of the user. Thus, it is important to ensure that these tokens are kept confidential. Please refer to the API section to identify the endpoints to obtain a token.

## Initialization

Make targets are provided for simple initialization. Please refer to the `Makefile` for the full list of supported commands. 

Because the application listens on port 80, please ensure you run the make command as root (sudo) to get the application started. The application requires a private key (.rsa) to sign the tokens with. Note that every initialization generates a new private and public key using openssl.

```shell
# To run the application
sudo make
```

## API

### GET /

Returns you basic information about the system.

```
Authentication: None
Status Codes: 200
Response: application/json
{
    "version"     string
    "description" string
}
```

### GET /health

Health endpoint to verify up status. Returns 200 if up.

```
Authentication: None
Status Codes: 200
Repsonse: None
```
### POST /oauth/token

Endpoint to retrieve a token. **Note**: Only password grants are supported at this time.

```
Authentication: Basic, or None
Payload: {
    "client_id"     string (only required if not using basic)
    "client_secret" string (only required if not using basic)
    "grant_type"    "password"
    "username"      string (email of the user)
}
Status Codes: 200, 400, 401, 500, 501
Repsonse:
{
    "access_token"  string (jwt token)
    "token_type"    string (grant type requested)
    "refresh_token" string (uuid used to refresh token)
    "expires_in"    string (jwt token ttl, seconds)
    "scope"         []string
    "jti"           string (token identifier)
}
```

## Developers

### Tools required:
1. [GO](https://golang.org/doc/install)
1. [git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
1. [openssl](https://www.openssl.org/source/)
1. [Running make on Windows](http://gnuwin32.sourceforge.net/packages/make.htm)

### 3rd Party Libraries:
1. [mux](https://github.com/gorilla/mux): URL Router
1. [negroni](https://github.com/urfave/negroni): HTTP middleware
1. [jwt-go](https://github.com/dgrijalva/jwt-go): JWT library for GO
1. [uuid](https://github.com/google/uuid): UUID library

## Testing

### Unit

The unit tests can be executed with the following make target. You can also directly run the go command as required. To find more information about GoLang's testing lib: [link](https://golang.org/pkg/testing/)

```shell
make test
```

### Acceptance

A rudimentary acceptance test is available as a Python3 executable. 

#### Tools

** NOTE**

Virtualenv is not being used as part of the requirements. By following the section below, you will be installing the required tooling to your system.

1. [python3](https://www.python.org/downloads/)

#### Summary

The acceptance test runs under the premise that the default credentials for the `client_id` and `client_secret` are used. The acceptance currently tests the supported API endpoints, and ensures that expected HTTP/1.1 status codes and responses are returned. To execute the acceptance test:

```shell
# Ensure you have the required python packages installed
pip install -r requirments.txt

# If not running, run the application
make

# Run Makefile target
make acceptance
```

## Future

The following outlines functionalities that are not yet supported:
1. Token revoke
1. Token renewal
1. Auto token expiry

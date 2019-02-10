# Fin 
Fin is an open source backend API to track personal finances, made with Go. _Very alpha_.

[![Documentation at Postman](https://img.shields.io/badge/Documentation-Postman-orange.svg)](https://documenter.getpostman.com/view/423288/RztoLTaX)
[![Build Status](https://travis-ci.org/jonatasbaldin/fin-backend.svg?branch=master)](https://travis-ci.org/jonatasbaldin/fin-backend)
[![Coverage Status](https://coveralls.io/repos/github/jonatasbaldin/fin-backend/badge.svg?branch=master)](https://coveralls.io/github/jonatasbaldin/fin-backend?branch=master)

## Considerations
I used this project to learn Go. Expect some messy code. Maybe some bugs. Definitely bugs.

So far this project has no "hosted version", so you need to deploy by yourself.

## Using it
Set the environment variables:
```
$ export DB=postgres://user:pass@host:port/dbame
$ export DB_TEST=postgres://user:pass@host:port/dbame
$ export PORT=5000
```

Getting the binary:
```
$ go get github.com/jonatasbaldin/fin
$ fin -migrate && fin -serve
```

With Docker:    
```
$ docker pull jonatsabaldin/fin
$ docker run -e DB=postgres://user:pass@host:port/dbame -p 5000:5000 fin
```

## Contributing
Building the project and running tests. Requires Go v1.11 or later.
```
$ git clone git@github.com:jonatasbaldin/fin
$ make test
$ make build
```

You may want to tackle some [issues](https://github.com/jonatasbaldin/fin/issues).

## Roadmap
- Add User profiles
- Add support for crypto currencies Rates
- Add logs/telemetry
- Add `fin` service to `docker-compose.yml`

## License
[MIT](https://github.com/jonatasbaldin/finblob/master/LICENSE).
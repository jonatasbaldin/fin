# Fin 
Fin is an open source backend API to track personal finances, made with Go. _Very alpha_.

[![Documentation at Postman](https://img.shields.io/badge/Documentation-Postman-orange.svg)](https://documenter.getpostman.com/view/423288/RztoLTaX)

## Considerations
I used this project to learn Go. Expect some messy code. Maybe some bugs. Definitely bugs.

So far this project has no "hosted version", so you need to deploy by yourself.

## Using it
Getting th binary:
```
$ go get github.com/jonatasbaldin/fin
$ export DB=postgres://user:pass@host:port/dbame
$ fin -migrate
$ fin -serve
```

With Docker:    
```
$ docker pull jonatsabaldin/fin
$ docker run -e DB=postgres://user:pass@host:port/dbame fin
```

## Contributing
Building the project and running tests. Requires Go v1.11 or later.
```
$ git clone git@github.com:jonatasbaldin/fin
$ go build
$ export DB_TEST=postgres://user:pass@host:port/dbame
$ go test
```

You may want to tackle some [issues](https://github.com/jonatasbaldin/fin/issues).

## Roadmap
- Add User profiles
- Add support for crypto currencies Rates
- Add logs/telemetry
- Add `fin` service to `docker-compose.yml`

## License
[MIT](https://github.com/jonatasbaldin/finblob/master/LICENSE).
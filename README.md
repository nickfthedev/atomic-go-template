# Project my-go-template

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## Shopify App?

Don't just run `make watch` to develop your app. Run `npm run dev` to develop your app. This will update the config, build the app, and run the app with live reload. Make sure to run `npm install` to install the dependencies and `npm run g:install` to install the global dependencies.

## MakeFile

run all make commands with clean tests

```bash
make all build
```

build the application

```bash
make build
```

run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB container

```bash
make docker-down
```

live reload the application

```bash
make watch
```

run the test suite

```bash
make test
```

clean up binary from the last build

```bash
make clean
```

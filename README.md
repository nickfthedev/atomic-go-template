# atomic-go-template

###### What is this? It is a template for a go project.

It comes with templ, htmx, auth, mail and jsx support preconfigured.

###### What is special about it?

This template is using a different approach to the usual go template.
Code for routes is organized in templ files. every templ route or component got its own GET, POST, ... functions inside. This way we got all the code for a route or component in one place.

This differs from the usual go template, where we have a MVC approach. In my opinion a MVC approach is not the best way to organize code in modern web applications. It was the best way back in the days of php, but i think this is outdated.

Heavily inspired by Remix.JS
Started with https://github.com/Melkeydev/go-blueprint but heavily modified it to fit my idea.

###### About JSX TSX support

Please see the example route for this. I followed the approach of having dynamic react islands.
Every react.ts or \*.react.ts file is compiled separately and embedded into the binary.
You find the compiled files in web/embed/assets/react/

###### Tech stack:

- Go (https://go.dev/)
- Templ (https://templ.guide/)
- HTMX (https://htmx.org/)
- Chi (https://github.com/go-chi/chi)
- Gorm (https://gorm.io/)
- Tailwind (https://tailwindcss.com/)
- DaisyUI (https://daisyui.com/)
- ESBuild (https://esbuild.github.io/)
- Resend (https://resend.com/)

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

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

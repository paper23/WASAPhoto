# WASAPhoto

**A full stack project in Go and Vue.js**

WASAPhoto is a social media platform based on **Instagram** and it is a **web app** created for the [Web And Software Architecture][WASA] course.

## Project structure

* `cmd/` contains all executables; Go programs here should only do "executable-stuff", like reading options from the CLI/env, etc.
	* `cmd/healthcheck` is an example of a daemon for checking the health of servers daemons; useful when the hypervisor is not providing HTTP readiness/liveness probes (e.g., Docker engine)
	* `cmd/webapi` contains an example of a web API server daemon
* `demo/` contains a demo config file
* `doc/` contains the documentation (usually, for APIs, this means an OpenAPI file)
* `service/` has all packages for implementing project-specific functionalities
	* `service/api` is the package with the common functionalities and types necessary to every other real controller or REST API endpoint
    * `service/database` contains all the database interactions written in Go
	* `service/globaltime` contains a wrapper package for `time.Time` (useful in unit testing)
* `vendor/` is managed by Go, and contains a copy of all dependencies
* `webui/` is the frontend code in Vue.js; it includes:
	* Bootstrap JavaScript framework
	* a customized version of "Bootstrap dashboard" template
	* feather icons as SVG
	* Go code for release embedding

Other project files include:
* `open-npm.sh` starts a new (temporary) container using `node:lts` image for safe web frontend development (you don't want to use `npm` in your system, do you?)


# How to build container images
## Backend

If you're not using the WebUI, or if you don't want to embed the WebUI into the final executable, then:
```shell
go build ./cmd/webapi
```

If you're using the WebUI and you want to embed it into the final executable:
```shell
./open-npm.sh
# (here you're inside the NPM container)
npm run build-embed
exit
# (outside the NPM container)
go build -tags webui ./cmd/webapi/
```

# How to run (in developement mode)
You can launch the backend only using:
```shell
go run ./cmd/webapi/
```
If you want to launch the WebUI, open a new tab and launch:
```shell
./open-npm.sh
# (here you're inside the NPM container)
npm run dev
```

# License
Based on [Enrico Bassetti][EBassetti]'s project [Fantastic Cofee Decaffeinated][FantasticCoffee]

See [LICENSE][License]

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2CA5E0?style=for-the-badge&logo=docker&logoColor=white)
![Swagger](https://img.shields.io/badge/Swagger-85EA2D?style=for-the-badge&logo=Swagger&logoColor=white)
![Vue.js](https://img.shields.io/badge/Vue.js-35495E?style=for-the-badge&logo=vuedotjs&logoColor=4FC08D)
![Node.js](https://img.shields.io/badge/Node.js-339933?style=for-the-badge&logo=nodedotjs&logoColor=white)
![SQLite](https://img.shields.io/badge/SQLite-07405E?style=for-the-badge&logo=sqlite&logoColor=white)


[WASA]: http://gamificationlab.uniroma1.it/en/wasa/
[FantasticCoffee]: https://github.com/sapienzaapps/fantastic-coffee-decaffeinated/
[EBassetti]: https://github.com/Enrico204
[License]: https://github.com/paper23/WASAPhoto/blob/main/LICENSE
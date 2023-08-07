# gomodtrace

Utility is intended to trace some dependency sub-graphs for Go projects.

# Installation

Run
`go install github.com/godepsresolve/gomodtrace/cmd/gomodtrace@latest`

# Usage

```
go mod graph | gomodtrace [OPTION]... PARENT_PACKAGE DEPENDENT_PACKAGE
-v	use verbose mode
```

Let's try to trace usages of vulnerable [gogo/protobuf](https://github.com/gogo/protobuf/)
library ([cve-2021-3121](https://nvd.nist.gov/vuln/detail/cve-2021-3121)) in [ory/hydra](https://github.com/ory/hydra)
well-known project.

## Usage of gomodtrace without modgraphviz

```shell
$ cd github.com/ory/hydra
$ go mod graph | gomodtrace github.com/ory/hydra/v2 github.com/gogo/protobuf@v1.1.1
github.com/ory/hydra/v2 github.com/ory/x@v0.0.574
github.com/ory/hydra/v2 github.com/prometheus/client_golang@v1.13.0
github.com/ory/hydra/v2 github.com/prometheus/common@v0.37.0
github.com/ory/x@v0.0.574 github.com/prometheus/client_golang@v1.13.0
github.com/ory/x@v0.0.574 github.com/prometheus/common@v0.37.0
github.com/prometheus/client_golang@v1.13.0 github.com/prometheus/common@v0.37.0
github.com/prometheus/common@v0.37.0 github.com/prometheus/client_golang@v1.12.1
github.com/prometheus/client_golang@v1.12.1 github.com/prometheus/common@v0.32.1
github.com/prometheus/common@v0.32.1 github.com/prometheus/client_golang@v1.11.0
github.com/prometheus/client_golang@v1.11.0 github.com/prometheus/common@v0.26.0
github.com/prometheus/common@v0.26.0 github.com/prometheus/client_golang@v1.7.1
github.com/prometheus/client_golang@v1.7.1 github.com/prometheus/common@v0.10.0
github.com/prometheus/common@v0.10.0 github.com/prometheus/client_golang@v1.0.0
github.com/prometheus/client_golang@v1.0.0 github.com/prometheus/common@v0.4.1
github.com/prometheus/common@v0.4.1 github.com/gogo/protobuf@v1.1.1

```

## Usage of gomodtrace with modgraphviz:

```shell
$ cd github.com/ory/hydra
$ go mod graph | gomodtrace github.com/ory/hydra/v2 github.com/gogo/protobuf@v1.1.1 | modgraphviz | dot -Tsvg -o graph_modg.svg
```

Then open graph_modg.svg with your favorite image viewer.
![modgraphviz graph image](/assets/images/graph_modg.svg)

So, it's really cool to see such compact output/image for such a big input graph:

```
$ cd github.com/ory/hydra
$ go mod graph | wc -l
4712
```

# Problem context

"Big" projects in Go could have a lot of dependencies.
Some of these dependencies are vulnerable, for example your security linter alerts you
(and it could be also false positive).
But it could be hard to determine how these dependencies have had come.
`go mod graph` command could return a lot of lines
(e.g. for investigation in one of the projects it returns more than 4000 lines).
So step-by-step "grep-ing" could become a hard deal.

There is a lot of visualization instruments like gomod and modgraphviz, you can read
about its usage in my another [repo](https://github.com/godepsresolve/dep_graph).
But rendered graph will take a lot of space and will be really overloaded and non human-readable.
So no help will come from there.

It also worth understanding that `go mod why` could return something like:

```
$ go mod why github.com/gogo/protobuf@v1.1.1
# github.com/gogo/protobuf@v1.1.1
(main module does not need package github.com/gogo/protobuf@v1.1.1)
```

because it not a direct dependency or MVS (minimal version selection) choose another version
of vulnerable lib. You can check this in widely known and widely
used [hydra project](https://github.com/ory/hydra/tree/339bf40e189e5285f7a8b9c7daa184ac00d0110f).

# Solution

Find path from your project to vulnerable dependency and filter out other non-relevant libraries/modules.
This utility is intended to do this.

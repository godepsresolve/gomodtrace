# gomodtrace

Utility is intended to trace some dependency sub-graphs for Go projects.

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
$ go mod why github.com/gogo/protobuf@v1.2.1
# github.com/gogo/protobuf@v1.2.1
(main module does not need package github.com/gogo/protobuf@v1.2.1)
```

because it not a direct dependency or MVS (minimal version selection) choose another version
of vulnerable lib. You can check this in widely known and widely
used [sql-migrate project](https://github.com/rubenv/sql-migrate/tree/9e20e0b824edc2b83a3ecc1c3219421f8b23516b).

# Solution

Find path from your project to vulnerable dependency and filter out other non-relevant libraries/modules.
This utility is intended to do this.

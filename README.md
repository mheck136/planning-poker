# Planning Poker

This is an experimental repository with the goal to create the prototype of a CQRS & EventSourcing toolkit.

The toolkit will provide location-transparent access to aggregate-roots with an api that allows sending Commands and
subscribing to Events.

The approach is influenced by Akka Persistence.

A registry is used to retrieve aggregate roots. The aggregate roots accept commands and guarantee conflict-free
execution of the commands. Every change to the aggregate leads to one or more new events that will be stored atomically
before they are applied to the aggregate. Only one instance of an aggregate (ID) is present in the application at avery
moment in time and all commands for that specific aggregate are processed by one goroutine.

## Case Study

As a first step I built an api for a scrum planning poker app that implements the described approach. The generic
functionality will be isolated into its own package and then further refined and refactored until it seems usable and
generic enough.

The planning poker app was used for two reasons:

- The lack of a user-friendly free application on the web
- The domain has enough substance to be used as a case study for the toolkit while it is simple enough to be implemented
  in short time

## Future

When the shape of the API is stable(ish), it will be translated into a generic api (either code generation or type
parameters aka. Generics which are due in Go 1.18). The toolkit will then be extended with a clustering mechanism
including _membership_, _consistent hashing_, _command routing_, and _fault tolerance_.

## Get Started

When you just want to use the api, just build it with `go build ./cmd/...` and then run the binary, e.g.
`./planning-poker-api`.



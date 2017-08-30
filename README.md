# errstack
[![GoDoc](https://godoc.org/github.com/robert-zaremba/errstack?status.png)](https://godoc.org/github.com/robert-zaremba/errstack)

stacked error wrappers

Package `errstack` provides a way to augment errors with a functional type and iteratively build new request error.

# Motivation

In software we have 3 types of errors:

* Infrastructure - this errors come from the requests to other services.
* Domain - this is the type of errors when we detect that our internal application state is wrong.
* Request - a user request error.

Request errors maps parameters (application inputs) to the cause description (why the value of wrong).

Domain errors usually arise form some model bug or unhandled condition.

Infrastructure errors arise from the system or connection error.

It is very handy to know the place (stack) where the error Infrastructure or Domain error arose - that's why they are threatened especially and are augmented with the stack trace.


# Credits

The initial work was done for the [AgFlow](http://www.agflow.com) project.

* [Krzysztof Dryś](https://github.com/krzysztofdrys)
* [Robert Zaremba](http://scale-it.pl)

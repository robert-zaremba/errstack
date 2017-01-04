/*
Package `errstack` provides a way to augment errors with a functional type and iteratively build new request error.

In software we have 3 types of errors:

* Infrastructure - this errors come from the requests to other services.

* Domain - this is the type of errors when we detect that our internal application state is wrong.

* Request - a user request error.

Request errors maps parameters (or input data fields) to the explanation (why the value of wrong).

Domain errors usually arise form some model bug or unhandled condition.

Infrastructure errors arise from the system or connection error.

It is very handy to know the place (stack) where the error Infrastructure or Domain error arose - that's why they are threatened especially and are augmented with the stack trace.
*/
package errstack

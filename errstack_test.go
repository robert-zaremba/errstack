package errstack

import (
	"errors"

	. "gopkg.in/check.v1"
)

type ESuite struct{}

func assertMarshal(err E, expected string, kind Kind, c *C) {
	b, errm := err.MarshalJSON()
	c.Assert(errm, IsNil)
	c.Assert(string(b), Equals, expected)
	c.Assert(err.Kind(), Equals, kind)
	c.Assert(err.IsReq(), Equals, kind == Request)
}

func (s *ESuite) TestNewReq(c *C) {
	// New request with no wrapping
	err := NewReq("error-text")
	c.Assert(err.Error(), Equals, "error-text")
	assertMarshal(err, `{"msg":"error-text"}`, Request, c)

	err2 := err.WithMsg("more_details")
	c.Assert(err2.Error(), Equals, "more_details [error-text]")
	assertMarshal(err2, `{"err":{"msg":"error-text"},"msg":"more_details"}`, Request, c)
}

func (s *ESuite) TestWrapAsReq(c *C) {
	err := errors.New("new error")

	werr := WrapAsReq(err, "one")
	assertMarshal(werr, `{"err":"new error","msg":"one"}`, Request, c)

	// Wrap wrapped error
	werr = WrapAsReq(werr, "two")
	c.Assert(werr.Error(), Equals, "two [one [new error]]")
	assertMarshal(werr, `{"err":{"err":"new error","msg":"one"},"msg":"two"}`, Request, c)

	// Wrap request
	err = NewReqDetails("key", "details", "message")
	werr = WrapAsReq(err, "two")
	assertMarshal(werr, `{"err":{"key":"details"},"msg":"two [message]"}`, Request, c)
}

func (s *ESuite) TestWrapAsReqF(c *C) {
	err := errors.New("new error")

	werr := WrapAsReqF(err, "%d", 1)
	assertMarshal(werr, `{"err":"new error","msg":"1"}`, Request, c)

	// Wrap wrapped error
	werr = WrapAsReqF(werr, "%d", 2)
	c.Assert(werr.Error(), Equals, "2 [1 [new error]]")
	assertMarshal(werr, `{"err":{"err":"new error","msg":"1"},"msg":"2"}`, Request, c)

	// Wrap request
	err = NewReqDetails("key", "details", "message")
	werr = WrapAsReqF(err, "%d", 3)
	assertMarshal(werr, `{"err":{"key":"details"},"msg":"3 [message]"}`, Request, c)
}

func (s *ESuite) TestWrapAsInf(c *C) {
	err := errors.New("new error")

	werr := WrapAsIOf(err, "one")
	assertMarshal(werr, `"Internal server error: one"`, IO, c)
}

func (s *ESuite) TestWrappingNil(c *C) {
	message := "message"
	c.Assert(WrapAsDomain(nil, message), IsNil)
	c.Assert(WrapAsDomainF(nil, message), IsNil)
	c.Assert(WrapAsIO(nil, message), IsNil)
	c.Assert(WrapAsIOf(nil, message), IsNil)
	c.Assert(WrapAsReq(nil, message), IsNil)
	c.Assert(WrapAsReqF(nil, message), IsNil)
}

func (s *ESuite) TestWrapping(c *C) {
	tests := []E{NewIO("a"), NewIOf("a"), NewDomain("a"), NewDomainF("a")}
	message := "message"

	var check = func(eOrig, eWrapped E, kind Kind) {
		c.Check(eWrapped.Error(), Equals, eOrig.WithMsg(message).Error())
		c.Check(eWrapped.IsReq(), Equals, kind == Request)
		c.Check(eWrapped.Kind(), Equals, kind)
	}

	for _, tt := range tests {
		check(tt, WrapAsDomain(tt, message), Domain)
		check(tt, WrapAsDomainF(tt, message), Domain)
		check(tt, WrapAsIO(tt, message), IO)
		check(tt, WrapAsIOf(tt, message), IO)
		check(tt, WrapAsReq(tt, message), Request)
		check(tt, WrapAsReqF(tt, message), Request)
	}
}

func (s *ESuite) TestRootErrAndCause(c *C) {
	err := errors.New("new Error")
	e2 := WrapAsDomain(err, "something more")
	e3 := WrapAsIO(e2, "something more 2")
	c.Assert(RootErr(e3), Equals, err)

	e2C := e2.(HasUnderlying)
	e3C := e3.(HasUnderlying)
	c.Assert(e3C.Cause(), Equals, e2)
	c.Assert(e2C.Cause(), Equals, err)

	// RootErr should return nil for errors build on top of nil
	e2 = WrapAsDomain(nil, "something more")
	e3 = WrapAsIO(e2, "something more 2")
	c.Assert(RootErr(e3), IsNil)

	// this will return nil, and interface E doesn't implement HasUnderlying
	_, ok := e2.(HasUnderlying)
	c.Assert(ok, Equals, false)
	_, ok = e3.(HasUnderlying)
	c.Assert(ok, Equals, false)
}

func checkKindFalse(c *C, err error) {
	c.Assert(IsKind(Other, err), Equals, false)
	c.Assert(IsKind(Request, err), Equals, false)
	c.Assert(IsKind(IO, err), Equals, false)
	c.Assert(IsKind(Domain, err), Equals, false)
	c.Assert(IsKind(Invalid, err), Equals, false)
}

func (s *ESuite) TestIsKind(c *C) {
	var err error = nil
	checkKindFalse(c, err)

	err = errors.New("new Error")
	checkKindFalse(c, err)

	err = NewReq("error")
	c.Assert(IsKind(Request, err), Equals, true)

	// err = NewIO("error")
	// c.Assert(IsKind(IO, err), Equals, true)

	err = NewDomain("error")
	c.Assert(IsKind(Domain, err), Equals, true)

	err = WrapAsReq(err, "hi")
	c.Assert(IsKind(Request, err), Equals, true)
}

func (s *ESuite) TestAdd(c *C) {
	var err = New(Request, "error1")
	err.Add("key1", "details1")
	err.Add("key1", "details2")

	c.Assert(err.Kind(), Equals, Request)
	c.Check(err.Error(), Equals, "error1")
	c.Check(err.Details(), DeepEquals, map[string]interface{}{"key1": "details2"})
}

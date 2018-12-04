package errstack

import (
	"errors"

	. "gopkg.in/check.v1"
)

type ESuite struct{}

func (s *ESuite) TestNewReq(c *C) {
	// New request with no wrapping
	err := NewReq("error-text")
	c.Assert(err.Error(), Equals, "error-text")

	// New request
	nerr := NewReq("one")
	b, errm := nerr.MarshalJSON()
	c.Assert(errm, IsNil)
	c.Assert(string(b), Equals, "{\"msg\":\"one\"}")

	nerr = NewReq("two")
	b, errm = nerr.MarshalJSON()
	c.Assert(errm, IsNil)
	c.Assert(string(b), Equals, "{\"msg\":\"two\"}")
}

func (s *ESuite) TestWrapAsReq(c *C) {
	err := errors.New("new error")

	werr := WrapAsReq(err, "one")
	b, errm := werr.MarshalJSON()
	c.Assert(errm, IsNil)
	c.Assert(string(b), Equals, `{"msg":"one [new error]"}`)

	// Wrap wrapped error
	werr = WrapAsReq(werr, "two")
	b, errm = werr.MarshalJSON()
	c.Assert(errm, IsNil)
	c.Assert(string(b), Equals, `{"msg":"two [one [new error]]"}`)

	// Wrap request
	err = NewReqDetails("key", "details")
	werr = WrapAsReq(err, "two")
	b, errm = werr.MarshalJSON()
	c.Assert(errm, IsNil)
	c.Assert(string(b), Equals, "{\"msg\":\"two [key: details\\n]\"}")
}

func (s *ESuite) TestWrapAsReqF(c *C) {
	err := errors.New("new error")

	werr := WrapAsReqF(err, "%d", 1)
	b, errm := werr.MarshalJSON()
	c.Assert(errm, IsNil)
	c.Assert(string(b), Equals, `{"msg":"1 [new error]"}`)

	// Wrap wrapped error
	werr = WrapAsReqF(werr, "%d", 2)
	b, errm = werr.MarshalJSON()
	c.Assert(errm, IsNil)
	c.Assert(string(b), Equals, `{"msg":"2 [1 [new error]]"}`)

	// Wrap request
	err = NewReqDetails("key", "details")
	werr = WrapAsReqF(err, "%d", 3)
	b, errm = werr.MarshalJSON()
	c.Assert(errm, IsNil)
	c.Assert(string(b), Equals, "{\"msg\":\"3 [key: details\\n]\"}")
}

func (s *ESuite) TestWrapAsInf(c *C) {
	err := errors.New("new error")

	werr := WrapAsInfF(err, "one")
	b, errm := werr.MarshalJSON()
	c.Assert(errm, IsNil)
	c.Assert(string(b), Equals, `"Internal server error"`)
}

func (s *ESuite) TestWrappingNil(c *C) {
	message := "message"
	c.Assert(WrapAsDomain(nil, message), IsNil)
	c.Assert(WrapAsDomainF(nil, message), IsNil)
	c.Assert(WrapAsInf(nil, message), IsNil)
	c.Assert(WrapAsInfF(nil, message), IsNil)
	c.Assert(WrapAsReq(nil, message), IsNil)
	c.Assert(WrapAsReqF(nil, message), IsNil)
}

func (s *ESuite) TestWrapping(c *C) {
	tests := []E{NewInf("a"), NewInfF("a"), NewDomain("a"), NewDomainF("a")}
	message := "message"

	for _, tt := range tests {
		c.Assert(WrapAsDomain(tt, message), DeepEquals, tt)
		c.Assert(WrapAsDomainF(tt, message), DeepEquals, tt)
		c.Assert(WrapAsInf(tt, message), DeepEquals, tt)
		c.Assert(WrapAsInfF(tt, message), DeepEquals, tt)
		c.Assert(WrapAsReq(tt, message), DeepEquals, tt)
		c.Assert(WrapAsReqF(tt, message), DeepEquals, tt)
	}
}

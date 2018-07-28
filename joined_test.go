package errstack

import (
	"errors"

	. "gopkg.in/check.v1"
)

type JoinSuite struct{}

func (s *BuilderSuite) TestJoinErrorsNil(c *C) {
	var err error

	c.Assert(Join(), IsNil)
	c.Assert(Join(nil), IsNil)
	c.Assert(Join(err), IsNil)
	c.Assert(Join(err, err), IsNil)
	c.Assert(Join(err, nil, err), IsNil)
}

func (s *BuilderSuite) TestJoinErrorsNotNil(c *C) {
	var err = errors.New("abc")
	var errNil error

	c.Assert(Join(err), Not(IsNil))
	c.Assert(Join(err, err), Not(IsNil))
	c.Assert(Join(errNil, err), Not(IsNil))
	c.Assert(Join(nil, err), Not(IsNil))
	c.Assert(Join(nil, err, nil), Not(IsNil))
}

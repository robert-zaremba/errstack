package errstack

import (
	. "github.com/scale-it/checkers"
	. "gopkg.in/check.v1"
)

type BuilderSuite struct{}

func (s *BuilderSuite) TestPut(c *C) {
	b := NewBuilder()
	c.Assert(b.NotNil(), IsFalse)
	c.Assert(b.ToReqErr(), IsNil)

	b.Put("k1", 1)
	c.Assert(b.NotNil(), IsTrue)
	errR := b.ToReqErr().(*request)
	c.Assert(errR, NotNil)
	c.Assert(errR.details, DeepEquals, errmap{"k1": 1})

	b.Put("k1", 3)
	b.Put("newkey", 4)
	errR = b.ToReqErr().(*request)
	c.Assert(errR, NotNil)
	c.Assert(errR.details, DeepEquals, errmap{"k1": chain{1, 3}, "newkey": 4})
}

func (s *BuilderSuite) TestFork(c *C) {
	b1 := NewBuilder()
	b1.Put("k", 1)
	c.Check(b1.Get("k"), Equals, 1)

	b2 := b1.Fork("b")
	b2.Put("k", 2)

	b3 := b1.Fork("c")
	b3.Put("k", 3)

	c.Check(b2.Get("k"), Equals, 2)
	c.Check(b3.Get("k"), Equals, 3)

	c.Assert(b1.NotNil(), IsTrue)
	b := b1.(builder)
	c.Assert(b.m, DeepEquals, errmap{"k": 1, "b|k": 2, "c|k": 3})
}

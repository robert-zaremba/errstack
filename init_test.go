package errstack

import (
	"testing"

	. "gopkg.in/check.v1"
)

func init() {
	//	logger = log15.New()
	Suite(&BuilderSuite{})
	Suite(&ESuite{})
	Suite(&JoinSuite{})
}

func Test(t *testing.T) { TestingT(t) }

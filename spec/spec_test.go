package spec

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"testing"

	. "github.com/essentialkaos/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

type SpecSuite struct{}

// ////////////////////////////////////////////////////////////////////////////////// //

var _ = Suite(&SpecSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *SpecSuite) TestGetDeps(c *C) {
	deps, err := GetDeps("../testdata/test.spec", []string{"_pg:99"})
	c.Assert(err, IsNil)
	c.Assert(len(deps), Equals, 10)

	c.Assert(deps.Names(false), DeepEquals, []string{
		"bash", "gcc", "make", "perl", "perl(ExtUtils::Embed)",
		"postgresql10-devel", "python", "python3dist(setuptools)",
		"readline-devel", "zlib-devel",
	})

	c.Assert(deps.Names(true), DeepEquals, []string{
		"bash = 4", "gcc", "make", "perl > 5", "perl(ExtUtils::Embed)",
		"postgresql10-devel = 10.23", "python < 3", "python3dist(setuptools)",
		"readline-devel <= 7", "zlib-devel >= 1.2.11",
	})

	c.Assert(deps[1].String(), Equals, "{gcc}")
	c.Assert(deps[9].String(), Equals, "{zlib-devel ≥ 1.2.11}")

	c.Assert(deps[0].Cond.String(), Equals, "=")
	c.Assert(deps[0].Cond.Clause(), Equals, "=")
	c.Assert(deps[3].Cond.String(), Equals, ">")
	c.Assert(deps[3].Cond.Clause(), Equals, ">")
	c.Assert(deps[6].Cond.String(), Equals, "<")
	c.Assert(deps[6].Cond.Clause(), Equals, "<")
	c.Assert(deps[8].Cond.String(), Equals, "≤")
	c.Assert(deps[8].Cond.Clause(), Equals, "<=")
	c.Assert(deps[9].Cond.String(), Equals, "≥")
	c.Assert(deps[9].Cond.Clause(), Equals, ">=")

	c.Assert(deps[1].Cond.String(), Equals, "")
	c.Assert(deps[1].Cond.Clause(), Equals, "")
}

func (s *SpecSuite) TestErrors(c *C) {
	_, err := GetDeps("/__unknown__", nil)
	c.Assert(err, ErrorMatches, "File /__unknown__ doesn't exist or not accessible")

	_, err = GetDeps("../testdata/broken.spec", nil)
	c.Assert(err, ErrorMatches, "Spec parsing error: line 1: Unknown tag: TEST1234!")

	c.Assert(Cond(0).String(), Equals, "")
	c.Assert(Cond(0).Clause(), Equals, "")

	c.Assert(parseCond("!"), Equals, Cond(0))
}

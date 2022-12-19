package e2e_test

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/stretchr/testify/suite"
)

type YdbSQLSuite struct {
	suite.Suite
}

func (s *YdbSQLSuite) Test_YdbSQL_MigrationSteps() {
	r := s.Require()

	c, err := pop.Connect("ydb")
	r.NoError(err)
	r.NoError(retryOpen(c))

	run(&s.Suite, c, runTestData(&s.Suite, c, true))
}

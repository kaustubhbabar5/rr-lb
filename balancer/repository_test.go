package balancer

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kaustubhbabar5/rr-lb/adapters/cache"
	"github.com/kaustubhbabar5/rr-lb/pkg/constants"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	repo *repository

	testUrl string
}

func TestRepository(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupSuite() {
	cache, err := cache.New("0.0.0.0:6379", "")
	if err != nil {
		log.Fatal(err)
	}

	repo := NewRepository(cache)

	s.repo = repo

	s.testUrl = "https://www.google.com/"
}

func (s *TestSuite) BeforeTest(suiteName, testName string) {
	switch testName {
	case "TestMarkReplicaUnhealthy":
		s.repo.cache.LPush(context.Background(), constants.HEALTHY_SERVERS, s.testUrl)
		// default:

	}
}

// The TearDownTest method will be run after every test in the suite.
func (s *TestSuite) TearDownTest() {
	s.repo.cache.FlushAll(context.Background())
}

func (s *TestSuite) TestMarkReplicaUnhealthy() {
	time.Sleep(10 * time.Second)
	err := s.repo.MarkReplicaUnhealthy(s.testUrl)
	s.Nil(err)
	time.Sleep(10 * time.Second)

	res := s.repo.cache.LPos(context.Background(), constants.UNHEALTHY_SERVERS, s.testUrl, redis.LPosArgs{})
	s.Nil(res.Err())
	s.Equal(int64(0), res.Val(), "index of value different")

}

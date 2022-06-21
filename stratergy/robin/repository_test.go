package robin

import (
	"context"
	"log"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/kaustubhbabar5/rr-lb/adapters/cache"
	"github.com/kaustubhbabar5/rr-lb/pkg/constants"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	repo *repository

	testUrl1 string
	testUrl2 string
	testUrl3 string
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

	s.testUrl1 = "https://www.google.com/"
	s.testUrl2 = "https://www.facebook.com/"
	s.testUrl3 = "https://www.amazon.com/"

}

func (s *TestSuite) BeforeTest(suiteName, testName string) {
	switch testName {
	case "TestGetServer":
		s.repo.cache.RPush(context.Background(), constants.HEALTHY_SERVERS, s.testUrl1, s.testUrl2, s.testUrl3)
		// default:

	}
}

func (s *TestSuite) AfterTest(suiteName, testName string) {
	switch testName {
	case "TestGetServer":
		s.repo.cache.FlushAll(context.Background())
		// default:

	}
}

func (s *TestSuite) TestGetServer() {

	url, err := s.repo.GetServer()
	s.Nil(err)
	s.Equal(s.testUrl1, url, "got different url")

	res := s.repo.cache.LPos(context.Background(), constants.HEALTHY_SERVERS, s.testUrl1, redis.LPosArgs{})
	s.Nil(res.Err())
	s.Equal(int64(2), res.Val(), "expected key at index 2 but got at: ", res.Val())
}

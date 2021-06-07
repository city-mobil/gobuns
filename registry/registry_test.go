package registry

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/sdk/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/city-mobil/gobuns/retry"
)

type regSuite struct {
	suite.Suite

	client       Client
	consulClient *api.Client
	consulServer *testutil.TestServer
}

func TestRegistry(t *testing.T) {
	s := new(regSuite)

	server, err := testutil.NewTestServerConfigT(t, nil)
	require.NoError(t, err)

	if server.Config.Bootstrap {
		server.WaitForLeader(t)
	}

	s.consulServer = server

	opts := Config{
		Addr:         server.HTTPAddr,
		QueryTimeout: 1 * time.Second,
		RetryConfig:  retry.NewDefRetryConfig(),
	}
	s.client, err = NewClient(opts)
	if assert.NoError(t, err) {
		s.consulClient = s.client.(*client).consulClient
		suite.Run(t, s)
	} else {
		_ = server.Stop()
	}
}

func (s *regSuite) TearDownSuite() {
	_ = s.consulServer.Stop()
}

func (s *regSuite) TestGetString() {
	t := s.T()

	key := "bb"
	want := []byte("Heisenberg")
	p := &api.KVPair{Key: key, Value: want}
	_, err := s.consulClient.KV().Put(p, nil)
	require.NoError(t, err)

	got, err := s.client.GetString(context.Background(), key)
	require.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func (s *regSuite) TestGetInt() {
	t := s.T()

	key := "bb"
	want := []byte("12673")
	p := &api.KVPair{Key: key, Value: want}
	_, err := s.consulClient.KV().Put(p, nil)
	require.NoError(t, err)

	got, err := s.client.GetInt(context.Background(), key)
	require.NoError(t, err)
	assert.Equal(t, 12673, got)
}

func (s *regSuite) TestGetBool() {
	t := s.T()

	key := "bb"
	want := []byte("true")
	p := &api.KVPair{Key: key, Value: want}
	_, err := s.consulClient.KV().Put(p, nil)
	require.NoError(t, err)

	got, err := s.client.GetBool(context.Background(), key)
	require.NoError(t, err)
	assert.True(t, got)
}

func (s *regSuite) TestGetNotExistKey() {
	t := s.T()

	got, err := s.client.GetString(context.Background(), "jessy")
	assert.EqualError(t, err, ErrKeyNotExist.Error())
	assert.Empty(t, got)
}

func (s *regSuite) TestQueryTimeout() {
	t := s.T()

	client, ok := s.client.(*client)
	require.True(t, ok)

	// Set tiny timeout and expect error.
	client.queryTimeout = 1 * time.Microsecond

	_, err := client.GetString(context.Background(), "bb")
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "context deadline exceeded")
	}
}

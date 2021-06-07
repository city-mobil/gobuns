package zwatch

import (
	"bytes"
	"testing"
	"time"

	"github.com/city-mobil/gobuns/registry"
	"github.com/city-mobil/gobuns/zlog"
	"github.com/city-mobil/gobuns/zlog/glog"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/sdk/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	tKey = "watch/logger/global"
)

type watchSuite struct {
	suite.Suite

	consulClient *api.Client
	consulServer *testutil.TestServer

	wc registry.WatchConfig
}

func (s *watchSuite) updateLevel(level zlog.Level) {
	p := &api.KVPair{Key: tKey, Value: []byte(level.String())}
	_, err := s.consulClient.KV().Put(p, nil)
	require.NoError(s.T(), err)
}

func TestWatcher(t *testing.T) {
	suite.Run(t, new(watchSuite))
}

func (s *watchSuite) SetupSuite() {
	t := s.T()

	server, err := testutil.NewTestServerConfigT(t, nil)
	require.NoError(t, err)

	if server.Config.Bootstrap {
		server.WaitForLeader(t)
	}

	server.WaitForSerfCheck(t)

	conf := api.DefaultConfig()
	conf.Address = server.HTTPAddr

	client, err := api.NewClient(conf)
	if err != nil {
		_ = server.Stop()
		require.NoError(t, err)
	}

	s.consulClient = client
	s.consulServer = server

	s.wc = registry.WatchConfig{
		Addr: s.consulServer.HTTPAddr,
		OnErr: func(err error) {
			assert.NoError(t, err)
		},
	}

	zlog.SetTimestampFunc(func() time.Time {
		return time.Date(2001, time.February, 3, 4, 5, 6, 7, time.UTC)
	})
}

func (s *watchSuite) AfterTest(_, _ string) {
	_, err := s.consulClient.KV().Delete(tKey, nil)
	require.NoError(s.T(), err)
}

func (s *watchSuite) TearDownSuite() {
	if s.consulServer != nil {
		_ = s.consulServer.Stop()
	}

	zlog.SetTimestampFunc(time.Now)
}

func (s *watchSuite) TestGlobalWatcher() {
	t := s.T()

	cancel, err := GlobalLevel(tKey, s.wc)
	require.NoError(t, err)
	defer cancel()

	out := &bytes.Buffer{}
	glog.Logger = zlog.New(out)

	zlog.SetGlobalLevel(zlog.FatalLevel)

	glog.Info().Msg("test")
	assert.Empty(t, out)

	s.updateLevel(zlog.DebugLevel)

	assert.Eventually(t, func() bool {
		return zlog.GlobalLevel() == zlog.DebugLevel
	}, 100*time.Millisecond, 5*time.Millisecond)

	glog.Info().Msg("test")
	assert.Equal(t, `{"level":"info","time":"2001-02-03T04:05:06Z","message":"test"}`+"\n", out.String())
}

func (s *watchSuite) TestLoggerWatcher() {
	t := s.T()
	out := &bytes.Buffer{}
	log := zlog.Raw(out).Level(zlog.FatalLevel)

	log.Info().Msg("test")
	assert.Empty(t, out)

	cancel, err := LoggerLevel(log, tKey, s.wc)
	require.NoError(t, err)
	defer cancel()

	s.updateLevel(zlog.DebugLevel)

	assert.Eventually(t, func() bool {
		return log.GetLevel() == zlog.DebugLevel
	}, 100*time.Millisecond, 5*time.Millisecond)

	log.Info().Msg("test")
	assert.Equal(t, `{"level":"info","message":"test"}`+"\n", out.String())
}

package rabbit

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/city-mobil/gobuns/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	cfgPath, err := filepath.Abs("testdata/config.yml")
	require.NoError(t, err)
	os.Args = append(os.Args, "--config="+cfgPath)

	cfgConnFn := NewConnectorConfig("rabbit.connector")
	cfgConsumerFn := NewConsumerConfig("rabbit.consumer")

	err = config.InitOnce()
	require.NoError(t, err)

	connCfg := cfgConnFn()
	assert.Equal(t, "t_addr", connCfg.Addr)
	assert.Equal(t, "t_login", connCfg.Login)
	assert.Equal(t, "t_password", connCfg.Password)
	assert.Equal(t, 140*time.Millisecond, connCfg.ReconnectDelay)
	assert.Equal(t, 17*time.Second, connCfg.ConnectTimeout)
	assert.Equal(t, 5*time.Second, connCfg.Heartbeat)

	consumerCfg := cfgConsumerFn()
	assert.Equal(t, "test_queue", consumerCfg.Queue)
	assert.Equal(t, "t_consumer", consumerCfg.Consumer)
	assert.Equal(t, int64(5), consumerCfg.ConsumersCount)
	assert.Equal(t, 14*time.Millisecond, consumerCfg.Heartbeat)
	assert.True(t, consumerCfg.AutoAck)
	assert.True(t, consumerCfg.NoWait)
	assert.True(t, consumerCfg.Exclusive)
}

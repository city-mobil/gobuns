package registry

import (
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/sdk/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWatch_EmptyKey_ShouldReturnErr(t *testing.T) {
	hd := WatchHandleFunc(func(data *string) {})
	wc := WatchConfig{
		Addr: "localhost:8500",
		OnErr: func(err error) {
			assert.NoError(t, err)
		},
	}

	cancel, err := Watch("", hd, wc)
	assert.EqualError(t, err, ErrNoKey.Error())
	assert.Nil(t, cancel)
}

func TestWatch_EmptyOnData_ShouldReturnErr(t *testing.T) {
	hd := WatchHandler(nil)
	wc := WatchConfig{
		Addr: "localhost:8500",
		OnErr: func(err error) {
			assert.NoError(t, err)
		},
	}

	cancel, err := Watch("key", hd, wc)
	assert.EqualError(t, err, ErrNoHandler.Error())
	assert.Nil(t, cancel)
}

func TestWatch(t *testing.T) {
	c, s := tMakeClient(t)
	defer func() {
		_ = s.Stop()
	}()

	s.WaitForSerfCheck(t)
	kv := c.KV()

	val := new(string)
	key := "watch/my/key"
	hd := WatchHandleFunc(func(data *string) {
		val = data
	})
	wc := WatchConfig{
		Addr: s.HTTPAddr,
		OnErr: func(err error) {
			assert.NoError(t, err)
		},
	}

	cancel, err := Watch(key, hd, wc)
	require.NoError(t, err)
	defer cancel()

	// Put the key.
	want := []byte("Heisenberg")
	p := &api.KVPair{Key: key, Value: want}
	_, err = kv.Put(p, nil)
	require.NoError(t, err)

	assert.Eventually(t, func() bool {
		return val != nil && *val == string(want)
	}, 100*time.Millisecond, 5*time.Millisecond)

	// Update the key.
	want = []byte("White")
	p = &api.KVPair{Key: key, Value: want}
	_, err = kv.Put(p, nil)
	require.NoError(t, err)

	assert.Eventually(t, func() bool {
		return val != nil && *val == string(want)
	}, 100*time.Millisecond, 5*time.Millisecond)

	// Delete the key.
	_, err = kv.Delete(key, nil)
	require.NoError(t, err)

	assert.Eventually(t, func() bool {
		return val == nil
	}, 100*time.Millisecond, 5*time.Millisecond)
}

func tMakeClient(t *testing.T) (*api.Client, *testutil.TestServer) {
	server, err := testutil.NewTestServerConfigT(t, nil)
	require.NoError(t, err)

	if server.Config.Bootstrap {
		server.WaitForLeader(t)
	}

	conf := api.DefaultConfig()
	conf.Address = server.HTTPAddr

	client, err := api.NewClient(conf)
	if err != nil {
		_ = server.Stop()
		require.NoError(t, err)
	}

	return client, server
}

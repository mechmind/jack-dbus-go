package jackdbus

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventRouterDelivery(t *testing.T) {
	client, err := New()
	require.NoError(t, err)

	var counter uint32
	handler := func(beat uint32) {
		assert.Equal(t, counter, beat)
		counter++
	}

	passiveHandler := func(beat uint32) {
		assert.Equal(t, counter, beat)
	}

	obj := client.conn.Object(client.conn.Names()[0], "/org/jackdbus/test")
	detach, err := client.eventRouter.addHandler(obj, "org.jackdbus.Emitter", "Heartbeat", handler)
	require.NoError(t, err)
	detachPassive, err := client.eventRouter.addHandler(obj, "org.jackdbus.Emitter", "Heartbeat", passiveHandler)
	require.NoError(t, err)

	for i := uint32(0); i < 3; i++ {
		client.conn.Emit("/org/jackdbus/test", "org.jackdbus.Emitter.Heartbeat", i)
		client.conn.Emit("/org/jackdbus/another-path", "org.jackdbus.Emitter.Heartbeat", i)
		client.conn.Emit("/org/jackdbus/test", "org.jackdbus.Emitter.Value", i)
		time.Sleep(50 * time.Millisecond)
	}

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, uint32(3), counter)

	require.NoError(t, detachPassive())

	client.conn.Emit("/org/jackdbus/test", "org.jackdbus.Emitter.Heartbeat", uint32(3))
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, uint32(4), counter)

	require.NoError(t, detach())
	client.conn.Emit("/org/jackdbus/test", "org.jackdbus.Emitter.Heartbeat", uint32(4))
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, uint32(4), counter)
}

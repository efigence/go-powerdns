package ipredir

import (
	"github.com/efigence/go-powerdns/backend/memdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var backend, _ = New(memdb.New())

func TestAdd(t *testing.T) {
	t.Run("Add IP redir", func(t *testing.T) {
		err := backend.AddRedirIp("1.2.3.4", "5.6.7.8")
		require.NoError(t, err)
		redirIp, err := backend.ListRedirIp()
		require.NoError(t, err)
		assert.Equal(t, "5.6.7.8", redirIp["1.2.3.4"])
	})
	t.Run("Batch set IP", func(t *testing.T) {
		backend.AddRedirIp("1.2.3.4", "5.6.7.8")
		err := backend.SetRedirIp(map[string]string{
			"2.2.2.2": "2.3.3.3",
			"3.2.2.2": "3.3.3.3",
			"4.2.2.2": "4.3.3.3",
		})
		require.NoError(t, err)
		redirIp, err := backend.ListRedirIp()
		require.NoError(t, err)
		assert.NotEqual(t, "5.6.7.8", redirIp["1.2.3.4"], "previous value removed")
		assert.Equal(t, "2.3.3.3", redirIp["2.2.2.2"])
		assert.Equal(t, "3.3.3.3", redirIp["3.2.2.2"])
		assert.Equal(t, "4.3.3.3", redirIp["4.2.2.2"])

	})
	t.Run("Delete IP", func(t *testing.T) {
		err := backend.SetRedirIp(map[string]string{
			"5.2.2.2": "5.3.3.3",
			"6.2.2.2": "6.3.3.3",
			"7.2.2.2": "7.3.3.3",
		})
		require.NoError(t, err)
		err = backend.DeleteRedirIp("6.2.2.2")
		require.NoError(t, err)
		err = backend.DeleteRedirIp("99.99.99.99")
		require.NoError(t, err)
		redirIp, err := backend.ListRedirIp()
		require.NoError(t, err)
		t.Run("Deleted IP should not exist", func(t *testing.T) {
			assert.NotContains(t, redirIp, "99.99.99.99")
			assert.NotContains(t, redirIp, "6.2.2.2")
		})
		t.Run("Non-deleted IPs should exist", func(t *testing.T) {
			assert.Equal(t, "5.3.3.3", redirIp["5.2.2.2"])
			assert.Equal(t, "7.3.3.3", redirIp["7.2.2.2"])
		})

	})
	t.Run("bad IP", func(t *testing.T) {
		assert.Error(t, backend.AddRedirIp("256.254.254.254", "1.2.3.4"))
		assert.Error(t, backend.AddRedirIp("1.2.3.4", "999.1.1.1"))
	})

}

package test

import (
	"github.com/crawlab-team/crawlab/db/redis"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRedisClient_Ping(t *testing.T) {
	var err error
	T.Setup(t)

	err = T.client.Ping()
	require.Nil(t, err)
}

func TestRedisClient_Get_Set(t *testing.T) {
	var err error
	T.Setup(t)

	err = T.client.Set(T.TestCollection, T.TestMessage)
	require.Nil(t, err)

	value, err := T.client.Get(T.TestCollection)
	require.Nil(t, err)
	require.Equal(t, T.TestMessage, value)
}

func TestRedisClient_Keys_AllKeys(t *testing.T) {
	var err error
	T.Setup(t)

	for _, key := range T.TestKeysAlpha {
		err = T.client.Set(key, key)
		require.Nil(t, err)
	}
	for _, key := range T.TestKeysBeta {
		err = T.client.Set(key, key)
		require.Nil(t, err)
	}

	keys, err := T.client.Keys("*alpha*")
	require.Nil(t, err)
	require.Len(t, keys, len(T.TestKeysAlpha))

	keys, err = T.client.Keys("*beta*")
	require.Nil(t, err)
	require.Len(t, keys, len(T.TestKeysBeta))

	keys, err = T.client.AllKeys()
	require.Nil(t, err)
	require.Len(t, keys, len(T.TestKeysAlpha)+len(T.TestKeysBeta))
}

func TestRedisClient_RPush_LPop_LLen(t *testing.T) {
	var err error
	T.Setup(t)

	for _, msg := range T.TestMessages {
		err = T.client.RPush(T.TestCollection, msg)
		require.Nil(t, err)
	}

	n, err := T.client.LLen(T.TestCollection)
	require.Nil(t, err)
	require.Equal(t, len(T.TestMessages), n)

	value, err := T.client.LPop(T.TestCollection)
	require.Nil(t, err)
	require.Equal(t, T.TestMessages[0], value)
}

func TestRedisClient_LPush_RPop(t *testing.T) {
	var err error
	T.Setup(t)

	for _, msg := range T.TestMessages {
		err = T.client.LPush(T.TestCollection, msg)
		require.Nil(t, err)
	}

	n, err := T.client.LLen(T.TestCollection)
	require.Nil(t, err)
	require.Equal(t, len(T.TestMessages), n)

	value, err := T.client.RPop(T.TestCollection)
	require.Nil(t, err)
	require.Equal(t, T.TestMessages[0], value)
}

func TestRedisClient_BRPop(t *testing.T) {
	var err error
	T.Setup(t)

	isErr := true
	go func(t *testing.T) {
		value, err := T.client.BRPop(T.TestCollection, 0)
		require.Nil(t, err)
		require.Equal(t, T.TestMessage, value)
		isErr = false
	}(t)

	err = T.client.LPush(T.TestCollection, T.TestMessage)
	require.Nil(t, err)
	time.Sleep(500 * time.Millisecond)
	require.False(t, isErr)
}

func TestRedisClient_BLPop(t *testing.T) {
	var err error
	T.Setup(t)

	isErr := true
	go func(t *testing.T) {
		value, err := T.client.BLPop(T.TestCollection, 0)
		require.Nil(t, err)
		require.Equal(t, T.TestMessage, value)
		isErr = false
	}(t)

	err = T.client.RPush(T.TestCollection, T.TestMessage)
	require.Nil(t, err)
	time.Sleep(500 * time.Millisecond)
	require.False(t, isErr)
}

func TestRedisClient_HSet_HGet_HDel(t *testing.T) {
	var err error
	T.Setup(t)

	for k, v := range T.TestMessagesMap {
		err = T.client.HSet(T.TestCollection, k, v)
		require.Nil(t, err)
	}

	for k, v := range T.TestMessagesMap {
		vr, err := T.client.HGet(T.TestCollection, k)
		require.Nil(t, err)
		require.Equal(t, v, vr)
	}

	for k := range T.TestMessagesMap {
		err = T.client.HDel(T.TestCollection, k)
		require.Nil(t, err)

		v, err := T.client.HGet(T.TestCollection, k)
		require.Nil(t, err)
		require.Empty(t, v)
	}
}

func TestRedisClient_HScan(t *testing.T) {
	var err error
	T.Setup(t)

	for k, v := range T.TestMessagesMap {
		err = T.client.HSet(T.TestCollection, k, v)
		require.Nil(t, err)
	}

	results, err := T.client.HScan(T.TestCollection)
	require.Nil(t, err)

	for k, vr := range results {
		v, ok := T.TestMessagesMap[k]
		require.True(t, ok)
		require.Equal(t, v, vr)
	}
}

func TestRedisClient_HKeys(t *testing.T) {
	var err error
	T.Setup(t)

	for k, v := range T.TestMessagesMap {
		err = T.client.HSet(T.TestCollection, k, v)
		require.Nil(t, err)
	}

	keys, err := T.client.HKeys(T.TestCollection)
	require.Nil(t, err)

	for _, k := range keys {
		_, ok := T.TestMessagesMap[k]
		require.True(t, ok)
	}
}

func TestRedisClient_ZAdd_ZCount_ZCountAll_ZPopMax_ZPopMin(t *testing.T) {
	var err error
	T.Setup(t)

	for i, v := range T.TestMessages {
		score := float32(i)
		err = T.client.ZAdd(T.TestCollection, score, v)
		require.Nil(t, err)
	}

	count, err := T.client.ZCountAll(T.TestCollection)
	require.Nil(t, err)
	require.Equal(t, len(T.TestMessages), count)

	value, err := T.client.ZPopMaxOne(T.TestCollection)
	require.Nil(t, err)
	require.Equal(t, T.TestMessages[len(T.TestMessages)-1], value)

	value, err = T.client.ZPopMinOne(T.TestCollection)
	require.Nil(t, err)
	require.Equal(t, T.TestMessages[0], value)
}

func TestRedisClient_BZPopMax_BZPopMin(t *testing.T) {
	var err error
	T.Setup(t)

	isErr := true
	go func(t *testing.T) {
		value, err := T.client.BZPopMax(T.TestCollection, 0)
		require.Nil(t, err)
		require.Equal(t, T.TestMessage, value)
		isErr = false
	}(t)

	err = T.client.ZAdd(T.TestCollection, 1, T.TestMessage)
	require.Nil(t, err)
	time.Sleep(500 * time.Millisecond)
	require.False(t, isErr)

	isErr = true
	go func(t *testing.T) {
		value, err := T.client.BZPopMin(T.TestCollection, 0)
		require.Nil(t, err)
		require.Equal(t, T.TestMessage, value)
		isErr = false
	}(t)

	err = T.client.ZAdd(T.TestCollection, 1, T.TestMessage)
	require.Nil(t, err)
	time.Sleep(500 * time.Millisecond)
	require.False(t, isErr)
}

func TestRedisClient_Lock_Unlock(t *testing.T) {
	var err error
	T.Setup(t)

	ts, err := T.client.Lock(T.TestLockKey)
	require.Nil(t, err)

	_, err = T.client.Lock(T.TestLockKey)
	require.NotNil(t, err)

	T.client.UnLock(T.TestLockKey, ts)

	ts, err = T.client.Lock(T.TestLockKey)
	require.Nil(t, err)

}

func TestRedisClient_MemoryStats(t *testing.T) {
	var err error
	T.Setup(t)

	stats, err := T.client.MemoryStats()
	require.Nil(t, err)

	for _, k := range redis.MemoryStatsMetrics {
		v, ok := stats[k]
		require.True(t, ok)
		require.Greater(t, v, int64(-1))
	}
}

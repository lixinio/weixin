package redis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRedisUrl(t *testing.T) {
	testcases := []struct {
		input  string
		expect struct {
			db       int
			host     string
			password string
		}
	}{
		{
			input: "redis://127.0.0.1:6379/0",
			expect: struct {
				db       int
				host     string
				password string
			}{
				db:       0,
				host:     "127.0.0.1:6379",
				password: "",
			},
		},
		{
			input: "redis://:secrets@127.0.0.1:6379/3",
			expect: struct {
				db       int
				host     string
				password string
			}{
				db:       3,
				host:     "127.0.0.1:6379",
				password: "secrets",
			},
		},
		{
			input: "redis://a.com:12345",
			expect: struct {
				db       int
				host     string
				password string
			}{
				db:       0,
				host:     "a.com:12345",
				password: "",
			},
		},
	}

	for _, testcase := range testcases {
		redisDB, redisHost, redisPwd, err := parseRedisURL(testcase.input)
		require.Equal(t, nil, err)
		require.Equal(t, testcase.expect.db, redisDB)
		require.Equal(t, testcase.expect.host, redisHost)
		require.Equal(t, testcase.expect.password, redisPwd)
	}
}

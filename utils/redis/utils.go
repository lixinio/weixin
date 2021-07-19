package redis

import (
	"net/url"
	"strconv"
	"strings"
)

func parseRedisURL(urlStr string) (int, string, string, error) {
	redisURL, err := url.Parse(urlStr)
	if err != nil {
		return 0, "", "", err
	}

	redisPwd := ""
	if redisURL.User != nil {
		if password, ok := redisURL.User.Password(); ok {
			redisPwd = password
		}
	}

	redisDb := 0
	if len(redisURL.Path) > 1 {
		db := strings.TrimPrefix(redisURL.Path, "/")
		intVar, err := strconv.Atoi(db)
		if err != nil {
			return 0, "", "", err
		}
		redisDb = intVar
	}

	return redisDb, redisURL.Host, redisPwd, nil
}

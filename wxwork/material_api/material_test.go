package material_api

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/wxwork"
	agentApi "github.com/lixinio/weixin/wxwork/agent"
	"github.com/stretchr/testify/require"
)

func TestMaterialUrl(t *testing.T) {
	cache := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	corp := wxwork.New(&wxwork.Config{
		Corpid: test.CorpID,
	})
	agent := agentApi.New(corp, cache, &agentApi.Config{
		AgentId: test.AgentID,
		Secret:  test.AgentSecret,
	})

	materialApi := NewAgentApi(agent)

	file, err := os.Open(test.ImagePath)
	require.Empty(t, err)
	defer file.Close()

	url, err := materialApi.UploadImg(test.ImagePath, file)
	require.Empty(t, err)
	fmt.Print(url)

	// 下载
	response, err := http.Get(url)
	require.Empty(t, err)
	defer response.Body.Close()

	_, _, err = image.Decode(response.Body)
	require.Empty(t, err)
}

func TestMaterialID(t *testing.T) {
	cache := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	corp := wxwork.New(&wxwork.Config{
		Corpid: test.CorpID,
	})
	agent := agentApi.New(corp, cache, &agentApi.Config{
		AgentId: test.AgentID,
		Secret:  test.AgentSecret,
	})

	materialApi := NewAgentApi(agent)

	file, err := os.Open(test.ImagePath)
	require.Empty(t, err)
	defer file.Close()

	result, err := materialApi.Upload(test.ImagePath, file, MediaTypeImage)
	require.Empty(t, err)
	require.Equal(t, result.Type, MediaTypeImage)
	fmt.Println(result.MediaID, result.CreatedAt, result.Type)

	// 计算源文件hash
	originHash := ""
	{
		file.Seek(0, 0)
		hasher := sha256.New()
		_, err = io.Copy(hasher, file)
		require.Empty(t, err)
		originHash = hex.EncodeToString(hasher.Sum(nil))
	}

	{
		resp, err := materialApi.Get(result.MediaID)
		require.Empty(t, err)

		// 计算hash
		hasher := sha256.New()
		_, err = io.Copy(hasher, resp.Body)
		require.Empty(t, err)

		thisHash := hex.EncodeToString(hasher.Sum(nil))
		require.Equal(t, thisHash, originHash)
	}

	{
		mediaID := fmt.Sprintf("0%s", result.MediaID)
		_, err := materialApi.Get(mediaID)
		require.NotEmpty(t, err)
	}
}

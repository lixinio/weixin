package material_api

import (
	"context"
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
	"github.com/lixinio/weixin/utils"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/weixin/official_account"
	"github.com/stretchr/testify/require"
)

func initOfficialAccount() *utils.Client {
	redis := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	officialAccount := official_account.New(redis, redis, &official_account.Config{
		Appid:  test.OfficialAccountAppid,
		Secret: test.OfficialAccountSecret,
	})
	return officialAccount.Client
}

func TestMaterialUrl(t *testing.T) {
	ctx := context.Background()
	for _, client := range []*utils.Client{
		initOfficialAccount(),
	} {
		materialApi := NewApi(client)

		file, err := os.Open(test.ImagePath)
		require.Empty(t, err)
		defer file.Close()

		fi, err := file.Stat()
		require.Empty(t, err)

		url, err := materialApi.UploadImg(ctx, test.ImagePath, fi.Size(), file)
		require.Empty(t, err)
		fmt.Print(url)

		// 下载
		response, err := http.Get(url)
		require.Empty(t, err)
		defer response.Body.Close()

		_, _, err = image.Decode(response.Body)
		require.Empty(t, err)
	}
}

func TestUploadMedia(t *testing.T) {
	ctx := context.Background()
	for _, client := range []*utils.Client{
		initOfficialAccount(),
	} {
		materialApi := NewApi(client)

		file, err := os.Open(test.ImagePath)
		require.Empty(t, err)
		defer file.Close()

		fi, err := file.Stat()
		require.Empty(t, err)

		result, err := materialApi.UploadMedia(ctx, test.ImagePath, fi.Size(), file, MediaTypeImage)
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
			resp, err := materialApi.GetMedia(ctx, result.MediaID)
			require.Empty(t, err)

			// 计算hash
			hasher := sha256.New()
			_, err = hasher.Write(resp)
			require.Empty(t, err)

			thisHash := hex.EncodeToString(hasher.Sum(nil))
			require.Equal(t, thisHash, originHash)
		}

		{
			r, w := io.Pipe()
			go func(t *testing.T) {
				defer w.Close()
				err := materialApi.SaveMedia(ctx, result.MediaID, w)
				require.Empty(t, err)
			}(t)

			// 计算hash
			hasher := sha256.New()
			_, err = io.Copy(hasher, r)
			require.Empty(t, err)

			thisHash := hex.EncodeToString(hasher.Sum(nil))
			require.Equal(t, thisHash, originHash)
		}

		{
			mediaID := fmt.Sprintf("0%s", result.MediaID)
			_, err := materialApi.GetMedia(ctx, mediaID)
			require.NotEmpty(t, err)
		}
	}
}

func TestUploadImageMaterial(t *testing.T) {
	ctx := context.Background()
	for _, client := range []*utils.Client{
		initOfficialAccount(),
	} {
		materialApi := NewApi(client)

		file, err := os.Open(test.ImagePath)
		require.Empty(t, err)
		defer file.Close()

		fi, err := file.Stat()
		require.Empty(t, err)

		result, err := materialApi.UploadMaterial(
			ctx,
			test.ImagePath,
			fi.Size(),
			file,
			MediaTypeImage,
		)
		require.Empty(t, err)
		fmt.Println(result.MediaID, result.URL)

		// 永久素材， 图片素材不一致
		{
			_, err := materialApi.GetMaterial(ctx, result.MediaID)
			require.Empty(t, err)
		}

		{
			_, w := io.Pipe()
			go func(t *testing.T) {
				defer w.Close()
				err := materialApi.SaveMaterial(ctx, result.MediaID, w)
				require.Empty(t, err)
			}(t)
		}

		{
			mediaID := fmt.Sprintf("0%s", result.MediaID)
			_, err := materialApi.GetMaterial(ctx, mediaID)
			require.NotEmpty(t, err)
		}

		{
			err = materialApi.DeleteMaterial(ctx, result.MediaID)
			require.Empty(t, err)
		}
	}
}

func TestUploadVoiceMaterial(t *testing.T) {
	ctx := context.Background()
	for _, client := range []*utils.Client{
		initOfficialAccount(),
	} {
		materialApi := NewApi(client)

		file, err := os.Open(test.AudioPath)
		require.Empty(t, err)
		defer file.Close()

		fi, err := file.Stat()
		require.Empty(t, err)

		result, err := materialApi.UploadMaterial(
			ctx,
			test.AudioPath,
			fi.Size(),
			file,
			MediaTypeVoice,
		)
		require.Empty(t, err)
		fmt.Println(result.MediaID, result.URL)

		// 永久素材， Voice素材不一致
		{
			_, err := materialApi.GetMaterial(ctx, result.MediaID)
			require.Empty(t, err)
		}

		{
			_, w := io.Pipe()
			go func(t *testing.T) {
				defer w.Close()
				err := materialApi.SaveMaterial(ctx, result.MediaID, w)
				require.Empty(t, err)
			}(t)
		}

		{
			mediaID := fmt.Sprintf("0%s", result.MediaID)
			_, err := materialApi.GetMaterial(ctx, mediaID)
			require.NotEmpty(t, err)
		}

		{
			err = materialApi.DeleteMaterial(ctx, result.MediaID)
			require.Empty(t, err)
		}
	}
}

func TestUploadVideoMaterial(t *testing.T) {
	ctx := context.Background()
	for _, client := range []*utils.Client{
		initOfficialAccount(),
	} {
		materialApi := NewApi(client)

		file, err := os.Open(test.VideoPath)
		require.Empty(t, err)
		defer file.Close()

		fi, err := file.Stat()
		require.Empty(t, err)

		result, err := materialApi.UploadVideoMaterial(
			ctx, test.VideoPath, "Title", "fjasdklfjasd", fi.Size(), file,
		)
		require.Empty(t, err)
		fmt.Println(result.MediaID, result.URL)

		{
			m, err := materialApi.GetVideoMaterial(ctx, result.MediaID)
			require.Empty(t, err)
			fmt.Println(m.Description, m.DownloadUrl, m.Title)
		}

		{
			err = materialApi.DeleteMaterial(ctx, result.MediaID)
			require.Empty(t, err)
		}
	}
}

func TestCountMaterial(t *testing.T) {
	ctx := context.Background()
	for _, client := range []*utils.Client{
		initOfficialAccount(),
	} {
		materialApi := NewApi(client)

		result, err := materialApi.GetMaterialStatistics(ctx)
		require.Empty(t, err)
		fmt.Println(result.ImageCount, result.VoiceCount, result.VideoCount, result.NewsCount)

		{
			materials, err := materialApi.ListMaterial(ctx, MediaTypeVoice, 0, 100)
			require.Empty(t, err)
			fmt.Println(materials.ItemCount, materials.TotalCount)
			for _, material := range materials.Items {
				fmt.Println(material.MediaID, material.Name, material.URL, material.UpdateTime)
			}
		}

		{
			materials, err := materialApi.ListMaterial(ctx, MediaTypeImage, 0, 100)
			require.Empty(t, err)
			fmt.Println(materials.ItemCount, materials.TotalCount)
			for _, material := range materials.Items {
				fmt.Println(material.MediaID, material.Name, material.URL, material.UpdateTime)
			}
		}

		{
			materials, err := materialApi.ListMaterial(ctx, MediaTypeVideo, 0, 100)
			require.Empty(t, err)
			fmt.Println(materials.ItemCount, materials.TotalCount)
			for _, material := range materials.Items {
				fmt.Println(material.MediaID, material.Name, material.URL, material.UpdateTime)
			}
		}

		{
			materials, err := materialApi.ListMpnewsMaterial(ctx, 0, 100)
			require.Empty(t, err)
			fmt.Println(materials.ItemCount, materials.TotalCount)
			for _, material := range materials.Items {
				fmt.Println(material.MediaID)
				for _, item := range materials.Items {
					fmt.Println(item.Content, item.MediaID)
				}
			}
		}
	}
}

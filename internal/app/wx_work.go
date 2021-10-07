package app

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/airdb/sailor/dbutil"
	"github.com/airdb/wxwork-kf/pkg/util"
	"github.com/go-redis/redis/v8"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/credential"
	wxContext "github.com/silenceper/wechat/v2/officialaccount/context"
	"github.com/silenceper/wechat/v2/officialaccount/material"
	"github.com/silenceper/wechat/v2/work"
	"github.com/silenceper/wechat/v2/work/config"
	"github.com/silenceper/wechat/v2/work/kf"
)

type wxWorkMedia interface {
	MediaUpload(mediaType material.MediaType, filename string) (media util.Media, err error)
}

var (
	WxWorkKF    *kf.Client
	WxWorkMedia wxWorkMedia
	Redis       *redis.Client
)

func InitWxWork() {
	var err error

	// Init Database.
	dbutil.InitDefaultDB()

	redisDB, err := strconv.Atoi(os.Getenv("WXKF_REDIS_DB"))
	if err != nil {
		redisDB = 2
	}

	// Redis
	Redis = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("WXKF_REDIS_ADDR"),
		Password: os.Getenv("WXKF_REDIS_PASSWD"),
		DB:       redisDB,
	})
	err = Redis.Ping(context.Background()).Err()
	if err != nil {
		log.Panicf("init redis err: %s", err.Error())
	}

	// SDK
	cfg := &config.Config{
		CorpID:     os.Getenv("WXKF_CORP_ID"),
		CorpSecret: os.Getenv("WXKF_SECRET"),
		// AgentID: "",
		Cache: cache.NewRedis(&cache.RedisOpts{
			Host:     os.Getenv("WXKF_REDIS_ADDR"),
			Password: os.Getenv("WXKF_REDIS_PASSWD"),
			Database: redisDB,
		}),
		RasPrivateKey: "",

		Token:          os.Getenv("WXKF_TOKEN"),
		EncodingAESKey: os.Getenv("WXKF_ENCODING_AES_KEY"),
	}

	clientWork := work.NewWork(cfg)

	if WxWorkKF, err = clientWork.GetKF(); err != nil {
		log.Panicf("init sdk err: %s", err.Error())
	}

	wxCtx := &wxContext.Context{
		Config: nil,
		AccessTokenHandle: credential.NewWorkAccessToken(
			cfg.CorpID, cfg.CorpSecret, credential.CacheKeyWorkPrefix, cfg.Cache),
	}
	WxWorkMedia = util.NewMaterial(wxCtx)
}

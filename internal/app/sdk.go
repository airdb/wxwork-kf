package app

import (
	"context"
	"log"
	"os"
	"strconv"

	sdk "github.com/NICEXAI/WeChatCustomerServiceSDK"
	"github.com/NICEXAI/WeChatCustomerServiceSDK/cache"
	"github.com/NICEXAI/WeChatCustomerServiceSDK/crypto"
	"github.com/go-redis/redis/v8"
)

var (
	wxOption sdk.Options
	WxCpt    *crypto.WXBizMsgCrypt
	WxClient *sdk.Client
	Redis    *redis.Client
)

func InitSdk() {
	var err error

	wxOption.CorpID = os.Getenv("WXKF_CORP_ID")
	wxOption.Secret = os.Getenv("WXKF_SECRET")
	wxOption.Token = os.Getenv("WXKF_SECRET")
	wxOption.EncodingAESKey = os.Getenv("WXKF_ENCODING_AES_KEY")

	var wxRedisOption cache.RedisOptions
	if redisDB, err := strconv.Atoi(os.Getenv("WXKF_REDIS_DB")); err == nil {
		// 初始化默认 redis
		Redis = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("WXKF_REDIS_ADDR"),
			Password: os.Getenv("WXKF_REDIS_PASSWD"),
			DB:       redisDB,
		})
		err = Redis.Ping(context.Background()).Err()
		if err != nil {
			log.Panicf("init redis err: %s", err.Error())
		}
		wxRedisOption.Addr = os.Getenv("WXKF_REDIS_ADDR")
		wxRedisOption.Password = os.Getenv("WXKF_REDIS_PASSWD")
		wxRedisOption.DB = redisDB
		wxOption.Cache = cache.NewRedis(wxRedisOption)
	} else {
		log.Panicf("init redis err: %s", err.Error())
	}

	log.Printf("get env from serverless.yaml: redis: %s,  corpID: %s, secret: %s, token: %s, encodingAESKey: %s\n",
		wxRedisOption.Addr,
		wxOption.CorpID,
		wxOption.Secret,
		wxOption.Token,
		wxOption.EncodingAESKey,
	)

	WxCpt = crypto.NewWXBizMsgCrypt(
		wxOption.Token, wxOption.EncodingAESKey, wxOption.CorpID, crypto.XmlType,
	)
	if WxClient, err = sdk.New(wxOption); err != nil {
		log.Panicf("init sdk err: %s", err.Error())
	}
}

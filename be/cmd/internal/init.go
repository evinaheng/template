package internal

import (
	"log"

	"github.com/template/be/lib/nsq"
	"github.com/template/be/lib/nsq/gonsq"
	"github.com/template/be/lib/panics"
	"github.com/template/be/lib/panics/tolog"
	"github.com/template/be/lib/panics/toslack"
	"github.com/template/be/lib/storage/database"
	"github.com/template/be/lib/storage/database/sqlt"
	"github.com/template/be/lib/storage/elastic"
	v5 "github.com/template/be/lib/storage/elastic/v5"
	"github.com/template/be/lib/storage/redis"
	"github.com/template/be/lib/storage/redis/redigo"
)

// init NSQ prodocuer module
func initNSQProducer(endpoint string) nsq.Producer {

	config := gonsq.Config{
		Endpoint: endpoint,
	}
	producer, err := gonsq.NewProducer(config)
	if err != nil {
		log.Fatalln("[FATAL] Can't connect to NSQ PRODUCER", err)
	}

	return producer
}

// init elastic module
func initElastic(cfg map[string]*ElasticConfig) map[string]elastic.Elastic {

	res := map[string]elastic.Elastic{}
	for k, v := range cfg {
		elasticModule := v5.New(v5.Config{
			Endpoint: v.Endpoint,
		})

		if elasticModule == nil {
			log.Fatalln("[FATAL] Can't connect to ELASTIC", k, v.Endpoint)
		}
		res[k] = elasticModule
	}
	return res
}

// init redis module
func initRedis(cfg map[string]*RedisConfig) map[string]redis.Redis {
	res := map[string]redis.Redis{}

	for k, v := range cfg {
		redisModule := redigo.New(redigo.Config{
			Endpoint: v.Endpoint,
			MaxIdle:  v.MaxIdle,
			Timeout:  v.Timeout,
		})

		if redisModule == nil {
			log.Fatalln("[FATAL] Can't connect to REDIS", k, v.Endpoint)
		}
		res[k] = redisModule
	}

	return res
}

// init database module
func initDatabase(cfg map[string]*DatabaseConfig) map[string]database.Database {
	res := map[string]database.Database{}

	for k, v := range cfg {
		dbModule := sqlt.New(sqlt.Config{
			Driver: v.Driver,
			Master: v.Master,
			Slave:  v.Slave,
		})

		if dbModule == nil {
			log.Fatalln("[FATAL] Can't connect to DATABASE", k, v.Master)
		}
		res[k] = dbModule
	}

	return res
}

// InitPanics wrapper
func InitPanics(currEnv, ipAddress string) panics.Panics {

	var isDevelopment = currEnv == Development
	var isProduction = currEnv == Production

	var panicModule panics.Panics
	// Init panics module
	if isDevelopment {
		panicModule = tolog.New()
	} else {
		panicModule = toslack.InitPanic(toslack.Config{
			Env:         currEnv,
			IPAddress:   ipAddress,
			WithMention: isProduction,
			SlackURL:    "",
		})
	}
	if panicModule == nil {
		log.Println("Can't load module PANICS")
	}

	return panicModule
}

# Redis

Initialization
```golang

	
	config := map[string]*redis.Config{}
	
	config["server_1"] = &redis.Config{
		Endpoint: 	"http:localhost:6379",
		Timeout:	time.Duration(10 * time.Seconds),
		MaxIdle: 	10,
	}

	// Create redis module
	// Initialize this ONLY ONCE
	redis.Init(cfg.Redis, isMocking)
	
	
```

Get redis value
```golang
	rds, err := redis.Connect("server_1")
	if err != nil {
		...
	}

	value, err := rds.Get("mykey").Int()
	if err != nil {
		...
	}

	fmt.Print(value)
```


Setex redis value
```golang
	rds, err := redis.Connect("server_1")
	if err != nil {
		...
	}

	expireTime := 3600
	if err := rds.Setex("mykey", expireTime, "testing"); err != nil {
		...
	}

```
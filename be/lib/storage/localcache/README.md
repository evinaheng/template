# Localcache

Local value with custom expiration time. All values are local and removed if the server is stopped/restarted.

Get value
```golang
	value := localcache.Get("mykey")
```

Set value
```golang
	localcache.Set("mykey", "value", 60) // Expired in 60 seconds
	localcache.Set("mykey", "value", -1) // No expiry
```

Delete value
```golang
	localcache.Delete("mykey")
```

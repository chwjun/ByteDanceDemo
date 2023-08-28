package redis

//var Ctx = context.Background()
//var RdbTest *redis.Client
//
//// UserFollowings 根据用户id找到他关注的人
//var UserFollowings *redis.Client
//
//// UserFollowers 根据用户id找到他的粉丝
//var UserFollowers *redis.Client
//
//// UserFriends 根据用户id找到他的好友
//var UserFriends *redis.Client
//
//// RdbVCid 存储video与comment的关系
//var RdbVCid *redis.Client
//
//// RdbCVid 根据commentId找videoId
//var RdbCVid *redis.Client
//
//// RdbCIdComment 根据commentId 找comment
//var RdbCIdComment *redis.Client
//
//var (
//	ProdRedisAddr string
//	ProRedisPwd   string
//	ExpireTime    time.Duration
//)
//
//// InitRedis 初始化 Redis 连接，redis 默认 16 个 DB
//func InitRedis() {
//	ProdRedisAddr = fmt.Sprintf("%s:%s",
//		viper.GetString("settings.redis.host"),
//		viper.GetString("settings.redis.port"),
//	)
//	ProRedisPwd = viper.GetString("settings.redis.password")
//	ExpireTime = viper.GetDuration("settings.redis.expirationTime") * time.Second
//	RdbTest = redis.NewClient(&redis.Options{
//		Addr:     ProdRedisAddr,
//		Password: ProRedisPwd,
//		DB:       0,
//	})
//	RdbVCid = redis.NewClient(&redis.Options{
//		Addr:     ProdRedisAddr,
//		Password: ProRedisPwd,
//		DB:       1,
//	})
//	RdbCVid = redis.NewClient(&redis.Options{
//		Addr:     ProdRedisAddr,
//		Password: ProRedisPwd,
//		DB:       2,
//	})
//	RdbCIdComment = redis.NewClient(&redis.Options{
//		Addr:     ProdRedisAddr,
//		Password: ProRedisPwd,
//		DB:       3,
//	})
//	UserFollowings = redis.NewClient(&redis.Options{
//		Addr:     ProdRedisAddr,
//		Password: ProRedisPwd,
//		DB:       11,
//	})
//	UserFollowers = redis.NewClient(&redis.Options{
//		Addr:     ProdRedisAddr,
//		Password: ProRedisPwd,
//		DB:       12,
//	})
//	UserFriends = redis.NewClient(&redis.Options{
//		Addr:     ProdRedisAddr,
//		Password: ProRedisPwd,
//		DB:       13,
//	})
//	_, err := RdbTest.Ping(Ctx).Result()
//	if err != nil {
//		log.Panicf("连接 redis 错误，错误信息: %v", err)
//	} else {
//		log.Println("Redis 连接成功！")
//	}
//}

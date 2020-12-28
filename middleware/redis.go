package middleware

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"

	"toolbox"
	"toolbox/stat"
)

const (
	TYPE_PROXY    = "proxy"
	TYPE_SENTINEL = "sentinel"
	TYPE_CLUSTER  = "cluster"
)

// 包装redis.Pool不让上层用户直接使用第三方库中的redis，
// 因为会造成用户code中也import上面的github路径，这样就
// 不能控制app的使用方式。
type RedisClient struct {
	*redis.Client
}

type RedisClusterClient struct {
	*redis.ClusterClient
}

type DpRedisClient interface {
	Pipeline() redis.Pipeliner
	Pipelined(fn func(redis.Pipeliner) error) ([]redis.Cmder, error)

	TxPipelined(fn func(redis.Pipeliner) error) ([]redis.Cmder, error)
	TxPipeline() redis.Pipeliner

	Command() *redis.CommandsInfoCmd
	ClientGetName() *redis.StringCmd
	Echo(message interface{}) *redis.StringCmd
	Ping() *redis.StatusCmd
	Quit() *redis.StatusCmd
	Del(keys ...string) *redis.IntCmd
	Unlink(keys ...string) *redis.IntCmd
	Dump(key string) *redis.StringCmd
	Exists(keys ...string) *redis.IntCmd
	Expire(key string, expiration time.Duration) *redis.BoolCmd
	ExpireAt(key string, tm time.Time) *redis.BoolCmd
	Keys(pattern string) *redis.StringSliceCmd
	Migrate(host, port, key string, db int64, timeout time.Duration) *redis.StatusCmd
	Move(key string, db int64) *redis.BoolCmd
	ObjectRefCount(key string) *redis.IntCmd
	ObjectEncoding(key string) *redis.StringCmd
	ObjectIdleTime(key string) *redis.DurationCmd
	Persist(key string) *redis.BoolCmd
	PExpire(key string, expiration time.Duration) *redis.BoolCmd
	PExpireAt(key string, tm time.Time) *redis.BoolCmd
	PTTL(key string) *redis.DurationCmd
	RandomKey() *redis.StringCmd
	Rename(key, newkey string) *redis.StatusCmd
	RenameNX(key, newkey string) *redis.BoolCmd
	Restore(key string, ttl time.Duration, value string) *redis.StatusCmd
	RestoreReplace(key string, ttl time.Duration, value string) *redis.StatusCmd
	Sort(key string, sort *redis.Sort) *redis.StringSliceCmd
	SortStore(key, store string, sort *redis.Sort) *redis.IntCmd
	SortInterfaces(key string, sort *redis.Sort) *redis.SliceCmd
	Touch(keys ...string) *redis.IntCmd
	TTL(key string) *redis.DurationCmd
	Type(key string) *redis.StatusCmd
	Scan(cursor uint64, match string, count int64) *redis.ScanCmd
	SScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd
	HScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd
	ZScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd
	Append(key, value string) *redis.IntCmd
	BitCount(key string, bitCount *redis.BitCount) *redis.IntCmd
	BitOpAnd(destKey string, keys ...string) *redis.IntCmd
	BitOpOr(destKey string, keys ...string) *redis.IntCmd
	BitOpXor(destKey string, keys ...string) *redis.IntCmd
	BitOpNot(destKey string, key string) *redis.IntCmd
	BitPos(key string, bit int64, pos ...int64) *redis.IntCmd
	Decr(key string) *redis.IntCmd
	DecrBy(key string, decrement int64) *redis.IntCmd
	Get(key string) *redis.StringCmd
	GetBit(key string, offset int64) *redis.IntCmd
	GetRange(key string, start, end int64) *redis.StringCmd
	GetSet(key string, value interface{}) *redis.StringCmd
	Incr(key string) *redis.IntCmd
	IncrBy(key string, value int64) *redis.IntCmd
	IncrByFloat(key string, value float64) *redis.FloatCmd
	MGet(keys ...string) *redis.SliceCmd
	MSet(pairs ...interface{}) *redis.StatusCmd
	MSetNX(pairs ...interface{}) *redis.BoolCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	SetBit(key string, offset int64, value int) *redis.IntCmd
	SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	SetXX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	SetRange(key string, offset int64, value string) *redis.IntCmd
	StrLen(key string) *redis.IntCmd
	HDel(key string, fields ...string) *redis.IntCmd
	HExists(key, field string) *redis.BoolCmd
	HGet(key, field string) *redis.StringCmd
	HGetAll(key string) *redis.StringStringMapCmd
	HIncrBy(key, field string, incr int64) *redis.IntCmd
	HIncrByFloat(key, field string, incr float64) *redis.FloatCmd
	HKeys(key string) *redis.StringSliceCmd
	HLen(key string) *redis.IntCmd
	HMGet(key string, fields ...string) *redis.SliceCmd
	HMSet(key string, fields map[string]interface{}) *redis.StatusCmd
	HSet(key, field string, value interface{}) *redis.BoolCmd
	HSetNX(key, field string, value interface{}) *redis.BoolCmd
	HVals(key string) *redis.StringSliceCmd
	BLPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd
	BRPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd
	BRPopLPush(source, destination string, timeout time.Duration) *redis.StringCmd
	LIndex(key string, index int64) *redis.StringCmd
	LInsert(key, op string, pivot, value interface{}) *redis.IntCmd
	LInsertBefore(key string, pivot, value interface{}) *redis.IntCmd
	LInsertAfter(key string, pivot, value interface{}) *redis.IntCmd
	LLen(key string) *redis.IntCmd
	LPop(key string) *redis.StringCmd
	LPush(key string, values ...interface{}) *redis.IntCmd
	LPushX(key string, value interface{}) *redis.IntCmd
	LRange(key string, start, stop int64) *redis.StringSliceCmd
	LRem(key string, count int64, value interface{}) *redis.IntCmd
	LSet(key string, index int64, value interface{}) *redis.StatusCmd
	LTrim(key string, start, stop int64) *redis.StatusCmd
	RPop(key string) *redis.StringCmd
	RPopLPush(source, destination string) *redis.StringCmd
	RPush(key string, values ...interface{}) *redis.IntCmd
	RPushX(key string, value interface{}) *redis.IntCmd
	SAdd(key string, members ...interface{}) *redis.IntCmd
	SCard(key string) *redis.IntCmd
	SDiff(keys ...string) *redis.StringSliceCmd
	SDiffStore(destination string, keys ...string) *redis.IntCmd
	SInter(keys ...string) *redis.StringSliceCmd
	SInterStore(destination string, keys ...string) *redis.IntCmd
	SIsMember(key string, member interface{}) *redis.BoolCmd
	SMembers(key string) *redis.StringSliceCmd
	SMembersMap(key string) *redis.StringStructMapCmd
	SMove(source, destination string, member interface{}) *redis.BoolCmd
	SPop(key string) *redis.StringCmd
	SPopN(key string, count int64) *redis.StringSliceCmd
	SRandMember(key string) *redis.StringCmd
	SRandMemberN(key string, count int64) *redis.StringSliceCmd
	SRem(key string, members ...interface{}) *redis.IntCmd
	SUnion(keys ...string) *redis.StringSliceCmd
	SUnionStore(destination string, keys ...string) *redis.IntCmd
	XAdd(a *redis.XAddArgs) *redis.StringCmd
	XLen(stream string) *redis.IntCmd
	XRange(stream, start, stop string) *redis.XMessageSliceCmd
	XRangeN(stream, start, stop string, count int64) *redis.XMessageSliceCmd
	XRevRange(stream string, start, stop string) *redis.XMessageSliceCmd
	XRevRangeN(stream string, start, stop string, count int64) *redis.XMessageSliceCmd
	XRead(a *redis.XReadArgs) *redis.XStreamSliceCmd
	XReadStreams(streams ...string) *redis.XStreamSliceCmd
	XGroupCreate(stream, group, start string) *redis.StatusCmd
	XGroupSetID(stream, group, start string) *redis.StatusCmd
	XGroupDestroy(stream, group string) *redis.IntCmd
	XGroupDelConsumer(stream, group, consumer string) *redis.IntCmd
	XReadGroup(a *redis.XReadGroupArgs) *redis.XStreamSliceCmd
	XAck(stream, group string, ids ...string) *redis.IntCmd
	XPending(stream, group string) *redis.XPendingCmd
	XPendingExt(a *redis.XPendingExtArgs) *redis.XPendingExtCmd
	XClaim(a *redis.XClaimArgs) *redis.XMessageSliceCmd
	XClaimJustID(a *redis.XClaimArgs) *redis.StringSliceCmd
	XTrim(key string, maxLen int64) *redis.IntCmd
	XTrimApprox(key string, maxLen int64) *redis.IntCmd
	ZAdd(key string, members ...redis.Z) *redis.IntCmd
	ZAddNX(key string, members ...redis.Z) *redis.IntCmd
	ZAddXX(key string, members ...redis.Z) *redis.IntCmd
	ZAddCh(key string, members ...redis.Z) *redis.IntCmd
	ZAddNXCh(key string, members ...redis.Z) *redis.IntCmd
	ZAddXXCh(key string, members ...redis.Z) *redis.IntCmd
	ZIncr(key string, member redis.Z) *redis.FloatCmd
	ZIncrNX(key string, member redis.Z) *redis.FloatCmd
	ZIncrXX(key string, member redis.Z) *redis.FloatCmd
	ZCard(key string) *redis.IntCmd
	ZCount(key, min, max string) *redis.IntCmd
	ZLexCount(key, min, max string) *redis.IntCmd
	ZIncrBy(key string, increment float64, member string) *redis.FloatCmd
	ZInterStore(destination string, store redis.ZStore, keys ...string) *redis.IntCmd
	ZRange(key string, start, stop int64) *redis.StringSliceCmd
	ZRangeWithScores(key string, start, stop int64) *redis.ZSliceCmd
	ZRangeByScore(key string, opt redis.ZRangeBy) *redis.StringSliceCmd
	ZRangeByLex(key string, opt redis.ZRangeBy) *redis.StringSliceCmd
	ZRangeByScoreWithScores(key string, opt redis.ZRangeBy) *redis.ZSliceCmd
	ZRank(key, member string) *redis.IntCmd
	ZRem(key string, members ...interface{}) *redis.IntCmd
	ZRemRangeByRank(key string, start, stop int64) *redis.IntCmd
	ZRemRangeByScore(key, min, max string) *redis.IntCmd
	ZRemRangeByLex(key, min, max string) *redis.IntCmd
	ZRevRange(key string, start, stop int64) *redis.StringSliceCmd
	ZRevRangeWithScores(key string, start, stop int64) *redis.ZSliceCmd
	ZRevRangeByScore(key string, opt redis.ZRangeBy) *redis.StringSliceCmd
	ZRevRangeByLex(key string, opt redis.ZRangeBy) *redis.StringSliceCmd
	ZRevRangeByScoreWithScores(key string, opt redis.ZRangeBy) *redis.ZSliceCmd
	ZRevRank(key, member string) *redis.IntCmd
	ZScore(key, member string) *redis.FloatCmd
	ZUnionStore(dest string, store redis.ZStore, keys ...string) *redis.IntCmd
	PFAdd(key string, els ...interface{}) *redis.IntCmd
	PFCount(keys ...string) *redis.IntCmd
	PFMerge(dest string, keys ...string) *redis.StatusCmd
	BgRewriteAOF() *redis.StatusCmd
	BgSave() *redis.StatusCmd
	ClientKill(ipPort string) *redis.StatusCmd
	ClientKillByFilter(keys ...string) *redis.IntCmd
	ClientList() *redis.StringCmd
	ClientPause(dur time.Duration) *redis.BoolCmd
	ConfigGet(parameter string) *redis.SliceCmd
	ConfigResetStat() *redis.StatusCmd
	ConfigSet(parameter, value string) *redis.StatusCmd
	ConfigRewrite() *redis.StatusCmd
	DBSize() *redis.IntCmd
	FlushAll() *redis.StatusCmd
	FlushAllAsync() *redis.StatusCmd
	FlushDB() *redis.StatusCmd
	FlushDBAsync() *redis.StatusCmd
	Info(section ...string) *redis.StringCmd
	LastSave() *redis.IntCmd
	Save() *redis.StatusCmd
	Shutdown() *redis.StatusCmd
	ShutdownSave() *redis.StatusCmd
	ShutdownNoSave() *redis.StatusCmd
	SlaveOf(host, port string) *redis.StatusCmd
	Time() *redis.TimeCmd
	Eval(script string, keys []string, args ...interface{}) *redis.Cmd
	EvalSha(sha1 string, keys []string, args ...interface{}) *redis.Cmd
	ScriptExists(hashes ...string) *redis.BoolSliceCmd
	ScriptFlush() *redis.StatusCmd
	ScriptKill() *redis.StatusCmd
	ScriptLoad(script string) *redis.StringCmd
	DebugObject(key string) *redis.StringCmd
	Publish(channel string, message interface{}) *redis.IntCmd
	PubSubChannels(pattern string) *redis.StringSliceCmd
	PubSubNumSub(channels ...string) *redis.StringIntMapCmd
	PubSubNumPat() *redis.IntCmd
	ClusterSlots() *redis.ClusterSlotsCmd
	ClusterNodes() *redis.StringCmd
	ClusterMeet(host, port string) *redis.StatusCmd
	ClusterForget(nodeID string) *redis.StatusCmd
	ClusterReplicate(nodeID string) *redis.StatusCmd
	ClusterResetSoft() *redis.StatusCmd
	ClusterResetHard() *redis.StatusCmd
	ClusterInfo() *redis.StringCmd
	ClusterKeySlot(key string) *redis.IntCmd
	ClusterCountFailureReports(nodeID string) *redis.IntCmd
	ClusterCountKeysInSlot(slot int) *redis.IntCmd
	ClusterDelSlots(slots ...int) *redis.StatusCmd
	ClusterDelSlotsRange(min, max int) *redis.StatusCmd
	ClusterSaveConfig() *redis.StatusCmd
	ClusterSlaves(nodeID string) *redis.StringSliceCmd
	ClusterFailover() *redis.StatusCmd
	ClusterAddSlots(slots ...int) *redis.StatusCmd
	ClusterAddSlotsRange(min, max int) *redis.StatusCmd
	GeoAdd(key string, geoLocation ...*redis.GeoLocation) *redis.IntCmd
	GeoPos(key string, members ...string) *redis.GeoPosCmd
	GeoRadius(key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd
	GeoRadiusRO(key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd
	GeoRadiusByMember(key, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd
	GeoRadiusByMemberRO(key, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd
	GeoDist(key string, member1, member2, unit string) *redis.FloatCmd
	GeoHash(key string, members ...string) *redis.StringSliceCmd
	ReadOnly() *redis.StatusCmd
	ReadWrite() *redis.StatusCmd
	MemoryUsage(key string, samples ...int) *redis.IntCmd
}

// redis库内部存在Options，这里根据最佳实践，筛选出需要
// app关注的，二次封装，虽然屏蔽细节，但经过打磨后应该会
// 节省app开发时间。
type redisOptions struct {
	masterName       string
	addr             string
	maxOpenConnCount int64
	idleTimeout      time.Duration
	maxConnAge       time.Duration
	maxRetries       int64
	password         string
	method           string

	// 下面3个选项，在初始化时时放在Dial方法中用的，这里app不需要关注
	connTimeout  time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
}

var defaultRedisOptions = redisOptions{
	masterName:       "mymaster",
	maxOpenConnCount: 50,
	idleTimeout:      240 * time.Second,
	maxRetries:       1,
	maxConnAge:       500 * time.Second,
	password:         "",
	method:           TYPE_PROXY,

	connTimeout:  300 * time.Millisecond,
	readTimeout:  1 * time.Second,
	writeTimeout: 1 * time.Second,
}

type redisOptionsFunc func(*redisOptions)

func RedisAddr(s string) redisOptionsFunc {
	return func(o *redisOptions) {
		o.addr = s
	}
}
func RedisMethod(s string) redisOptionsFunc {
	return func(o *redisOptions) {
		o.method = s
	}
}

func RedisMaxOpenConnCount(s int64) redisOptionsFunc {
	return func(o *redisOptions) {
		o.maxOpenConnCount = s
	}
}

func RedisMasterName(s string) redisOptionsFunc {
	return func(o *redisOptions) {
		o.masterName = s
	}
}

func RedisPassword(s string) redisOptionsFunc {
	return func(o *redisOptions) {
		o.password = s
	}
}

func RedisIdleTimeout(s int64) redisOptionsFunc {
	return func(o *redisOptions) {
		o.idleTimeout = time.Duration(s) * time.Millisecond
	}
}
func RedisMaxConnAge(s int64) redisOptionsFunc {
	return func(o *redisOptions) {
		o.maxConnAge = time.Duration(s) * time.Millisecond
	}
}
func RedisMaxRetries(s int64) redisOptionsFunc {
	return func(o *redisOptions) {
		o.maxRetries = s
	}
}

func RedisConnTimeout(s int64) redisOptionsFunc {
	return func(o *redisOptions) {
		o.connTimeout = time.Duration(s) * time.Millisecond
	}
}

func RedisReadTimeout(s int64) redisOptionsFunc {
	return func(o *redisOptions) {
		o.readTimeout = time.Duration(s) * time.Millisecond
	}
}

func RedisWriteTimeout(s int64) redisOptionsFunc {
	return func(o *redisOptions) {
		o.writeTimeout = time.Duration(s) * time.Millisecond
	}
}
func initRedisSentinel(opts redisOptions) (DpRedisClient, error) {

	//不支持sfns
	addrs := strings.Split(strings.Split(opts.addr, "//")[1], ",")

	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    opts.masterName,
		SentinelAddrs: addrs,
		DialTimeout:   opts.connTimeout,
		ReadTimeout:   opts.readTimeout,
		WriteTimeout:  opts.writeTimeout,
		MaxRetries:    int(opts.maxRetries),
		Password:      opts.password,

		PoolSize:    int(opts.maxOpenConnCount),
		MaxConnAge:  opts.maxConnAge,
		IdleTimeout: opts.idleTimeout,
	})

	redis.SetLogger(log.New(os.Stdout, "redis: ", log.LstdFlags|log.Lshortfile))

	return &RedisClient{client}, nil
}
func initRedisProxy(opts redisOptions) (DpRedisClient, error) {
	//foxns, err := ns.New(ns.WithService(opts.addr), ns.WithConnTimeout(opts.connTimeout))
	//if err != nil {
	//	return nil, err
	//}

	client := redis.NewClient(&redis.Options{
		Addr:         opts.addr,
		DialTimeout:  opts.connTimeout,
		ReadTimeout:  opts.readTimeout,
		WriteTimeout: opts.writeTimeout,
		MaxRetries:   int(opts.maxRetries),
		Password:     opts.password,

		PoolSize:    int(opts.maxOpenConnCount),
		MaxConnAge:  opts.maxConnAge,
		IdleTimeout: opts.idleTimeout,

		// 支持sfns
		//Dialer: foxns.Dial,
	})

	redis.SetLogger(log.New(os.Stdout, "redis: ", log.LstdFlags|log.Lshortfile))

	return &RedisClient{client}, nil
}
func initRedisCluster(opts redisOptions) (DpRedisClient, error) {

	//不支持sfns
	addrs := strings.Split(strings.Split(opts.addr, "//")[1], ",")

	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        addrs,
		Password:     opts.password,
		DialTimeout:  opts.connTimeout,
		ReadTimeout:  opts.readTimeout,
		WriteTimeout: opts.writeTimeout,
		MaxRetries:   int(opts.maxRetries),

		PoolSize:    int(opts.maxOpenConnCount),
		MaxConnAge:  opts.maxConnAge,
		IdleTimeout: opts.idleTimeout,
	})

	redis.SetLogger(log.New(os.Stdout, "redis: ", log.LstdFlags|log.Lshortfile))

	return &RedisClusterClient{clusterClient}, nil
}

func InitRedis(opt ...redisOptionsFunc) (DpRedisClient, error) {
	opts := defaultRedisOptions
	for _, o := range opt {
		o(&opts)
	}
	switch opts.method {
	case TYPE_SENTINEL:
		return initRedisSentinel(opts)
	case TYPE_PROXY:
		return initRedisProxy(opts)
	case TYPE_CLUSTER:
		return initRedisCluster(opts)
	default:
		return nil, fmt.Errorf("type errr")
	}

}

func (redis *RedisClient) Pipeline() redis.Pipeliner {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Pipeline()
}

func (redis *RedisClient) Pipelined(fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Pipelined(fn)
}

func (redis *RedisClient) TxPipelined(fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.TxPipelined(fn)
}

func (redis *RedisClient) TxPipeline() redis.Pipeliner {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.TxPipeline()
}

func (redis *RedisClient) Command() *redis.CommandsInfoCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Command()
}

func (redis *RedisClient) ClientGetName() *redis.StringCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ClientGetName()
}

func (redis *RedisClient) Echo(message interface{}) *redis.StringCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Echo(message)
}

func (redis *RedisClient) Ping() *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Ping()
}

func (redis *RedisClient) Quit() *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Quit()
}

func (redis *RedisClient) Del(keys ...string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Del(keys...)
}

func (redis *RedisClient) Unlink(keys ...string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Unlink(keys...)
}

func (redis *RedisClient) Dump(key string) *redis.StringCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Dump(key)
}

func (redis *RedisClient) Exists(keys ...string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Exists(keys...)
}

func (redis *RedisClient) Expire(key string, expiration time.Duration) *redis.BoolCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Expire(key, expiration)
}

func (redis *RedisClient) ExpireAt(key string, tm time.Time) *redis.BoolCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ExpireAt(key, tm)
}

func (redis *RedisClient) Keys(pattern string) *redis.StringSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Keys(pattern)
}

func (redis *RedisClient) Migrate(host string, port string, key string, db int64, timeout time.Duration) *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Migrate(host, port, key, db, timeout)
}

func (redis *RedisClient) Move(key string, db int64) *redis.BoolCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Move(key, db)
}

func (redis *RedisClient) ObjectRefCount(key string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ObjectRefCount(key)
}

func (redis *RedisClient) ObjectEncoding(key string) *redis.StringCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ObjectEncoding(key)
}

func (redis *RedisClient) ObjectIdleTime(key string) *redis.DurationCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ObjectIdleTime(key)
}

func (redis *RedisClient) Persist(key string) *redis.BoolCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Persist(key)
}

func (redis *RedisClient) PExpire(key string, expiration time.Duration) *redis.BoolCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.PExpire(key, expiration)
}

func (redis *RedisClient) PExpireAt(key string, tm time.Time) *redis.BoolCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.PExpireAt(key, tm)
}

func (redis *RedisClient) PTTL(key string) *redis.DurationCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.PTTL(key)
}

func (redis *RedisClient) RandomKey() *redis.StringCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.RandomKey()
}

func (redis *RedisClient) Rename(key string, newkey string) *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Rename(key, newkey)
}

func (redis *RedisClient) RenameNX(key string, newkey string) *redis.BoolCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.RenameNX(key, newkey)
}

func (redis *RedisClient) Restore(key string, ttl time.Duration, value string) *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Restore(key, ttl, value)
}

func (redis *RedisClient) RestoreReplace(key string, ttl time.Duration, value string) *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.RestoreReplace(key, ttl, value)
}

func (redis *RedisClient) Sort(key string, sort *redis.Sort) *redis.StringSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Sort(key, sort)
}

func (redis *RedisClient) SortStore(key string, store string, sort *redis.Sort) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SortStore(key, store, sort)
}

func (redis *RedisClient) SortInterfaces(key string, sort *redis.Sort) *redis.SliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SortInterfaces(key, sort)
}

func (redis *RedisClient) Touch(keys ...string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Touch(keys...)
}

func (redis *RedisClient) TTL(key string) *redis.DurationCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.TTL(key)
}

func (redis *RedisClient) Type(key string) *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Type(key)
}

func (redis *RedisClient) Scan(cursor uint64, match string, count int64) *redis.ScanCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Scan(cursor, match, count)
}

func (redis *RedisClient) SScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SScan(key, cursor, match, count)
}

func (redis *RedisClient) HScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.HScan(key, cursor, match, count)
}

func (redis *RedisClient) ZScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZScan(key, cursor, match, count)
}

func (redis *RedisClient) Append(key string, value string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Append(key, value)
}

func (redis *RedisClient) BitCount(key string, bitCount *redis.BitCount) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.BitCount(key, bitCount)
}

func (redis *RedisClient) BitOpAnd(destKey string, keys ...string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.BitOpAnd(destKey, keys...)
}

func (redis *RedisClient) BitOpOr(destKey string, keys ...string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.BitOpOr(destKey, keys...)
}

func (redis *RedisClient) BitOpXor(destKey string, keys ...string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.BitOpXor(destKey, keys...)
}

func (redis *RedisClient) BitOpNot(destKey string, key string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.BitOpNot(destKey, key)
}

func (redis *RedisClient) BitPos(key string, bit int64, pos ...int64) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.BitPos(key, bit, pos...)
}

func (redis *RedisClient) Decr(key string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Decr(key)
}

func (redis *RedisClient) DecrBy(key string, decrement int64) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.DecrBy(key, decrement)
}

func (redis *RedisClient) Get(key string) *redis.StringCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Get(key)
}

func (redis *RedisClient) GetBit(key string, offset int64) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.GetBit(key, offset)
}

func (redis *RedisClient) GetRange(key string, start int64, end int64) *redis.StringCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.GetRange(key, start, end)
}

func (redis *RedisClient) GetSet(key string, value interface{}) *redis.StringCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.GetSet(key, value)
}

func (redis *RedisClient) Incr(key string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Incr(key)
}

func (redis *RedisClient) IncrBy(key string, value int64) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.IncrBy(key, value)
}

func (redis *RedisClient) IncrByFloat(key string, value float64) *redis.FloatCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.IncrByFloat(key, value)
}

func (redis *RedisClient) MGet(keys ...string) *redis.SliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.MGet(keys...)
}

func (redis *RedisClient) MSet(pairs ...interface{}) *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.MSet(pairs...)
}

func (redis *RedisClient) MSetNX(pairs ...interface{}) *redis.BoolCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.MSetNX(pairs...)
}

func (redis *RedisClient) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Set(key, value, expiration)
}

func (redis *RedisClient) SetBit(key string, offset int64, value int) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SetBit(key, offset, value)
}

func (redis *RedisClient) SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SetNX(key, value, expiration)
}

func (redis *RedisClient) SetXX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SetXX(key, value, expiration)
}

func (redis *RedisClient) SetRange(key string, offset int64, value string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SetRange(key, offset, value)
}

func (redis *RedisClient) StrLen(key string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.StrLen(key)
}

func (redis *RedisClient) HDel(key string, fields ...string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.HDel(key, fields...)
}

func (redis *RedisClient) HExists(key string, field string) *redis.BoolCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.HExists(key, field)
}

func (redis *RedisClient) HGet(key string, field string) *redis.StringCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.HGet(key, field)
}

func (redis *RedisClient) HGetAll(key string) *redis.StringStringMapCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.HGetAll(key)
}

func (redis *RedisClient) HIncrBy(key string, field string, incr int64) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.HIncrBy(key, field, incr)
}

func (redis *RedisClient) HIncrByFloat(key string, field string, incr float64) *redis.FloatCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.HIncrByFloat(key, field, incr)
}

func (redis *RedisClient) HKeys(key string) *redis.StringSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.HKeys(key)
}

func (redis *RedisClient) HLen(key string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.HLen(key)
}

func (redis *RedisClient) HMGet(key string, fields ...string) *redis.SliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.HMGet(key, fields...)
}

func (redis *RedisClient) HMSet(key string, fields map[string]interface{}) *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.HMSet(key, fields)
}

func (redis *RedisClient) HSet(key string, field string, value interface{}) *redis.BoolCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.HSet(key, field, value)
}

func (redis *RedisClient) HSetNX(key string, field string, value interface{}) *redis.BoolCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.HSetNX(key, field, value)
}

func (redis *RedisClient) HVals(key string) *redis.StringSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.HVals(key)
}

func (redis *RedisClient) BLPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.BLPop(timeout, keys...)
}

func (redis *RedisClient) BRPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.BRPop(timeout, keys...)
}

func (redis *RedisClient) BRPopLPush(source string, destination string, timeout time.Duration) *redis.StringCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.BRPopLPush(source, destination, timeout)
}

func (redis *RedisClient) LIndex(key string, index int64) *redis.StringCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.LIndex(key, index)
}

func (redis *RedisClient) LInsert(key string, op string, pivot interface{}, value interface{}) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.LInsert(key, op, pivot, value)
}

func (redis *RedisClient) LInsertBefore(key string, pivot interface{}, value interface{}) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.LInsertBefore(key, pivot, value)
}

func (redis *RedisClient) LInsertAfter(key string, pivot interface{}, value interface{}) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.LInsertAfter(key, pivot, value)
}

func (redis *RedisClient) LLen(key string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.LLen(key)
}

func (redis *RedisClient) LPop(key string) *redis.StringCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.LPop(key)
}

func (redis *RedisClient) LPush(key string, values ...interface{}) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.LPush(key, values...)
}

func (redis *RedisClient) LPushX(key string, value interface{}) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.LPushX(key, value)
}

func (redis *RedisClient) LRange(key string, start int64, stop int64) *redis.StringSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.LRange(key, start, stop)
}

func (redis *RedisClient) LRem(key string, count int64, value interface{}) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.LRem(key, count, value)
}

func (redis *RedisClient) LSet(key string, index int64, value interface{}) *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.LSet(key, index, value)
}

func (redis *RedisClient) LTrim(key string, start int64, stop int64) *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.LTrim(key, start, stop)
}

func (redis *RedisClient) RPop(key string) *redis.StringCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.RPop(key)
}

func (redis *RedisClient) RPopLPush(source string, destination string) *redis.StringCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.RPopLPush(source, destination)
}

func (redis *RedisClient) RPush(key string, values ...interface{}) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.RPush(key, values...)
}

func (redis *RedisClient) RPushX(key string, value interface{}) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.RPushX(key, value)
}

func (redis *RedisClient) SAdd(key string, members ...interface{}) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SAdd(key, members...)
}

func (redis *RedisClient) SCard(key string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SCard(key)
}

func (redis *RedisClient) SDiff(keys ...string) *redis.StringSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SDiff(keys...)
}

func (redis *RedisClient) SDiffStore(destination string, keys ...string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SDiffStore(destination, keys...)
}

func (redis *RedisClient) SInter(keys ...string) *redis.StringSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SInter(keys...)
}

func (redis *RedisClient) SInterStore(destination string, keys ...string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SInterStore(destination, keys...)
}

func (redis *RedisClient) SIsMember(key string, member interface{}) *redis.BoolCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SIsMember(key, member)
}

func (redis *RedisClient) SMembers(key string) *redis.StringSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SMembers(key)
}

func (redis *RedisClient) SMembersMap(key string) *redis.StringStructMapCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SMembersMap(key)
}

func (redis *RedisClient) SMove(source string, destination string, member interface{}) *redis.BoolCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SMove(source, destination, member)
}

func (redis *RedisClient) SPop(key string) *redis.StringCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SPop(key)
}

func (redis *RedisClient) SPopN(key string, count int64) *redis.StringSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SPopN(key, count)
}

func (redis *RedisClient) SRandMember(key string) *redis.StringCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SRandMember(key)
}

func (redis *RedisClient) SRandMemberN(key string, count int64) *redis.StringSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SRandMemberN(key, count)
}

func (redis *RedisClient) SRem(key string, members ...interface{}) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SRem(key, members...)
}

func (redis *RedisClient) SUnion(keys ...string) *redis.StringSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SUnion(keys...)
}

func (redis *RedisClient) SUnionStore(destination string, keys ...string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SUnionStore(destination, keys...)
}

func (redis *RedisClient) XAdd(a *redis.XAddArgs) *redis.StringCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.XAdd(a)
}

func (redis *RedisClient) XLen(stream string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.XLen(stream)
}

func (redis *RedisClient) XRange(stream string, start string, stop string) *redis.XMessageSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.XRange(stream, start, stop)
}

func (redis *RedisClient) XRangeN(stream string, start string, stop string, count int64) *redis.XMessageSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.XRangeN(stream, start, stop, count)
}

func (redis *RedisClient) XRevRange(stream string, start string, stop string) *redis.XMessageSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.XRevRange(stream, start, stop)
}

func (redis *RedisClient) XRevRangeN(stream string, start string, stop string, count int64) *redis.XMessageSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.XRevRangeN(stream, start, stop, count)
}

func (redis *RedisClient) XRead(a *redis.XReadArgs) *redis.XStreamSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.XRead(a)
}

func (redis *RedisClient) XReadStreams(streams ...string) *redis.XStreamSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.XReadStreams(streams...)
}

func (redis *RedisClient) XGroupCreate(stream string, group string, start string) *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.XGroupCreate(stream, group, start)
}

func (redis *RedisClient) XGroupSetID(stream string, group string, start string) *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.XGroupSetID(stream, group, start)
}

func (redis *RedisClient) XGroupDestroy(stream string, group string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.XGroupDestroy(stream, group)
}

func (redis *RedisClient) XGroupDelConsumer(stream string, group string, consumer string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.XGroupDelConsumer(stream, group, consumer)
}

func (redis *RedisClient) XReadGroup(a *redis.XReadGroupArgs) *redis.XStreamSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.XReadGroup(a)
}

func (redis *RedisClient) XAck(stream string, group string, ids ...string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.XAck(stream, group, ids...)
}

func (redis *RedisClient) XPending(stream string, group string) *redis.XPendingCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.XPending(stream, group)
}

func (redis *RedisClient) XPendingExt(a *redis.XPendingExtArgs) *redis.XPendingExtCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.XPendingExt(a)
}

func (redis *RedisClient) XClaim(a *redis.XClaimArgs) *redis.XMessageSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.XClaim(a)
}

func (redis *RedisClient) XClaimJustID(a *redis.XClaimArgs) *redis.StringSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.XClaimJustID(a)
}

func (redis *RedisClient) XTrim(key string, maxLen int64) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.XTrim(key, maxLen)
}

func (redis *RedisClient) XTrimApprox(key string, maxLen int64) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.XTrimApprox(key, maxLen)
}

func (redis *RedisClient) ZAdd(key string, members ...redis.Z) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZAdd(key, members...)
}

func (redis *RedisClient) ZAddNX(key string, members ...redis.Z) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZAddNX(key, members...)
}

func (redis *RedisClient) ZAddXX(key string, members ...redis.Z) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZAddXX(key, members...)
}

func (redis *RedisClient) ZAddCh(key string, members ...redis.Z) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZAddCh(key, members...)
}

func (redis *RedisClient) ZAddNXCh(key string, members ...redis.Z) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZAddNXCh(key, members...)
}

func (redis *RedisClient) ZAddXXCh(key string, members ...redis.Z) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZAddXXCh(key, members...)
}

func (redis *RedisClient) ZIncr(key string, member redis.Z) *redis.FloatCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZIncr(key, member)
}

func (redis *RedisClient) ZIncrNX(key string, member redis.Z) *redis.FloatCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZIncrNX(key, member)
}

func (redis *RedisClient) ZIncrXX(key string, member redis.Z) *redis.FloatCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZIncrXX(key, member)
}

func (redis *RedisClient) ZCard(key string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZCard(key)
}

func (redis *RedisClient) ZCount(key string, min string, max string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZCount(key, min, max)
}

func (redis *RedisClient) ZLexCount(key string, min string, max string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZLexCount(key, min, max)
}

func (redis *RedisClient) ZIncrBy(key string, increment float64, member string) *redis.FloatCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZIncrBy(key, increment, member)
}

func (redis *RedisClient) ZInterStore(destination string, store redis.ZStore, keys ...string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZInterStore(destination, store, keys...)
}

func (redis *RedisClient) ZRange(key string, start int64, stop int64) *redis.StringSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZRange(key, start, stop)
}

func (redis *RedisClient) ZRangeWithScores(key string, start int64, stop int64) *redis.ZSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZRangeWithScores(key, start, stop)
}

func (redis *RedisClient) ZRangeByScore(key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZRangeByScore(key, opt)
}

func (redis *RedisClient) ZRangeByLex(key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZRangeByLex(key, opt)
}

func (redis *RedisClient) ZRangeByScoreWithScores(key string, opt redis.ZRangeBy) *redis.ZSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZRangeByScoreWithScores(key, opt)
}

func (redis *RedisClient) ZRank(key string, member string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZRank(key, member)
}

func (redis *RedisClient) ZRem(key string, members ...interface{}) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZRem(key, members...)
}

func (redis *RedisClient) ZRemRangeByRank(key string, start int64, stop int64) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZRemRangeByRank(key, start, stop)
}

func (redis *RedisClient) ZRemRangeByScore(key string, min string, max string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZRemRangeByScore(key, min, max)
}

func (redis *RedisClient) ZRemRangeByLex(key string, min string, max string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZRemRangeByLex(key, min, max)
}

func (redis *RedisClient) ZRevRange(key string, start int64, stop int64) *redis.StringSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZRevRange(key, start, stop)
}

func (redis *RedisClient) ZRevRangeWithScores(key string, start int64, stop int64) *redis.ZSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZRevRangeWithScores(key, start, stop)
}

func (redis *RedisClient) ZRevRangeByScore(key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZRevRangeByScore(key, opt)
}

func (redis *RedisClient) ZRevRangeByLex(key string, opt redis.ZRangeBy) *redis.StringSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZRevRangeByLex(key, opt)
}

func (redis *RedisClient) ZRevRangeByScoreWithScores(key string, opt redis.ZRangeBy) *redis.ZSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZRevRangeByScoreWithScores(key, opt)
}

func (redis *RedisClient) ZRevRank(key string, member string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZRevRank(key, member)
}

func (redis *RedisClient) ZScore(key string, member string) *redis.FloatCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZScore(key, member)
}

func (redis *RedisClient) ZUnionStore(dest string, store redis.ZStore, keys ...string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ZUnionStore(dest, store, keys...)
}

func (redis *RedisClient) PFAdd(key string, els ...interface{}) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.PFAdd(key, els...)
}

func (redis *RedisClient) PFCount(keys ...string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.PFCount(keys...)
}

func (redis *RedisClient) PFMerge(dest string, keys ...string) *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.PFMerge(dest, keys...)
}

func (redis *RedisClient) BgRewriteAOF() *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.BgRewriteAOF()
}

func (redis *RedisClient) BgSave() *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.BgSave()
}

func (redis *RedisClient) ClientKill(ipPort string) *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ClientKill(ipPort)
}

func (redis *RedisClient) ClientKillByFilter(keys ...string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ClientKillByFilter(keys...)
}

func (redis *RedisClient) ClientList() *redis.StringCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ClientList()
}

func (redis *RedisClient) ClientPause(dur time.Duration) *redis.BoolCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ClientPause(dur)
}

func (redis *RedisClient) ConfigGet(parameter string) *redis.SliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ConfigGet(parameter)
}

func (redis *RedisClient) ConfigResetStat() *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ConfigResetStat()
}

func (redis *RedisClient) ConfigSet(parameter string, value string) *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ConfigSet(parameter, value)
}

func (redis *RedisClient) ConfigRewrite() *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ConfigRewrite()
}

func (redis *RedisClient) DBSize() *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.DBSize()
}

func (redis *RedisClient) FlushAll() *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.FlushAll()
}

func (redis *RedisClient) FlushAllAsync() *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.FlushAllAsync()
}

func (redis *RedisClient) FlushDB() *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.FlushDB()
}

func (redis *RedisClient) FlushDBAsync() *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.FlushDBAsync()
}

func (redis *RedisClient) Info(section ...string) *redis.StringCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Info(section...)
}

func (redis *RedisClient) LastSave() *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.LastSave()
}

func (redis *RedisClient) Save() *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Save()
}

func (redis *RedisClient) Shutdown() *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Shutdown()
}

func (redis *RedisClient) ShutdownSave() *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ShutdownSave()
}

func (redis *RedisClient) ShutdownNoSave() *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ShutdownNoSave()
}

func (redis *RedisClient) SlaveOf(host string, port string) *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.SlaveOf(host, port)
}

func (redis *RedisClient) Time() *redis.TimeCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Time()
}

func (redis *RedisClient) Eval(script string, keys []string, args ...interface{}) *redis.Cmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Eval(script, keys, args...)
}

func (redis *RedisClient) EvalSha(sha1 string, keys []string, args ...interface{}) *redis.Cmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.EvalSha(sha1, keys, args...)
}

func (redis *RedisClient) ScriptExists(hashes ...string) *redis.BoolSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ScriptExists(hashes...)
}

func (redis *RedisClient) ScriptFlush() *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ScriptFlush()
}

func (redis *RedisClient) ScriptKill() *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ScriptKill()
}

func (redis *RedisClient) ScriptLoad(script string) *redis.StringCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ScriptLoad(script)
}

func (redis *RedisClient) DebugObject(key string) *redis.StringCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.DebugObject(key)
}

func (redis *RedisClient) Publish(channel string, message interface{}) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.Publish(channel, message)
}

func (redis *RedisClient) PubSubChannels(pattern string) *redis.StringSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.PubSubChannels(pattern)
}

func (redis *RedisClient) PubSubNumSub(channels ...string) *redis.StringIntMapCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.PubSubNumSub(channels...)
}

func (redis *RedisClient) PubSubNumPat() *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.PubSubNumPat()
}

func (redis *RedisClient) ClusterSlots() *redis.ClusterSlotsCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ClusterSlots()
}

func (redis *RedisClient) ClusterNodes() *redis.StringCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ClusterNodes()
}

func (redis *RedisClient) ClusterMeet(host string, port string) *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ClusterMeet(host, port)
}

func (redis *RedisClient) ClusterForget(nodeID string) *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ClusterForget(nodeID)
}

func (redis *RedisClient) ClusterReplicate(nodeID string) *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ClusterReplicate(nodeID)
}

func (redis *RedisClient) ClusterResetSoft() *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ClusterResetSoft()
}

func (redis *RedisClient) ClusterResetHard() *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ClusterResetHard()
}

func (redis *RedisClient) ClusterInfo() *redis.StringCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ClusterInfo()
}

func (redis *RedisClient) ClusterKeySlot(key string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ClusterKeySlot(key)
}

func (redis *RedisClient) ClusterCountFailureReports(nodeID string) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ClusterCountFailureReports(nodeID)
}

func (redis *RedisClient) ClusterCountKeysInSlot(slot int) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ClusterCountKeysInSlot(slot)
}

func (redis *RedisClient) ClusterDelSlots(slots ...int) *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ClusterDelSlots(slots...)
}

func (redis *RedisClient) ClusterDelSlotsRange(min int, max int) *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ClusterDelSlotsRange(min, max)
}

func (redis *RedisClient) ClusterSaveConfig() *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ClusterSaveConfig()
}

func (redis *RedisClient) ClusterSlaves(nodeID string) *redis.StringSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ClusterSlaves(nodeID)
}

func (redis *RedisClient) ClusterFailover() *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ClusterFailover()
}

func (redis *RedisClient) ClusterAddSlots(slots ...int) *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ClusterAddSlots(slots...)
}

func (redis *RedisClient) ClusterAddSlotsRange(min int, max int) *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ClusterAddSlotsRange(min, max)
}

func (redis *RedisClient) GeoAdd(key string, geoLocation ...*redis.GeoLocation) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.GeoAdd(key, geoLocation...)
}

func (redis *RedisClient) GeoPos(key string, members ...string) *redis.GeoPosCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.GeoPos(key, members...)
}

func (redis *RedisClient) GeoRadius(key string, longitude float64, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.GeoRadius(key, longitude, latitude, query)
}

func (redis *RedisClient) GeoRadiusRO(key string, longitude float64, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.GeoRadiusRO(key, longitude, latitude, query)
}

func (redis *RedisClient) GeoRadiusByMember(key string, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.GeoRadiusByMember(key, member, query)
}

func (redis *RedisClient) GeoRadiusByMemberRO(key string, member string, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.GeoRadiusByMemberRO(key, member, query)
}

func (redis *RedisClient) GeoDist(key string, member1 string, member2 string, unit string) *redis.FloatCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.GeoDist(key, member1, member2, unit)
}

func (redis *RedisClient) GeoHash(key string, members ...string) *redis.StringSliceCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.GeoHash(key, members...)
}

func (redis *RedisClient) ReadOnly() *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ReadOnly()
}

func (redis *RedisClient) ReadWrite() *redis.StatusCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.ReadWrite()
}

func (redis *RedisClient) MemoryUsage(key string, samples ...int) *redis.IntCmd {
	startTime := time.Now()
	defer stat.ClientStat(toolbox.StatMetrix(stat.Redis), startTime)
	return redis.Client.MemoryUsage(key, samples...)
}

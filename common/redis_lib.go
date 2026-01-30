// redis_lib
package common

import (
	"context"
	//"fmt"
	"log"
	//"os"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-redis/redis/v9" //https://pkg.go.dev/github.com/go-redis/redis/v8#XAddArgs
)

// var rdb *redis.UniversalClient
var Rdb *redis.Client

//var rdb *redis.ClusterClient

//var rdb *redis.Ring

var Ctx = context.Background()

func init() {
	redisDb, err := beego.AppConfig.String("redisDb")
	redisPass, err := beego.AppConfig.String("redisPasswd")
	if err != nil {
		log.Printf("[WARN] Get redisDb config Fail:%s\n", err)
		return
	}
	if redisDb != "" {
		Rdb = redis.NewClient(&redis.Options{
			Addr:     redisDb,
			Password: redisPass, // no password set
			DB:       0,         // use default DB
		})
	}
	//defer Rdb.Close() //加上会提示error:redis: client is closed
	// rdb = redis.NewClusterClient(&redis.ClusterOptions{
	// 	Addrs: []string{":7000", ":7001", ":7002", ":7003", ":7004", ":7005"},
	// 	// To route commands by latency or randomly, enable one of the following.
	// 	//RouteByLatency: true,
	// 	//RouteRandomly: true,
	// })

	// rdb = redis.NewFailoverClient(&redis.FailoverOptions{
	// 	MasterName:    "master-name",
	// 	SentinelAddrs: []string{":9126", ":9127", ":9128"},
	// })

	// rdb = redis.NewRing(&redis.RingOptions{
	// 	Addrs: map[string]string{
	// 		"shard1": ":7000",
	// 		"shard2": ":7001",
	// 		"shard3": ":7002",
	// 	},
	// })

	// rdb is *redis.Client.
	//rdb = redis.NewUniversalClient(&redis.UniversalOptions{
	//	Addrs: []string{"192.168.10.171:6379"},
	//})

	// rdb is *redis.ClusterClient.
	//rdb := NewUniversalClient(&redis.UniversalOptions{
	//	Addrs: []string{":6379", ":6380"},
	//})

	// rdb is *redis.FailoverClient.
	//rdb := NewUniversalClient(&redis.UniversalOptions{
	//	Addrs:      []string{":6379"},
	//	MasterName: "mymaster",
	//})
}

func Get(key string) string {
	reply, err := Rdb.Get(Ctx, key).Result()
	if err == redis.Nil {
		log.Printf("[INFO] Get %s not exist\n", key)
		return ""
	} else if err != nil {
		log.Printf("[ERROR] Get %s error:%s\n", key, err)
		return ""
	} else {
		//log.Printf("[INFO] Get result:%s", reply)
		return reply
	}
}

func Ttl(key string) (int64, error) {
	reply, err := Rdb.TTL(Ctx, key).Result()
	return int64(reply / time.Second), err
}

func Persist(key string) (bool, error) {
	reply, err := Rdb.Persist(Ctx, key).Result()
	return reply, err
}

func Del(key string) int64 {
	reply, err := Rdb.Del(Ctx, key).Result()
	if err == redis.Nil {
		log.Printf("[INFO] Del %s not exist\n", key)
		return 0
	} else if err != nil {
		log.Printf("[ERROR] Del %s error:%s\n", key, err)
		return 0
	} else {
		//log.Printf("[INFO] Get result:%s", reply)
		return reply
	}
}

func Exists(key string) (int64, error) {
	num, err := Rdb.Exists(Ctx, key).Result()
	if err == redis.Nil {
		//log.Printf("[WARN] Exists %s not exist\n", key)
		return 0, err
	} else if err != nil {
		//log.Printf("[ERROR] Exists %s error:%s\n", key, err)
		return 0, err
	} else {
		return num, nil
	}
}

func Set(key, value string) error {
	err := Rdb.Set(Ctx, key, value, 0).Err()
	if err != nil {
		log.Printf("[ERROR] Set %s:%s error:%s\n", key, value, err)
		return err
	}
	log.Printf("[INFO] SET %s:%s ok", key, value)
	return nil
}

func IncrBy(key string, value int64) (int64, error) {
	num, err := Rdb.IncrBy(Ctx, key, value).Result()
	if err != nil {
		log.Printf("[ERROR] IncrBy %s:%d error:%s\n", key, value, err)
		return 0, err
	}
	log.Printf("[INFO] IncrBy %s:%d ok", key, value)
	return num, nil
}

func Keys(key string) []string {
	respArry, err := Rdb.Keys(Ctx, key).Result()
	if err != nil {
		log.Printf("[ERROR] Keys %s error:%s\n", key, err)
	}
	return respArry
}

func SetEx(key, value string, expire int64) error {
	err := Rdb.SetEx(Ctx, key, value, time.Duration(expire)*time.Second).Err()
	if err != nil {
		log.Printf("[ERROR] SetEx %s:%s,expire:%d error:%s\n", key, value, expire, err)
		return err
	}
	//log.Printf("[INFO] SetEx %s:%s,expire:%d ok", key, value, expire)
	return nil
}

func HSet(htable, key, value string) error {
	_, err := Rdb.HSet(Ctx, htable, key, value).Result()
	if err != nil {
		log.Printf("[ERROR] HSet %s,%s:%s error:%s\n", htable, key, value, err)
		return err
	}
	//log.Printf("[INFO] HSet %s,%s:%s success,result:%d\n", htable, key, value, resp)
	return nil
}

func HGet(htable, key string) (string, error) {
	resp, err := Rdb.HGet(Ctx, htable, key).Result()
	if err != nil {
		log.Printf("[ERROR] HGet %s:%s error:%s\n", htable, key, err)
		return "", err
	}
	if err == redis.Nil {
		log.Printf("[INFO] HGet %s:%s not exist\n", htable, key)
		return "", err
	}
	//log.Printf("[INFO] HGet %s:%s success,result:%s \n", htable, key, resp)
	return resp, nil
}

func Hexists(htable, key string) (bool, error) {
	resp, err := Rdb.HExists(Ctx, htable, key).Result()
	if err != nil {
		log.Printf("[ERROR] Hexists %s:%s error:%s\n", htable, key, err)
		return false, err
	}
	return resp, nil
}

func HDel(htable, key string) error {
	resp, err := Rdb.HDel(Ctx, htable, key).Result()
	if err != nil {
		log.Printf("[ERROR] HDel %s:%s error:%s\n", htable, key, err)
		return err
	}
	log.Printf("[INFO] HDel %s:%s success,result:%d\n", htable, key, resp)
	return nil
}

func HDels(htable string, key []string) error {
	resp, err := Rdb.HDel(Ctx, htable, key...).Result()
	if err != nil {
		log.Printf("[ERROR] HDels %s:%s error:%s\n", htable, key, err)
		return err
	}
	log.Printf("[INFO] HDels %s:%s success,result:%d\n", htable, key, resp)
	return nil
}

type kv struct {
	Key   string
	Value string
}

func HGetAll(htable string) []kv {
	hMap := Rdb.HGetAll(Ctx, htable).Val()
	kvArry := make([]kv, 0)
	for k, v := range hMap {
		kvArry = append(kvArry, *&kv{k, v})
	}
	//log.Printf("[INFO] HGetAll success,result:%v \n", kvArry)
	return kvArry
}

func HKeys(htable string) []string {
	respArry, err := Rdb.HKeys(Ctx, htable).Result()
	if err != nil {
		log.Printf("[ERROR] HKeys %s error:%s\n", htable, err)
	}
	return respArry
}

func HLen(htable string) (int64, error) {
	resp, err := Rdb.HLen(Ctx, htable).Result()
	if err != nil {
		log.Printf("[ERROR] HKeys %s error:%s\n", htable, err)
		return 0, err
	}
	return resp, nil
}

// 写redis hash表 多个值
func HMset(htable string, kv []interface{}) (bool, error) {
	resp, err := Rdb.HMSet(Ctx, htable, kv).Result()
	if err != nil {
		log.Printf("[ERROR] HMSET %s:%s error:%s\n", htable, kv, err)
		return false, err
	}
	log.Printf("[INFO] HMSET result:%v", resp)
	return true, nil
}

func SAdd(keys string, kv []interface{}) (int64, error) {
	resp, err := Rdb.SAdd(Ctx, keys, kv).Result()
	if err != nil {
		log.Printf("[ERROR] SAdd %s:%s error:%s\n", keys, kv, err)
		return 0, err
	}
	log.Printf("[INFO] SAdd result:%d", resp)
	return resp, nil
}

func Smembers(keys string) ([]string, error) {
	var tt = []string{}
	resp, err := Rdb.SMembers(Ctx, keys).Result()
	if err != nil {
		log.Printf("[ERROR] Smembers %serror:%s\n", keys, err)
		return tt, err
	}
	log.Printf("[INFO] Smembers result:%s", resp)
	return resp, nil
}

// func Subscribe(channel string) {
// 	pubsub := rdb.Subscribe(ctx, channel)
// 	defer pubsub.Close()
// }

// xadd single field
func XAdd(keys string, field, value string) (string, error) {
	resp, err := Rdb.XAdd(Ctx, &redis.XAddArgs{
		Stream:     keys,
		NoMkStream: false,
		MaxLen:     0,
		ID:         "",
		Values:     []string{field, value},
	}).Result()
	if err != nil {
		log.Printf("[ERROR] XAdd stream:%s %s:%s error:%s\n", keys, field, value, err)
		return "", err
	}
	//log.Printf("[INFO] XAdd stream:%s %s:%s result:%s", keys, field, value, resp)
	return resp, nil
}

func XAddMap(keys string, kv map[string]string) (string, error) {
	resp, err := Rdb.XAdd(Ctx, &redis.XAddArgs{
		Stream:     keys,
		NoMkStream: false,
		MaxLen:     0,
		ID:         "",
		Values:     kv,
	}).Result()
	if err != nil {
		log.Printf("[ERROR] XAddMap stream:%s %v error:%s\n", keys, kv, err)
		return "", err
	}
	//log.Printf("[INFO] XAddMap stream:%s %v result:%s", keys, kv, resp)
	return resp, nil
}

// 创建组，如果不存在stream 就创建stream
func XGroupCreateMkStream(keys, group string) (string, error) {
	resp, err := Rdb.XGroupCreateMkStream(Ctx, keys, group, "$").Result()
	if err != nil {
		log.Printf("[ERROR] XGroupCreateMkStream stream:%s group:%s error:%s\n", keys, group, err)
		return "", err
	}
	log.Printf("[INFO] XGroupCreateMkStream stream:%s group:%s result:%s", keys, group, resp)
	return resp, nil
}

// 存在stream的情况下创建组，否则报错
func XGroupCreate(keys, group string) (string, error) {
	resp, err := Rdb.XGroupCreate(Ctx, keys, group, "$").Result()
	if err != nil {
		log.Printf("[ERROR] XGroupCreate stream:%s group:%s error:%s\n", keys, group, err)
		return "", err
	}
	log.Printf("[INFO] XGroupCreate stream:%s group:%s result:%s", keys, group, resp)
	return resp, nil
}

func XAck(keys, group, msgId string) error {
	_, err := Rdb.XAck(Ctx, keys, group, msgId).Result()
	//resp, err := Rdb.XAck(Ctx, keys, group, []msgId).Result() //确认多个消息
	if err != nil {
		log.Printf("[ERROR] XAck stream:%s group:%s error:%s\n", keys, group, err)
		return err
	}
	//log.Printf("[INFO] XAck stream:%s group:%s result:%d", keys, group, resp)
	return nil
}

// https://pkg.go.dev/github.com/go-redis/redis/v8#XReadGroupArgs
// https://www.jianshu.com/p/c3cee95c3a34
func XReadGroup(keys, group, consumerId, startId string) ([]redis.XStream, error) {

	var resp = make([]redis.XStream, 0)
	var err error

	resp, err = Rdb.XReadGroup(Ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumerId,
		Streams:  []string{keys, startId}, //id 为>时 读取最新未被传递到其他消费者的消息，0是读取所有未被xack的消息,在这种情况下，BLOCK 和 NOACK 都被忽略。。//当streams这里是监听多个streams时，resp数组的长度就是streams的长度。
		Count:    5,                       //读取消息的数量
		Block:    100 * time.Second,       //0为无限时长阻塞侦听消息time.Duration,目前发现设置为0时不生效，只有10s中就断开，设置100000毫秒，不加time.*时只有一分钟。注意需要加上time.Second才会生效。
		NoAck:    false,                   //为true时，读取消息后就自动xack 了。此时不需要再执行autoclaim了
	}).Result()
	if err == redis.Nil {
		log.Printf("[WARN] TaskWorker queue noTask:%s\n", err)
		return resp, nil
	} else if err != nil {
		log.Printf("[ERROR] TaskWorker XReadGroup error:%s\n", err)
		return resp, err
	}
	// for i := 0; i < len(resp[0].Messages); i++ {
	// 	messageID := resp[0].Messages[i].ID
	// 	values := resp[0].Messages[i].Values
	// 	tskId := fmt.Sprintf("%v", values["taskId"])
	// 	if taskId != "" {
	// 		log.Printf("taskId:%s", taskId)
	// 		XAck(keys, group, messageID)
	// 	}
	// }
	return resp, nil
}

// 查看未ack的消息
func XPending(key, group string) (int64, error) {
	//var resp = make([]redis.XStream, 0)
	//var err error
	resp, err := Rdb.XPending(Ctx, key, group).Result()
	if err != nil {
		//log.Printf("[ERROR] XPending error:%s\n", err)
		return 0, err
	}
	return resp.Count, nil
}

// 将所有已分派，还没有xack的消息进行重新分配
// XAUTOCLAIM taskIdQueue taskIdConsumerGroup ccc3 100 0-0 COUNT 10
// 将超过所有100毫秒的未xack的10条消息 转移到ccc3 consumer 上去。
func XAutoClaim(key, group, newConsumer string) ([]string, error) {
	var taskId []string
	resp, _, err := Rdb.XAutoClaim(Ctx, &redis.XAutoClaimArgs{
		Stream:   key,
		Group:    group,
		MinIdle:  15 * time.Second,
		Start:    "0-0",
		Count:    5,
		Consumer: newConsumer,
	}).Result()
	if err != nil {
		//log.Printf("[ERROR] XAutoClaim error:%s\n", err)
		return taskId, err
	}
	if len(resp) > 0 {
		for i := 0; i < len(resp); i++ {
			//log.Printf("[INFO] XAutoClaim id:%s,taskid:%s", resp[i].ID, resp[i].Values["taskId"])
			taskId = append(taskId, resp[i].Values["taskId"].(string))
		}
	}
	return taskId, nil
}

// 添加锁
func AddLock(key, value string, expire int64) bool {
	resp, err := Rdb.SetNX(Ctx, key, value, time.Duration(expire)*time.Second).Result()
	if err != nil {
		log.Printf("[INFO] AddLock %s:%s,Fail:%s", key, value, err)
		return false
	}
	//log.Printf("[INFO] SET %s:%s ok", key, value)
	return resp
}

// 获取锁
func GetLock(key string) bool {
	resp, err := Rdb.Exists(Ctx, key).Result()
	if err != nil {
		return false
	}
	if resp > 0 {
		return true
	} else {
		return false
	}
}

// 删除锁
func UnLock(key string) bool {
	_, err := Rdb.Del(Ctx, key).Result() //返回结果1表示删除成功、0表示不存在key。所依只要不存在key就任务解锁成功
	if err != nil {
		return false
	}
	return true
}

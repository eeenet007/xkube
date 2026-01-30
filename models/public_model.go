// public_model.go
package models

import (
	"fmt"
	"xkube/common"
)

func ClearCache(clusterId string) int {
	//删除资源列表缓存
	delKeyArry := common.Keys("*List" + clusterId + "*")
	for _, vv := range delKeyArry {
		_ = common.Del(vv)
	}
	num := len(delKeyArry)

	//删除根据集群+节点统计的pod数量
	delKeyArry2 := common.Keys("podCount" + clusterId + "*")
	for _, vv2 := range delKeyArry2 {
		_ = common.Del(vv2)
	}
	num += len(delKeyArry2)

	//删除首页统计的数量
	_ = common.Del("count_" + clusterId)

	fmt.Printf("[INFO]%s:cache clear success:%d", clusterId, num)
	return num
}

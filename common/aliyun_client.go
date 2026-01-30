// aliyun_client.go
package common

import (
	//"fmt"
	//"log"
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	devops "github.com/alibabacloud-go/devops-20210625/v4/client"

	//util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

//var AccessKeyId = ""
//var AccessKeySecret = ""
var Endpint = "devops.cn-hangzhou.aliyuncs.com"

func GetAkbyAccountId(aliyun_id string) (string, string) {
	akstr, err := GetAliyunAk(aliyun_id)
	if err != nil {
		return "", ""
	}
	//log.Println(akstr)
	akArry := strings.Split(akstr, ",")
	return akArry[0], akArry[1]
}

func AliClient(aliyun_id string) (_result *devops.Client, _err error) {
	akId, akSecret := GetAkbyAccountId(aliyun_id)
	config := &openapi.Config{
		AccessKeyId:     &akId,
		AccessKeySecret: &akSecret,
	}
	// Endpoint 请参考 https://api.aliyun.com/product/devops
	config.Endpoint = tea.String(Endpint)
	_result = &devops.Client{}
	_result, _err = devops.NewClient(config)
	return _result, _err
}

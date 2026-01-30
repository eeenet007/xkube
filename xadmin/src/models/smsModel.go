// smsModel.go
package models

import (
	"fmt"
	//"lib"
	"log"
	common "xkube/common"
	lib "xkube/xadmin/src/lib"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	credential "github.com/aliyun/credentials-go/credentials"
)

// VerifyTelCode 验证手机号码验证码 / Verify mobile phone verification code
// username: 用户名 / username
// code: 验证码 / verification code
// bool: 验证结果，true表示验证成功，false表示验证失败 / Verification result, true means verification successful, false means verification failed
func VerifyTelCode(username, code string) bool {
	user, _ := GetUserByUsername(username)
	if user.Id == 0 {
		return false
	}
	if common.Get(user.Telphone) == code {
		return true
	}
	return false
}

// CreateClient 创建阿里云短信服务客户端 / Create Alibaba Cloud SMS service client
// _result: 阿里云短信服务客户端 / Alibaba Cloud SMS service client
// _err: 错误信息 / error message
func CreateClient() (_result *dysmsapi20170525.Client, _err error) {

	var AliyunAkId = lib.GetConfigString("AliyunAkId")
	var AliyunAkSecret = lib.GetConfigString("AliyunAkSecret")

	akConfig := new(credential.Config).
		SetType("access_key").
		SetAccessKeyId(AliyunAkId).
		SetAccessKeySecret(AliyunAkSecret)

	akCredential, err := credential.NewCredential(akConfig)
	if err != nil {
		return _result, err
	}

	config := &openapi.Config{
		Credential: akCredential,
	}

	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

// SendSmsV2 发送短信验证码 / Send SMS verification code
// username: 用户名 / username
// error: 错误信息 / error message
func SendSmsV2(username string) error {
	client, _err := CreateClient()
	if _err != nil {
		log.Printf("[ERROR] SendSmsV2 CreateClient error:%s\n", _err)
		return _err
	}

	//var err error
	user, _ := GetUserByUsername(username)
	if user.Id == 0 {
		return fmt.Errorf("NoMobilePhone")
	}

	//set In redis
	code := lib.GetRandomNumber(5)
	expireTime, _ := common.Ttl(user.Telphone)
	if expireTime > 0 {
		return fmt.Errorf("%d秒后重试", expireTime)
	}

	tmpParam := fmt.Sprintf(`{"code":"%s"}`, code)
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(user.Telphone),
		SignName:      tea.String(lib.GetConfigString("SignName")),
		TemplateCode:  tea.String(lib.GetConfigString("TemplateCode")),
		TemplateParam: tea.String(tmpParam),
	}

	runtime := &util.RuntimeOptions{}
	_, _err = client.SendSmsWithOptions(sendSmsRequest, runtime)
	if _err != nil {
		log.Printf("[ERROR] SendSmsV2 SendSmsWithOptions error:%s\n", _err)
		return _err
	}
	_ = common.SetEx(user.Telphone, code, 180)

	return nil
}
// cluster.go
package controllers

import (
	//"fmt"

	"log"
	//"strings"
	//"time"
	//"encoding/json"
	//"net/http"
	"strconv"
	"xkube/common"
	m "xkube/models"

	//beego "github.com/beego/beego/v2/server/web"
	"github.com/tidwall/gjson"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	devops "github.com/alibabacloud-go/devops-20210625/v4/client"
	aliutil "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

//xxx网络：5f20ed2
//xxx2020团队：5f6b
//此处查看https://devops.aliyun.com/organization/
//流水线ID：https://flow.aliyun.com/pipelines/【PipelineId】/current
//jobid：通过GetPipelineRun来获取https://help.aliyun.com/document_detail/460565.html

//type CicdController struct {
//	beego.Controller
//}

// 运行流水线
func (this *CicdController) Start() {
	aliyun_id := this.GetString("aliyun_id")
	organizationId := this.GetString("organization_id")
	pipelineId := this.GetString("pipeline_id")
	cicdId, _ := this.GetInt64("cicdId")
	if aliyun_id == "" || organizationId == "" || pipelineId == "" {
		datas, err := m.GetPipelinesByCicdId(cicdId)
		if err != nil {
			this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": err.Error()}
			this.ServeJSON()
			return
		}
		aliyun_id = datas.AliyunId
		organizationId = datas.OrganizationId
		pipelineId = datas.PipelineId
	}

	client, _err := common.AliClient(aliyun_id)
	if _err != nil {
		log.Printf("[ERROR] pipeline start err:%s\n", _err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": _err.Error()}
		this.ServeJSON()
		return
	}
	startPipelineRunRequest := &devops.StartPipelineRunRequest{}
	runtime := &aliutil.RuntimeOptions{}
	headers := make(map[string]*string)

	resp, _err := client.StartPipelineRunWithOptions(tea.String(organizationId), tea.String(pipelineId), startPipelineRunRequest, headers, runtime)
	if _err != nil {
		log.Printf("[ERROR] pipeline start err:%s\n", _err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": _err.Error()}
		this.ServeJSON()
		return
	}
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	this.Ctx.WriteString(resp.Body.String())
}

func ListPipelineRuns(aliyun_id, organizationId, pipelineId string, num int64) (*devops.ListPipelineRunsResponse, error) {
	var resp *devops.ListPipelineRunsResponse
	client, _err := common.AliClient(aliyun_id)
	if _err != nil {
		log.Printf("[ERROR] pipeline start err:%s\n", _err)
		return resp, _err
	}
	listPipelineRunsRequest := &devops.ListPipelineRunsRequest{
		MaxResults: tea.Int64(num),
	}
	runtime := &aliutil.RuntimeOptions{}
	headers := make(map[string]*string)
	resp2, err := client.ListPipelineRunsWithOptions(tea.String(organizationId), tea.String(pipelineId), listPipelineRunsRequest, headers, runtime)
	return resp2, err
}

// 获取流水线运行实例列表
func (this *CicdController) ListRun() {
	aliyun_id := this.GetString("aliyun_id")
	organizationId := this.GetString("organization_id")
	pipelineId := this.GetString("pipeline_id")
	num, _ := this.GetInt64("num")

	cicdId, _ := this.GetInt64("cicdId")
	if aliyun_id == "" || organizationId == "" || pipelineId == "" {
		datas, err := m.GetPipelinesByCicdId(cicdId)
		if err != nil {
			this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": err.Error()}
			this.ServeJSON()
			return
		}
		aliyun_id = datas.AliyunId
		organizationId = datas.OrganizationId
		pipelineId = datas.PipelineId
	}
	resp, err := ListPipelineRuns(aliyun_id, organizationId, pipelineId, num)
	if err != nil {
		log.Printf("[ERROR] pipeline start err2:%s\n", err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": err.Error()}
		this.ServeJSON()
		return
	}

	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	this.Ctx.WriteString(resp.Body.String())
}

// 终止流水线运行
func (this *CicdController) Stop() {
	aliyun_id := this.GetString("aliyun_id")
	organizationId := this.GetString("organization_id")
	pipelineId := this.GetString("pipeline_id")
	pipelineRunId := this.GetString("pipelineRunId")

	cicdId, _ := this.GetInt64("cicdId")
	if aliyun_id == "" || organizationId == "" || pipelineId == "" {
		datas, err := m.GetPipelinesByCicdId(cicdId)
		if err != nil {
			this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": err.Error()}
			this.ServeJSON()
			return
		}
		aliyun_id = datas.AliyunId
		organizationId = datas.OrganizationId
		pipelineId = datas.PipelineId
	}

	client, _err := common.AliClient(aliyun_id)
	if _err != nil {
		log.Printf("[ERROR] pipeline Stop err2:%s\n", _err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": _err.Error()}
		this.ServeJSON()
		return
	}

	runtime := &aliutil.RuntimeOptions{}
	headers := make(map[string]*string)

	resp, _err := client.StopPipelineRunWithOptions(tea.String(organizationId), tea.String(pipelineId), tea.String(pipelineRunId), headers, runtime)
	if _err != nil {
		log.Printf("[ERROR] pipeline Stop err2:%s\n", _err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": _err.Error()}
		this.ServeJSON()
		return
	}
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	this.Ctx.WriteString(resp.Body.String())
}

// 重试流水线的job
func (this *CicdController) Retry() {
	aliyun_id := this.GetString("aliyun_id")
	organizationId := this.GetString("organization_id")
	pipelineId := this.GetString("pipeline_id")
	pipelineRunId := this.GetString("pipelineRunId")
	jobId := this.GetString("jobId")

	cicdId, _ := this.GetInt64("cicdId")
	if aliyun_id == "" || organizationId == "" || pipelineId == "" {
		datas, err := m.GetPipelinesByCicdId(cicdId)
		if err != nil {
			this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": err.Error()}
			this.ServeJSON()
			return
		}
		aliyun_id = datas.AliyunId
		organizationId = datas.OrganizationId
		pipelineId = datas.PipelineId
	}

	client, _err := common.AliClient(aliyun_id)
	if _err != nil {
		log.Printf("[ERROR] pipeline Retry err:%s\n", _err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": _err.Error()}
		this.ServeJSON()
		return
	}

	runtime := &aliutil.RuntimeOptions{}
	headers := make(map[string]*string)
	resp, _err := client.RetryPipelineJobRunWithOptions(tea.String(organizationId), tea.String(pipelineId), tea.String(pipelineRunId), tea.String(jobId), headers, runtime)
	if _err != nil {
		log.Printf("[ERROR] pipeline Retry err2:%s\n", _err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": _err.Error()}
		this.ServeJSON()
		return
	}

	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	this.Ctx.WriteString(resp.Body.String())
}

// 停止job
func (this *CicdController) StopJob() {
	aliyun_id := this.GetString("aliyun_id")
	organizationId := this.GetString("organization_id")
	pipelineId := this.GetString("pipeline_id")
	pipelineRunId := this.GetString("pipelineRunId")
	jobId := this.GetString("jobId")

	cicdId, _ := this.GetInt64("cicdId")
	if aliyun_id == "" || organizationId == "" || pipelineId == "" {
		datas, err := m.GetPipelinesByCicdId(cicdId)
		if err != nil {
			this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": err.Error()}
			this.ServeJSON()
			return
		}
		aliyun_id = datas.AliyunId
		organizationId = datas.OrganizationId
		pipelineId = datas.PipelineId
	}

	client, _err := common.AliClient(aliyun_id)
	if _err != nil {
		log.Printf("[ERROR] pipeline Retry err:%s\n", _err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": _err.Error()}
		this.ServeJSON()
		return
	}

	runtime := &aliutil.RuntimeOptions{}
	headers := make(map[string]*string)
	resp, _err := client.StopPipelineJobRunWithOptions(tea.String(organizationId), tea.String(pipelineId), tea.String(pipelineRunId), tea.String(jobId), headers, runtime)
	if _err != nil {
		log.Printf("[ERROR] pipeline Retry err2:%s\n", _err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": _err.Error()}
		this.ServeJSON()
		return
	}

	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	this.Ctx.WriteString(resp.Body.String())
}

// 跳过job
func (this *CicdController) SkipJob() {
	aliyun_id := this.GetString("aliyun_id")
	organizationId := this.GetString("organization_id")
	pipelineId := this.GetString("pipeline_id")
	pipelineRunId := this.GetString("pipelineRunId")
	jobId := this.GetString("jobId")

	cicdId, _ := this.GetInt64("cicdId")
	if aliyun_id == "" || organizationId == "" || pipelineId == "" {
		datas, err := m.GetPipelinesByCicdId(cicdId)
		if err != nil {
			this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": err.Error()}
			this.ServeJSON()
			return
		}
		aliyun_id = datas.AliyunId
		organizationId = datas.OrganizationId
		pipelineId = datas.PipelineId
	}

	client, _err := common.AliClient(aliyun_id)
	if _err != nil {
		log.Printf("[ERROR] pipeline Retry err:%s\n", _err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": _err.Error()}
		this.ServeJSON()
		return
	}

	runtime := &aliutil.RuntimeOptions{}
	headers := make(map[string]*string)
	resp, _err := client.SkipPipelineJobRunWithOptions(tea.String(organizationId), tea.String(pipelineId), tea.String(pipelineRunId), tea.String(jobId), headers, runtime)
	if _err != nil {
		log.Printf("[ERROR] pipeline Retry err2:%s\n", _err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": _err.Error()}
		this.ServeJSON()
		return
	}

	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	this.Ctx.WriteString(resp.Body.String())
}

// 通过人工卡点
func (this *CicdController) PassValidate() {
	aliyun_id := this.GetString("aliyun_id")
	organizationId := this.GetString("organization_id")
	pipelineId := this.GetString("pipeline_id")
	pipelineRunId := this.GetString("pipelineRunId")
	jobId := this.GetString("jobId")

	cicdId, _ := this.GetInt64("cicdId")
	if aliyun_id == "" || organizationId == "" || pipelineId == "" {
		datas, err := m.GetPipelinesByCicdId(cicdId)
		if err != nil {
			this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": err.Error()}
			this.ServeJSON()
			return
		}
		aliyun_id = datas.AliyunId
		organizationId = datas.OrganizationId
		pipelineId = datas.PipelineId
	}

	client, _err := common.AliClient(aliyun_id)
	if _err != nil {
		log.Printf("[ERROR] pipeline PassValidate err:%s\n", _err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": _err.Error()}
		this.ServeJSON()
		return
	}

	runtime := &aliutil.RuntimeOptions{}
	headers := make(map[string]*string)

	resp, _err := client.PassPipelineValidateWithOptions(tea.String(organizationId), tea.String(pipelineId), tea.String(pipelineRunId), tea.String(jobId), headers, runtime)
	if _err != nil {
		log.Printf("[ERROR] pipeline PassValidate err2:%s\n", _err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": _err.Error()}
		this.ServeJSON()
		return
	}
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	this.Ctx.WriteString(resp.Body.String())
}

// 拒绝人工卡点
func (this *CicdController) RefuseValidate() {
	aliyun_id := this.GetString("aliyun_id")
	organizationId := this.GetString("organization_id")
	pipelineId := this.GetString("pipeline_id")
	pipelineRunId := this.GetString("pipelineRunId")
	jobId := this.GetString("jobId")

	cicdId, _ := this.GetInt64("cicdId")
	if aliyun_id == "" || organizationId == "" || pipelineId == "" {
		datas, err := m.GetPipelinesByCicdId(cicdId)
		if err != nil {
			this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": err.Error()}
			this.ServeJSON()
			return
		}
		aliyun_id = datas.AliyunId
		organizationId = datas.OrganizationId
		pipelineId = datas.PipelineId
	}

	client, _err := common.AliClient(aliyun_id)
	if _err != nil {
		log.Printf("[ERROR] pipeline RefuseValidate err:%s\n", _err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": _err.Error()}
		this.ServeJSON()
		return
	}

	runtime := &aliutil.RuntimeOptions{}
	headers := make(map[string]*string)

	resp, _err := client.RefusePipelineValidateWithOptions(tea.String(organizationId), tea.String(pipelineId), tea.String(pipelineRunId), tea.String(jobId), headers, runtime)
	if _err != nil {
		log.Printf("[ERROR] pipeline RefuseValidate err2:%s\n", _err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": _err.Error()}
		this.ServeJSON()
		return
	}
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	this.Ctx.WriteString(resp.Body.String())
}

// 获取运行实例
func (this *CicdController) GetRun() {
	aliyun_id := this.GetString("aliyun_id")
	organizationId := this.GetString("organization_id")
	pipelineId := this.GetString("pipeline_id")
	pipelineRunId := this.GetString("pipelineRunId")
	cicdId, _ := this.GetInt64("cicdId")
	if aliyun_id == "" || organizationId == "" || pipelineId == "" {
		datas, err := m.GetPipelinesByCicdId(cicdId)
		if err != nil {
			this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": err.Error()}
			this.ServeJSON()
			return
		}
		aliyun_id = datas.AliyunId
		organizationId = datas.OrganizationId
		pipelineId = datas.PipelineId
	}

	if pipelineRunId == "" {
		resp, err := ListPipelineRuns(aliyun_id, organizationId, pipelineId, 1)
		if err != nil {
			log.Printf("[ERROR] ListPipelineRuns err:%s\n", err)
			this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": err.Error()}
			this.ServeJSON()
			return
		}
		gp := gjson.Parse(resp.Body.String())
		if gp.Get("pipelineRuns").Exists() {
			pipelineRunId = strconv.FormatInt(gp.Get("pipelineRuns.0.pipelineRunId").Int(), 10)
		} else {
			log.Println("[ERROR] ListPipelineRuns Get Fail")
		}
	}
	if pipelineRunId == "" {
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": "pipelineRunId null"}
		this.ServeJSON()
		return
	}
	log.Println(aliyun_id, organizationId, pipelineId, pipelineRunId)
	client, _err := common.AliClient(aliyun_id)
	if _err != nil {
		log.Printf("[ERROR] pipeline GetRun err:%s\n", _err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": _err.Error()}
		this.ServeJSON()
		return
	}

	runtime := &aliutil.RuntimeOptions{}
	headers := make(map[string]*string)
	resp, _err := client.GetPipelineRunWithOptions(tea.String(organizationId), tea.String(pipelineId), tea.String(pipelineRunId), headers, runtime)
	if _err != nil {
		log.Printf("[ERROR] pipeline GetRun err2:%s\n", _err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": _err.Error()}
		this.ServeJSON()
		return
	}

	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	this.Ctx.WriteString(resp.Body.String())
}

// 获取运行实例的日志
func (this *CicdController) GetJobLog() {
	aliyun_id := this.GetString("aliyun_id")
	organizationId := this.GetString("organization_id")
	pipelineId := this.GetString("pipeline_id")
	jobId := this.GetString("jobId")
	pipelineRunId := this.GetString("pipelineRunId")

	cicdId, _ := this.GetInt64("cicdId")
	if aliyun_id == "" || organizationId == "" || pipelineId == "" {
		datas, err := m.GetPipelinesByCicdId(cicdId)
		if err != nil {
			this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": err.Error()}
			this.ServeJSON()
			return
		}
		aliyun_id = datas.AliyunId
		organizationId = datas.OrganizationId
		pipelineId = datas.PipelineId
	}

	client, _err := common.AliClient(aliyun_id)
	if _err != nil {
		log.Printf("[ERROR] pipeline GetJobLog err:%s\n", _err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": _err.Error()}
		this.ServeJSON()
		return
	}

	runtime := &aliutil.RuntimeOptions{}
	headers := make(map[string]*string)
	resp, _err := client.LogPipelineJobRunWithOptions(tea.String(organizationId), tea.String(pipelineId), tea.String(jobId), tea.String(pipelineRunId), headers, runtime)
	if _err != nil {
		log.Printf("[ERROR] pipeline GetJobLog err2:%s\n", _err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": _err.Error()}
		this.ServeJSON()
		return
	}
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	this.Ctx.WriteString(resp.Body.String())
}

func (this *CicdController) ListOrganizations() {
	aliyun_id := this.GetString("aliyun_id")
	client, _err := common.AliClient(aliyun_id)
	if _err != nil {
		log.Printf("[ERROR] ListOrganizations err:%s\n", _err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": _err.Error()}
		this.ServeJSON()
		return
	}

	runtime := &aliutil.RuntimeOptions{}
	headers := make(map[string]*string)
	resp, _err := client.ListJoinedOrganizationsWithOptions(headers, runtime)
	if _err != nil {
		log.Printf("[ERROR] ListOrganizations err2:%s\n", _err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": _err.Error()}
		this.ServeJSON()
		return
	}

	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	this.Ctx.WriteString(resp.Body.String())
}

func (this *CicdController) GetOrganizationsByAk() {
	gp := gjson.ParseBytes(this.Ctx.Input.RequestBody)
	akId := gp.Get("accesskey_id").String()
	akSecret := gp.Get("accesskey_secret").String()
	config := &openapi.Config{
		AccessKeyId:     &akId,
		AccessKeySecret: &akSecret,
	}
	config.Endpoint = tea.String(common.Endpint)
	//client = &devops.Client{}
	client, _err := devops.NewClient(config)
	if _err != nil {
		log.Printf("[ERROR] GetOrganizationsByAk err:%s\n", _err)
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": _err.Error()}
		this.ServeJSON()
		return
	}
	runtime := &aliutil.RuntimeOptions{}
	headers := make(map[string]*string)
	resp, _err2 := client.ListJoinedOrganizationsWithOptions(headers, runtime)
	if _err2 != nil {
		log.Printf("[ERROR] GetOrganizationsByAk err2:%s\n", _err2)
		this.Data["json"] = &map[string]interface{}{"code": -1, "success": false, "msg": _err2.Error()}
		this.ServeJSON()
		return
	}
	this.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	this.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json;charset=utf-8")
	this.Ctx.WriteString(resp.Body.String())
}

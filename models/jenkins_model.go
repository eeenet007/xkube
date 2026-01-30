// jenkins_models.go
package models

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	"github.com/bndr/gojenkins"
	_ "github.com/go-sql-driver/mysql"
)

type BuildInfo struct {
	Id        int64   `json:"id"`
	BuildTime string  `json:"buildTime"`
	Duration  float64 `json:"duration"`
	Status    string  `json:"status"`
	Log       string  `json:"log"`
}

var JksIdMap = make(map[string]Xkb_cicd_jks)

func NewClient(jksId string) (*gojenkins.Jenkins, error) {
	var jenkins *gojenkins.Jenkins
	ctx := context.Background()
	jks, err := GetJksConfig(jksId)
	if err != nil {
		log.Printf("GetJksConfig失败, %v\n", err)
		return jenkins, err
	}
	jenkins = gojenkins.CreateJenkins(nil, jks.JksUrl, jks.JksUser, jks.JksPasswd)
	_, err = jenkins.Init(ctx)
	if err != nil {
		log.Printf("连接Jenkins失败, %v\n", err)
		return jenkins, err
	}
	//log.Println("Jenkins连接成功")
	return jenkins, nil
}

func GetJksConfig(jksId string) (Xkb_cicd_jks, error) {
	o := orm.NewOrm()
	var jkscfg Xkb_cicd_jks
	if vv, ok := JksIdMap[jksId]; ok {
		return vv, nil
	} else {
		sqlstr := fmt.Sprintf("SELECT * FROM xkb_cicd_jks WHERE jks_id = '%s'", jksId)
		err := o.Raw(sqlstr).QueryRow(&jkscfg)
		if err != nil {
			return jkscfg, err
		}
		if jkscfg.JksId != "" {
			JksIdMap[jksId] = jkscfg
			return jkscfg, nil
		} else {
			return jkscfg, fmt.Errorf("emptyJksId")
		}
	}
}

func RunJobBuild(jksId, jobName string, m map[string]string) (int64, error) {
	Jks, err := NewClient(jksId)
	if err != nil {
		log.Printf("[ERROR] RunJobBuild NewClient Fail:%s\n", err)
		return 0, err
	}

	job, err := Jks.GetJob(context.Background(), jobName)
	if err != nil {
		log.Printf("[ERROR] GetJobInfo Fail:%s\n", err)
		return 0, err
	}
	lastBuild, _ := job.GetLastBuild(context.Background())
	buildId := lastBuild.GetBuildNumber()

	_, err2 := Jks.BuildJob(context.Background(), jobName, m)
	if err2 != nil {
		log.Printf("[ERROR] RunJobBuild Fail:%s\n", err2)
		return 0, err2
	}
	//当前执行的是在上一次构建的ID基础上+1
	return buildId + 1, err2
}

func GetBuildList(jksId, jobName string) ([]BuildInfo, error) {
	var buildArry = make([]BuildInfo, 0)
	Jks, err := NewClient(jksId)
	if err != nil {
		log.Printf("[ERROR] RunJobBuild NewClient Fail:%s\n", err)
		return buildArry, err
	}
	builds, err := Jks.GetAllBuildIds(context.Background(), jobName)
	if err != nil {
		log.Println(err)
		return buildArry, err
	}
	for _, vv := range builds {
		//log.Println(vv.Number, vv.URL)
		data, err := Jks.GetBuild(context.Background(), jobName, vv.Number)
		if err != nil {
			continue
		}
		buildArry = append(buildArry, *&BuildInfo{
			Id:        vv.Number,
			BuildTime: data.GetTimestamp().Format("2006-01-02T15:04:05"),
			Duration:  data.GetDuration(),
			Status:    data.GetResult(),
			Log:       "",
		})
		//log.Printf("ID:%d,时间:%s,耗时:%f,状态:%s\n", vv.Number, data.GetTimestamp().Format("2006-01-02 15:04:05"), data.GetDuration(), data.GetResult())
	}
	return buildArry, nil
}

func GetBuildLog(jksId, jobName string, buildId int64) (string, error) {
	//buildIdInt, _ := strconv.ParseInt(buildId, 10, 64)
	Jks, err := NewClient(jksId)
	if err != nil {
		log.Printf("[ERROR] RunJobBuild NewClient Fail:%s\n", err)
		return "", err
	}

	if buildId == 0 {
		job, err := Jks.GetJob(context.Background(), jobName)
		if err != nil {
			log.Printf("[ERROR] GetJobInfo Fail:%s\n", err)
			return "GetJobInfo Fail", err
		}
		lastBuild, _ := job.GetLastBuild(context.Background())
		buildId = lastBuild.GetBuildNumber()
	}

	build, err := Jks.GetBuild(context.Background(), jobName, buildId)
	if err != nil {
		log.Printf("[ERROR] GetBuildLog GetBuild Fail:%s", err)
		return "", err
	}
	//build.GetResult()
	logtext := build.GetConsoleOutput(context.Background())
	logtext = strings.Replace(logtext, `"`, `\"`, -1)
	return logtext, nil
}

func GetBuildState(jksId, jobName string, buildId int64) (BuildInfo, error) {
	var bi BuildInfo
	Jks, err := NewClient(jksId)
	if err != nil {
		log.Printf("[ERROR] RunJobBuild NewClient Fail:%s\n", err)
		return bi, err
	}

	if buildId == 0 {
		job, err := Jks.GetJob(context.Background(), jobName)
		if err != nil {
			log.Printf("[ERROR] GetJobInfo Fail:%s\n", err)
			return bi, err
		}
		lastBuild, _ := job.GetLastBuild(context.Background())
		buildId = lastBuild.GetBuildNumber()
	}

	data, err := Jks.GetBuild(context.Background(), jobName, buildId)
	if err != nil {
		log.Printf("[ERROR] GetBuildLog GetBuild Fail:%s", err)
		return bi, err
	}
	logtext := data.GetConsoleOutput(context.Background())
	return BuildInfo{
		Id:        buildId,
		BuildTime: data.GetTimestamp().Format("2006-01-02T15:04:05"),
		Duration:  data.GetDuration() / 1000,
		Status:    data.GetResult(),
		Log:       logtext,
	}, nil
}

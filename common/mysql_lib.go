// mysql_lib.go
package common

import (
	//"database/sql"
	"fmt"
	//"log"
	//"os"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/client/orm"
	//beego "github.com/beego/beego/v2/server/web"

	//. "github.com/beego/admin/src/lib"

	_ "github.com/go-sql-driver/mysql"
)

var ClusterMap = make(orm.Params)
var ClusterTokenMap = make(orm.Params)
var ClusterVersionMap = make(orm.Params)

var AliyunIdMap = make(map[string]string)

func GetKubeConfigByClusterId(clusterid string) (string, error) {
	o := orm.NewOrm()
	if clusterid == "undefined" || clusterid == "" {
		return "", fmt.Errorf("error clusterId")
	}
	if vv, ok := ClusterMap[clusterid]; ok {
		return vv.(string), nil
	} else {
		sqlstr := fmt.Sprintf("SELECT cluster_id,kube_config FROM xkb_cluster WHERE cluster_id = '%s'", clusterid)
		_, err := o.Raw(sqlstr).RowsToMap(&ClusterMap, "cluster_id", "kube_config")
		if err != nil {
			return "", err
		}
		kubecfg := ClusterMap[clusterid]
		if kubecfg == nil {
			return "", fmt.Errorf("errClusterId")
		}
		return kubecfg.(string), nil
	}
	// if appname != "" {
	// 	sqlstr = fmt.Sprintf("select cluster_id,kube_config FROM xkb_cluster", appname)
	// }
	// num, err2 := o.Raw(sqlstr).QueryRows(&xlr)
}

// 获取k8s的版本，来进行兼容api版本
func GetKubeVersionByClusterId(clusterid string) (float64, error) {
	o := orm.NewOrm()
	if vv, ok := ClusterVersionMap[clusterid]; ok {
		vArry := strings.Split(vv.(string), ".")
		if len(vArry) >= 2 {
			versionStr := vArry[0] + "." + vArry[1]
			version, err := strconv.ParseFloat(versionStr, 64)
			if err != nil {
				return 0.0, err
			}
			return version, nil
		}
		return 0, nil
	} else {
		sqlstr := fmt.Sprintf("SELECT cluster_id,kube_version FROM xkb_cluster WHERE cluster_id = '%s'", clusterid)
		_, err := o.Raw(sqlstr).RowsToMap(&ClusterVersionMap, "cluster_id", "kube_version")
		if err != nil {
			return 0, err
		}
		kubecfg := ClusterVersionMap[clusterid]
		vArry := strings.Split(kubecfg.(string), ".")
		if len(vArry) >= 2 {
			versionStr := vArry[0] + "." + vArry[1]
			version, err := strconv.ParseFloat(versionStr, 64)
			if err != nil {
				return 0, err
			}
			return version, nil
		}
		return 0, fmt.Errorf("k8s version format error")
	}
}

//获取k8s的BearerToken，用于k8s进行二次认证这种情况：例如aws的eks
//aws获取token:aws eks get-token --cluster-name pcauto-internation-eks| jq -r .status.token
//其他集群 获得方式，未测试成功
//# TOKEN=$(kubectl describe secret $(kubectl get secrets | grep default | cut -f1 -d ' ') | grep -E '^token' | cut -f2 -d':' | tr -d '\t')
//TOKEN=$(kubectl describe secret $(kubectl get secrets | grep default | cut -f1 -d ' ') | grep -E '^token' | cut -f2 -d':' | tr -d ' ')
//参考：https://m.bufeishi.cn/99201.html
//https://blog.csdn.net/chrboy/article/details/122540792
//https://www.jianshu.com/p/d23c5ebd6695
//免刷新token生成方法
//通过以下方式能生成一个不受影响的Kube config：
//$ kubectl -n kube-system create serviceaccount kubeconfig-sa
//$ kubectl create clusterrolebinding add-on-cluster-admin --clusterrole=cluster-admin --serviceaccount=kube-system:kubeconfig-sa
//$ TOKENNAME=`kubectl -n kube-system get serviceaccount/kubeconfig-sa -o jsonpath='{.secrets[0].name}'`
//$ TOKEN=`kubectl -n kube-system get secret $TOKENNAME -o jsonpath='{.data.token}'| base64 --decode`
//将生成的TOKEN填充到kube config中的${TOKEN}即可实现免刷新。
//https://aws.amazon.com/cn/blogs/china/talk-about-aws-eks-identity-authentication-processing-in-simple-terms/

func GetBearerTokenByClusterId(clusterid string) (string, error) {
	o := orm.NewOrm()
	if vv, ok := ClusterTokenMap[clusterid]; ok {
		if vv == nil {
			return "", nil
		}
		return vv.(string), nil
	} else {
		sqlstr := fmt.Sprintf("SELECT cluster_id,bearer_token FROM xkb_cluster WHERE cluster_id = '%s'", clusterid)
		_, err := o.Raw(sqlstr).RowsToMap(&ClusterTokenMap, "cluster_id", "bearer_token")
		if err != nil {
			return "", err
		}
		bearerToken := ClusterTokenMap[clusterid]
		if bearerToken == nil {
			return "", nil
		}
		return bearerToken.(string), nil
	}
}

type AliAk struct {
	AliyunId        string
	AccesskeyId     string
	AccesskeySecret string
}

func GetAliyunAk(aliyunId string) (string, error) {
	o := orm.NewOrm()
	if vv, ok := AliyunIdMap[aliyunId]; ok {
		return vv, nil
	} else {
		var aliak AliAk
		sqlstr := fmt.Sprintf("SELECT aliyun_id,accesskey_id,accesskey_secret FROM xkb_cicd_ak WHERE aliyun_id = '%s'", aliyunId)
		err := o.Raw(sqlstr).QueryRow(&aliak)
		if err != nil {
			return "", err
		}
		//fmt.Println(aliak)
		if aliak.AccesskeyId != "" && aliak.AccesskeySecret != "" {
			AliyunIdMap[aliak.AliyunId] = aliak.AccesskeyId + "," + aliak.AccesskeySecret
			return aliak.AccesskeyId + "," + aliak.AccesskeySecret, nil
		} else {
			return "", fmt.Errorf("emptyAk")
		}
	}
}

// terminalws.go
package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"xkube/common"
	m "xkube/models"
	adm "xkube/xadmin/src/lib"

	beegolog "github.com/beego/beego/v2/core/logs"
	//beego "github.com/beego/beego/v2/server/web"
	"log"
	"strings"

	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"k8s.io/api/core/v1"

	//"k8s.io/apimachinery/pkg/runtime"
	//"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"

	//"k8s.io/client-go/rest"
	//"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

var sessionName = adm.GetConfigString("SessionName")

func (self TerminalSockjs) Read(p []byte) (int, error) {
	var reply string
	var msg map[string]uint16
	reply, err := self.conn.Recv()
	if err != nil {
		return 0, err
	}
	if err := json.Unmarshal([]byte(reply), &msg); err != nil {
		return copy(p, reply), nil
	} else {
		self.sizeChan <- &remotecommand.TerminalSize{
			msg["cols"],
			msg["rows"],
		}
		return 0, nil
	}
}

func (self TerminalSockjs) Write(p []byte) (int, error) {
	err := self.conn.Send(string(p))
	return len(p), err
}

type TerminalSockjs struct {
	conn      sockjs.Session
	sizeChan  chan *remotecommand.TerminalSize
	context   string
	clusterId string
	namespace string
	pod       string
	container string
}

// 实现tty size queue
func (self *TerminalSockjs) Next() *remotecommand.TerminalSize {
	size := <-self.sizeChan
	//beegolog.Debug(fmt.Sprintf("terminal size to width: %d height: %d", size.Width, size.Height))
	return size
}

// 处理输入输出与sockjs 交互
func Handler(t *TerminalSockjs, cmd string) error {
	//restclient, config := common.RestClient(t.clusterId)
	clientset, config := common.ClientSetConfig(t.clusterId)
	req := clientset.CoreV1().RESTClient().Post().
		//req := restclient.Post().
		Resource("pods").
		Name(t.pod).
		Namespace(t.namespace).
		SubResource("exec").
		Param("container", t.container).
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("command", cmd).Param("tty", "true")
	req.VersionedParams(
		&v1.PodExecOptions{
			Container: t.container,
			Command:   []string{},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		},
		scheme.ParameterCodec,
	)
	executor, err := remotecommand.NewSPDYExecutor(
		config, http.MethodPost, req.URL(),
	)
	if err != nil {
		return err
	}
	return executor.Stream(remotecommand.StreamOptions{
		Stdin:             t,
		Stdout:            t,
		Stderr:            t,
		Tty:               true,
		TerminalSizeQueue: t,
	})
}

// 实现http.handler 接口获取入参
func (self TerminalSockjs) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// cookie, err := r.Cookie(sessionName)
	// if err != nil || cookie.Value == "" {
	// 	log.Printf("[ERROR] TerminalSockjs no session cookie found\n")
	// 	http.Error(w, "Authentication failed - no session", http.StatusUnauthorized)
	// 	return
	// }
	// sessionData := common.Get(cookie.Value)
	// if len(sessionData) < 100 {
	// 	log.Printf("[ERROR] TerminalSockjs session key not found in Redis: %s\n", cookie.Value)
	// 	http.Error(w, "Authentication failed - invalid session", http.StatusUnauthorized)
	// 	return
	// }

	context := r.FormValue("context")
	clusterId := r.FormValue("clusterId")
	namespace := r.FormValue("nameSpace")
	pod := r.FormValue("podName")
	container := r.FormValue("container")

	containerList, err := m.PodContainerList(clusterId, namespace, pod)
	if err != nil {
		log.Printf("[ERROR] TerminalSockjs PodContainerList error:%s\n", err)
		fmt.Fprintf(w, "clusterId,namespace,container Null")
		return
	}
	var containerStr string
	for _, vv := range containerList {
		containerStr += fmt.Sprintf("%s,", vv.ContainerName)
	}

	if !strings.Contains(containerStr, container) {
		log.Printf("[ERROR] TerminalSockjs ServeHTTP noThisContainer")
		fmt.Fprintf(w, "noThisContainer")
		return
	}

	Sockjshandler := func(session sockjs.Session) {
		t := &TerminalSockjs{session, make(chan *remotecommand.TerminalSize),
			context, clusterId, namespace, pod, container}
		if err := Handler(t, "/bin/sh"); err != nil {
			//beegolog.Error(err)
			log.Printf("[ERROR] TerminalSockjs error:%s\n", err)
			beegolog.Error(Handler(t, "/bin/bash"))
		}
	}

	sockjs.NewHandler("/terminal/ws", sockjs.DefaultOptions, Sockjshandler).ServeHTTP(w, r)
}

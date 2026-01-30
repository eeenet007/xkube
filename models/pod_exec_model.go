// pod_exec_model.go
package models

import (
	"net/http"
	//"strings"
	"bytes"
	"context"
	"fmt"
	"time"
	"xkube/common"

	corev1 "k8s.io/api/core/v1"

	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

func PodExec(clusterId, nameSpace, podName, containerName string, command string, timeout int64) (string, string, error) {
	clientset, config := common.ClientSetConfig(clusterId)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	req := clientset.CoreV1().RESTClient().Post().
		//req := restclient.Post().
		Resource("pods").
		Name(podName).
		Namespace(nameSpace).
		SubResource("exec").
		VersionedParams(
			&corev1.PodExecOptions{
				Container: containerName,
				Command:   []string{"sh", "-c", command},
				Stdout:    true,
				Stderr:    true,
			},
			scheme.ParameterCodec,
		)
	executor, err := remotecommand.NewSPDYExecutor(
		config, http.MethodPost, req.URL(),
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to create executor: %v", err)
	}
	var stdout, stderr bytes.Buffer
	err = executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if ctx.Err() == context.DeadlineExceeded {
		return "", "", fmt.Errorf("command execution timeout after %v", timeout)
	}
	return stdout.String(), stderr.String(), err
}

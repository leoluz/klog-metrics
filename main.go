package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/ericchiang/k8s"
	corev1 "github.com/ericchiang/k8s/apis/core/v1"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func main() {
	fmt.Println(".: Kube Log Metrics")
	LoadConfiguration()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(strconv.Itoa(GetServerPort()), nil))
}

func publishPodsLogs() {
	client, _ := k8s.NewInClusterClient()

	var pods corev1.PodList
	err := client.List(context.Background(), client.Namespace, &pods)
	if err != nil {
		fmt.Printf("Error listing pods: %v", err)
		os.Exit(1)
	}

	for _, pod := range pods.GetItems() {
		podName := pod.GetMetadata().GetName()
		if value, ok := pod.GetMetadata().GetLabels()["ad-app"]; ok {
			fmt.Printf("\nPod: %s ad-app: %s\n", podName, value)
		} else {
			fmt.Printf("\nPod: %s no ad-app\n", podName)
		}
		for _, container := range pod.GetSpec().GetContainers() {
			handleLog(client, podName, container.GetName())
		}
	}
}

func handleLog(client *k8s.Client, pod, container string) error {
	url := getLogResourceUrl(client.Endpoint, client.Namespace, pod, container)
	var body io.Reader
	req, err := http.NewRequest("GET", url, body)
	if err != nil {
		return fmt.Errorf("Error creating the request: %v\n", err)
	}
	if client.SetHeaders != nil {
		client.SetHeaders(req.Header)
	}
	resp, err := client.Client.Do(req)
	if err != nil {
		fmt.Errorf("Request error: %v\n", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Read body error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\tContainer: %s, Log size: %d\n", container, len(respBody))
	return nil
}

func getLogResourceUrl(baseurl, namespace, pod, container string) string {
	return fmt.Sprintf("%s/api/v1/namespaces/%s/pods/%s/log?container=%s", baseurl, namespace, pod, container)
}

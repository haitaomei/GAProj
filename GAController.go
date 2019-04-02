package main

import (
	"fmt"
	"k8s.io/client-go/rest"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var kubeClient *kubernetes.Clientset

const nameSpace = "default"

func main() {
	fmt.Println("Hi")

	//init kubernetes client
	var err error
	kubeClient, err = initInClusterKubeClient()
	if err != nil {
		panic("Can not create kubernetes client...")
	}

}

func initInClusterKubeClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println("Failed to create InClusterConfig...")
	}
	client, err := kubernetes.NewForConfig(config)
	return client, err
}

func deployIsland(islandName string, id string) {
	pr := islandPodRec(islandName, id)
	pod, err := kubeClient.CoreV1().Pods(nameSpace).Create(pr)
	if err != nil {
		fmt.Println("Failed to create Pod in Kubernetes", err)
	} else {
		fmt.Println("Created Pod", islandName, "in Kubernetes cluster. PodName =", pod.ObjectMeta.Name)
	}
}

func deleteIsland(id string) {
	pods, err := kubeClient.CoreV1().Pods(nameSpace).List(metav1.ListOptions{
		LabelSelector: "islandid=" + id,
	})
	if err != nil {
		fmt.Println("Failed to list all the pods in Kubernetes", err)
		return
	}

	for _, pod := range pods.Items {
		podName := pod.ObjectMeta.Name
		islandName := pod.ObjectMeta.Labels["islandName"]

		var delBackGround = metav1.DeletePropagationBackground
		err = kubeClient.CoreV1().Pods(nameSpace).Delete(podName, &metav1.DeleteOptions{
			PropagationPolicy: &delBackGround,
		})

		if err != nil {
			fmt.Println("Failed to delete pod with id =", id, "name =", islandName)
		} else {
			fmt.Println("Successfully deleted pod with id =", id, "name =", islandName)
		}
	}
}

func islandPodRec(podName string, id string) *corev1.Pod {

	dockerRegistry := "zs1517"
	imageName := "gacontroller"
	if img, ok := os.LookupEnv("GA_CONTROLLER_IMAGE_NAME"); ok {
		imageName = img
	}

	pr := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: nameSpace,
			Labels: map[string]string{
				"islandName": podName,
				"creator":    "ODGA",
				"islandid":   id,
			},
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyOnFailure,
			Volumes: []corev1.Volume{
				corev1.Volume{
					Name: "island-storage",
					//VolumeSource: corev1.VolumeSource{
					//	PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					//		ClaimName: "PVC-Name",
					//	},
					//},
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
			},
			Containers: []corev1.Container{
				{
					Name:            podName,
					Image:           dockerRegistry + "/" + imageName,
					ImagePullPolicy: "Always",
					VolumeMounts: []corev1.VolumeMount{
						corev1.VolumeMount{
							Name:      "island-storage",
							MountPath: "/islandsdatastore",
						},
					},
					Env: []corev1.EnvVar{
						{
							Name:  "ISLAND_ID",
							Value: id,
						},
					},
				},
			},
		},
	}

	return pr
}

package main

import (
	"fmt"
	ocsv1 "github.com/openshift/ocs-operator/pkg/apis/ocs/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	codecs := serializer.NewCodecFactory(scheme.Scheme)
	parameterCodec := runtime.NewParameterCodec(scheme.Scheme)

	//replace with the required kubeconfig path
	ocsConfig, err := clientcmd.BuildConfigFromFlags("", "/home/k7/kvellalo-spike1/auth/kubeconfig")
	if err != nil {
		fmt.Println(err)
	}
	ocsConfig.GroupVersion = &ocsv1.SchemeGroupVersion
	ocsConfig.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: codecs}
	ocsConfig.APIPath = "/apis"
	ocsConfig.ContentType = runtime.ContentTypeJSON
	if ocsConfig.UserAgent == "" {
		ocsConfig.UserAgent = rest.DefaultKubernetesUserAgent()
	}
	metav1.AddToGroupVersion(scheme.Scheme, *ocsConfig.GroupVersion)

	ocsClient, err := rest.RESTClientFor(ocsConfig)
	if err != nil {
		fmt.Println(err)
	}

	sc := &ocsv1.StorageCluster{}
	err = ocsClient.Get().
		Resource("storageclusters").
		Namespace("openshift-storage").
		Name("ocs-storagecluster").
		VersionedParams(&metav1.GetOptions{}, parameterCodec).
		Do().
		Into(sc)
	if err != nil {
		fmt.Println(err)
	}

	//Prints Storage cluster creation date
	fmt.Println(sc.CreationTimestamp.Date())
}

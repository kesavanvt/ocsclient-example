package main

import (
	"fmt"
	ocsv1 "github.com/openshift/ocs-operator/pkg/apis/ocs/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

func main() {
	//adding objects to the scheme
	ocsv1.SchemeBuilder.Register(&ocsv1.StorageCluster{})
	ocsScheme, err := ocsv1.SchemeBuilder.Build()

	codecs := serializer.NewCodecFactory(ocsScheme)
	parameterCodec := runtime.NewParameterCodec(ocsScheme)

	//replace with the required kubeconfig path
	ocsConfig, err := clientcmd.BuildConfigFromFlags("", "/home/k7/kvellalo-drain5/auth/kubeconfig")
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


	//creating watchlist
	sc = &ocsv1.StorageCluster{}
	watchlist := cache.NewListWatchFromClient(ocsClient, "storageclusters", "openshift-storage", fields.Everything())
	_, controller := cache.NewInformer(watchlist, sc, time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				fmt.Println("StorageCluster/pod added")
			},
			DeleteFunc: func(obj interface{}) {
				fmt.Println("StorageCluster/pod deleted")
			},
			UpdateFunc: nil,
		},
	)
	stop := make(chan struct{})
	go controller.Run(stop)
	for {
		time.Sleep(time.Second)
	}
}

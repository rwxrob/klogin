/*
Package kubeconf handles interactions with Kubernetes configuration as encapsulated by the kubectl config command which is the only appropriate way of handing the many possibilities of combined Kubernetes configuration files. This is also most reliable than using the client-go library as evidenced by the Kubernetes project locking specific third-party package dependencies in that library. Therefore, many of the values are taken from kubectl config view -o json and marshaled appropriately ensuring maximum longevity and compatibility (which using client-go unfortunately does not provide).
*/
package kubeconf

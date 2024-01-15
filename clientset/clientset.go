package clientset

import (
	"fmt"
	"github.com/martin-helmich/kubernetes-crd-example/clientset/v1alpha1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/flowcontrol"
	"net/http"
)

// clientset represents a {groupversion}.{resources}.{action} look at line 52 in main.go file
type Clientset struct {
	*discovery.DiscoveryClient
	customResource *v1alpha1.CustomResourceV1Client
}

// represents an api group
type Interface interface {
	ProjectsCrdV1Alpha1() v1alpha1.CustomResourceV1ClientInterface
}

// this represents the api group
func (c *Clientset) ProjectsCrdV1Alpha1() *v1alpha1.CustomResourceV1Client {
	return c.customResource
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {

	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

func NewForConfigWithHTTP(c *rest.Config) (*Clientset, error) {
	config := *c

	httpClient, err := rest.HTTPClientFor(&config)
	if err != nil {
		return nil, err
	}
	return NewForConfigAndClient(&config, httpClient)
}

func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c

	if configShallowCopy.UserAgent == "" {
		configShallowCopy.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	// share the transport between all clients
	httpClient, err := rest.HTTPClientFor(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	return NewForConfigAndClient(&configShallowCopy, httpClient)
}

// Note the http client provided takes precedence over the configured transport values.
// If config's RateLimiter is not set and QPS and Burst are acceptable,
// NewForConfigAndClient will generate a rate-limiter in configShallowCopy.
func NewForConfigAndClient(c *rest.Config, httpClient *http.Client) (*Clientset, error) {

	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		if configShallowCopy.Burst <= 0 {
			return nil, fmt.Errorf("burst is required to be greater than 0 when RateLimiter is not set and QPS is set to greater than 0")
		}
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}

	var cs Clientset
	var err error

	cs.customResource, err = v1alpha1.NewForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfigAndClient(&configShallowCopy, httpClient)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

func NewForConfigOrDie(c *rest.Config) *Clientset {
	cs, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return cs
}

func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.customResource = v1alpha1.New(c)
	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}

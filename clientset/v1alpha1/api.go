package v1alpha1

import (
	"github.com/martin-helmich/kubernetes-crd-example/api/types/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"net/http"
)

type CustomResourceV1ClientInterface interface {
	RESTClient() rest.Interface
	ProjectsGetter
}

type CustomResourceV1Client struct {
	restClient rest.Interface
}

// CustomResourceV1Client implements CustomResourceV1ClientInterface.ProjectsGetter interface. Projects represent the resource
func (c *CustomResourceV1Client) Projects(namespace string) ProjectsInterface {
	return newProjects(c, namespace)
}

// CustomResourceV1Client implements CustomResourceV1Interface interface
func (c *CustomResourceV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.RESTClient()
}

func NewForConfig(c *rest.Config) (*CustomResourceV1Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha1.GroupName, Version: v1alpha1.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &CustomResourceV1Client{restClient: client}, nil
}

func NewForConfigWithHTTP(c *rest.Config) (*CustomResourceV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	httpClient, err := rest.HTTPClientFor(&config)
	if err != nil {
		return nil, err
	}
	return NewForConfigAndClient(&config, httpClient)
}

// NewForConfigAndClient creates a new AppsV1Client for the given config and http client.
// Note the http client provided takes precedence over the configured transport values.
func NewForConfigAndClient(c *rest.Config, h *http.Client) (*CustomResourceV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientForConfigAndClient(&config, h)
	if err != nil {
		return nil, err
	}
	return &CustomResourceV1Client{client}, nil
}

// NewForConfigOrDie creates a new AppsV1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *CustomResourceV1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

func setConfigDefaults(config *rest.Config) error {

	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha1.GroupName, Version: v1alpha1.GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}
	return nil
}

func New(c rest.Interface) *CustomResourceV1Client {
	return &CustomResourceV1Client{c}
}

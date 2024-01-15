package v1alpha1

import (
	"context"
	"github.com/martin-helmich/kubernetes-crd-example/api/types/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ProjectsGetter interface {
	Projects(namespace string) ProjectsInterface
}

type ProjectsInterface interface {
	List(opts metav1.ListOptions) (*v1alpha1.ProjectList, error)
	Get(name string, options metav1.GetOptions) (*v1alpha1.Project, error)
	Create(*v1alpha1.Project) (*v1alpha1.Project, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

// projects implements the ProjectInterface
type projects struct {
	client rest.Interface
	ns     string
}

func newProjects(c *CustomResourceV1Client, namespace string) *projects {
	return &projects{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// implement the functions of ProjectInterface interface

func (p *projects) List(opts metav1.ListOptions) (*v1alpha1.ProjectList, error) {
	result := v1alpha1.ProjectList{}
	err := p.client.
		Get().
		Namespace(p.ns).
		Resource("projects").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (p *projects) Get(name string, opts metav1.GetOptions) (*v1alpha1.Project, error) {
	result := v1alpha1.Project{}
	err := p.client.
		Get().
		Namespace(p.ns).
		Resource("projects").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (p *projects) Create(project *v1alpha1.Project) (*v1alpha1.Project, error) {
	result := v1alpha1.Project{}
	err := p.client.
		Post().
		Namespace(p.ns).
		Resource("projects").
		Body(project).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (p *projects) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return p.client.
		Get().
		Namespace(p.ns).
		Resource("projects").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}

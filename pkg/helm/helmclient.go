package helm

import (
	"fmt"
	"github.com/weapons97/helm-api/pkg/helm/adapter"
	protos "github.com/weapons97/helm-api/pkg/protos"
	"github.com/weapons97/helm-api/pkg/utils"
	"time"
)

const HelmClientInContextName = `helm_context`

type ContextRequest interface {
	GetContextName() string
}
type TemporaryContextRequest interface {
	TemporaryRepoRequest
	TemporaryKubeRequest
}
type TemporaryRepoRequest interface {
	GetRepoinfo() *protos.RepoInfo
}
type TemporaryKubeRequest interface {
	GetKubeinfo() *protos.KubeInfo
}

type KubeNamespaceRequest interface {
	GetNamespace() string
}

type WithHelmContextOpt func(hc *HelmContext)

func WithHelmContextNamespace(namespace string) WithHelmContextOpt {
	return func(hc *HelmContext) {
		if namespace != "" {
			hc.KubeNamespace = namespace
		}
	}
}
func WithHelmContextNamespaceOpts(r interface{}) WithHelmContextOpt {
	switch req := r.(type) {
	case KubeNamespaceRequest:
		return WithHelmContextNamespace(req.GetNamespace())
	}
	return func(hc *HelmContext) {}
}

func NewHelmClient(r interface{}, opts ...WithHelmContextOpt) (client adapter.Adapter, e error) {

	ContextName := ""
	switch req := r.(type) {
	case ContextRequest:
		ContextName = req.GetContextName()
	default:
		client, e = NewHelmClientTemporary(r, opts...)
		if e == NoTemporaryContextInfo{
			return nil, NoContextInfo
		}
		if e != nil {
			return nil, e
		}
		return
	}
	hc, e := LoadHelmContext(ContextName)
	if e != nil {
		return
	}
	for i := range opts {
		if opts[i] == nil {
			continue
		}
		opts[i](hc)
	}
	client, e = hc.NewHelmClient()
	return
}

var NoTemporaryContextInfo = fmt.Errorf(`NoTemporaryContextInfo`)
var NoContextInfo = fmt.Errorf(`NoContextInfo`)

type TemporaryContext struct {
	Name   string
	Cancel func()
}

func NewHelmClientTemporary(r interface{}, opts ...WithHelmContextOpt) (client adapter.Adapter, e error) {
	var kubeinfo *protos.KubeInfo
	var repoinfo *protos.RepoInfo

	switch req := r.(type) {
	case TemporaryContextRequest:
		kubeinfo = req.GetKubeinfo()
		repoinfo = req.GetRepoinfo()
	case TemporaryKubeRequest:
		kubeinfo = req.GetKubeinfo()
	case TemporaryRepoRequest:
		repoinfo = req.GetRepoinfo()
	}
	if kubeinfo == nil && repoinfo == nil {
		return nil, NoTemporaryContextInfo
	}
	entrys := repoinfo.GetEntrys()
	rentrys := make([]adapter.Entry, len(entrys))
	for i, pe := range entrys {
		entry := adapter.ConvertEntry(pe)
		rentrys[i] = entry
	}
	hc, e := NewHelmContext(utils.RandomString(32), false,
		[]byte(kubeinfo.GetKubeconfig()), kubeinfo.GetContext(), kubeinfo.GetNamespace(),
		time.Now(),
		rentrys...,
	)
	if e != nil {
		return nil, e
	}
	for i := range opts {
		if opts[i] == nil {
			continue
		}
		opts[i](hc)
	}
	return hc.NewHelmClient()
}

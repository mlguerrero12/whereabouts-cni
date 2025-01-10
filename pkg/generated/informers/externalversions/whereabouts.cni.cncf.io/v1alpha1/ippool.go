/*
Copyright 2025 The Kubernetes Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	whereaboutscnicncfiov1alpha1 "github.com/k8snetworkplumbingwg/whereabouts/pkg/api/whereabouts.cni.cncf.io/v1alpha1"
	versioned "github.com/k8snetworkplumbingwg/whereabouts/pkg/generated/clientset/versioned"
	internalinterfaces "github.com/k8snetworkplumbingwg/whereabouts/pkg/generated/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/k8snetworkplumbingwg/whereabouts/pkg/generated/listers/whereabouts.cni.cncf.io/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// IPPoolInformer provides access to a shared informer and lister for
// IPPools.
type IPPoolInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.IPPoolLister
}

type iPPoolInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewIPPoolInformer constructs a new informer for IPPool type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewIPPoolInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredIPPoolInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredIPPoolInformer constructs a new informer for IPPool type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredIPPoolInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.WhereaboutsV1alpha1().IPPools(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.WhereaboutsV1alpha1().IPPools(namespace).Watch(context.TODO(), options)
			},
		},
		&whereaboutscnicncfiov1alpha1.IPPool{},
		resyncPeriod,
		indexers,
	)
}

func (f *iPPoolInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredIPPoolInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *iPPoolInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&whereaboutscnicncfiov1alpha1.IPPool{}, f.defaultInformer)
}

func (f *iPPoolInformer) Lister() v1alpha1.IPPoolLister {
	return v1alpha1.NewIPPoolLister(f.Informer().GetIndexer())
}

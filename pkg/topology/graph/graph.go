package graph

import (
	"log/slog"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/rancher/opni/pkg/logger"
	"github.com/rancher/opni/pkg/util"
	"github.com/steveteuber/kubectl-graph/pkg/graph"
	kgraph "github.com/steveteuber/kubectl-graph/pkg/graph"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	// Import to initialize client auth plugins
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
)

func NewRuntimeFactory() cmdutil.Factory {
	return cmdutil.NewFactory(&genericclioptions.ConfigFlags{})
}

// panics if run outside of a kubernetes environment
func TraverseTopology(lg *slog.Logger, f cmdutil.Factory) (*kgraph.Graph, error) {
	clientSet, err := kubernetes.NewForConfig(util.Must(rest.InClusterConfig()))
	if err != nil {
		return nil, err
	}
	objs := []*unstructured.Unstructured{}
	r := f.NewBuilder().
		Unstructured().
		NamespaceParam("").DefaultNamespace().AllNamespaces(true).
		//FIXME: reimplement these at a later date
		// FilenameParam(o.ExplicitNamespace, &o.FilenameOptions).
		// LabelSelectorParam(o.LabelSelector).
		// FieldSelectorParam(o.FieldSelector).
		// RequestChunksOf(o.ChunkSize).
		ResourceTypeOrNameArgs(true, "all").
		ContinueOnError().
		Latest().
		Flatten().
		Do()

	if err := r.Err(); err != nil {
		lg.Error("hit an error in the kubernetes runtime", logger.Err(err))
		return nil, err
	}

	infos, err := r.Infos() // doesn't use error types
	if err != nil {
		lg.Warn("hit an error while collecting kubernetes topology", logger.Err(err))
	}
	if len(infos) == 0 && err != nil { // should only exit in this case
		return nil, err
	}

	for _, info := range infos {
		objs = append(objs, info.Object.(*unstructured.Unstructured))
	}

	kubegraph, err := graph.NewGraph(clientSet, objs, func() {})
	if err != nil {
		lg.Error("error", logger.Err(err))
		return nil, err
	}
	return kubegraph, nil
}

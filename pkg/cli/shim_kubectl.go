package cli

import (
	"os"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	kclientcmd "k8s.io/client-go/tools/clientcmd"
	kclientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	kcmdset "k8s.io/kubectl/pkg/cmd/set"
	kcmdcreate "k8s.io/kubectl/pkg/cmd/create"
	describeversioned "k8s.io/kubectl/pkg/describe"
	"k8s.io/kubectl/pkg/generate/versioned"
	"k8s.io/kubectl/pkg/polymorphichelpers"

	"github.com/openshift/oc/pkg/helpers/clientcmd"
	"github.com/openshift/oc/pkg/helpers/originpolymorphichelpers"
)

func shimKubectlForOc() {
	// we only need this change for `oc`.  `kubectl` should behave as close to `kubectl` as we can
	// if we call this factory construction method, we want the openshift style config loading
	kclientcmd.ErrEmptyConfig = genericclioptions.ErrEmptyConfig
	kclientcmd.ClusterDefaults = kclientcmdapi.Cluster{Server: os.Getenv("KUBERNETES_MASTER")}
	kcmdset.ParseDockerImageReferenceToStringFunc = clientcmd.ParseDockerImageReferenceToStringFunc
	kcmdcreate.AddSpecialVerb("use", schema.GroupResource{
		Group:    "security.openshift.io",
		Resource: "securitycontextconstraints",
	})
	kclientcmd.UseModifyConfigLock = false

	// update polymorphic helpers
	polymorphichelpers.AttachablePodForObjectFn = originpolymorphichelpers.NewAttachablePodForObjectFn(polymorphichelpers.AttachablePodForObjectFn)
	polymorphichelpers.CanBeExposedFn = originpolymorphichelpers.NewCanBeExposedFn(polymorphichelpers.CanBeExposedFn)
	polymorphichelpers.HistoryViewerFn = originpolymorphichelpers.NewHistoryViewerFn(polymorphichelpers.HistoryViewerFn)
	polymorphichelpers.LogsForObjectFn = originpolymorphichelpers.NewLogsForObjectFn(polymorphichelpers.LogsForObjectFn)
	polymorphichelpers.MapBasedSelectorForObjectFn = originpolymorphichelpers.NewMapBasedSelectorForObjectFn(polymorphichelpers.MapBasedSelectorForObjectFn)
	polymorphichelpers.ObjectPauserFn = originpolymorphichelpers.NewObjectPauserFn(polymorphichelpers.ObjectPauserFn)
	polymorphichelpers.ObjectResumerFn = originpolymorphichelpers.NewObjectResumerFn(polymorphichelpers.ObjectResumerFn)
	polymorphichelpers.PortsForObjectFn = originpolymorphichelpers.NewPortsForObjectFn(polymorphichelpers.PortsForObjectFn)
	polymorphichelpers.ProtocolsForObjectFn = originpolymorphichelpers.NewProtocolsForObjectFn(polymorphichelpers.ProtocolsForObjectFn)
	polymorphichelpers.RollbackerFn = originpolymorphichelpers.NewRollbackerFn(polymorphichelpers.RollbackerFn)
	polymorphichelpers.StatusViewerFn = originpolymorphichelpers.NewStatusViewerFn(polymorphichelpers.StatusViewerFn)
	polymorphichelpers.UpdatePodSpecForObjectFn = originpolymorphichelpers.NewUpdatePodSpecForObjectFn(polymorphichelpers.UpdatePodSpecForObjectFn)

	// update some functions we inject
	versioned.GeneratorFn = originpolymorphichelpers.NewGeneratorsFn(versioned.GeneratorFn)
	describeversioned.DescriberFn = originpolymorphichelpers.NewDescriberFn(describeversioned.DescriberFn)
}

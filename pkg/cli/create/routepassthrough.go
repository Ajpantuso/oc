package create

import (
	"context"

	"github.com/spf13/cobra"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	kcmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/scheme"
	"k8s.io/kubectl/pkg/util"
	"k8s.io/kubectl/pkg/util/templates"

	routev1 "github.com/openshift/api/route/v1"
	"github.com/openshift/oc/pkg/cli/create/route"
)

var (
	passthroughRouteLong = templates.LongDesc(`
		Create a route that uses passthrough TLS termination.

		Specify the service (either just its name or using type/name syntax) that the
		generated route should expose via the --service flag.
	`)

	passthroughRouteExample = templates.Examples(`
		# Create a passthrough route named "my-route" that exposes the frontend service
		oc create route passthrough my-route --service=frontend

		# Create a passthrough route that exposes the frontend service and specify
		# a host name. If the route name is omitted, the service name will be used
		oc create route passthrough --service=frontend --hostname=www.example.com
	`)
)

type CreatePassthroughRouteOptions struct {
	CreateRouteSubcommandOptions *CreateRouteSubcommandOptions

	Hostname       string
	Port           string
	InsecurePolicy string
	Service        string
	WildcardPolicy string
}

// NewCmdCreatePassthroughRoute is a macro command to create a passthrough route.
func NewCmdCreatePassthroughRoute(f kcmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	o := &CreatePassthroughRouteOptions{
		CreateRouteSubcommandOptions: NewCreateRouteSubcommandOptions(streams),
	}
	cmd := &cobra.Command{
		Use:     "passthrough [NAME] --service=SERVICE",
		Short:   "Create a route that uses passthrough TLS termination",
		Long:    passthroughRouteLong,
		Example: passthroughRouteExample,
		Run: func(cmd *cobra.Command, args []string) {
			kcmdutil.CheckErr(o.Complete(f, cmd, args))
			kcmdutil.CheckErr(o.Run())
		},
	}
	cmd.Flags().StringVar(&o.Hostname, "hostname", o.Hostname, "Set a hostname for the new route")
	cmd.Flags().StringVar(&o.Port, "port", o.Port, "Name of the service port or number of the container port the route will route traffic to")
	cmd.Flags().StringVar(&o.InsecurePolicy, "insecure-policy", o.InsecurePolicy, "Set an insecure policy for the new route")
	cmd.Flags().StringVar(&o.Service, "service", o.Service, "Name of the service that the new route is exposing")
	cmd.MarkFlagRequired("service")
	cmd.Flags().StringVar(&o.WildcardPolicy, "wildcard-policy", o.WildcardPolicy, "Sets the WilcardPolicy for the hostname, the default is \"None\". valid values are \"None\" and \"Subdomain\"")

	kcmdutil.AddValidateFlags(cmd)
	o.CreateRouteSubcommandOptions.AddFlags(cmd)
	kcmdutil.AddDryRunFlag(cmd)

	return cmd
}

func (o *CreatePassthroughRouteOptions) Complete(f kcmdutil.Factory, cmd *cobra.Command, args []string) error {
	return o.CreateRouteSubcommandOptions.Complete(f, cmd, args)
}

func (o *CreatePassthroughRouteOptions) Run() error {
	serviceName, err := resolveServiceName(o.CreateRouteSubcommandOptions.Mapper, o.Service)
	if err != nil {
		return err
	}
	route, err := route.UnsecuredRoute(o.CreateRouteSubcommandOptions.CoreClient, o.CreateRouteSubcommandOptions.Namespace, o.CreateRouteSubcommandOptions.Name, serviceName, o.Port, false, o.CreateRouteSubcommandOptions.EnforceNamespace)
	if err != nil {
		return err
	}

	if len(o.WildcardPolicy) > 0 {
		route.Spec.WildcardPolicy = routev1.WildcardPolicyType(o.WildcardPolicy)
	}

	route.Spec.Host = o.Hostname
	route.Spec.TLS = new(routev1.TLSConfig)
	route.Spec.TLS.Termination = routev1.TLSTerminationPassthrough

	if len(o.InsecurePolicy) > 0 {
		route.Spec.TLS.InsecureEdgeTerminationPolicy = routev1.InsecureEdgeTerminationPolicyType(o.InsecurePolicy)
	}

	if err := util.CreateOrUpdateAnnotation(o.CreateRouteSubcommandOptions.CreateAnnotation, route, scheme.DefaultJSONEncoder()); err != nil {
		return err
	}

	if o.CreateRouteSubcommandOptions.DryRunStrategy != kcmdutil.DryRunClient {
		route, err = o.CreateRouteSubcommandOptions.Client.Routes(o.CreateRouteSubcommandOptions.Namespace).Create(context.TODO(), route, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}

	return o.CreateRouteSubcommandOptions.Printer.PrintObj(route, o.CreateRouteSubcommandOptions.Out)
}

/*
SPDX-FileCopyrightText: 2026 SAP SE or an SAP affiliate company and redis-operator-cop contributors
SPDX-License-Identifier: Apache-2.0
*/

package operator

import (
	"embed"
	"flag"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/sap/component-operator-runtime/pkg/component"
	helmgenerator "github.com/sap/component-operator-runtime/pkg/manifests/helm"
	"github.com/sap/component-operator-runtime/pkg/operator"

	operatorv1alpha1 "github.com/sap/redis-operator-cop/api/v1alpha1"
	"github.com/sap/redis-operator-cop/internal/transformer"
)

const Name = "redis-operator-cop.cs.sap.com"

//go:embed all:data
var data embed.FS

type Options struct {
	Name                  string
	DefaultServiceAccount string
	FlagPrefix            string
}

type Operator struct {
	options Options
}

var defaultOperator operator.Operator = New()

func GetName() string {
	return defaultOperator.GetName()
}

func InitScheme(scheme *runtime.Scheme) {
	defaultOperator.InitScheme(scheme)
}

func InitFlags(flagset *flag.FlagSet) {
	defaultOperator.InitFlags(flagset)
}

func ValidateFlags() error {
	return defaultOperator.ValidateFlags()
}

func GetUncacheableTypes() []client.Object {
	return defaultOperator.GetUncacheableTypes()
}

func Setup(mgr ctrl.Manager) error {
	return defaultOperator.Setup(mgr)
}

func New() *Operator {
	return NewWithOptions(Options{})
}

func NewWithOptions(options Options) *Operator {
	operator := &Operator{options: options}
	if operator.options.Name == "" {
		operator.options.Name = Name
	}
	return operator
}

func (o *Operator) GetName() string {
	return o.options.Name
}

func (o *Operator) InitScheme(scheme *runtime.Scheme) {
	utilruntime.Must(operatorv1alpha1.AddToScheme(scheme))
}

func (o *Operator) InitFlags(flagset *flag.FlagSet) {
	flagset.StringVar(&o.options.DefaultServiceAccount, "default-service-account", o.options.DefaultServiceAccount, "Default service account name")
}

func (o *Operator) ValidateFlags() error {
	return nil
}

func (o *Operator) GetUncacheableTypes() []client.Object {
	return []client.Object{&operatorv1alpha1.RedisOperator{}}
}

func (o *Operator) Setup(mgr ctrl.Manager) error {
	resourceGenerator, err := helmgenerator.NewHelmGeneratorWithParameterTransformer(
		data,
		"data/charts/redis-operator",
		mgr.GetClient(),
		transformer.NewParameterTransformer(),
	)
	if err != nil {
		return errors.Wrap(err, "error initializing resource generator")
	}

	if err := component.NewReconciler[*operatorv1alpha1.RedisOperator](
		o.options.Name,
		resourceGenerator,
		component.ReconcilerOptions{
			DefaultServiceAccount: &o.options.DefaultServiceAccount,
		},
	).SetupWithManager(mgr); err != nil {
		return errors.Wrapf(err, "unable to create controller")
	}

	return nil
}

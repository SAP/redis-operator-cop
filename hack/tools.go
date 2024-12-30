//go:build tools
// +build tools

/*
SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and redis-operator-cop contributors
SPDX-License-Identifier: Apache-2.0
*/

package tools

import (
	_ "k8s.io/code-generator"
	_ "sigs.k8s.io/controller-runtime/tools/setup-envtest"
	_ "sigs.k8s.io/controller-tools/cmd/controller-gen"
)

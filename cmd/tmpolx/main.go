/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2020 Red Hat, Inc.
 */

package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"

	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"

	"github.com/fromanirh/tmpolx/pkg/tmpolx"
)

func main() {
	var numaNodes string
	var policyName string
	flag.StringVarP(&numaNodes, "numa", "N", "0-7", "set NUMA configuration")
	flag.StringVarP(&policyName, "policy", "P", "none", "set Topology manager Policy")
	flag.Parse()

	numaConf, err := cpuset.Parse(numaNodes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bad format for NUMA configuration: %v", err)
		os.Exit(1)
	}

	params := tmpolx.Params{
		PolicyName: policyName,
		NUMANodes:  numaConf.ToSlice(),
		RawHints:   flag.Args(),
	}

	tmpx, err := tmpolx.NewFromParams(params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating TMPolx object: %v", err)
		os.Exit(2)
	}

	fmt.Fprintf(os.Stderr, "%s\n", tmpx.String())

	bestHint, admit := tmpx.Run()
	fmt.Printf("admit=%v hint=%v\n", admit, bestHint)
}
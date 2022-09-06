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

package tmpolx

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"k8s.io/kubernetes/pkg/kubelet/cm/topologymanager"

	"github.com/fromanirh/cpumgrx/pkg/tmutils"
)

const (
	MaxNUMANodes = 8 // TODO keep in sync with TM sources
)

type Params struct {
	PolicyName   string
	NUMANodes    []int
	RawHints     []string
	UseJSONHints bool
}

type TMPolx struct {
	policy topologymanager.Policy
	hints  map[string][]topologymanager.TopologyHint
}

func (tmpx *TMPolx) GetPolicyName() string {
	return tmpx.policy.Name()
}

func (tmpx *TMPolx) GetHints(resName string) []topologymanager.TopologyHint {
	var ret []topologymanager.TopologyHint
	for _, hint := range tmpx.hints[resName] {
		ret = append(ret, hint)
	}
	return ret
}

func (tmpx *TMPolx) String() string {
	var buf strings.Builder
	tw := tabwriter.NewWriter(&buf, 0, 8, 0, '\t', tabwriter.AlignRight)
	fmt.Fprintf(tw, ".\tresource\thints\t\n")
	for res, hints := range tmpx.hints {
		fmt.Fprintf(tw, ".\t%s\t%v\t\n", res, hints)
	}
	tw.Flush()
	return fmt.Sprintf("using policy %q\n%s", tmpx.policy.Name(), buf.String())
}

/*
> From: https://github.com/kubernetes/kubernetes/issues/84597#issuecomment-548414942

The restricted policy operates by limiting preferred alignments to the minimum possible alignment for the given request size on the given machine.

For your machine, this means that:

    Request sizes <= 6 will be restricted to a single NUMA node.
    Request sizes 7-12 will be restricted to 2 NUMA nodes.
    Request sizes 12-18 will be restricted to 3 NUMA nodes.

For your exact example, since there exists a way to allocate 3 CPUs from a single NUMA node on your machine (e.g. when no other pods are running), then requests of size 3 are restricted to single NUMA alignment for all pods.

This differs from the single-numa-node policy in that, no matter what the machine configuration looks like you must have alignment on a single NUMA node in order for the pod to be admitted. In your setup, this would mean that requests of sizes 7-18 would never have a path to admission.

The semantics you seem to be expecting are part of the best-effort policy, which will attempt to align on as few NUMA nodes as possible, only spilling over to another one if necessary.

---

> From: https://kubernetes.slack.com/archives/C0BP8PW9G/p1661761032814389?thread_ts=1661680145.406899&cid=C0BP8PW9G

The three policies are:
single-numa-node: only allow allocations from a single NUMA node, fail otherwise. Even if one of the requested resource requires more than one NUMA node to be satisfied.
restricted: only allow allocations from the minimum number of NUMA nodes. Look at each resource request, see what the minimum number of NUMA nodes are required to satisfy that resource request. Allow alignment to that number of NUMA nodes for all resources. Fail otherwise.
best-effort: Run as restricted, but never fail the allocation. Fall back to allocating from any remaining NUMA nodes as necessary. (edited)
*/

func NewFromParams(params Params) (*TMPolx, error) {
	if len(params.NUMANodes) > MaxNUMANodes {
		return nil, fmt.Errorf("TM currently supports up to %d NUMA nodes (got %d)", MaxNUMANodes, len(params.NUMANodes))
	}

	var policy topologymanager.Policy
	switch params.PolicyName {

	case topologymanager.PolicyNone:
		policy = topologymanager.NewNonePolicy()

	case topologymanager.PolicyBestEffort:
		policy = topologymanager.NewBestEffortPolicy(params.NUMANodes)

	case topologymanager.PolicyRestricted:
		policy = topologymanager.NewRestrictedPolicy(params.NUMANodes)

	case topologymanager.PolicySingleNumaNode:
		policy = topologymanager.NewSingleNumaNodePolicy(params.NUMANodes)

	default:
		return nil, fmt.Errorf("unknown policy: %q", params.PolicyName)
	}

	var err error
	hints := make(map[string][]topologymanager.TopologyHint)
	if params.UseJSONHints {
		hints, err = tmutils.ParseJSONHints(params.RawHints)
	} else {
		hints, err = tmutils.ParseGOHints(params.RawHints)
	}

	if err != nil {
		return nil, err
	}

	tmpx := &TMPolx{
		hints:  hints,
		policy: policy,
	}
	return tmpx, nil
}

func (tmpx *TMPolx) Run() (string, bool) {
	allHints := []map[string][]topologymanager.TopologyHint{tmpx.hints}
	bestHint, admit := tmpx.policy.Merge(allHints)
	return fmt.Sprintf("%v", bestHint), admit
}

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

// keep in this file the code which mimics k8s' TM

import (
	"fmt"

	"k8s.io/kubernetes/pkg/kubelet/cm/topologymanager"
)

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

	tmpx := &TMPolx{
		policy: policy,
		hints:  make(map[string][]topologymanager.TopologyHint),
	}
	if params.UseJSONHints {
		if err := tmpx.ParseJSONHints(params.RawHints); err != nil {
			return nil, err
		}
	} else {
		if err := tmpx.ParseGOHints(params.RawHints); err != nil {
			return nil, err
		}
	}
	return tmpx, nil
}

func (tmpx *TMPolx) Run() (string, bool) {
	allHints := []map[string][]topologymanager.TopologyHint{tmpx.hints}
	bestHint, admit := tmpx.policy.Merge(allHints)
	return fmt.Sprintf("%v", bestHint), admit
}

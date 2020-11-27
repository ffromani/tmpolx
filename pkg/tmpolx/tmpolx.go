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

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
	"encoding/json"
	"strings"

	"k8s.io/kubernetes/pkg/kubelet/cm/topologymanager"
	"k8s.io/kubernetes/pkg/kubelet/cm/topologymanager/bitmask"
)

type ResHints struct {
	Resource string `json:"R"`
	Hints    []Hint `json:"H"`
}

type Hint struct {
	Mask      string `json:"M"`
	Preferred bool   `json:"P"`
}

func (ht Hint) ToTM() topologymanager.TopologyHint {
	tmht := topologymanager.TopologyHint{
		Preferred:        ht.Preferred,
		NUMANodeAffinity: bitmask.NewEmptyBitMask(),
	}
	num := len(ht.Mask)
	for idx := 0; idx < num; idx++ {
		if ht.Mask[idx] == '1' {
			tmht.NUMANodeAffinity.Add(num - 1 - idx)
		}
	}
	return tmht
}

func (tmpx *TMPolx) addHint(rh ResHints) {
	hints := tmpx.hints[rh.Resource]
	for _, ht := range rh.Hints {
		hints = append(hints, ht.ToTM())
	}
	tmpx.hints[rh.Resource] = hints
}

func (tmpx *TMPolx) ParseJSONHints(rawHints []string) error {
	var err error
	for _, rawHint := range rawHints {
		var rh ResHints
		err = json.Unmarshal([]byte(rawHint), &rh)
		if err != nil {
			return err
		}

		tmpx.addHint(rh)
	}
	return nil
}

// cpu:[{01 true} {10 true} {11 false}]
func (tmpx *TMPolx) ParseGOHints(rawHints []string) error {
	for _, rawHint := range rawHints {
		data := strings.SplitN(rawHint, ":", 2)
		rh := ResHints{
			Resource: strings.TrimSpace(data[0]),
		}
		items := strings.Split(data[1], "} ")
		for _, item := range items {
			hintData := strings.SplitN(strings.Trim(item, "{}"), " ", 2)
			rh.Hints = append(rh.Hints, Hint{
				Mask:      hintData[0],
				Preferred: hintData[1] == "true",
			})
		}

		tmpx.addHint(rh)
	}
	return nil
}

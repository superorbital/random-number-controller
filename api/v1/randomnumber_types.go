/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Entropy defines how random the number should be
type Entropy int

const (
	// Low not very random
	Low Entropy = 0
	// Medium a little bit random
	Medium Entropy = 1
	// High the most random
	High Entropy = 2
)

// RandomNumberSpec defines the desired state of RandomNumber
type RandomNumberSpec struct {
	Entropy `json:"entropy,omitempty"`
}

// RandomNumberStatus defines the observed state of RandomNumber
type RandomNumberStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RandomNumber is the Schema for the randomnumbers API
type RandomNumber struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RandomNumberSpec   `json:"spec,omitempty"`
	Status RandomNumberStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RandomNumberList contains a list of RandomNumber
type RandomNumberList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RandomNumber `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RandomNumber{}, &RandomNumberList{})
}

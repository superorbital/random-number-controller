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

package controllers

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	randomv1 "github.com/superorbital/random-number-controller/api/v1"
)

var _ = Describe("RandomNumber controller", func() {
	Context("When creating a random number", func() {
		It("Should create a configmap with a random value", func() {
			By("By choosing entropy value")
			ctx := context.Background()
			rn := &randomv1.RandomNumber{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-random",
					Namespace: "default",
				},
				Spec: randomv1.RandomNumberSpec{
					Entropy: randomv1.High,
				},
			}
			Expect(k8sClient.Create(ctx, rn)).Should(Succeed())
			By("Verifying the resultant configmap exists")
		})
	})
})

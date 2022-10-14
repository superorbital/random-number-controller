# random-number-controller

This project is a sample Kubebuilder controller that demonstrates a few different
testing approaches.

The controller itself is straightforward, it defines a CRD that outlines how
much "entropy" you'd like in your random number and then generates a ConfigMap
with a random number in it that is guaranteed to be random because it 
was the result of a fair die roll.

## Testing Strategies

There are three types of tests in here:

1. Tests using a fake client from the Kubernetes controller-runtime project
2. Tests using an API server spun up by `envtest`
3. End to end tests run against a kind or colima cluster provisioned by a script in the repo

The goal of this repository is to help outline the drawbacks/benefits of each test type 
and enable you to write more tests against your own custom controllers.

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

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


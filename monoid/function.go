// Copyright (c) 2023 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package monoid

import (
	F "github.com/IBM/fp-go/function"
	S "github.com/IBM/fp-go/semigroup"
)

// FunctionMonoid forms a monoid as long as you can provide a monoid for the codomain.
func FunctionMonoid[A, B any](M Monoid[B]) Monoid[func(A) B] {
	return MakeMonoid(
		S.FunctionSemigroup[A, B](M).Concat,
		F.Constant1[A](M.Empty()),
	)
}

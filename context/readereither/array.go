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

package readereither

import (
	RE "github.com/IBM/fp-go/readereither/generic"
)

// TraverseArray transforms an array
func TraverseArray[A, B any](f func(A) ReaderEither[B]) func([]A) ReaderEither[[]B] {
	return RE.TraverseArray[ReaderEither[B], ReaderEither[[]B], []A](f)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[A any](ma []ReaderEither[A]) ReaderEither[[]A] {
	return RE.SequenceArray[ReaderEither[A], ReaderEither[[]A]](ma)
}

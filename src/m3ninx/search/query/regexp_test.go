// Copyright (c) 2018 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package query

import (
	"testing"

	"github.com/m3db/m3ninx/index"
	"github.com/m3db/m3ninx/search"

	"github.com/stretchr/testify/require"
)

func TestRegexpQuery(t *testing.T) {
	tests := []struct {
		name          string
		field, regexp []byte
		expectErr     bool
	}{
		{
			name:      "valid field and regexp should not return an error",
			field:     []byte("fruit"),
			regexp:    []byte(".*ple"),
			expectErr: false,
		},
		{
			name:      "invalid regexp should return an error",
			field:     []byte("fruit"),
			regexp:    []byte("(*]ple"),
			expectErr: true,
		},
	}

	rs := index.Readers{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			q, err := NewRegexpQuery(test.field, test.regexp)

			if test.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			_, err = q.Searcher(rs)
			require.NoError(t, err)
		})
	}
}

func TestRegexpQueryEqual(t *testing.T) {
	tests := []struct {
		name        string
		left, right search.Query
		expected    bool
	}{
		{
			name:     "same field and regexp",
			left:     newRegexpQuery(t, []byte("fruit"), []byte(".*ple")),
			right:    newRegexpQuery(t, []byte("fruit"), []byte(".*ple")),
			expected: true,
		},
		{
			name: "singular conjunction query",
			left: newRegexpQuery(t, []byte("fruit"), []byte(".*ple")),
			right: NewConjuctionQuery([]search.Query{
				newRegexpQuery(t, []byte("fruit"), []byte(".*ple")),
			}),
			expected: true,
		},
		{
			name: "singular disjunction query",
			left: newRegexpQuery(t, []byte("fruit"), []byte(".*ple")),
			right: NewDisjuctionQuery([]search.Query{
				newRegexpQuery(t, []byte("fruit"), []byte(".*ple")),
			}),
			expected: true,
		},
		{
			name:     "different field",
			left:     newRegexpQuery(t, []byte("fruit"), []byte(".*ple")),
			right:    newRegexpQuery(t, []byte("food"), []byte(".*ple")),
			expected: false,
		},
		{
			name:     "different regexp",
			left:     newRegexpQuery(t, []byte("fruit"), []byte(".*ple")),
			right:    newRegexpQuery(t, []byte("fruit"), []byte(".*na")),
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expected, test.left.Equal(test.right))
		})
	}
}

func newRegexpQuery(t *testing.T, field, regexp []byte) search.Query {
	q, err := NewRegexpQuery(field, regexp)
	require.NoError(t, err)
	return q
}

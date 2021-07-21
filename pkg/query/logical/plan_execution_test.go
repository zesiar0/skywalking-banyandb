// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package logical_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	apiv1 "github.com/apache/skywalking-banyandb/api/fbs/v1"
	"github.com/apache/skywalking-banyandb/banyand/series"
	"github.com/apache/skywalking-banyandb/pkg/posting/roaring"
	"github.com/apache/skywalking-banyandb/pkg/query/executor"
	logical2 "github.com/apache/skywalking-banyandb/pkg/query/logical"
)

func TestPlanExecution_Limit(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, s := prepareSchema(assert)

	tests := []struct {
		name           string
		unresolvedPlan logical2.UnresolvedPlan
		wantLength     int
	}{
		{
			name:           "Limit 1",
			unresolvedPlan: logical2.Limit(NewMockDataFactory(ctrl, m, s, 20).MockParentPlan(), 1),
			wantLength:     1,
		},
		{
			name:           "Limit 10",
			unresolvedPlan: logical2.Limit(NewMockDataFactory(ctrl, m, s, 20).MockParentPlan(), 10),
			wantLength:     10,
		},
		{
			name:           "Limit 50",
			unresolvedPlan: logical2.Limit(NewMockDataFactory(ctrl, m, s, 20).MockParentPlan(), 50),
			wantLength:     20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := require.New(t)
			plan, err := tt.unresolvedPlan.Analyze(s)
			assert.NoError(err)
			assert.NotNil(plan)

			ec := executor.NewMockExecutionContext(ctrl)
			entities, err := plan.Execute(ec)
			assert.NoError(err)
			assert.Len(entities, tt.wantLength)
		})
	}
}

func TestPlanExecution_Offset(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, s := prepareSchema(assert)

	tests := []struct {
		name           string
		unresolvedPlan logical2.UnresolvedPlan
		wantLength     int
	}{
		{
			name:           "Offset 0",
			unresolvedPlan: logical2.Offset(NewMockDataFactory(ctrl, m, s, 20).MockParentPlan(), 0),
			wantLength:     20,
		},
		{
			name:           "Offset 10",
			unresolvedPlan: logical2.Offset(NewMockDataFactory(ctrl, m, s, 20).MockParentPlan(), 10),
			wantLength:     10,
		},
		{
			name:           "Limit 50",
			unresolvedPlan: logical2.Offset(NewMockDataFactory(ctrl, m, s, 20).MockParentPlan(), 50),
			wantLength:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := require.New(t)
			plan, err := tt.unresolvedPlan.Analyze(s)
			assert.NoError(err)
			assert.NotNil(plan)

			ec := executor.NewMockExecutionContext(ctrl)
			entities, err := plan.Execute(ec)
			assert.NoError(err)
			assert.Len(entities, tt.wantLength)
		})
	}
}

func TestPlanExecution_TraceIDFetch(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, s := prepareSchema(assert)

	traceID := "asdf1234"

	p := logical2.TraceIDFetch(traceID, m, s)
	assert.NotNil(p)
	f := NewMockDataFactory(ctrl, m, s, 10)
	entities, err := p.Execute(f.MockTraceIDFetch(traceID))
	assert.NoError(err)
	assert.Len(entities, 10)
}

func TestPlanExecution_IndexScan(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, s := prepareSchema(assert)

	st, et := time.Now().Add(-3*time.Hour), time.Now()

	tests := []struct {
		name           string
		unresolvedPlan logical2.UnresolvedPlan
		wantLength     int
		indexMatchers  []*indexMatcher
	}{
		{
			name: "Single Index Search",
			unresolvedPlan: logical2.IndexScan(uint64(st.UnixNano()), uint64(et.UnixNano()), m, []logical2.Expr{
				logical2.Eq(logical2.NewFieldRef("http.method"), logical2.Str("GET")),
			}, series.TraceStateDefault),
			indexMatchers: []*indexMatcher{NewIndexMatcher("http.method", roaring.NewPostingListWithInitialData(1, 2, 3, 4, 5, 6, 7, 8, 9, 10))},
			wantLength:    10,
		},
		{
			name: "Multiple Index Search",
			unresolvedPlan: logical2.IndexScan(uint64(st.UnixNano()), uint64(et.UnixNano()), m, []logical2.Expr{
				logical2.Eq(logical2.NewFieldRef("http.method"), logical2.Str("GET")),
				logical2.Eq(logical2.NewFieldRef("status_code"), logical2.Str("200")),
			}, series.TraceStateDefault),
			indexMatchers: []*indexMatcher{
				NewIndexMatcher("http.method", roaring.NewPostingListWithInitialData(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)),
				NewIndexMatcher("status_code", roaring.NewPostingListWithInitialData(1, 3, 5, 7, 9)),
			},
			wantLength: 5,
		},
		{
			name: "Multiple Index With One Empty Result(ChunkID)",
			unresolvedPlan: logical2.IndexScan(uint64(st.UnixNano()), uint64(et.UnixNano()), m, []logical2.Expr{
				logical2.Eq(logical2.NewFieldRef("http.method"), logical2.Str("GET")),
				logical2.Eq(logical2.NewFieldRef("status_code"), logical2.Str("200")),
			}, series.TraceStateDefault),
			indexMatchers: []*indexMatcher{
				NewIndexMatcher("http.method", roaring.NewPostingListWithInitialData(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)),
				NewIndexMatcher("status_code", roaring.NewPostingList()),
			},
			wantLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := require.New(t)
			p, err := tt.unresolvedPlan.Analyze(s)
			assert.NoError(err)
			assert.NotNil(p)

			f := NewMockDataFactory(ctrl, m, s, 0)
			entities, err := p.Execute(f.MockIndexScan(st, et, tt.indexMatchers...))
			assert.NoError(err)
			assert.NotNil(entities)
			assert.Len(entities, tt.wantLength)
		})
	}
}

func TestPlanExecution_OrderBy(t *testing.T) {
	assert := require.New(t)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m, s := prepareSchema(assert)

	tests := []struct {
		name        string
		targetField string
		// TODO: avoid hardcoded index?
		targetFieldIdx int
		sortDirection  apiv1.Sort
	}{
		{
			name:           "Sort By trace_id ASC",
			targetField:    "trace_id",
			targetFieldIdx: 0,
			sortDirection:  apiv1.SortASC,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := require.New(t)
			p, err := logical2.OrderBy(NewMockDataFactory(ctrl, m, s, 20).MockParentPlan(), tt.targetField, tt.sortDirection).Analyze(s)
			assert.NoError(err)
			assert.NotNil(p)

			ec := executor.NewMockExecutionContext(ctrl)
			entities, err := p.Execute(ec)
			assert.NoError(err)
			assert.NotNil(entities)

			assert.True(logical2.Sorted(entities, tt.targetFieldIdx, tt.sortDirection))
			assert.False(logical2.Sorted(entities, tt.targetFieldIdx, reverseSortDirection(tt.sortDirection)))
		})
	}
}

func reverseSortDirection(sort apiv1.Sort) apiv1.Sort {
	if sort == apiv1.SortDESC {
		return apiv1.SortASC
	}
	return apiv1.SortDESC
}
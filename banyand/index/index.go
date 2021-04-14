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

package index

import (
	"github.com/apache/skywalking-banyandb/banyand/internal/bus"
	"github.com/apache/skywalking-banyandb/banyand/storage"
	"github.com/apache/skywalking-banyandb/pkg/logger"
	"github.com/apache/skywalking-banyandb/pkg/run"
)

const name = "index"

var (
	_ bus.MessageListener    = (*Index)(nil)
	_ run.PreRunner          = (*Index)(nil)
	_ storage.DataSubscriber = (*Index)(nil)
)

type Index struct {
	log *logger.Logger
}

func (s *Index) ComponentName() string {
	return name
}

func (s *Index) Sub(subscriber bus.Subscriber) error {
	return subscriber.Subscribe(storage.TraceIndex, s)
}

func (s Index) Name() string {
	return name
}

func (s Index) PreRun() error {
	s.log = logger.GetLogger(name)
	return nil
}

func (s Index) Rev(message bus.Message) {
	s.log.Info("rev", logger.Any("msg", message.Data()))
}
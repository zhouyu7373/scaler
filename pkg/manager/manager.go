/*
Copyright 2023 The Alibaba Cloud Serverless Authors.
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

package manager

import (
	"fmt"
	"log"
	"sync"

	"github.com/AliyunContainerService/scaler/pkg/config"
	"github.com/AliyunContainerService/scaler/pkg/model"
	"github.com/AliyunContainerService/scaler/pkg/scaler"
)

type Manager struct {
	mutex      sync.Mutex
	schedulers map[string]scaler.Scaler
	config     *config.Config
}

func New(config *config.Config) *Manager {
	return &Manager{
		mutex:      sync.Mutex{},
		schedulers: make(map[string]scaler.Scaler),
		config:     config,
	}
}

func (m *Manager) GetOrCreate(metaData *model.Meta) scaler.Scaler {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if scheduler := m.schedulers[metaData.Key]; scheduler != nil {
		return scheduler
	}

	log.Printf("Create new scaler for app %s", metaData.Key)
	scheduler := scaler.New(metaData, m.config)
	m.schedulers[metaData.Key] = scheduler
	return scheduler
}

func (m *Manager) Get(metaKey string) (scaler.Scaler, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if scheduler := m.schedulers[metaKey]; scheduler != nil {
		return scheduler, nil
	}
	return nil, fmt.Errorf("scaler of app: %s not found", metaKey)
}

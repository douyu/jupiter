// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mockserver

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"
)

type notification struct {
	NamespaceName  string `json:"namespaceName,omitempty"`
	NotificationID int    `json:"notificationId,omitempty"`
}

type result struct {
	NamespaceName  string            `json:"namespaceName"`
	Configurations map[string]string `json:"configurations"`
	ReleaseKey     string            `json:"releaseKey"`
}

type mockServer struct {
	server http.Server

	lock          sync.Mutex
	notifications map[string]int
	config        map[string]map[string]string
}

// NotificationHandler ...
func (s *mockServer) NotificationHandler(rw http.ResponseWriter, req *http.Request) {
	s.lock.Lock()
	defer s.lock.Unlock()
	req.ParseForm()
	var notifications []notification
	if err := json.Unmarshal([]byte(req.FormValue("notifications")), &notifications); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	var changes []notification
	for _, noti := range notifications {
		if currentID := s.notifications[noti.NamespaceName]; currentID != noti.NotificationID {
			changes = append(changes, notification{NamespaceName: noti.NamespaceName, NotificationID: currentID})
		}
	}

	if len(changes) == 0 {
		rw.WriteHeader(http.StatusNotModified)
		return
	}
	bts, err := json.Marshal(&changes)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Write(bts)
}

// ConfigHandler ...
func (s *mockServer) ConfigHandler(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	strs := strings.Split(req.RequestURI, "/")
	var namespace, releaseKey = strings.Split(strs[4], "?")[0], req.FormValue("releaseKey")
	config := s.Get(namespace)

	var result = result{NamespaceName: namespace, Configurations: config, ReleaseKey: releaseKey}
	bts, err := json.Marshal(&result)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Write(bts)
}

var server *mockServer

// Set ...
func (s *mockServer) Set(namespace, key, value string) {
	server.lock.Lock()
	defer server.lock.Unlock()

	notificationID := s.notifications[namespace]
	notificationID++
	s.notifications[namespace] = notificationID

	if kv, ok := s.config[namespace]; ok {
		kv[key] = value
		return
	}
	kv := map[string]string{key: value}
	s.config[namespace] = kv
}

// Get ...
func (s *mockServer) Get(namespace string) map[string]string {
	server.lock.Lock()
	defer server.lock.Unlock()

	retM := make(map[string]string)
	for k, v := range s.config[namespace] {
		retM[k] = v
	}
	return retM
}

// Delete ...
func (s *mockServer) Delete(namespace, key string) {
	server.lock.Lock()
	defer server.lock.Unlock()

	if kv, ok := s.config[namespace]; ok {
		delete(kv, key)
	}

	notificationID := s.notifications[namespace]
	notificationID++
	s.notifications[namespace] = notificationID
}

// Set namespace's key value
func Set(namespace, key, value string) {
	server.Set(namespace, key, value)
}

// Delete namespace's key
func Delete(namespace, key string) {
	server.Delete(namespace, key)
}

// Run mock server
func Run() error {
	return server.server.ListenAndServe()
}

func init() {
	server = &mockServer{
		notifications: map[string]int{},
		config:        map[string]map[string]string{},
	}
	mux := http.NewServeMux()
	mux.Handle("/notifications/", http.HandlerFunc(server.NotificationHandler))
	mux.Handle("/configs/", http.HandlerFunc(server.ConfigHandler))
	server.server.Handler = mux
	server.server.Addr = ":16852"
}

// Close mock server
func Close() error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	defer cancel()

	return server.server.Shutdown(ctx)
}

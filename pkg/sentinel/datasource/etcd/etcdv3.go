package etcd

import (
	"context"
	"github.com/douyu/jupiter/pkg/client/etcdv3"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"

	"github.com/alibaba/sentinel-golang/ext/datasource"
	"github.com/alibaba/sentinel-golang/util"
	"github.com/pkg/errors"
)

type DataSource struct {
	datasource.Base
	propertyKey         string
	lastUpdatedRevision int64
	client              *etcdv3.Client
	// cancel is the func, call cancel will stop watching on the propertyKey
	cancel context.CancelFunc
	// closed indicate whether continuing to watch on the propertyKey
	closed util.AtomicBool
}

// NewDatasource new a DataSource instance.
// client is the etcd client, it must be useful and should be release by User.
func NewDatasource(client *etcdv3.Client, key string, handlers ...datasource.PropertyHandler) (*DataSource, error) {
	if client == nil {
		return nil, errors.New("The etcd client is nil.")
	}
	ds := &DataSource{
		client:      client,
		propertyKey: key,
	}
	for _, h := range handlers {
		ds.AddPropertyHandler(h)
	}
	return ds, nil
}

func (s *DataSource) Initialize() error {
	err := s.doReadAndUpdate()
	if err != nil {
		xlog.Errorf("Fail to update data for key[%s] when execute Initialize function, err: %+v", s.propertyKey, err)
	}
	go util.RunWithRecover(s.watch)
	return err
}

func (s *DataSource) ReadSource() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := s.client.Get(ctx, s.propertyKey)
	if err != nil {
		return nil, errors.Errorf("Fail to get value for property key[%s]", s.propertyKey)
	}
	if resp.Count == 0 {
		return nil, errors.Errorf("The key[%s] is not existed in etcd server.", s.propertyKey)
	}
	s.lastUpdatedRevision = resp.Header.GetRevision()
	xlog.Infof("Get the newest data for key:%s, revision: %d, value: %s", s.propertyKey, resp.Header.GetRevision(), resp.Kvs[0].Value)
	return resp.Kvs[0].Value, nil
}

func (s *DataSource) doReadAndUpdate() error {
	src, err := s.ReadSource()
	if err != nil {
		return err
	}
	return s.Handle(src)
}

func (s *DataSource) processWatchResponse(resp *clientv3.WatchResponse) {
	if resp.CompactRevision > s.lastUpdatedRevision {
		s.lastUpdatedRevision = resp.CompactRevision
	}
	if resp.Header.GetRevision() > s.lastUpdatedRevision {
		s.lastUpdatedRevision = resp.Header.GetRevision()
	}

	if err := resp.Err(); err != nil {
		xlog.Errorf("Watch on etcd endpoints(%+v) occur error, err: %+v", s.client.Endpoints(), err)
		return
	}

	for _, ev := range resp.Events {
		if ev.Type == mvccpb.PUT {
			err := s.doReadAndUpdate()
			if err != nil {
				xlog.Errorf("Fail to execute doReadAndUpdate for PUT event, err: %+v", err)
			}
		}
		if ev.Type == mvccpb.DELETE {
			updateErr := s.Handle(nil)
			if updateErr != nil {
				xlog.Errorf("Fail to execute doReadAndUpdate for DELETE event, err: %+v", updateErr)
			}
		}
	}
}

func (s *DataSource) watch() {
	// Add watch for propertyKey from lastUpdatedRevision updated after Initializing
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	rch := s.client.Watch(ctx, s.propertyKey, clientv3.WithCreatedNotify(), clientv3.WithRev(s.lastUpdatedRevision))
	for {
		for resp := range rch {
			s.processWatchResponse(&resp)
		}
		// Stop watching if datasource had been closed.
		if s.closed.Get() {
			return
		}
		time.Sleep(time.Duration(1) * time.Second)
		ctx, cancel = context.WithCancel(context.Background())
		s.cancel = cancel
		if s.lastUpdatedRevision > 0 {
			rch = s.client.Watch(ctx, s.propertyKey, clientv3.WithCreatedNotify(), clientv3.WithRev(s.lastUpdatedRevision))
		} else {
			rch = s.client.Watch(ctx, s.propertyKey, clientv3.WithCreatedNotify())
		}
	}
}

func (s *DataSource) Close() error {
	// stop to watch property key.
	s.closed.Set(true)
	s.cancel()

	return nil
}

package apollo

import (
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/philchia/agollo"
	"sync"
)

type apolloDataSource struct {
	client      *agollo.Client
	namespace   string
	propertyKey string
	changed     chan struct{}
	quit        chan struct{}
	sync.Once
}

// NewDataSource create an apolloDataSource
func NewDataSource(conf *agollo.Conf, namespace string, key string) conf.DataSource {
	client := agollo.NewClient(conf)
	ap := &apolloDataSource{
		client:      client,
		namespace:   namespace,
		propertyKey: key,
		changed:     make(chan struct{}),
		quit:        make(chan struct{}),
	}
	return ap
}

// ReadConfig read config content from apollo
func (ap *apolloDataSource) ReadConfig() ([]byte, error) {
	ap.Once.Do(func() {
		ap.client.Start()
		changedEvent := ap.client.WatchUpdate()
		go func() {
			ap.watch(changedEvent)
		}()
	})
	value := ap.client.GetStringValueWithNameSpace(ap.namespace, ap.propertyKey, "")
	return []byte(value), nil
}

// IsConfigChanged returns a chanel for notification when the config changed
func (ap *apolloDataSource) IsConfigChanged() <-chan struct{} {
	return ap.changed
}

func (ap *apolloDataSource) watch(changedEvent <-chan *agollo.ChangeEvent) {
	for {
		select {
		case <-changedEvent:
			ap.changed <- struct{}{}
		case <-ap.quit:
			close(ap.changed)
			return
		}
	}
}

// Close stop watching the config changed
func (ap *apolloDataSource) Close() error {
	ap.quit <- struct{}{}
	ap.client.Stop()
	return nil
}

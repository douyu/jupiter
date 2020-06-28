package apollo

import (
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/util/xgo"
	"github.com/philchia/agollo"
)

type apolloDataSource struct {
	client      *agollo.Client
	namespace   string
	propertyKey string
	changed     chan struct{}
	quit        chan struct{}
}

// NewDataSource creates an apolloDataSource
func NewDataSource(conf *agollo.Conf, namespace string, key string) conf.DataSource {
	client := agollo.NewClient(conf)
	ap := &apolloDataSource{
		client:      client,
		namespace:   namespace,
		propertyKey: key,
		changed:     make(chan struct{}),
		quit:        make(chan struct{}),
	}
	ap.client.Start()
	changedEvent := ap.client.WatchUpdate()
	xgo.Go(func() {
		ap.watch(changedEvent)
	})
	return ap
}

// ReadConfig reads config content from apollo
func (ap *apolloDataSource) ReadConfig() ([]byte, error) {
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
			ap.client.Stop()
			close(ap.changed)
			return
		}
	}
}

// Close stops watching the config changed
func (ap *apolloDataSource) Close() error {
	ap.quit <- struct{}{}
	return nil
}

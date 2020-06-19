package apollo

import (
	"github.com/douyu/jupiter/pkg/datasource/apollo/mockserver"
	"github.com/philchia/agollo"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	setup()
	time.Sleep(time.Second)
	code := m.Run()
	time.Sleep(time.Second * 3)
	mockserver.Close()
	os.Exit(code)
}

func setup() {
	go func() {
		if err := mockserver.Run(); err != nil {
			log.Fatal(err)
		}
	}()
}

func TestReadConfig(t *testing.T) {
	ds := NewDataSource(&agollo.Conf{
		AppID:          "SampleApp",
		Cluster:        "default",
		NameSpaceNames: []string{"application"},
		IP:             "localhost:16852",
	}, "application", "key_test")

	mockserver.Set("application", "key_test", "value1")
	value, err := ds.ReadConfig()
	assert.Nil(t, err)
	assert.Equal(t, "value1", string(value))
	t.Logf("read: %s", value)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		mockserver.Set("application", "key_test", "value2")
		time.Sleep(time.Second * 3)
		ds.Close()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range ds.IsConfigChanged() {
			value, err := ds.ReadConfig()
			assert.Nil(t, err)
			assert.Equal(t, "value2", string(value))
			t.Logf("read: %s", value)
		}
	}()
	wg.Wait()
}

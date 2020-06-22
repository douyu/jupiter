package apollo

import (
	"github.com/douyu/jupiter/pkg/datasource/apollo/mockserver"
	"github.com/philchia/agollo"
	"github.com/stretchr/testify/assert"
	"log"
	"sync"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	setup()
	m.Run()
	teardown()
}

func setup() {
	go func() {
		if err := mockserver.Run(); err != nil {
			log.Println(err)
		}
	}()
	// wait for mock server to run
	time.Sleep(time.Second)
}

func teardown() {
	mockserver.Close()
}

func TestReadConfig(t *testing.T) {
	testData := []string{"value1", "value2"}
	ds := NewDataSource(&agollo.Conf{
		AppID:          "SampleApp",
		Cluster:        "default",
		NameSpaceNames: []string{"application"},
		IP:             "localhost:16852",
	}, "application", "key_test")

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		mockserver.Set("application", "key_test", testData[0])
		time.Sleep(time.Second * 3)
		mockserver.Set("application", "key_test", testData[1])
		time.Sleep(time.Second * 3)
		ds.Close()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second)
		index := 0

		for range ds.IsConfigChanged() {
			value, err := ds.ReadConfig()
			assert.Nil(t, err)
			assert.Equal(t, testData[index], string(value))
			index++
			t.Logf("read: %s", value)
		}
	}()
	wg.Wait()
}

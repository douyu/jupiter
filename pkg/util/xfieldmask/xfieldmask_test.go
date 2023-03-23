package xfieldmask

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNestedFieldMask(t *testing.T) {

	t.Run("test Filter method", func(t *testing.T) {
		msg := &SayHelloResponse{
			Error: 0,
			Msg:   "请求正常",
			Data: &SayHelloResponse_Data{
				Name:      "foo",
				AgeNumber: 18,
				Sex:       Sex_SEX_MALE,
				Metadata: map[string]string{
					"Bar": "bar",
				},
			},
		}
		nfm := New([]string{"data", "data.name", "data.age_number", "data.metadata"})
		nfm.Filter(msg)
		assert.Equal(t, 0, int(msg.GetError()))
		assert.Equal(t, "请求正常", msg.GetMsg())
		assert.Equal(t, "foo", msg.GetData().GetName())
		assert.Equal(t, 18, int(msg.GetData().GetAgeNumber()))
		assert.Equal(t, Sex_SEX_UNSPECIFIED, msg.GetData().GetSex())
		assert.Equal(t, map[string]string{"Bar": "bar"}, msg.GetData().GetMetadata())
	})

	t.Run("test Prune method", func(t *testing.T) {
		msg := &SayHelloResponse{
			Error: 0,
			Msg:   "请求正常",
			Data: &SayHelloResponse_Data{
				Name:      "foo",
				AgeNumber: 18,
				Sex:       Sex_SEX_MALE,
				Metadata: map[string]string{
					"Bar": "bar",
				},
			},
		}
		nfm := New([]string{"data", "data.name", "data.age_number", "data.metadata"})
		nfm.Prune(msg)
		assert.Equal(t, 0, int(msg.GetError()))
		assert.Equal(t, "请求正常", msg.GetMsg())
		assert.Equal(t, "", msg.GetData().GetName())
		assert.Equal(t, 0, int(msg.GetData().GetAgeNumber()))
		assert.Equal(t, Sex_SEX_MALE, msg.GetData().GetSex())
		assert.Equal(t, map[string]string(nil), msg.GetData().GetMetadata())

	})

	t.Run("test Masked method", func(t *testing.T) {
		nfm := New([]string{"data", "data.name", "data.age_number", "data.metadata"})
		assert.False(t, nfm.Masked("error"))
		assert.False(t, nfm.Masked("msg"))
		assert.True(t, nfm.Masked("data"))
		assert.True(t, nfm.Masked("data.name"))
		assert.True(t, nfm.Masked("data.age_number"))
		assert.False(t, nfm.Masked("data.sex"))
		assert.True(t, nfm.Masked("data.metadata"))
	})
}

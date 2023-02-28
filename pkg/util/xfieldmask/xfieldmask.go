package xfieldmask

import (
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// NestedFieldMask helps to build a nested field mask.
type NestedFieldMask map[string]NestedFieldMask

type MaskMode uint32

const (
	// FiledMaskMode_Filter indicates that the only fields in fieldmask should be kept.
	MaskMode_Filter MaskMode = 0
	// FiledMaskMode_Prune indicates that the field in fieldmask should be removed.
	MaskMode_Prune MaskMode = 1
)

// New returns a new NestedFieldMask with the given field mask.
func New(path []string) NestedFieldMask {
	if len(path) == 0 {
		return nil
	}
	return NewWithPaths(path...)
}

// NewWithPaths returns a new NestedFieldMask with the given paths.
func NewWithPaths(paths ...string) NestedFieldMask {
	if len(paths) == 0 {
		return nil
	}

	nfm := make(NestedFieldMask, len(paths))
	for i := 0; i < len(paths); i++ {
		p := paths[i]
		if arr := strings.Split(p, "."); len(arr) > 1 {
			cur := nfm
			ok := false

			for _, p1 := range arr {
				if _, ok = cur[p1]; !ok {
					cur[p1] = make(NestedFieldMask, 1)
				}
				cur = cur[p1]
			}
		}

		if _, ok := nfm[p]; !ok {
			nfm[p] = make(NestedFieldMask, 1)
		}
	}

	return nfm
}

func (mask NestedFieldMask) Filter(m proto.Message) {
	if len(mask) == 0 {
		return
	}

	pr := m.ProtoReflect()
	pr.Range(func(fd protoreflect.FieldDescriptor, _ protoreflect.Value) bool {
		if string(fd.Name()) == "error" || string(fd.Name()) == "msg" {
			return true
		}
		nfm, ok := mask[string(fd.Name())]
		if !ok {
			pr.Clear(fd)
			return true
		}

		if len(nfm) == 0 {
			return true
		}

		// repeated field type
		if fd.IsList() && fd.Kind() == protoreflect.MessageKind {
			l := pr.Get(fd).List()
			for i := 0; i < l.Len(); i++ {
				nfm.Filter(l.Get(i).Message().Interface())
			}
			return true
		}

		if fd.Kind() == protoreflect.MessageKind {
			nfm.Filter(pr.Get(fd).Message().Interface())
		}

		return true
	})
}

func (mask NestedFieldMask) Prune(m proto.Message) {
	if len(mask) == 0 {
		return
	}

	pr := m.ProtoReflect()
	pr.Range(func(fd protoreflect.FieldDescriptor, _ protoreflect.Value) bool {
		if string(fd.Name()) == "error" || string(fd.Name()) == "msg" {
			return true
		}
		nfm, ok := mask[string(fd.Name())]
		if !ok {
			return true
		}

		if len(nfm) == 0 {
			pr.Clear(fd)
			return true
		}

		// repeated field type
		if fd.IsList() {
			switch fd.Kind() {
			case protoreflect.MessageKind:
				l := pr.Get(fd).List()
				for i := 0; i < l.Len(); i++ {
					nfm.Prune(l.Get(i).Message().Interface())
				}
			default:
				pr.Clear(fd)
			}
			return true
		}

		if fd.Kind() == protoreflect.MessageKind {
			nfm.Prune(pr.Get(fd).Message().Interface())
		}

		return true
	})
}

func (mask NestedFieldMask) Masked(path string) bool {
	arr := strings.Split(path, ".")
	cur := mask
	for _, p := range arr {
		if _, ok := cur[p]; !ok {
			return false
		}
		cur = cur[p]
	}

	return true
}

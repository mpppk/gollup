package ast

import (
	"fmt"
	"go/types"
)

func findObject(objects []types.Object, object types.Object) (types.Object, bool) {
	for _, o := range objects {
		if isSameObject(o, object) {
			return o, true
		}
	}
	return nil, false
}

func isSameObject(obj1, obj2 types.Object) bool {
	return getObjectUniqueStr(obj1) == getObjectUniqueStr(obj2)
}

func getObjectUniqueStr(obj types.Object) string {
	recv := ""
	switch t := obj.Type().(type) {
	case *types.Signature:
		if t.Recv() != nil {
			recv = t.Recv().Type().String()
		}
	}
	return fmt.Sprintf("%s:%s:%s", obj.Pkg(), recv, obj.Name())
}

func distinctObjects(objects []types.Object) (newObjects []types.Object) {
	m := map[string]types.Object{}
	for _, object := range objects {
		m[getObjectUniqueStr(object)] = object
	}
	for _, object := range m {
		newObjects = append(newObjects, object)
	}
	return
}

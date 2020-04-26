package ast

import "go/types"

func findObject(objects []types.Object, object types.Object) (types.Object, bool) {
	for _, o := range objects {
		if isSameObject(o, object) {
			return o, true
		}
	}
	return nil, false
}

func isSameObject(obj1, obj2 types.Object) bool {
	return obj1.Pkg().Path() == obj2.Pkg().Path() &&
		obj1.Name() == obj2.Name()
}

func distinctObjects(objects []types.Object) (newObjects []types.Object) {
	m := map[string]types.Object{}
	for _, object := range objects {
		m[object.Pkg().Path()+object.Name()] = object
	}
	for _, object := range m {
		newObjects = append(newObjects, object)
	}
	return
}

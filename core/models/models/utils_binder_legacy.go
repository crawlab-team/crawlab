package models

//func AssignFields(d interface{}, fieldIds ...interfaces.ModelId) (res interface{}, err error) {
//	return assignFields(d, fieldIds...)
//}
//
//func AssignListFields(list interface{}, fieldIds ...interfaces.ModelId) (res arraylist.List, err error) {
//	return assignListFields(list, fieldIds...)
//}
//
//func AssignListFieldsAsPtr(list interface{}, fieldIds ...interfaces.ModelId) (res arraylist.List, err error) {
//	return assignListFieldsAsPtr(list, fieldIds...)
//}
//
//func assignFields(d interface{}, fieldIds ...interfaces.ModelId) (res interface{}, err error) {
//	doc, ok := d.(interfaces.Model)
//	if !ok {
//		return nil, errors.ErrorModelInvalidType
//	}
//	if len(fieldIds) == 0 {
//		return doc, nil
//	}
//	for _, fid := range fieldIds {
//		switch fid {
//		case interfaces.ModelIdTag:
//			// convert interface
//			d, ok := doc.(interfaces.ModelWithTags)
//			if !ok {
//				return nil, errors.ErrorModelInvalidType
//			}
//
//			// attempt to get artifact
//			a, err := doc.GetArtifact()
//			if err != nil {
//				return nil, err
//			}
//
//			// skip if no artifact found
//			if a == nil {
//				return d, nil
//			}
//
//			// assign tags
//			tags, err := a.GetTags()
//			if err != nil {
//				return nil, err
//			}
//			d.SetTags(tags)
//
//			return d, nil
//		}
//	}
//	return doc, nil
//}
//
//func _assignListFields(asPtr bool, list interface{}, fieldIds ...interfaces.ModelId) (res arraylist.List, err error) {
//	vList := reflect.ValueOf(list)
//	if vList.Kind() != reflect.Array &&
//		vList.Kind() != reflect.Slice {
//		return res, errors.ErrorModelInvalidType
//	}
//	for i := 0; i < vList.Len(); i++ {
//		vItem := vList.Index(i)
//		var item interface{}
//		if vItem.CanAddr() {
//			item = vItem.Addr().Interface()
//		} else {
//			item = vItem.Interface()
//		}
//		doc, ok := item.(interfaces.Model)
//		if !ok {
//			return res, errors.ErrorModelInvalidType
//		}
//		ptr, err := assignFields(doc, fieldIds...)
//		if err != nil {
//			return res, err
//		}
//		v := reflect.ValueOf(ptr)
//		if !asPtr {
//			// non-pointer item
//			res.Add(v.Elem().Interface())
//		} else {
//			// pointer item
//			res.Add(v.Interface())
//		}
//	}
//	return res, nil
//}
//
//func assignListFields(list interface{}, fieldIds ...interfaces.ModelId) (res arraylist.List, err error) {
//	return _assignListFields(false, list, fieldIds...)
//}
//
//func assignListFieldsAsPtr(list interface{}, fieldIds ...interfaces.ModelId) (res arraylist.List, err error) {
//	return _assignListFields(true, list, fieldIds...)
//}

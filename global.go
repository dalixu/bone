/*Package bone 一个跨平台的MVVM框架 可以方便的和native窗口结合起来
windows 可以参考bone/gwin来实现
*/
package bone

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
)

//Component 元信息 根据Component创建1个Bone
type Component struct {
	Template    string
	TemplateURI string
	Provider    []reflect.Type //用于依赖注入 只能是结构的指针类型 可以实现Constructor方法来作为构造方法
	DataContext reflect.Type   //可以实现Constructor方法来作为构造方法
}

type descriptor struct {
	component Component
	doc       xmlDoc
}

var boneStore map[string]*descriptor

//Import 导入组件名
func Import(key string, component Component) error {
	template := component.Template
	if template == "" && component.TemplateURI != "" {
		t, err := ioutil.ReadFile(component.TemplateURI)
		if err != nil {
			return err
		}
		template = string(t)
	}
	if template == "" {
		return errors.New("empty template or templateuri")
	}
	if key == "" {
		return errors.New("empty key")
	}
	if boneStore == nil {
		boneStore = make(map[string]*descriptor)
	}
	des := &descriptor{component: component}
	if err := des.doc.Parse(template); err != nil {
		return err
	}
	boneStore[key] = des
	return nil
}

func createBoneByTag(key string, ij *Injector) (*Bone, error) {
	des := boneStore[key]
	if des == nil {
		return nil, fmt.Errorf("unknown tag %s", key)
	}

	be, err := createBone(des, ij)
	if err != nil {
		return nil, err
	}
	//是否把be添加到父Bone 再考虑一下
	return be, err
}

func createBone(des *descriptor, ij *Injector) (*Bone, error) {
	//根据依赖注入创建ijInjector
	var provider []interface{}
	if len(des.component.Provider) > 0 {
		provider = make([]interface{}, len(des.component.Provider))
		for k, v := range des.component.Provider {
			if v.Kind() != reflect.Ptr {
				return nil, fmt.Errorf("%v must be Pointer", v)
			}
			ptr := reflect.New(v.Elem())
			ctor := ptr.Elem().MethodByName("Constructor")
			if ctor.IsValid() {
				_, err := ij.Call(ctor)
				if err != nil {
					return nil, err
				}
			}
			provider[k] = ptr
		}
	}
	thisIJ := &Injector{Parent: ij, Provider: provider}
	//执行vm的构造函数
	rt := des.component.DataContext
	var vm reflect.Value
	//var ctor reflect.Value
	if rt.Kind() == reflect.Ptr {
		vm = reflect.New(rt.Elem())
	} else {
		vm = reflect.New(rt).Elem()
	}
	// ctor = vm.MethodByName()
	// if ctor.IsValid() {
	// 	thisIJ.Call(ctor)
	// }
	be := NewBone(thisIJ, vm.Interface(), des.doc)
	be.CallVM("Constructor")
	return be, nil
}

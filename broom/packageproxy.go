package broom

import (
	"fmt"
	"reflect"
)

type PackageProxy struct {
	name  string
	funcs map[string]reflect.Value
}

func (pp *PackageProxy) Name() string {
	return pp.name
}

func (pp *PackageProxy) Find(name string) (reflect.Value, bool) {
	val, ok := pp.funcs[name]
	return val, ok
}

func (pp *PackageProxy) register(name string, f interface{}) {
	v := reflect.ValueOf(f)
	pp.funcs[name] = v
}

func NewPackageProxy(name string) *PackageProxy {
	p := &PackageProxy{name: name}
	p.funcs = make(map[string]reflect.Value)
	return p
}

func MakeReflectPackage() *PackageProxy {
	p := NewPackageProxy("reflect")
	p.register("Copy", reflect.Copy)
	p.register("ValueOf", reflect.ValueOf)
	p.register("TypeOf", reflect.TypeOf)
	p.register("DeepEqual", reflect.DeepEqual)

	p.register("ChanOf", reflect.ChanOf)
	p.register("MakeChan", reflect.MakeChan)

	p.register("MapOf", reflect.MapOf)
	p.register("MakeMap", reflect.MakeMap)

	p.register("SliceOf", reflect.SliceOf)
	p.register("MakeSlice", reflect.MakeSlice)

	p.register("PtrTo", reflect.PtrTo)
	p.register("Indirect", reflect.Indirect)
	p.register("New", reflect.New)
	//NewAt.
	p.register("Zero", reflect.Zero)

	p.register("Append", reflect.Append)
	p.register("AppendSlice", reflect.AppendSlice)

	p.register("Select", reflect.Select)

	p.register("ChanDir", func(recv, send bool) reflect.ChanDir {
		switch {
		case recv && send:
			return reflect.BothDir
		case recv:
			return reflect.RecvDir
		case send:
			return reflect.SendDir
		default:
			panic("never reach")
		}
	})

	p.register("SelectCase",
		func(Dir reflect.SelectDir, Chan reflect.Value, Send reflect.Value) *reflect.SelectCase {
			return &reflect.SelectCase{Dir: Dir, Chan: Chan, Send: Send}
		})

	p.register("SliceHeader",
		func(Data uintptr, Len int, Cap int) *reflect.SliceHeader {
			return &reflect.SliceHeader{Data: Data, Len: Len, Cap: Cap}
		})
	p.register("StringHeader",
		func(Data uintptr, Len int) *reflect.StringHeader {
			return &reflect.StringHeader{Data: Data, Len: Len}
		})
	p.register("StructField",
		func(Name string,
			PkgPath string,
			Type reflect.Type,
			Tag reflect.StructTag,
			Offset uintptr,
			Index []int,
			Anonymous bool) *reflect.StructField {
			return &reflect.StructField{
				Name:      Name,
				PkgPath:   PkgPath,
				Type:      Type,
				Tag:       Tag,
				Offset:    Offset,
				Index:     Index,
				Anonymous: Anonymous}
		})

	p.register("Bool", func() reflect.Kind { return reflect.Bool })
	p.register("Int", func() reflect.Kind { return reflect.Int })
	p.register("Int8", func() reflect.Kind { return reflect.Int8 })
	p.register("Int16", func() reflect.Kind { return reflect.Int16 })
	p.register("Int32", func() reflect.Kind { return reflect.Int32 })
	p.register("Int64", func() reflect.Kind { return reflect.Int64 })
	p.register("Uint", func() reflect.Kind { return reflect.Uint })
	p.register("Uint8", func() reflect.Kind { return reflect.Uint8 })
	p.register("Uint16", func() reflect.Kind { return reflect.Uint16 })
	p.register("Uint32", func() reflect.Kind { return reflect.Uint32 })
	p.register("Uint64", func() reflect.Kind { return reflect.Uint64 })
	p.register("Uintptr", func() reflect.Kind { return reflect.Uintptr })
	p.register("Float32", func() reflect.Kind { return reflect.Float32 })
	p.register("Float64", func() reflect.Kind { return reflect.Float64 })
	p.register("Complex64", func() reflect.Kind { return reflect.Complex64 })
	p.register("Complex128", func() reflect.Kind { return reflect.Complex128 })
	p.register("Array", func() reflect.Kind { return reflect.Array })
	p.register("Chan", func() reflect.Kind { return reflect.Chan })
	p.register("Func", func() reflect.Kind { return reflect.Func })
	p.register("Interface", func() reflect.Kind { return reflect.Interface })
	p.register("Map", func() reflect.Kind { return reflect.Map })
	p.register("Ptr", func() reflect.Kind { return reflect.Ptr })
	p.register("Slice", func() reflect.Kind { return reflect.Slice })
	p.register("String", func() reflect.Kind { return reflect.String })
	p.register("Struct", func() reflect.Kind { return reflect.Struct })
	p.register("UnsafePointer", func() reflect.Kind { return reflect.UnsafePointer })

	p.register("TypeByKind", func(k reflect.Kind) reflect.Type {
		switch k {
		case reflect.Bool:
			return reflect.TypeOf(true)
		case reflect.Int:
			return reflect.TypeOf(0)
		case reflect.Int8:
			return reflect.TypeOf(int8(0))
		case reflect.Int16:
			return reflect.TypeOf(int16(0))
		case reflect.Int32:
			return reflect.TypeOf(int32(0))
		case reflect.Int64:
			return reflect.TypeOf(int64(0))
		case reflect.Uint:
			return reflect.TypeOf(uint(0))
		case reflect.Uint8:
			return reflect.TypeOf(uint8(0))
		case reflect.Uint16:
			return reflect.TypeOf(uint16(0))
		case reflect.Uint32:
			return reflect.TypeOf(uint32(0))
		case reflect.Uint64:
			return reflect.TypeOf(uint64(0))
		case reflect.Uintptr:
			return reflect.TypeOf(uintptr(0))
		case reflect.Float32:
			return reflect.TypeOf(float32(0))
		case reflect.Float64:
			return reflect.TypeOf(float64(0))
		case reflect.Complex64:
			return reflect.TypeOf(complex64(0))
		case reflect.Complex128:
			return reflect.TypeOf(complex128(0))
		default:
			panic(fmt.Sprintf("bad reflect.Kind: %d", k))
		}
		panic("never reach")
	})
	return p
}

func MakeOSPackage() *PackageProxy {
	p := NewPackageProxy("os")

	return p
}

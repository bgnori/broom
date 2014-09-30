package broom

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
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
		case reflect.String:
			return reflect.TypeOf("abc")
		default:
			panic(fmt.Sprintf("bad reflect.Kind: %d", k))
		}
		panic("never reach")
	})
	return p
}

func MakeOSPackage() *PackageProxy {
	p := NewPackageProxy("os")

	// CONSTANTS
	p.register("O_RDONLY", func() int { return os.O_RDONLY })
	p.register("O_WRONLY", func() int { return os.O_WRONLY })
	p.register("O_RDWR  ", func() int { return os.O_RDWR })
	p.register("O_APPEND", func() int { return os.O_APPEND })
	p.register("O_CREATE", func() int { return os.O_CREATE })
	p.register("O_EXCL  ", func() int { return os.O_EXCL })
	p.register("O_SYNC  ", func() int { return os.O_SYNC })
	p.register("O_TRUNC ", func() int { return os.O_TRUNC })
	p.register("SEEK_SET", func() int { return os.SEEK_SET })
	p.register("SEEK_CUR", func() int { return os.SEEK_CUR })
	p.register("SEEK_END", func() int { return os.SEEK_END })
	p.register("PathSeparator", func() rune { return os.PathSeparator })
	p.register("PathListSeparator", func() rune { return os.PathListSeparator })
	p.register("DevNull", func() string { return os.DevNull })

	//VARIABLES
	p.register("ErrInvalid", func() error { return os.ErrInvalid })
	p.register("ErrPermission", func() error { return os.ErrPermission })
	p.register("ErrExist", func() error { return os.ErrExist })
	p.register("ErrNotExist", func() error { return os.ErrNotExist })
	p.register("Stdin", func() *os.File { return os.Stdin })
	p.register("Stdout", func() *os.File { return os.Stdout })
	p.register("Stderr", func() *os.File { return os.Stderr })
	p.register("Args", func() []string { return os.Args })

	//FUNCTIONS
	p.register("Chdir", os.Chdir)
	p.register("Chmod", os.Chmod)
	p.register("Chown", os.Chown)
	p.register("Chtimes", os.Chtimes)
	p.register("Clearenv", os.Clearenv)
	p.register("Environ", os.Environ)
	p.register("Exit", os.Exit)
	p.register("Expand", os.Expand)
	p.register("ExpandEnv", os.ExpandEnv)
	p.register("Getegid", os.Getegid)
	p.register("Getenv", os.Getenv)
	p.register("Geteuid", os.Geteuid)
	p.register("Getgid", os.Getgid)
	p.register("Getgroups", os.Getgroups)
	p.register("Getpagesize", os.Getpagesize)
	p.register("Getpid", os.Getpid)
	p.register("Getppid", os.Getppid)
	p.register("Getuid", os.Getuid)
	p.register("Getwd", os.Getwd)
	p.register("Hostname", os.Hostname)
	p.register("IsExist", os.IsExist)
	p.register("IsNotExist", os.IsNotExist)
	p.register("IsPathSeparator", os.IsPathSeparator)
	p.register("IsPermission", os.IsPermission)
	p.register("Lchown", os.Lchown)
	p.register("Link", os.Link)
	p.register("Mkdir", os.Mkdir)
	p.register("MkdirAll", os.MkdirAll)
	p.register("NewSyscallError", os.NewSyscallError)
	p.register("Readlink", os.Readlink)
	p.register("Remove", os.Remove)
	p.register("RemoveAll", os.RemoveAll)
	p.register("Rename", os.Rename)
	p.register("SameFile", os.SameFile)
	p.register("Setenv", os.Setenv)
	p.register("Symlink", os.Symlink)
	p.register("TempDir", os.TempDir)
	p.register("Truncate", os.Truncate)

	//TYPES
	//File
	p.register("Create", os.Create)
	p.register("NewFile", os.NewFile)
	p.register("Open", os.Open)
	p.register("OpenFile", os.OpenFile)
	p.register("Pipe", os.Pipe)
	p.register("Lstat", os.Lstat)
	p.register("Stat", os.Stat)

	//type FileMode uint32
	p.register("FileModeX", func(v int) os.FileMode { return os.FileMode(v) }) // Avoid stupod typing,
	p.register("FileMode", func(v uint32) os.FileMode { return os.FileMode(v) })
	p.register("ModeDir", func() os.FileMode { return os.ModeDir })
	p.register("ModeAppend", func() os.FileMode { return os.ModeAppend })
	p.register("ModeExclusive", func() os.FileMode { return os.ModeExclusive })
	p.register("ModeTemporary", func() os.FileMode { return os.ModeTemporary })
	p.register("ModeSymlink", func() os.FileMode { return os.ModeSymlink })
	p.register("ModeDevice", func() os.FileMode { return os.ModeDevice })
	p.register("ModeNamedPipe", func() os.FileMode { return os.ModeNamedPipe })
	p.register("ModeSocket", func() os.FileMode { return os.ModeSocket })
	p.register("ModeSetuid", func() os.FileMode { return os.ModeSetuid })
	p.register("ModeSetgid", func() os.FileMode { return os.ModeSetgid })
	p.register("ModeCharDevice", func() os.FileMode { return os.ModeCharDevice })
	p.register("ModeSticky", func() os.FileMode { return os.ModeSticky })
	p.register("ModeType", func() os.FileMode { return os.ModeType })
	p.register("ModePerm", func() os.FileMode { return os.ModePerm })

	//Process
	p.register("FindProcess", os.FindProcess)
	p.register("StartProcess", os.StartProcess)

	p.register("Interrupt", func() os.Signal { return os.Interrupt })
	p.register("Kill", func() os.Signal { return os.Kill })

	return p
}

func MakeRuntimePackage() *PackageProxy {
	p := NewPackageProxy("runtime")

	//Const
	p.register("Compiler", func() string { return runtime.Compiler })
	p.register("GOARCH", func() string { return runtime.GOARCH })
	p.register("GOOS", func() string { return runtime.GOOS })

	//variables
	p.register("MemProfileRate", func() int { return runtime.MemProfileRate })

	p.register("MakeArrayOfBlockProfileRecord", func(n int) []runtime.BlockProfileRecord {
		return make([]runtime.BlockProfileRecord, n)
	})
	p.register("BlockProfile", func(p []runtime.BlockProfileRecord) (n int, ok bool) {
		return runtime.BlockProfile(p)
	})
	p.register("Breakpoint", func() { runtime.Breakpoint() })
	p.register("CPUProfile", func() []byte {
		return runtime.CPUProfile()
	})
	p.register("Caller", func(skip int) (pc uintptr, file string, line int, ok bool) {
		return runtime.Caller(skip)
	})

	p.register("Callers", func(skip int, pc []uintptr) int {
		return runtime.Callers(skip, pc)
	})
	p.register("GC", func() {
		runtime.GC()
		return
	})
	p.register("GOMAXPROCS", func(n int) int {
		return runtime.GOMAXPROCS(n)
	})
	p.register("GOROOT", func() string {
		return runtime.GOROOT()
	})
	p.register("Goexit", func() {
		runtime.Goexit()
		return
	})

	p.register("MakeArrayOfStackRecord", func(n int) []runtime.StackRecord {
		return make([]runtime.StackRecord, n)
	})

	p.register("GoroutineProfile", func(p []runtime.StackRecord) (n int, ok bool) {
		return runtime.GoroutineProfile(p)
	})

	p.register("Gosched", func() {
		runtime.Gosched()
		return
	})

	p.register("LockOSThread", func() {
		runtime.LockOSThread()
		return
	})

	p.register("UnlockOSThread", func() {
		runtime.UnlockOSThread()
		return
	})

	p.register("MakeArrayOfMemProfileRecord", func(n int) []runtime.MemProfileRecord {
		return make([]runtime.MemProfileRecord, n)
	})
	p.register("MemProfile", func(p []runtime.MemProfileRecord, inuseZero bool) (n int, ok bool) {
		return runtime.MemProfile(p, inuseZero)
	})

	p.register("NumCPU", func() int {
		return runtime.NumCPU()
	})

	p.register("NumCgoCall", func() int64 {
		return runtime.NumCgoCall()
	})

	p.register("NumGoroutine", func() int {
		return runtime.NumGoroutine()
	})

	p.register("NewMemStats", func() *runtime.MemStats {
		return &runtime.MemStats{}
	})

	p.register("ReadMemStats", func(m *runtime.MemStats) {
		runtime.ReadMemStats(m)
		return
	})

	p.register("SetBlockProfileRate", func(rate int) {
		runtime.SetBlockProfileRate(rate)
		return
	})

	p.register("SetCPUProfileRate", func(hz int) {
		runtime.SetCPUProfileRate(hz)
		return
	})

	p.register("SetFinalizer", func(x, f interface{}) {
		runtime.SetFinalizer(x, f)
		return
	})

	p.register("Stack", func(buf []byte, all bool) int {
		return runtime.Stack(buf, all)
	})

	p.register("ThreadCreateProfile", func(p []runtime.StackRecord) (n int, ok bool) {
		return runtime.ThreadCreateProfile(p)
	})

	p.register("Version", func() string {
		return runtime.Version()
	})

	p.register("FuncForPC", func(pc uintptr) *runtime.Func {
		return runtime.FuncForPC(pc)
	})

	//methods
	//func (f *Func) Entry() uintptr
	//func (f *Func) FileLine(pc uintptr) (file string, line int)
	//func (f *Func) Name() string
	//func (r *MemProfileRecord) InUseBytes() int64
	//func (r *MemProfileRecord) InUseObjects() int64
	//func (r *MemProfileRecord) Stack() []uintptr
	//func (e *TypeAssertionError) Error() string
	//func (*TypeAssertionError) RuntimeError()

	return p
}

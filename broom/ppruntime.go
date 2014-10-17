package broom

import (
	"runtime"
)

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

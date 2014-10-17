package broom

import (
	"os"
)

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

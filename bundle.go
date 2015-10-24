package main

import (
	"archive/tar"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"encoding/json"
	"github.com/opencontainers/specs"
	"golang.org/x/sys/unix"
)

// (g GraphTool) Bundle  ...
func (g *GraphTool) Bundle(imageName string, dst string) error {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "dg-bundle")
	if err != nil {
		return err
	}

	tmpMount := tmpDir + "/mount"
	tmpRoot := tmpDir + "/rootfs"

	if err := os.MkdirAll(tmpMount, 0755); err != nil {
		g.logger.Error(err.Error())
	}

	if err := g.Mount(imageName, tmpMount, []string{"ro", "nosuid"}); err != nil {
		g.logger.Error(err.Error())
	}

	if err := os.MkdirAll(tmpRoot, 0755); err != nil {
		g.logger.Error(err.Error())
	}

	var bytesCopied int64

	tarFile, err := os.Create(dst)
	if err != nil {
		return err
	}

	tarArchive := tar.NewWriter(tarFile)
	// Order is important
	defer g.Unmount(tmpMount)
	defer os.RemoveAll(tmpDir)
	defer tarArchive.Close()

	err = g.specFiles(tarArchive)
	if err != nil {
		g.logger.Error(err.Error())
	}

	if err := filepath.Walk(tmpMount, func(path string, info os.FileInfo, err error) error {
		dstFile := filepath.Join("rootfs", strings.TrimPrefix(path, tmpMount))
		var link string
		if info.Mode()&os.ModeSymlink != 0 {
			link, err = os.Readlink(path)
			if err != nil {
				return err
			}
		}

		hdr, err := tar.FileInfoHeader(info, link)
		if err != nil {
			return err
		}

		hdr.Name = dstFile

		if err := tarArchive.WriteHeader(hdr); err != nil {
			return err
		}

		if info.Mode().IsRegular() {
			if n, err := g.tarCp(path, tarArchive); err != nil {
				return err
			} else {
				bytesCopied += n
			}
		}
		return nil
	}); err != nil {
		g.logger.Error(err.Error())
	}

	g.logger.Infof("%d MB copied", bytesCopied/1024)
	return nil
}

func (g *GraphTool) tarCp(srcName string, tw *tar.Writer) (int64, error) {
	var (
		src *os.File
		err error
	)

	if src, err = os.Open(srcName); err != nil {
		return 0, err
	}
	defer src.Close()

	srcStat, err := src.Stat()
	if err != nil {
		g.logger.Error(err.Error())
	} else if err := unix.Fadvise(int(src.Fd()), 0, srcStat.Size(), unix.MADV_SEQUENTIAL); err != nil {
		g.logger.Error(err.Error())
	}

	if n, err := io.Copy(tw, src); err != nil {
		g.logger.Error(err.Error())
	} else {
		return n, nil
	}

	return 0, nil
}

func (g *GraphTool) specFiles(tw *tar.Writer) error {
	// shameless copy from https://github.com/opencontainers/runc/blob/master/spec.go
	spec := specs.LinuxSpec{
		Spec: specs.Spec{
			Version: specs.Version,
			Platform: specs.Platform{
				OS:   runtime.GOOS,
				Arch: runtime.GOARCH,
			},
			Root: specs.Root{
				Path:     "rootfs",
				Readonly: true,
			},
			Process: specs.Process{
				Terminal: true,
				User:     specs.User{},
				Args: []string{
					"sh",
				},
				Env: []string{
					"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
					"TERM=xterm",
				},
			},
			Hostname: "shell",
			Mounts: []specs.MountPoint{
				{
					Name: "proc",
					Path: "/proc",
				},
				{
					Name: "dev",
					Path: "/dev",
				},
				{
					Name: "devpts",
					Path: "/dev/pts",
				},
				{
					Name: "shm",
					Path: "/dev/shm",
				},
				{
					Name: "mqueue",
					Path: "/dev/mqueue",
				},
				{
					Name: "sysfs",
					Path: "/sys",
				},
				{
					Name: "cgroup",
					Path: "/sys/fs/cgroup",
				},
			},
		},
		Linux: specs.Linux{
			Capabilities: []string{
				"CAP_AUDIT_WRITE",
				"CAP_KILL",
				"CAP_NET_BIND_SERVICE",
			},
		},
	}
	rspec := specs.LinuxRuntimeSpec{
		RuntimeSpec: specs.RuntimeSpec{
			Mounts: map[string]specs.Mount{
				"proc": {
					Type:    "proc",
					Source:  "proc",
					Options: nil,
				},
				"dev": {
					Type:    "tmpfs",
					Source:  "tmpfs",
					Options: []string{"nosuid", "strictatime", "mode=755", "size=65536k"},
				},
				"devpts": {
					Type:    "devpts",
					Source:  "devpts",
					Options: []string{"nosuid", "noexec", "newinstance", "ptmxmode=0666", "mode=0620", "gid=5"},
				},
				"shm": {
					Type:    "tmpfs",
					Source:  "shm",
					Options: []string{"nosuid", "noexec", "nodev", "mode=1777", "size=65536k"},
				},
				"mqueue": {
					Type:    "mqueue",
					Source:  "mqueue",
					Options: []string{"nosuid", "noexec", "nodev"},
				},
				"sysfs": {
					Type:    "sysfs",
					Source:  "sysfs",
					Options: []string{"nosuid", "noexec", "nodev"},
				},
				"cgroup": {
					Type:    "cgroup",
					Source:  "cgroup",
					Options: []string{"nosuid", "noexec", "nodev", "relatime", "ro"},
				},
			},
		},
		Linux: specs.LinuxRuntime{
			Namespaces: []specs.Namespace{
				{
					Type: "pid",
				},
				{
					Type: "network",
				},
				{
					Type: "ipc",
				},
				{
					Type: "uts",
				},
				{
					Type: "mount",
				},
			},
			Rlimits: []specs.Rlimit{
				{
					Type: "RLIMIT_NOFILE",
					Hard: uint64(1024),
					Soft: uint64(1024),
				},
			},
			Devices: []specs.Device{
				{
					Type:        'c',
					Path:        "/dev/null",
					Major:       1,
					Minor:       3,
					Permissions: "rwm",
					FileMode:    0666,
					UID:         0,
					GID:         0,
				},
				{
					Type:        'c',
					Path:        "/dev/random",
					Major:       1,
					Minor:       8,
					Permissions: "rwm",
					FileMode:    0666,
					UID:         0,
					GID:         0,
				},
				{
					Type:        'c',
					Path:        "/dev/full",
					Major:       1,
					Minor:       7,
					Permissions: "rwm",
					FileMode:    0666,
					UID:         0,
					GID:         0,
				},
				{
					Type:        'c',
					Path:        "/dev/tty",
					Major:       5,
					Minor:       0,
					Permissions: "rwm",
					FileMode:    0666,
					UID:         0,
					GID:         0,
				},
				{
					Type:        'c',
					Path:        "/dev/zero",
					Major:       1,
					Minor:       5,
					Permissions: "rwm",
					FileMode:    0666,
					UID:         0,
					GID:         0,
				},
				{
					Type:        'c',
					Path:        "/dev/urandom",
					Major:       1,
					Minor:       9,
					Permissions: "rwm",
					FileMode:    0666,
					UID:         0,
					GID:         0,
				},
			},
			Resources: &specs.Resources{
				Memory: specs.Memory{
					Swappiness: -1,
				},
			},
			Seccomp: specs.Seccomp{
				DefaultAction: "SCMP_ACT_ALLOW",
				Syscalls:      []*specs.Syscall{},
			},
		},
	}
	specData, err := json.MarshalIndent(spec, "", "\t")
	if err != nil {
		return nil
	}

	rspecData, err := json.MarshalIndent(rspec, "", "\t")
	if err != nil {
		return err
	}

	hdrConfig := &tar.Header{
		Name: "config.json",
		Size: int64(len(specData)),
	}
	tw.WriteHeader(hdrConfig)
	if _, err := tw.Write(specData); err != nil {
		return err
	}

	hdrRun := &tar.Header{
		Name: "runtime.json",
		Size: int64(len(rspecData)),
	}
	tw.WriteHeader(hdrRun)
	if _, err := tw.Write(rspecData); err != nil {
		return err
	}

	return nil
}

package fs

import (
	"io/ioutil"
	"testing"

	_ "crypto/sha256"

	"github.com/docker/containerd/fs/fstest"
	"github.com/pkg/errors"
)

// TODO: Create copy directory which requires privilege
//  chown
//  mknod
//  setxattr fstest.SetXAttr("/home", "trusted.overlay.opaque", "y"),

func TestCopyDirectory(t *testing.T) {
	apply := fstest.MultiApply(
		fstest.CreateDirectory("/etc/", 0755),
		fstest.NewTestFile("/etc/hosts", []byte("localhost 127.0.0.1"), 0644),
		fstest.Link("/etc/hosts", "/etc/hosts.allow"),
		fstest.CreateDirectory("/usr/local/lib", 0755),
		fstest.NewTestFile("/usr/local/lib/libnothing.so", []byte{0x00, 0x00}, 0755),
		fstest.Symlink("libnothing.so", "/usr/local/lib/libnothing.so.2"),
		fstest.CreateDirectory("/home", 0755),
	)

	if err := testCopy(apply); err != nil {
		t.Fatalf("Copy test failed: %+v", err)
	}
}

func testCopy(apply fstest.Applier) error {
	t1, err := ioutil.TempDir("", "test-copy-src-")
	if err != nil {
		return errors.Wrap(err, "failed to create temporary directory")
	}
	t2, err := ioutil.TempDir("", "test-copy-dst-")
	if err != nil {
		return errors.Wrap(err, "failed to create temporary directory")
	}

	if err := apply(t1); err != nil {
		return errors.Wrap(err, "failed to apply changes")
	}

	if err := CopyDirectory(t2, t1); err != nil {
		return errors.Wrap(err, "failed to copy")
	}

	return fstest.CheckDirectoryEqual(t1, t2)
}
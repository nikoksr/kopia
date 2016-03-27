package dir

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/kopia/kopia/content"
)

// EntryType describes the type of an backup entry.
type EntryType string

const (
	// EntryTypeFile represents a regular file.
	EntryTypeFile EntryType = "f"

	// EntryTypeDirectory represents a directory entry which is a subdirectory.
	EntryTypeDirectory EntryType = "d"

	// EntryTypeSymlink represents a symbolic link.
	EntryTypeSymlink EntryType = "l"

	// EntryTypeSocket represents a UNIX socket.
	EntryTypeSocket EntryType = "s"

	// EntryTypeDevice represents a device.
	EntryTypeDevice EntryType = "v"

	// EntryTypeNamedPipe represents a named pipe.
	EntryTypeNamedPipe EntryType = "n"
)

// FileModeToType converts os.FileMode into EntryType.
func FileModeToType(mode os.FileMode) EntryType {
	switch mode & os.ModeType {
	case os.ModeDir:
		return EntryTypeDirectory

	case os.ModeDevice:
		return EntryTypeDevice

	case os.ModeSocket:
		return EntryTypeSocket

	case os.ModeSymlink:
		return EntryTypeSymlink

	case os.ModeNamedPipe:
		return EntryTypeNamedPipe

	default:
		return EntryTypeFile
	}
}

// Entry stores attributes of a single entry in a directory.
type Entry struct {
	Name          string
	Size          int64
	Type          EntryType
	ModTime       time.Time
	Mode          int16 // 0000 .. 0777
	UserID        uint32
	GroupID       uint32
	ObjectID      content.ObjectID
	MetadataCRC32 uint32

	List func() (Listing, error)
	Open func() (io.ReadCloser, error)
}

func (e *Entry) String() string {
	return fmt.Sprintf(
		"name: '%v' type: %v modTime: %v size: %v oid: '%v' uid: %v gid: %v",
		e.Name, e.Type, e.ModTime, e.Size, e.ObjectID, e.UserID, e.GroupID,
	)
}

// Directory contains access to contents of directory, both in original order and indexed by name.
type Directory struct {
	Ordered []*Entry
	ByName  map[string]*Entry
}
package filesync

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/v3/disk"
)

type Files []File

func (f Files) SpaceNeeded() int64 {
	var space int64

	for _, v := range f {
		switch v.Operation {
		case OperationCreate:
			space += v.Source.Info.Size()
		case OperationUpdate:
			space += v.Source.Info.Size() - v.Destination.Info.Size()
			break
		case OperationDelete:
			space -= v.Destination.Info.Size()
		}
	}

	return space
}

type DiskSpace struct {
	Free   string
	Needed string
	Left   string
}

func CalculateDiskSpace(files Files, destination string) (DiskSpace, error) {
	di, err := disk.Usage(destination)
	if err != nil {
		return DiskSpace{}, err
	}

	needed := files.SpaceNeeded()
	left := int64(di.Free) - needed

	neededDisplay := humanize.Bytes(uint64(needed))
	if needed < 0 {
		neededDisplay = "-" + humanize.Bytes(uint64(needed*-1))
	}

	leftDisplay := humanize.Bytes(uint64(left))
	if left < 0 {
		leftDisplay = "-" + humanize.Bytes(uint64(left*-1))
	}

	diskSpace := DiskSpace{
		Free:   humanize.Bytes(di.Free),
		Needed: neededDisplay,
		Left:   leftDisplay,
	}

	if left < 1 {
		return diskSpace, fmt.Errorf("not enough free disk space, need %d more bytes", left*-1)
	}

	return diskSpace, nil
}

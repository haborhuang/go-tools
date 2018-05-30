package du

import "syscall"

func DiskUsage(diskPath string) (int64, int64, int64, error) {
	statfs := &syscall.Statfs_t{}
	err := syscall.Statfs(diskPath, statfs)
	if err != nil {
		return 0, 0, 0, err
	}

	// Available is blocks available * fragment size
	available := int64(statfs.Bavail) * int64(statfs.Bsize)

	// Capacity is total block count * fragment size
	capacity := int64(statfs.Blocks) * int64(statfs.Bsize)

	// Usage is block being used * fragment size (aka block size).
	usage := (int64(statfs.Blocks) - int64(statfs.Bfree)) * int64(statfs.Bsize)

	return available, capacity, usage, nil
}
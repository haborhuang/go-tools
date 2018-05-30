package du

import "testing"

func TestDiskUsage(t *testing.T) {
	t.Log(DiskUsage("/"))
}
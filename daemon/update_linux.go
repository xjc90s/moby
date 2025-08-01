package daemon

import (
	"time"

	"github.com/moby/moby/api/types/container"
	libcontainerdtypes "github.com/moby/moby/v2/daemon/internal/libcontainerd/types"
	"github.com/opencontainers/runtime-spec/specs-go"
)

func toContainerdResources(resources container.Resources) *libcontainerdtypes.Resources {
	var r libcontainerdtypes.Resources

	if resources.BlkioWeight != 0 {
		r.BlockIO = &specs.LinuxBlockIO{
			Weight: &resources.BlkioWeight,
		}
	}

	cpu := specs.LinuxCPU{
		Cpus: resources.CpusetCpus,
		Mems: resources.CpusetMems,
	}
	if resources.CPUShares != 0 {
		shares := uint64(resources.CPUShares)
		cpu.Shares = &shares
	}

	var (
		period uint64
		quota  int64
	)
	if resources.NanoCPUs != 0 {
		period = uint64(100 * time.Millisecond / time.Microsecond)
		quota = resources.NanoCPUs * int64(period) / 1e9
	}
	if quota == 0 && resources.CPUQuota != 0 {
		quota = resources.CPUQuota
	}
	if period == 0 && resources.CPUPeriod != 0 {
		period = uint64(resources.CPUPeriod)
	}

	if period != 0 {
		cpu.Period = &period
	}
	if quota != 0 {
		cpu.Quota = &quota
	}

	if cpu != (specs.LinuxCPU{}) {
		r.CPU = &cpu
	}

	var memory specs.LinuxMemory
	if resources.Memory != 0 {
		memory.Limit = &resources.Memory
	}
	if resources.MemoryReservation != 0 {
		memory.Reservation = &resources.MemoryReservation
	}
	if resources.KernelMemory != 0 { //nolint:staticcheck // ignore SA1019: memory.Kernel is deprecated: kernel-memory limits are not supported in cgroups v2, and were obsoleted in [kernel v5.4]. This field should no longer be used, as it may be ignored by runtimes.
		memory.Kernel = &resources.KernelMemory //nolint:staticcheck // ignore SA1019: memory.Kernel is deprecated: kernel-memory limits are not supported in cgroups v2, and were obsoleted in [kernel v5.4]. This field should no longer be used, as it may be ignored by runtimes.
	}
	if resources.MemorySwap > 0 {
		memory.Swap = &resources.MemorySwap
	}

	if memory != (specs.LinuxMemory{}) {
		r.Memory = &memory
	}

	r.Pids = getPidsLimit(resources)
	return &r
}

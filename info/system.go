package info

import (
	"fmt"
	"time"

	"github.com/mickep76/runshit/color"
	"github.com/mickep76/runshit/ts"
)

func (s *System) FmtStringColor(addr string) string {
	f := fmt.Sprintf("\n\t%s%%-24s%s : %s%%v%s", color.Cyan, color.Reset, color.Yellow, color.Reset)

	txt := fmt.Sprintf("%s%s%s", color.White, addr, color.Reset)
	txt += fmt.Sprintf(f, "UUID", s.Uuid)

	if s.Created != nil {
		txt += fmt.Sprintf(f, "Created", ts.Timestamp(*s.Created))
	}

	if s.Updated != nil {
		txt += fmt.Sprintf(f, "Updated", ts.Timestamp(*s.Updated))
	}

	if s.LastSeen != nil {
		lastSeen := time.Now().Sub(ts.Timestamp(*s.LastSeen).Time)
		txt += fmt.Sprintf("%s %sago%s", fmt.Sprintf(f, "Last Seen", lastSeen.Truncate(time.Second)), color.Cyan, color.Reset)
	}

	txt += fmt.Sprintf(f, "Hostname", s.Hostname)
	txt += fmt.Sprintf(f, "Manufacturer", s.Manufacturer)
	txt += fmt.Sprintf(f, "Product", s.Product)
	txt += fmt.Sprintf(f, "Product Version", s.ProductVersion)
	txt += fmt.Sprintf(f, "Serial Number", s.SerialNumber)

	switch s.Kernel {
	case "Linux":
		txt += fmt.Sprintf(f, "BIOS Vendor", s.BiosVendor)
		txt += fmt.Sprintf(f, "BIOS Date", s.BiosDate)
		txt += fmt.Sprintf(f, "BIOS Version", s.BiosVersion)
	case "Darwin":
		txt += fmt.Sprintf(f, "Boot ROM Version", s.BootRomVersion)
		txt += fmt.Sprintf(f, "SMC Version", s.SmcVersion)
	}

	txt += fmt.Sprintf("%s %sGB%s", fmt.Sprintf(f, "Memory", s.MemoryGb), color.Cyan, color.Reset)
	txt += fmt.Sprintf(f, "CPU Model", s.CpuModel)
	txt += fmt.Sprintf(f, "CPU Flags", s.CpuFlags)
	txt += fmt.Sprintf(f, "CPU Cores Per Socket", s.CpuCoresPerSocket)
	txt += fmt.Sprintf(f, "CPU Physical Cores", s.CpuPhysicalCores)
	txt += fmt.Sprintf(f, "CPU Logical Cores", s.CpuLogicalCores)
	txt += fmt.Sprintf(f, "CPU Sockets", s.CpuSockets)
	txt += fmt.Sprintf(f, "CPU Threads Per Core", s.CpuThreadsPerCore)
	txt += fmt.Sprintf(f, "Operating System", s.Os)
	txt += fmt.Sprintf(f, "Operating System Version", s.OsVersion)
	txt += fmt.Sprintf(f, "Operating System Build", s.OsBuild)
	txt += fmt.Sprintf(f, "Kernel", s.Kernel)
	txt += fmt.Sprintf(f, "Kernel Version", s.KernelVersion)
	txt += fmt.Sprintf(f, "Kernel Release", s.KernelRelease)
	txt += "\n\n"
	return txt
}

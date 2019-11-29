package scaleio

import "github.com/netdata/go.d.plugin/modules/scaleio/client"

func (s *ScaleIO) collectSystemOverview(mx *metrics, stats client.SelectedStatistics) {
	collectSystemCapacity(mx, stats)
	collectSystemComponents(mx, stats)
	collectSystemIOWorkload(mx, stats)
	collectSystemRebalance(mx, stats)
	collectSystemRebuild(mx, stats)
}

func collectSystemCapacity(mx *metrics, stats client.SelectedStatistics) {
	m := &mx.SystemOverview.Capacity
	s := stats.System

	// Physical Capacity Calculation (as in the ScaleIO GUI)
	{
		m.AvailableForVolumeAllocation = s.CapacityAvailableForVolumeAllocationInKb
		m.MaxCapacity = s.MaxCapacityInKb

		m.Decreased = sum(s.MaxCapacityInKb, -s.CapacityLimitInKb) // TODO: probably wrong
		m.Degraded = sum(s.DegradedFailedCapacityInKb, s.DegradedHealthyCapacityInKb)
		m.Failed = s.FailedCapacityInKb
		m.InMaintenance = s.InMaintenanceCapacityInKb
		m.Protected = s.ProtectedCapacityInKb
		m.Spare = s.SpareCapacityInKb
		m.UnreachableUnused = s.UnreachableUnusedCapacityInKb
		// Note: can't use 'UnusedCapacityInKb' directly, dashboard shows calculated value
		used := sum(
			s.ProtectedCapacityInKb,
			s.InMaintenanceCapacityInKb,
			m.Decreased,
			m.Degraded,
			s.FailedCapacityInKb,
			s.SpareCapacityInKb,
			s.UnreachableUnusedCapacityInKb,
		)
		m.Unused = sum(s.MaxCapacityInKb, -used)
	}

	m.ThickInUse = s.ThickCapacityInUseInKb
	m.ThinAllocated = s.ThinCapacityAllocatedInKb
	m.ThinInUse = s.ThinCapacityInUseInKb
	m.ThinFree = sum(s.ThinCapacityAllocatedInKb, -s.ThinCapacityInUseInKb)
}

func collectSystemComponents(mx *metrics, stats client.SelectedStatistics) {
	m := &mx.SystemOverview.Components
	s := stats.System

	m.Devices = s.NumOfDevices
	m.FaultSets = s.NumOfFaultSets
	m.MappedToAllVolumes = s.NumOfMappedToAllVolumes
	m.ProtectionDomains = s.NumOfProtectionDomains
	m.RfcacheDevices = s.NumOfRfcacheDevices
	m.Sdc = s.NumOfSdc
	m.Sds = s.NumOfSds
	m.Snapshots = s.NumOfSnapshots
	m.StoragePools = s.NumOfStoragePools
	m.VTrees = s.NumOfVtrees
	m.Volumes = s.NumOfVolumes

	m.ThickBaseVolumes = s.NumOfThickBaseVolumes
	m.ThinBaseVolumes = s.NumOfThinBaseVolumes
	m.UnmappedVolumes = s.NumOfUnmappedVolumes
	m.MappedVolumes = sum(s.NumOfVolumes, -s.NumOfUnmappedVolumes)
}

func collectSystemIOWorkload(mx *metrics, stats client.SelectedStatistics) {
	m := &mx.SystemOverview.IOWorkload
	s := stats.System

	m.Total.BW.set(
		calcBW(s.TotalReadBwc),
		calcBW(s.TotalWriteBwc),
	)
	m.Frontend.BW.set(
		calcBW(s.UserDataReadBwc),
		calcBW(s.UserDataWriteBwc),
	)
	m.Backend.Primary.BW.set(
		calcBW(s.PrimaryReadBwc),
		calcBW(s.PrimaryWriteBwc),
	)
	m.Backend.Secondary.BW.set(
		calcBW(s.SecondaryReadBwc),
		calcBW(s.SecondaryWriteBwc),
	)
	m.Backend.Total.BW.set(
		sum(m.Backend.Primary.BW.Read, m.Backend.Secondary.BW.Read),
		sum(m.Backend.Primary.BW.Write, m.Backend.Secondary.BW.Write),
	)

	m.Total.IOPS.set(
		calcIOPS(s.TotalReadBwc),
		calcIOPS(s.TotalWriteBwc),
	)
	m.Frontend.IOPS.set(
		calcIOPS(s.UserDataReadBwc),
		calcIOPS(s.UserDataWriteBwc),
	)
	m.Backend.Primary.IOPS.set(
		calcIOPS(s.PrimaryReadBwc),
		calcIOPS(s.PrimaryWriteBwc),
	)
	m.Backend.Secondary.IOPS.set(
		calcIOPS(s.SecondaryReadBwc),
		calcIOPS(s.SecondaryWriteBwc),
	)
	m.Backend.Total.IOPS.set(
		sum(m.Backend.Primary.IOPS.Read, m.Backend.Secondary.IOPS.Read),
		sum(m.Backend.Primary.IOPS.Write, m.Backend.Secondary.IOPS.Write),
	)

	m.Total.IOSize.set(
		calcIOSize(s.TotalReadBwc),
		calcIOSize(s.TotalWriteBwc),
	)
	m.Frontend.IOSize.set(
		calcIOSize(s.UserDataReadBwc),
		calcIOSize(s.UserDataWriteBwc),
	)
	m.Backend.Primary.IOSize.set(
		calcIOSize(s.PrimaryReadBwc),
		calcIOSize(s.PrimaryWriteBwc),
	)
	m.Backend.Secondary.IOSize.set(
		calcIOSize(s.SecondaryReadBwc),
		calcIOSize(s.SecondaryWriteBwc),
	)
	m.Backend.Total.IOSize.set(
		sum(m.Backend.Primary.IOSize.Read, m.Backend.Secondary.IOSize.Read),
		sum(m.Backend.Primary.IOSize.Write, m.Backend.Secondary.IOSize.Write),
	)
}

func collectSystemRebuild(mx *metrics, stats client.SelectedStatistics) {
	m := &mx.SystemOverview.Rebuild
	s := stats.System

	m.Forward.BW.set(
		calcBW(s.FwdRebuildReadBwc),
		calcBW(s.FwdRebuildWriteBwc),
	)
	m.Backward.BW.set(
		calcBW(s.BckRebuildReadBwc),
		calcBW(s.BckRebuildWriteBwc),
	)
	m.Normal.BW.set(
		calcBW(s.NormRebuildReadBwc),
		calcBW(s.NormRebuildWriteBwc),
	)
	m.Total.BW.set(
		sum(m.Forward.BW.Read, m.Backward.BW.Read, m.Normal.BW.Read),
		sum(m.Forward.BW.Write, m.Backward.BW.Write, m.Normal.BW.Write),
	)

	m.Forward.IOPS.set(
		calcIOPS(s.FwdRebuildReadBwc),
		calcIOPS(s.FwdRebuildWriteBwc),
	)
	m.Backward.IOPS.set(
		calcIOPS(s.BckRebuildReadBwc),
		calcIOPS(s.BckRebuildWriteBwc),
	)
	m.Normal.IOPS.set(
		calcIOPS(s.NormRebuildReadBwc),
		calcIOPS(s.NormRebuildWriteBwc),
	)
	m.Total.IOPS.set(
		sum(m.Forward.IOPS.Read, m.Backward.IOPS.Read, m.Normal.IOPS.Read),
		sum(m.Forward.IOPS.Write, m.Backward.IOPS.Write, m.Normal.IOPS.Write),
	)

	m.Forward.IOSize.set(
		calcIOSize(s.FwdRebuildReadBwc),
		calcIOSize(s.FwdRebuildWriteBwc),
	)
	m.Backward.IOSize.set(
		calcIOSize(s.BckRebuildReadBwc),
		calcIOSize(s.BckRebuildWriteBwc),
	)
	m.Normal.IOSize.set(
		calcIOSize(s.NormRebuildReadBwc),
		calcIOSize(s.NormRebuildWriteBwc),
	)
	m.Total.IOSize.set(
		sum(m.Forward.IOSize.Read, m.Backward.IOSize.Read, m.Normal.IOSize.Read),
		sum(m.Forward.IOSize.Write, m.Backward.IOSize.Write, m.Normal.IOSize.Write),
	)

	// --Pending Capacity--
	m.Forward.Pending = s.PendingFwdRebuildCapacityInKb
	m.Backward.Pending = s.PendingBckRebuildCapacityInKb
	m.Normal.Pending = s.PendingNormRebuildCapacityInKb
	m.Total.Pending = sum(m.Forward.Pending, m.Backward.Pending, m.Normal.Pending)
}

func collectSystemRebalance(mx *metrics, stats client.SelectedStatistics) {
	m := &mx.SystemOverview.Rebalance
	s := stats.System

	m.BW.set(
		calcBW(s.RebalanceReadBwc),
		calcBW(s.RebalanceWriteBwc),
	)

	m.IOPS.set(
		calcIOPS(s.RebalanceReadBwc),
		calcIOPS(s.RebalanceWriteBwc),
	)

	m.IOSize.set(
		calcIOSize(s.RebalanceReadBwc),
		calcIOSize(s.RebalanceWriteBwc),
	)

	m.Pending = s.PendingRebalanceCapacityInKb
}

func calcBW(bwc client.Bwc) float64     { return div(bwc.TotalWeightInKb, bwc.NumSeconds) }
func calcIOPS(bwc client.Bwc) float64   { return div(bwc.NumOccured, bwc.NumSeconds) }
func calcIOSize(bwc client.Bwc) float64 { return div(bwc.TotalWeightInKb, bwc.NumOccured) }

func sum(vs ...float64) (res float64) {
	for _, v := range vs {
		res += v
	}
	return res
}

func div(a, b int64) float64 {
	if b == 0 {
		return 0
	}
	return float64(a) / float64(b)
}

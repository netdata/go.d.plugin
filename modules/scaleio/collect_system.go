package scaleio

import "github.com/netdata/go.d.plugin/modules/scaleio/client"

func (s *ScaleIO) collectSystemOverview(mx *metrics, stats client.SelectedStatistics) {
	collectSystemCapacity(mx, stats)
	collectSystemWorkload(mx, stats)
	collectSystemRebalance(mx, stats)
	collectSystemRebuild(mx, stats)
	collectSystemComponents(mx, stats)
}

func collectSystemCapacity(mx *metrics, stats client.SelectedStatistics) {
	collectCapacity(&mx.SystemOverview.Capacity, stats.System.CapacityStatistics)
}

func collectCapacity(m *capacity, s client.CapacityStatistics) {
	// Health
	m.Protected = s.ProtectedCapacityInKb
	m.InMaintenance = s.InMaintenanceCapacityInKb
	m.Degraded = sum(s.DegradedFailedCapacityInKb, s.DegradedHealthyCapacityInKb)
	m.Failed = s.FailedCapacityInKb
	m.UnreachableUnused = s.UnreachableUnusedCapacityInKb

	// Capacity
	m.MaxCapacity = s.MaxCapacityInKb
	m.ThickInUse = s.ThickCapacityInUseInKb
	m.ThinInUse = s.ThinCapacityInUseInKb
	m.Snapshot = s.SnapCapacityInUseOccupiedInKb
	m.Spare = s.SpareCapacityInKb
	m.Decreased = sum(s.MaxCapacityInKb, -s.CapacityLimitInKb) // TODO: probably wrong
	// Note: can't use 'UnusedCapacityInKb' directly, dashboard shows calculated value
	used := sum(
		s.ProtectedCapacityInKb,
		s.InMaintenanceCapacityInKb,
		m.Decreased,
		m.Degraded,
		s.FailedCapacityInKb,
		s.SpareCapacityInKb,
		s.UnreachableUnusedCapacityInKb,
		s.SnapCapacityInUseOccupiedInKb,
	)
	m.Unused = sum(s.MaxCapacityInKb, -used)

	// Other
	m.InUse = s.CapacityInUseInKb
	m.AvailableForVolumeAllocation = s.CapacityAvailableForVolumeAllocationInKb
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

func collectSystemWorkload(mx *metrics, stats client.SelectedStatistics) {
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
		sumFloat(m.Backend.Primary.BW.Read, m.Backend.Secondary.BW.Read),
		sumFloat(m.Backend.Primary.BW.Write, m.Backend.Secondary.BW.Write),
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
		sumFloat(m.Backend.Primary.IOPS.Read, m.Backend.Secondary.IOPS.Read),
		sumFloat(m.Backend.Primary.IOPS.Write, m.Backend.Secondary.IOPS.Write),
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
		sumFloat(m.Backend.Primary.IOSize.Read, m.Backend.Secondary.IOSize.Read),
		sumFloat(m.Backend.Primary.IOSize.Write, m.Backend.Secondary.IOSize.Write),
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
		sumFloat(m.Forward.BW.Read, m.Backward.BW.Read, m.Normal.BW.Read),
		sumFloat(m.Forward.BW.Write, m.Backward.BW.Write, m.Normal.BW.Write),
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
		sumFloat(m.Forward.IOPS.Read, m.Backward.IOPS.Read, m.Normal.IOPS.Read),
		sumFloat(m.Forward.IOPS.Write, m.Backward.IOPS.Write, m.Normal.IOPS.Write),
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
		sumFloat(m.Forward.IOSize.Read, m.Backward.IOSize.Read, m.Normal.IOSize.Read),
		sumFloat(m.Forward.IOSize.Write, m.Backward.IOSize.Write, m.Normal.IOSize.Write),
	)

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
	m.TimeUntilFinish = divFloat(float64(m.Pending), m.BW.ReadWrite)
}

func calcBW(bwc client.Bwc) float64     { return div(bwc.TotalWeightInKb, bwc.NumSeconds) }
func calcIOPS(bwc client.Bwc) float64   { return div(bwc.NumOccured, bwc.NumSeconds) }
func calcIOSize(bwc client.Bwc) float64 { return div(bwc.TotalWeightInKb, bwc.NumOccured) }

func sum(a, b int64, others ...int64) (res int64) {
	for _, v := range others {
		res += v
	}
	return res + a + b
}

func sumFloat(a, b float64, others ...float64) (res float64) {
	for _, v := range others {
		res += v
	}
	return res + a + b
}

func div(a, b int64) float64 {
	return divFloat(float64(a), float64(b))
}

func divFloat(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

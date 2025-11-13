package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

type ovsDPCollector struct {
	// PMD stats
	missWithSuccessUpcallMetric      *prometheus.Desc
	missWithFailedUpcallMetric       *prometheus.Desc
	avgSubtableLookupsMegaflowMetric *prometheus.Desc
	processingCyclesMetric           *prometheus.Desc
	idleCyclesMetric                 *prometheus.Desc
	// Memory/show stats
	memoryHandlersMetric            *prometheus.Desc
	memoryIdlCellsOpenVSwitchMetric *prometheus.Desc
	memoryOfconnsMetric             *prometheus.Desc
	memoryPortsMetric               *prometheus.Desc
	memoryRevalidatorsMetric        *prometheus.Desc
	memoryRulesMetric               *prometheus.Desc
	memoryUdpifKeysMetric           *prometheus.Desc
	// Offload stats
	offloadEnqueuedMetric           *prometheus.Desc
	offloadInsertedMetric           *prometheus.Desc
	offloadCtUniDirConnMetric       *prometheus.Desc
	offloadCtBiDirConnMetric        *prometheus.Desc
	offloadCumAvgLatencyUsMetric    *prometheus.Desc
	offloadCumLatencyStddevUsMetric *prometheus.Desc
	offloadCumLatencyMaxUsMetric    *prometheus.Desc
	offloadCumLatencyMinUsMetric    *prometheus.Desc
	offloadExpAvgLatencyUsMetric    *prometheus.Desc
	offloadExpLatencyStddevUsMetric *prometheus.Desc
	// Drop reasons
	upcallDropsMetric                      *prometheus.Desc
	upcallDropsLockErrorMetric             *prometheus.Desc
	rxDropsInvalidPacketMetric             *prometheus.Desc
	datapathDropMeterMetric                *prometheus.Desc
	datapathDropUserspaceActionErrorMetric *prometheus.Desc
	datapathDropTunnelPushErrorMetric      *prometheus.Desc
	datapathDropTunnelPopErrorMetric       *prometheus.Desc
	datapathDropRecircErrorMetric          *prometheus.Desc
	datapathDropInvalidPortMetric          *prometheus.Desc
	datapathDropInvalidTnlPortMetric       *prometheus.Desc
	datapathDropSampleErrorMetric          *prometheus.Desc
	datapathDropNshDecapErrorMetric        *prometheus.Desc
	dropActionOfPipelineMetric             *prometheus.Desc
	dropActionBridgeNotFoundMetric         *prometheus.Desc
	dropActionRecursionTooDeepMetric       *prometheus.Desc
	dropActionTooManyResubmitMetric        *prometheus.Desc
	dropActionStackTooDeepMetric           *prometheus.Desc
	dropActionNoRecirculationContextMetric *prometheus.Desc
	dropActionRecirculationConflictMetric  *prometheus.Desc
	dropActionTooManyMplsLabelsMetric      *prometheus.Desc
	dropActionInvalidTunnelMetadataMetric  *prometheus.Desc
	dropActionUnsupportedPacketTypeMetric  *prometheus.Desc
	dropActionCongestionMetric             *prometheus.Desc
	dropActionForwardingDisabledMetric     *prometheus.Desc
	// Drop reasons new
	netdevVxlanTsoDropsMetric         *prometheus.Desc
	netdevGeneveTsoDropsMetric        *prometheus.Desc
	netdevPushHeaderDropsMetric       *prometheus.Desc
	netdevSoftSegDropsMetric          *prometheus.Desc
	datapathDropTunnelTsoRecircMetric *prometheus.Desc
	datapathDropInvalidBondMetric     *prometheus.Desc
	datapathDropHwMissRecoverMetric   *prometheus.Desc
	// DOCA
	ovsDocaNoMarkMetric              *prometheus.Desc
	ovsDocaInvalidClassifyPortMetric *prometheus.Desc
	docaQueueEmptyMetric             *prometheus.Desc
	docaQueueNoneProcessedMetric     *prometheus.Desc
	docaResizeBlockMetric            *prometheus.Desc
	docaPipeResizeMetric             *prometheus.Desc
	docaPipeResizeOver10MsMetric     *prometheus.Desc
	// Upcall Flow Limit
	UpcallFlowLimitGrewMetric    *prometheus.Desc
	UpcallFlowLimitHitMetric     *prometheus.Desc
	UpcallFlowLimitKillMetric    *prometheus.Desc
	UpcallFlowLimitReducedMetric *prometheus.Desc
	UpcallFlowLimitScaledMetric  *prometheus.Desc
	// DOCA Pipe Group
	docaUniqueItemTemplatesMetric *prometheus.Desc
}

// allow tests to stub metric fetching
var fetchOvsMetrics = getOvsMetric

func isValidMetric(value float64) bool {
	return value != -1
}

func newOvsDPCollector() *ovsDPCollector {
	return &ovsDPCollector{
		// PMD stats
		missWithSuccessUpcallMetric: prometheus.NewDesc("ovsdp_miss_with_success_upcall",
			"Cache miss with successuful upcall",
			nil, nil,
		),
		missWithFailedUpcallMetric: prometheus.NewDesc("ovsdp_miss_with_failed_upcall",
			"Cache miss with failed upcall",
			nil, nil,
		),
		processingCyclesMetric: prometheus.NewDesc("ovsdp_processing_cycles",
			"CPU cycles spent actively checking for packets in a loop",
			nil, nil,
		),
		idleCyclesMetric: prometheus.NewDesc("ovsdp_idle_cycles",
			"Idle cycles waiting for packets",
			nil, nil,
		),
		avgSubtableLookupsMegaflowMetric: prometheus.NewDesc("ovsdp_avg_subtable_lookups_megaflow",
			"Average of subtable lookups per megaflow hit",
			nil, nil,
		),
		// Memory/show stats
		memoryHandlersMetric: prometheus.NewDesc("ovsdp_memory_handlers",
			"Number of OVS handler threads handling OpenFlow connections and upcalls",
			nil, nil,
		),
		memoryIdlCellsOpenVSwitchMetric: prometheus.NewDesc("ovsdp_memory_idl_cells_open_vswitch",
			"OVSDB cells in use for Open_vSwitch table (transaction/monitor memory)",
			nil, nil,
		),
		memoryOfconnsMetric: prometheus.NewDesc("ovsdp_memory_ofconns",
			"Active OpenFlow controller connections",
			nil, nil,
		),
		memoryPortsMetric: prometheus.NewDesc("ovsdp_memory_ports",
			"Configured datapath ports (physical, virtual, and internal)",
			nil, nil,
		),
		memoryRevalidatorsMetric: prometheus.NewDesc("ovsdp_memory_revalidators",
			"Revalidator threads that periodically revalidate userspace datapath flows",
			nil, nil,
		),
		memoryRulesMetric: prometheus.NewDesc("ovsdp_memory_rules",
			"Installed OpenFlow rules (software and hardware offloaded)",
			nil, nil,
		),
		memoryUdpifKeysMetric: prometheus.NewDesc("ovsdp_memory_udpif_keys",
			"Unique userspace datapath (udpif) flow keys handled in software",
			nil, nil,
		),
		// Offload stats
		offloadEnqueuedMetric: prometheus.NewDesc("ovsdp_offload_enqueued",
			"Enqueued offloads total",
			nil, nil,
		),
		offloadInsertedMetric: prometheus.NewDesc("ovsdp_offload_inserted",
			"Inserted offloads total",
			nil, nil,
		),
		offloadCtUniDirConnMetric: prometheus.NewDesc("ovsdp_offload_ct_unidir_connections",
			"CT uni-dir connections offloaded",
			nil, nil,
		),
		offloadCtBiDirConnMetric: prometheus.NewDesc("ovsdp_offload_ct_bidir_connections",
			"CT bi-dir connections offloaded",
			nil, nil,
		),
		offloadCumAvgLatencyUsMetric: prometheus.NewDesc("ovsdp_offload_cum_avg_latency_us",
			"Cumulative average latency (us)",
			nil, nil,
		),
		offloadCumLatencyStddevUsMetric: prometheus.NewDesc("ovsdp_offload_cum_latency_stddev_us",
			"Cumulative latency stddev (us)",
			nil, nil,
		),
		offloadCumLatencyMaxUsMetric: prometheus.NewDesc("ovsdp_offload_cum_latency_max_us",
			"Cumulative latency max (us)",
			nil, nil,
		),
		offloadCumLatencyMinUsMetric: prometheus.NewDesc("ovsdp_offload_cum_latency_min_us",
			"Cumulative latency min (us)",
			nil, nil,
		),
		offloadExpAvgLatencyUsMetric: prometheus.NewDesc("ovsdp_offload_exp_avg_latency_us",
			"Exponential average latency (us)",
			nil, nil,
		),
		offloadExpLatencyStddevUsMetric: prometheus.NewDesc("ovsdp_offload_exp_latency_stddev_us",
			"Exponential latency stddev (us)",
			nil, nil,
		),
		// Drop reasons
		upcallDropsMetric: prometheus.NewDesc("ovsdp_datapath_drop_upcall_error",
			"Drop packet due to error in the Upcall process",
			nil, nil,
		),
		upcallDropsLockErrorMetric: prometheus.NewDesc("ovsdp_datapath_drop_lock_error",
			"Drop packet due to Upcall lock contention",
			nil, nil,
		),
		rxDropsInvalidPacketMetric: prometheus.NewDesc("ovsdp_datapath_drop_rx_invalid_packet",
			"Drop invalid packet having size lower than what wrote in the Ethernet header",
			nil, nil,
		),
		datapathDropMeterMetric: prometheus.NewDesc("ovsdp_datapath_drop_meter",
			"Drop packet in the OpenFlow (1.3+) Meter Table",
			nil, nil,
		),
		datapathDropUserspaceActionErrorMetric: prometheus.NewDesc("ovsdp_datapath_drop_userspace_action_error",
			"Drop packet due to generic error executing the action",
			nil, nil,
		),
		datapathDropTunnelPushErrorMetric: prometheus.NewDesc("ovsdp_datapath_drop_tunnel_push_error",
			"Drop packet due to error executing the tunnel push (aka encapsulation) action",
			nil, nil,
		),
		datapathDropTunnelPopErrorMetric: prometheus.NewDesc("ovsdp_datapath_drop_tunnel_pop_error",
			"Drop packet due to error executing the tunnel pop (aka decapsulation) action",
			nil, nil,
		),
		datapathDropRecircErrorMetric: prometheus.NewDesc("ovsdp_datapath_drop_recirc_error",
			"Drop packet due to error in the recirculation (this can also happen in the tunnel pop action)",
			nil, nil,
		),
		datapathDropInvalidPortMetric: prometheus.NewDesc("ovsdp_datapath_drop_invalid_port",
			"Drop packet due to invalid port",
			nil, nil,
		),
		datapathDropInvalidTnlPortMetric: prometheus.NewDesc("ovsdp_datapath_drop_invalid_tnl_port",
			"Drop packet due to invalid tunnel port executing the pop action",
			nil, nil,
		),
		datapathDropSampleErrorMetric: prometheus.NewDesc("ovsdp_datapath_drop_sample_error",
			"Drop packet due to sampling error",
			nil, nil,
		),
		datapathDropNshDecapErrorMetric: prometheus.NewDesc("ovsdp_datapath_drop_nsh_decap_error",
			"Drop packet due to invalid NSH pop (aka decapsulation)",
			nil, nil,
		),
		dropActionOfPipelineMetric: prometheus.NewDesc("ovsdp_drop_action_of_pipeline",
			"Drop packet due to pipeline errors, e.g., error parsing datapath actions",
			nil, nil,
		),
		dropActionBridgeNotFoundMetric: prometheus.NewDesc("ovsdp_drop_action_bridge_not_found",
			"Drop packet due to bridge not found but, at time of translation, existing",
			nil, nil,
		),
		dropActionRecursionTooDeepMetric: prometheus.NewDesc("ovsdp_drop_action_recursion_too_deep",
			"Drop packet due to too many translations, system limit to protect from excessive time/space usage",
			nil, nil,
		),
		dropActionTooManyResubmitMetric: prometheus.NewDesc("ovsdp_drop_action_too_many_resubmit",
			"Drop packet due to too many resubmitted, system limit to protect from excessive time/space usage",
			nil, nil,
		),
		dropActionStackTooDeepMetric: prometheus.NewDesc("ovsdp_drop_action_stack_too_deep",
			"Drop packet due to the stack consuming more than 64 kB, system limit to protect from excessive time/space usage",
			nil, nil,
		),
		dropActionNoRecirculationContextMetric: prometheus.NewDesc("ovsdp_drop_action_no_recirculation_context",
			"Drop packet due to missing recirculation context",
			nil, nil,
		),
		dropActionRecirculationConflictMetric: prometheus.NewDesc("ovsdp_drop_action_recirculation_conflict",
			"Drop packet due to conflict in the recirculation",
			nil, nil,
		),
		dropActionTooManyMplsLabelsMetric: prometheus.NewDesc("ovsdp_drop_action_too_many_mpls_labels",
			"Drop packet due to MPLS pop action can't be performed as it has more labels than supported (in OVS 2.13 up to 3)",
			nil, nil,
		),
		dropActionInvalidTunnelMetadataMetric: prometheus.NewDesc("ovsdp_drop_action_invalid_tunnel_metadata",
			"Drop packet due to invalid GENEVE tunnel metadata",
			nil, nil,
		),
		dropActionUnsupportedPacketTypeMetric: prometheus.NewDesc("ovsdp_drop_action_unsupported_packet_type",
			"Drop packet due to unsupported packet type (e.g. Ethernet VLAN encapsulation)",
			nil, nil,
		),
		dropActionCongestionMetric: prometheus.NewDesc("ovsdp_drop_action_congestion",
			"Drop packet due to congestion ECN (Explicit Congestion Notification) mismatch",
			nil, nil,
		),
		dropActionForwardingDisabledMetric: prometheus.NewDesc("ovsdp_drop_action_forwarding_disabled",
			"Drop packet when forwarding for a port is disabled (e.g. when port is admin down)",
			nil, nil,
		),
		// Drop reasons new
		netdevVxlanTsoDropsMetric: prometheus.NewDesc("ovsdp_netdev_vxlan_tso_drops",
			"Drop packet due to VXLAN TSO (TCP Segmentation Offload) issues",
			nil, nil,
		),
		netdevGeneveTsoDropsMetric: prometheus.NewDesc("ovsdp_netdev_geneve_tso_drops",
			"Drop packet due to Geneve TSO (TCP Segmentation Offload) issues",
			nil, nil,
		),
		netdevPushHeaderDropsMetric: prometheus.NewDesc("ovsdp_netdev_push_header_drops",
			"Drop packet due to push header errors",
			nil, nil,
		),
		netdevSoftSegDropsMetric: prometheus.NewDesc("ovsdp_netdev_soft_seg_drops",
			"Drop packet due to soft segmentation issues",
			nil, nil,
		),
		datapathDropTunnelTsoRecircMetric: prometheus.NewDesc("ovsdp_datapath_drop_tunnel_tso_recirc",
			"Drop packet due to tunnel TSO recirculation errors",
			nil, nil,
		),
		datapathDropInvalidBondMetric: prometheus.NewDesc("ovsdp_datapath_drop_invalid_bond",
			"Drop packet due to invalid bond configuration",
			nil, nil,
		),
		datapathDropHwMissRecoverMetric: prometheus.NewDesc("ovsdp_datapath_drop_hw_miss_recover",
			"Drop packet due to hardware miss recovery failure",
			nil, nil,
		),
		// DOCA
		ovsDocaNoMarkMetric: prometheus.NewDesc("ovsdp_ovs_doca_no_mark",
			"Number of packets dropped due to missing mark in OVS-DOCA",
			nil, nil,
		),
		ovsDocaInvalidClassifyPortMetric: prometheus.NewDesc("ovsdp_ovs_doca_invalid_classify_port",
			"Number of packets dropped due to invalid classify port in OVS-DOCA",
			nil, nil,
		),
		docaQueueEmptyMetric: prometheus.NewDesc("ovsdp_doca_queue_empty",
			"Number of times an offload queue is found empty during completion operations",
			nil, nil,
		),
		docaQueueNoneProcessedMetric: prometheus.NewDesc("ovsdp_doca_queue_none_processed",
			"Number of times no entries were processed from a queue despite pending entries",
			nil, nil,
		),
		docaResizeBlockMetric: prometheus.NewDesc("ovsdp_doca_resize_block",
			"Number of times queue processing is blocked due to pipeline resizing when no entries are processed",
			nil, nil,
		),
		docaPipeResizeMetric: prometheus.NewDesc("ovsdp_doca_pipe_resize",
			"Number of times a pipe resize operation begins",
			nil, nil,
		),
		docaPipeResizeOver10MsMetric: prometheus.NewDesc("ovsdp_doca_pipe_resize_over_10_ms",
			"Number of times a pipe resize operation takes longer than 10ms",
			nil, nil,
		),
		// Upcall Flow Limit
		UpcallFlowLimitGrewMetric: prometheus.NewDesc("ovsdp_upcall_flow_limit_grew",
			"Number of times the flow_limit increased due to fast processing time",
			nil, nil,
		),
		UpcallFlowLimitHitMetric: prometheus.NewDesc("ovsdp_upcall_flow_limit_hit",
			"Number of times the flow_limit was hit during upcall processing",
			nil, nil,
		),
		UpcallFlowLimitKillMetric: prometheus.NewDesc("ovsdp_upcall_flow_limit_kill",
			"Number of times flows were killed due to exceeding flow_limit",
			nil, nil,
		),
		UpcallFlowLimitReducedMetric: prometheus.NewDesc("ovsdp_upcall_flow_limit_reduced",
			"Number of times the flow_limit was reduced due to high processing time",
			nil, nil,
		),
		UpcallFlowLimitScaledMetric: prometheus.NewDesc("ovsdp_upcall_flow_limit_scaled",
			"Number of times the flow_limit was scaled down proportionally due to very long processing time",
			nil, nil,
		),
		// DOCA Pipe Group
		docaUniqueItemTemplatesMetric: prometheus.NewDesc("ovsdp_doca_unique_item_templates",
			"Number of unique item templates created from doca-pipe-group/dump",
			nil, nil,
		),
	}
}

func (collector *ovsDPCollector) Describe(ch chan<- *prometheus.Desc) {
	// PMD stats
	ch <- collector.missWithSuccessUpcallMetric
	ch <- collector.missWithFailedUpcallMetric
	ch <- collector.processingCyclesMetric
	ch <- collector.idleCyclesMetric
	ch <- collector.avgSubtableLookupsMegaflowMetric
	ch <- collector.dropActionOfPipelineMetric
	// Memory/show stats
	ch <- collector.memoryHandlersMetric
	ch <- collector.memoryIdlCellsOpenVSwitchMetric
	ch <- collector.memoryOfconnsMetric
	ch <- collector.memoryPortsMetric
	ch <- collector.memoryRevalidatorsMetric
	ch <- collector.memoryRulesMetric
	ch <- collector.memoryUdpifKeysMetric
	// Offload stats
	ch <- collector.offloadEnqueuedMetric
	ch <- collector.offloadInsertedMetric
	ch <- collector.offloadCtUniDirConnMetric
	ch <- collector.offloadCtBiDirConnMetric
	ch <- collector.offloadCumAvgLatencyUsMetric
	ch <- collector.offloadCumLatencyStddevUsMetric
	ch <- collector.offloadCumLatencyMaxUsMetric
	ch <- collector.offloadCumLatencyMinUsMetric
	ch <- collector.offloadExpAvgLatencyUsMetric
	ch <- collector.offloadExpLatencyStddevUsMetric
	// Drop reasons
	ch <- collector.upcallDropsMetric
	ch <- collector.upcallDropsLockErrorMetric
	ch <- collector.rxDropsInvalidPacketMetric
	ch <- collector.datapathDropMeterMetric
	ch <- collector.datapathDropUserspaceActionErrorMetric
	ch <- collector.datapathDropTunnelPushErrorMetric
	ch <- collector.datapathDropTunnelPopErrorMetric
	ch <- collector.datapathDropRecircErrorMetric
	ch <- collector.datapathDropInvalidPortMetric
	ch <- collector.datapathDropInvalidTnlPortMetric
	ch <- collector.datapathDropSampleErrorMetric
	ch <- collector.datapathDropNshDecapErrorMetric
	ch <- collector.dropActionOfPipelineMetric
	ch <- collector.dropActionBridgeNotFoundMetric
	ch <- collector.dropActionRecursionTooDeepMetric
	ch <- collector.dropActionTooManyResubmitMetric
	ch <- collector.dropActionStackTooDeepMetric
	ch <- collector.dropActionNoRecirculationContextMetric
	ch <- collector.dropActionRecirculationConflictMetric
	ch <- collector.dropActionTooManyMplsLabelsMetric
	ch <- collector.dropActionInvalidTunnelMetadataMetric
	ch <- collector.dropActionUnsupportedPacketTypeMetric
	ch <- collector.dropActionCongestionMetric
	ch <- collector.dropActionForwardingDisabledMetric
	// Drop reasons new
	ch <- collector.netdevVxlanTsoDropsMetric
	ch <- collector.netdevGeneveTsoDropsMetric
	ch <- collector.netdevPushHeaderDropsMetric
	ch <- collector.netdevSoftSegDropsMetric
	ch <- collector.datapathDropTunnelTsoRecircMetric
	ch <- collector.datapathDropInvalidBondMetric
	ch <- collector.datapathDropHwMissRecoverMetric
	// DOCA
	ch <- collector.ovsDocaNoMarkMetric
	ch <- collector.ovsDocaInvalidClassifyPortMetric
	ch <- collector.docaQueueEmptyMetric
	ch <- collector.docaQueueNoneProcessedMetric
	ch <- collector.docaResizeBlockMetric
	ch <- collector.docaPipeResizeMetric
	ch <- collector.docaPipeResizeOver10MsMetric
	// Upcall Flow Limit
	ch <- collector.UpcallFlowLimitGrewMetric
	ch <- collector.UpcallFlowLimitHitMetric
	ch <- collector.UpcallFlowLimitKillMetric
	ch <- collector.UpcallFlowLimitReducedMetric
	ch <- collector.UpcallFlowLimitScaledMetric
	// DOCA Pipe Group
	ch <- collector.docaUniqueItemTemplatesMetric
}

func (collector *ovsDPCollector) Collect(ch chan<- prometheus.Metric) {
	ovsMetric, successCount := fetchOvsMetrics()
	if successCount == 0 {
		ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("ovsdp_scrape_error", "Scrape failed: all ovs-appctl commands failed", nil, nil), fmt.Errorf("all ovs-appctl commands failed"))
		return
	}
	// PMD stats
	if isValidMetric(ovsMetric.MissWithSuccessUpcall) {
		ch <- prometheus.MustNewConstMetric(collector.missWithSuccessUpcallMetric, prometheus.CounterValue, float64(ovsMetric.MissWithSuccessUpcall))
	}
	if isValidMetric(ovsMetric.MissWithFailedUpcall) {
		ch <- prometheus.MustNewConstMetric(collector.missWithFailedUpcallMetric, prometheus.CounterValue, float64(ovsMetric.MissWithFailedUpcall))
	}
	if isValidMetric(ovsMetric.ProcessingCycles) {
		ch <- prometheus.MustNewConstMetric(collector.processingCyclesMetric, prometheus.GaugeValue, float64(ovsMetric.ProcessingCycles))
	}
	if isValidMetric(ovsMetric.IdleCycles) {
		ch <- prometheus.MustNewConstMetric(collector.idleCyclesMetric, prometheus.GaugeValue, float64(ovsMetric.IdleCycles))
	}
	if isValidMetric(ovsMetric.AvgSubtableLookupsMegaflow) {
		ch <- prometheus.MustNewConstMetric(collector.avgSubtableLookupsMegaflowMetric, prometheus.CounterValue, float64(ovsMetric.AvgSubtableLookupsMegaflow))
	}
	// Memory/show stats
	if isValidMetric(ovsMetric.MemoryHandlers) {
		ch <- prometheus.MustNewConstMetric(collector.memoryHandlersMetric, prometheus.GaugeValue, float64(ovsMetric.MemoryHandlers))
	}
	if isValidMetric(ovsMetric.MemoryIdlCellsOpenVSwitch) {
		ch <- prometheus.MustNewConstMetric(collector.memoryIdlCellsOpenVSwitchMetric, prometheus.GaugeValue, float64(ovsMetric.MemoryIdlCellsOpenVSwitch))
	}
	if isValidMetric(ovsMetric.MemoryOfconns) {
		ch <- prometheus.MustNewConstMetric(collector.memoryOfconnsMetric, prometheus.GaugeValue, float64(ovsMetric.MemoryOfconns))
	}
	if isValidMetric(ovsMetric.MemoryPorts) {
		ch <- prometheus.MustNewConstMetric(collector.memoryPortsMetric, prometheus.GaugeValue, float64(ovsMetric.MemoryPorts))
	}
	if isValidMetric(ovsMetric.MemoryRevalidators) {
		ch <- prometheus.MustNewConstMetric(collector.memoryRevalidatorsMetric, prometheus.GaugeValue, float64(ovsMetric.MemoryRevalidators))
	}
	if isValidMetric(ovsMetric.MemoryRules) {
		ch <- prometheus.MustNewConstMetric(collector.memoryRulesMetric, prometheus.GaugeValue, float64(ovsMetric.MemoryRules))
	}
	if isValidMetric(ovsMetric.MemoryUdpifKeys) {
		ch <- prometheus.MustNewConstMetric(collector.memoryUdpifKeysMetric, prometheus.GaugeValue, float64(ovsMetric.MemoryUdpifKeys))
	}
	// Offload stats
	if isValidMetric(ovsMetric.OffloadEnqueued) {
		ch <- prometheus.MustNewConstMetric(collector.offloadEnqueuedMetric, prometheus.CounterValue, float64(ovsMetric.OffloadEnqueued))
	}
	if isValidMetric(ovsMetric.OffloadInserted) {
		ch <- prometheus.MustNewConstMetric(collector.offloadInsertedMetric, prometheus.CounterValue, float64(ovsMetric.OffloadInserted))
	}
	if isValidMetric(ovsMetric.OffloadCtUniDirConnections) {
		ch <- prometheus.MustNewConstMetric(collector.offloadCtUniDirConnMetric, prometheus.CounterValue, float64(ovsMetric.OffloadCtUniDirConnections))
	}
	if isValidMetric(ovsMetric.OffloadCtBiDirConnections) {
		ch <- prometheus.MustNewConstMetric(collector.offloadCtBiDirConnMetric, prometheus.CounterValue, float64(ovsMetric.OffloadCtBiDirConnections))
	}
	if isValidMetric(ovsMetric.OffloadCumAvgLatencyUs) {
		ch <- prometheus.MustNewConstMetric(collector.offloadCumAvgLatencyUsMetric, prometheus.GaugeValue, float64(ovsMetric.OffloadCumAvgLatencyUs))
	}
	if isValidMetric(ovsMetric.OffloadCumLatencyStddevUs) {
		ch <- prometheus.MustNewConstMetric(collector.offloadCumLatencyStddevUsMetric, prometheus.GaugeValue, float64(ovsMetric.OffloadCumLatencyStddevUs))
	}
	if isValidMetric(ovsMetric.OffloadCumLatencyMaxUs) {
		ch <- prometheus.MustNewConstMetric(collector.offloadCumLatencyMaxUsMetric, prometheus.GaugeValue, float64(ovsMetric.OffloadCumLatencyMaxUs))
	}
	if isValidMetric(ovsMetric.OffloadCumLatencyMinUs) {
		ch <- prometheus.MustNewConstMetric(collector.offloadCumLatencyMinUsMetric, prometheus.GaugeValue, float64(ovsMetric.OffloadCumLatencyMinUs))
	}
	if isValidMetric(ovsMetric.OffloadExpAvgLatencyUs) {
		ch <- prometheus.MustNewConstMetric(collector.offloadExpAvgLatencyUsMetric, prometheus.GaugeValue, float64(ovsMetric.OffloadExpAvgLatencyUs))
	}
	if isValidMetric(ovsMetric.OffloadExpLatencyStddevUs) {
		ch <- prometheus.MustNewConstMetric(collector.offloadExpLatencyStddevUsMetric, prometheus.GaugeValue, float64(ovsMetric.OffloadExpLatencyStddevUs))
	}
	// Drop reasons
	if isValidMetric(ovsMetric.UpcallDrops) {
		ch <- prometheus.MustNewConstMetric(collector.upcallDropsMetric, prometheus.CounterValue, float64(ovsMetric.UpcallDrops))
	}
	if isValidMetric(ovsMetric.UpcallDropsLockError) {
		ch <- prometheus.MustNewConstMetric(collector.upcallDropsLockErrorMetric, prometheus.CounterValue, float64(ovsMetric.UpcallDropsLockError))
	}
	if isValidMetric(ovsMetric.RxDropsInvalidPacket) {
		ch <- prometheus.MustNewConstMetric(collector.rxDropsInvalidPacketMetric, prometheus.CounterValue, float64(ovsMetric.RxDropsInvalidPacket))
	}
	if isValidMetric(ovsMetric.DatapathDropMeter) {
		ch <- prometheus.MustNewConstMetric(collector.datapathDropMeterMetric, prometheus.CounterValue, float64(ovsMetric.DatapathDropMeter))
	}
	if isValidMetric(ovsMetric.DatapathDropUserspaceActionError) {
		ch <- prometheus.MustNewConstMetric(collector.datapathDropUserspaceActionErrorMetric, prometheus.CounterValue, float64(ovsMetric.DatapathDropUserspaceActionError))
	}
	if isValidMetric(ovsMetric.DatapathDropTunnelPushError) {
		ch <- prometheus.MustNewConstMetric(collector.datapathDropTunnelPushErrorMetric, prometheus.CounterValue, float64(ovsMetric.DatapathDropTunnelPushError))
	}
	if isValidMetric(ovsMetric.DatapathDropTunnelPopError) {
		ch <- prometheus.MustNewConstMetric(collector.datapathDropTunnelPopErrorMetric, prometheus.CounterValue, float64(ovsMetric.DatapathDropTunnelPopError))
	}
	if isValidMetric(ovsMetric.DatapathDropRecircError) {
		ch <- prometheus.MustNewConstMetric(collector.datapathDropRecircErrorMetric, prometheus.CounterValue, float64(ovsMetric.DatapathDropRecircError))
	}
	if isValidMetric(ovsMetric.DatapathDropInvalidPort) {
		ch <- prometheus.MustNewConstMetric(collector.datapathDropInvalidPortMetric, prometheus.CounterValue, float64(ovsMetric.DatapathDropInvalidPort))
	}
	if isValidMetric(ovsMetric.DatapathDropInvalidTnlPort) {
		ch <- prometheus.MustNewConstMetric(collector.datapathDropInvalidTnlPortMetric, prometheus.CounterValue, float64(ovsMetric.DatapathDropInvalidTnlPort))
	}
	if isValidMetric(ovsMetric.DatapathDropSampleError) {
		ch <- prometheus.MustNewConstMetric(collector.datapathDropSampleErrorMetric, prometheus.CounterValue, float64(ovsMetric.DatapathDropSampleError))
	}
	if isValidMetric(ovsMetric.DatapathDropNshDecapError) {
		ch <- prometheus.MustNewConstMetric(collector.datapathDropNshDecapErrorMetric, prometheus.CounterValue, float64(ovsMetric.DatapathDropNshDecapError))
	}
	if isValidMetric(ovsMetric.DropActionOfPipeline) {
		ch <- prometheus.MustNewConstMetric(collector.dropActionOfPipelineMetric, prometheus.CounterValue, float64(ovsMetric.DropActionOfPipeline))
	}
	if isValidMetric(ovsMetric.DropActionBridgeNotFound) {
		ch <- prometheus.MustNewConstMetric(collector.dropActionBridgeNotFoundMetric, prometheus.CounterValue, float64(ovsMetric.DropActionBridgeNotFound))
	}
	if isValidMetric(ovsMetric.DropActionRecursionTooDeep) {
		ch <- prometheus.MustNewConstMetric(collector.dropActionRecursionTooDeepMetric, prometheus.CounterValue, float64(ovsMetric.DropActionRecursionTooDeep))
	}
	if isValidMetric(ovsMetric.DropActionTooManyResubmit) {
		ch <- prometheus.MustNewConstMetric(collector.dropActionTooManyResubmitMetric, prometheus.CounterValue, float64(ovsMetric.DropActionTooManyResubmit))
	}
	if isValidMetric(ovsMetric.DropActionStackTooDeep) {
		ch <- prometheus.MustNewConstMetric(collector.dropActionStackTooDeepMetric, prometheus.CounterValue, float64(ovsMetric.DropActionStackTooDeep))
	}
	if isValidMetric(ovsMetric.DropActionNoRecirculationContext) {
		ch <- prometheus.MustNewConstMetric(collector.dropActionNoRecirculationContextMetric, prometheus.CounterValue, float64(ovsMetric.DropActionNoRecirculationContext))
	}
	if isValidMetric(ovsMetric.DropActionRecirculationConflict) {
		ch <- prometheus.MustNewConstMetric(collector.dropActionRecirculationConflictMetric, prometheus.CounterValue, float64(ovsMetric.DropActionRecirculationConflict))
	}
	if isValidMetric(ovsMetric.DropActionTooManyMplsLabels) {
		ch <- prometheus.MustNewConstMetric(collector.dropActionTooManyMplsLabelsMetric, prometheus.CounterValue, float64(ovsMetric.DropActionTooManyMplsLabels))
	}
	if isValidMetric(ovsMetric.DropActionInvalidTunnelMetadata) {
		ch <- prometheus.MustNewConstMetric(collector.dropActionInvalidTunnelMetadataMetric, prometheus.CounterValue, float64(ovsMetric.DropActionInvalidTunnelMetadata))
	}
	if isValidMetric(ovsMetric.DropActionUnsupportedPacketType) {
		ch <- prometheus.MustNewConstMetric(collector.dropActionUnsupportedPacketTypeMetric, prometheus.CounterValue, float64(ovsMetric.DropActionUnsupportedPacketType))
	}
	if isValidMetric(ovsMetric.DropActionCongestion) {
		ch <- prometheus.MustNewConstMetric(collector.dropActionCongestionMetric, prometheus.CounterValue, float64(ovsMetric.DropActionCongestion))
	}
	if isValidMetric(ovsMetric.DropActionForwardingDisabled) {
		ch <- prometheus.MustNewConstMetric(collector.dropActionForwardingDisabledMetric, prometheus.CounterValue, float64(ovsMetric.DropActionForwardingDisabled))
	}
	// Drop reasons new
	if isValidMetric(ovsMetric.NetdevVxlanTsoDrops) {
		ch <- prometheus.MustNewConstMetric(collector.netdevVxlanTsoDropsMetric, prometheus.CounterValue, float64(ovsMetric.NetdevVxlanTsoDrops))
	}
	if isValidMetric(ovsMetric.NetdevGeneveTsoDrops) {
		ch <- prometheus.MustNewConstMetric(collector.netdevGeneveTsoDropsMetric, prometheus.CounterValue, float64(ovsMetric.NetdevGeneveTsoDrops))
	}
	if isValidMetric(ovsMetric.NetdevPushHeaderDrops) {
		ch <- prometheus.MustNewConstMetric(collector.netdevPushHeaderDropsMetric, prometheus.CounterValue, float64(ovsMetric.NetdevPushHeaderDrops))
	}
	if isValidMetric(ovsMetric.NetdevSoftSegDrops) {
		ch <- prometheus.MustNewConstMetric(collector.netdevSoftSegDropsMetric, prometheus.CounterValue, float64(ovsMetric.NetdevSoftSegDrops))
	}
	if isValidMetric(ovsMetric.DatapathDropTunnelTsoRecirc) {
		ch <- prometheus.MustNewConstMetric(collector.datapathDropTunnelTsoRecircMetric, prometheus.CounterValue, float64(ovsMetric.DatapathDropTunnelTsoRecirc))
	}
	if isValidMetric(ovsMetric.DatapathDropInvalidBond) {
		ch <- prometheus.MustNewConstMetric(collector.datapathDropInvalidBondMetric, prometheus.CounterValue, float64(ovsMetric.DatapathDropInvalidBond))
	}
	if isValidMetric(ovsMetric.DatapathDropHwMissRecover) {
		ch <- prometheus.MustNewConstMetric(collector.datapathDropHwMissRecoverMetric, prometheus.CounterValue, float64(ovsMetric.DatapathDropHwMissRecover))
	}
	// DOCA
	if isValidMetric(ovsMetric.OvsDocaNoMark) {
		ch <- prometheus.MustNewConstMetric(collector.ovsDocaNoMarkMetric, prometheus.CounterValue, float64(ovsMetric.OvsDocaNoMark))
	}
	if isValidMetric(ovsMetric.OvsDocaInvalidClassifyPort) {
		ch <- prometheus.MustNewConstMetric(collector.ovsDocaInvalidClassifyPortMetric, prometheus.CounterValue, float64(ovsMetric.OvsDocaInvalidClassifyPort))
	}
	if isValidMetric(ovsMetric.DocaQueueEmpty) {
		ch <- prometheus.MustNewConstMetric(collector.docaQueueEmptyMetric, prometheus.CounterValue, float64(ovsMetric.DocaQueueEmpty))
	}
	if isValidMetric(ovsMetric.DocaQueueNoneProcessed) {
		ch <- prometheus.MustNewConstMetric(collector.docaQueueNoneProcessedMetric, prometheus.CounterValue, float64(ovsMetric.DocaQueueNoneProcessed))
	}
	if isValidMetric(ovsMetric.DocaResizeBlock) {
		ch <- prometheus.MustNewConstMetric(collector.docaResizeBlockMetric, prometheus.CounterValue, float64(ovsMetric.DocaResizeBlock))
	}
	if isValidMetric(ovsMetric.DocaPipeResize) {
		ch <- prometheus.MustNewConstMetric(collector.docaPipeResizeMetric, prometheus.CounterValue, float64(ovsMetric.DocaPipeResize))
	}
	if isValidMetric(ovsMetric.DocaPipeResizeOver10Ms) {
		ch <- prometheus.MustNewConstMetric(collector.docaPipeResizeOver10MsMetric, prometheus.CounterValue, float64(ovsMetric.DocaPipeResizeOver10Ms))
	}
	// Upcall Flow Limit
	if isValidMetric(ovsMetric.UpcallFlowLimitGrew) {
		ch <- prometheus.MustNewConstMetric(collector.UpcallFlowLimitGrewMetric, prometheus.CounterValue, float64(ovsMetric.UpcallFlowLimitGrew))
	}
	if isValidMetric(ovsMetric.UpcallFlowLimitHit) {
		ch <- prometheus.MustNewConstMetric(collector.UpcallFlowLimitHitMetric, prometheus.CounterValue, float64(ovsMetric.UpcallFlowLimitHit))
	}
	if isValidMetric(ovsMetric.UpcallFlowLimitKill) {
		ch <- prometheus.MustNewConstMetric(collector.UpcallFlowLimitKillMetric, prometheus.CounterValue, float64(ovsMetric.UpcallFlowLimitKill))
	}
	if isValidMetric(ovsMetric.UpcallFlowLimitReduced) {
		ch <- prometheus.MustNewConstMetric(collector.UpcallFlowLimitReducedMetric, prometheus.CounterValue, float64(ovsMetric.UpcallFlowLimitReduced))
	}
	if isValidMetric(ovsMetric.UpcallFlowLimitScaled) {
		ch <- prometheus.MustNewConstMetric(collector.UpcallFlowLimitScaledMetric, prometheus.CounterValue, float64(ovsMetric.UpcallFlowLimitScaled))
	}
	// DOCA Pipe Group
	if isValidMetric(ovsMetric.DocaUniqueItemTemplates) {
		ch <- prometheus.MustNewConstMetric(collector.docaUniqueItemTemplatesMetric, prometheus.GaugeValue, float64(ovsMetric.DocaUniqueItemTemplates))
	}

	// Collect parsed metrics from ovs-appctl metrics/show
	for _, mf := range ovsMetric.ParsedMetrics {
		for _, m := range mf.Metric {
			labels := make([]string, 0, len(m.Label))
			labelValues := make([]string, 0, len(m.Label))
			for _, label := range m.Label {
				labels = append(labels, *label.Name)
				labelValues = append(labelValues, *label.Value)
			}

			desc := prometheus.NewDesc(*mf.Name, *mf.Help, labels, nil)

			var value float64
			var valueType prometheus.ValueType

			switch {
			case m.Counter != nil:
				value = *m.Counter.Value
				valueType = prometheus.CounterValue
			case m.Gauge != nil:
				value = *m.Gauge.Value
				valueType = prometheus.GaugeValue
			case m.Histogram != nil:
				value = float64(*m.Histogram.SampleCount)
				valueType = prometheus.CounterValue
			case m.Summary != nil:
				value = float64(*m.Summary.SampleCount)
				valueType = prometheus.CounterValue
			default:
				continue // Skip unknown metric types
			}

			ch <- prometheus.MustNewConstMetric(desc, valueType, value, labelValues...)
		}
	}
}

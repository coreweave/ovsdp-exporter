package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

type OvsMetric struct {
	// PMD stats
	MissWithSuccessUpcall      float64
	MissWithFailedUpcall       float64
	AvgSubtableLookupsMegaflow float64
	ProcessingCycles           float64
	IdleCycles                 float64
	// Drop reasons
	UpcallDrops                      float64
	UpcallDropsLockError             float64
	RxDropsInvalidPacket             float64
	DatapathDropMeter                float64
	DatapathDropUserspaceActionError float64
	DatapathDropTunnelPushError      float64
	DatapathDropTunnelPopError       float64
	DatapathDropRecircError          float64
	DatapathDropInvalidPort          float64
	DatapathDropInvalidTnlPort       float64
	DatapathDropSampleError          float64
	DatapathDropNshDecapError        float64
	DropActionOfPipeline             float64
	DropActionBridgeNotFound         float64
	DropActionRecursionTooDeep       float64
	DropActionTooManyResubmit        float64
	DropActionStackTooDeep           float64
	DropActionNoRecirculationContext float64
	DropActionRecirculationConflict  float64
	DropActionTooManyMplsLabels      float64
	DropActionInvalidTunnelMetadata  float64
	DropActionUnsupportedPacketType  float64
	DropActionCongestion             float64
	DropActionForwardingDisabled     float64
	// Drop reasons new
	NetdevVxlanTsoDrops         float64
	NetdevGeneveTsoDrops        float64
	NetdevPushHeaderDrops       float64
	NetdevSoftSegDrops          float64
	DatapathDropTunnelTsoRecirc float64
	DatapathDropInvalidBond     float64
	DatapathDropHwMissRecover   float64
	// DOCA
	OvsDocaNoMark              float64
	OvsDocaInvalidClassifyPort float64
	DocaQueueEmpty             float64
	DocaQueueNoneProcessed     float64
	DocaResizeBlock            float64
	DocaPipeResize             float64
	DocaPipeResizeOver10Ms     float64
}

func getOvsMetric() *OvsMetric {
	var ovsMetric OvsMetric

	cmd := exec.Command("ovs-appctl", "dpif-netdev/pmd-stats-show")
	pmdStatsOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
	} else {
		parsePMDStats(&ovsMetric, string(pmdStatsOutput))
	}

	cmd = exec.Command("ovs-appctl", "coverage/show")
	coverageOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
	} else {
		parseCoverageDropReasons(&ovsMetric, string(coverageOutput))
		parseCoverageDoca(&ovsMetric, string(coverageOutput))
	}

	return &ovsMetric
}

func parseCoverageDoca(metrics *OvsMetric, coverageStats string) {
	// DOCA
	ovsDocaNoMarkRegexp := regexp.MustCompile(`(?m)^[ \t]*ovs_doca_no_mark.*total:\s*(\d+)`)
	ovsDocaNoMarkMatch := ovsDocaNoMarkRegexp.FindStringSubmatch(coverageStats)
	metrics.OvsDocaNoMark = -1
	if len(ovsDocaNoMarkMatch) > 1 {
		v, err := strconv.ParseFloat(ovsDocaNoMarkMatch[1], 64)
		if err == nil {
			metrics.OvsDocaNoMark = v
		}
	}

	ovsDocaInvalidClassifyPortRegexp := regexp.MustCompile(`(?m)^[ \t]*ovs_doca_invalid_classify_port.*total:\s*(\d+)`)
	ovsDocaInvalidClassifyPortMatch := ovsDocaInvalidClassifyPortRegexp.FindStringSubmatch(coverageStats)
	metrics.OvsDocaInvalidClassifyPort = -1
	if len(ovsDocaInvalidClassifyPortMatch) > 1 {
		v, err := strconv.ParseFloat(ovsDocaInvalidClassifyPortMatch[1], 64)
		if err == nil {
			metrics.OvsDocaInvalidClassifyPort = v
		}
	}

	docaQueueEmptyRegexp := regexp.MustCompile(`(?m)^[ \t]*doca_queue_empty.*total:\s*(\d+)`)
	docaQueueEmptyMatch := docaQueueEmptyRegexp.FindStringSubmatch(coverageStats)
	metrics.DocaQueueEmpty = -1
	if len(docaQueueEmptyMatch) > 1 {
		v, err := strconv.ParseFloat(docaQueueEmptyMatch[1], 64)
		if err == nil {
			metrics.DocaQueueEmpty = v
		}
	}

	docaQueueNoneProcessedRegexp := regexp.MustCompile(`(?m)^[ \t]*doca_queue_none_processed.*total:\s*(\d+)`)
	docaQueueNoneProcessedMatch := docaQueueNoneProcessedRegexp.FindStringSubmatch(coverageStats)
	metrics.DocaQueueNoneProcessed = -1
	if len(docaQueueNoneProcessedMatch) > 1 {
		v, err := strconv.ParseFloat(docaQueueNoneProcessedMatch[1], 64)
		if err == nil {
			metrics.DocaQueueNoneProcessed = v
		}
	}

	docaResizeBlockRegexp := regexp.MustCompile(`(?m)^[ \t]*doca_resize_block.*total:\s*(\d+)`)
	docaResizeBlockMatch := docaResizeBlockRegexp.FindStringSubmatch(coverageStats)
	metrics.DocaResizeBlock = -1
	if len(docaResizeBlockMatch) > 1 {
		v, err := strconv.ParseFloat(docaResizeBlockMatch[1], 64)
		if err == nil {
			metrics.DocaResizeBlock = v
		}
	}

	docaPipeResizeRegexp := regexp.MustCompile(`(?m)^[ \t]*doca_pipe_resize.*total:\s*(\d+)`)
	docaPipeResizeMatch := docaPipeResizeRegexp.FindStringSubmatch(coverageStats)
	metrics.DocaPipeResize = -1
	if len(docaPipeResizeMatch) > 1 {
		v, err := strconv.ParseFloat(docaPipeResizeMatch[1], 64)
		if err == nil {
			metrics.DocaPipeResize = v
		}
	}

	docaPipeResizeOver10MsRegexp := regexp.MustCompile(`(?m)^[ \t]*doca_pipe_resize_over_10_ms.*total:\s*(\d+)`)
	docaPipeResizeOver10MsMatch := docaPipeResizeOver10MsRegexp.FindStringSubmatch(coverageStats)
	metrics.DocaPipeResizeOver10Ms = -1
	if len(docaPipeResizeOver10MsMatch) > 1 {
		v, err := strconv.ParseFloat(docaPipeResizeOver10MsMatch[1], 64)
		if err == nil {
			metrics.DocaPipeResizeOver10Ms = v
		}
	}
}

func parseCoverageDropReasons(metrics *OvsMetric, coverageStats string) {
	// Drop reasons
	// Upcall drops
	upcallDropsRegexp := regexp.MustCompile(`(?m)^[ \t]*datapath_drop_upcall_error.*total:\s*(\d+)`)
	upcallDropsMatch := upcallDropsRegexp.FindStringSubmatch(coverageStats)
	metrics.UpcallDrops = -1
	if len(upcallDropsMatch) > 1 {
		v, err := strconv.ParseFloat(upcallDropsMatch[1], 64)
		if err == nil {
			metrics.UpcallDrops = v
		}
	}

	// Upcall drops lock error
	upcallDropsLockErrorRegexp := regexp.MustCompile(`(?m)^[ \t]*datapath_drop_lock_error.*total:\s*(\d+)`)
	upcallDropsLockErrorMatch := upcallDropsLockErrorRegexp.FindStringSubmatch(coverageStats)
	metrics.UpcallDropsLockError = -1
	if len(upcallDropsLockErrorMatch) > 1 {
		v, err := strconv.ParseFloat(upcallDropsLockErrorMatch[1], 64)
		if err == nil {
			metrics.UpcallDropsLockError = v
		}
	}

	// RX drops invalid packet
	rxDropsInvalidPacketRegexp := regexp.MustCompile(`(?m)^[ \t]*datapath_drop_rx_invalid_packet.*total:\s*(\d+)`)
	rxDropsInvalidPacketMatch := rxDropsInvalidPacketRegexp.FindStringSubmatch(coverageStats)
	metrics.RxDropsInvalidPacket = -1
	if len(rxDropsInvalidPacketMatch) > 1 {
		v, err := strconv.ParseFloat(rxDropsInvalidPacketMatch[1], 64)
		if err == nil {
			metrics.RxDropsInvalidPacket = v
		}
	}

	// Datapath drop meter
	datapathDropMeterRegexp := regexp.MustCompile(`(?m)^[ \t]*datapath_drop_meter.*total:\s*(\d+)`)
	datapathDropMeterMatch := datapathDropMeterRegexp.FindStringSubmatch(coverageStats)
	metrics.DatapathDropMeter = -1
	if len(datapathDropMeterMatch) > 1 {
		v, err := strconv.ParseFloat(datapathDropMeterMatch[1], 64)
		if err == nil {
			metrics.DatapathDropMeter = v
		}
	}

	// Datapath drop userspace action error
	datapathDropUserspaceActionErrorRegexp := regexp.MustCompile(`(?m)^[ \t]*datapath_drop_userspace_action_error.*total:\s*(\d+)`)
	datapathDropUserspaceActionErrorMatch := datapathDropUserspaceActionErrorRegexp.FindStringSubmatch(coverageStats)
	metrics.DatapathDropUserspaceActionError = -1
	if len(datapathDropUserspaceActionErrorMatch) > 1 {
		v, err := strconv.ParseFloat(datapathDropUserspaceActionErrorMatch[1], 64)
		if err == nil {
			metrics.DatapathDropUserspaceActionError = v
		}
	}

	// Datapath drop tunnel push error
	datapathDropTunnelPushErrorRegexp := regexp.MustCompile(`(?m)^[ \t]*datapath_drop_tunnel_push_error.*total:\s*(\d+)`)
	datapathDropTunnelPushErrorMatch := datapathDropTunnelPushErrorRegexp.FindStringSubmatch(coverageStats)
	metrics.DatapathDropTunnelPushError = -1
	if len(datapathDropTunnelPushErrorMatch) > 1 {
		v, err := strconv.ParseFloat(datapathDropTunnelPushErrorMatch[1], 64)
		if err == nil {
			metrics.DatapathDropTunnelPushError = v
		}
	}

	// Datapath drop tunnel pop error
	datapathDropTunnelPopErrorRegexp := regexp.MustCompile(`(?m)^[ \t]*datapath_drop_tunnel_pop_error.*total:\s*(\d+)`)
	datapathDropTunnelPopErrorMatch := datapathDropTunnelPopErrorRegexp.FindStringSubmatch(coverageStats)
	metrics.DatapathDropTunnelPopError = -1
	if len(datapathDropTunnelPopErrorMatch) > 1 {
		v, err := strconv.ParseFloat(datapathDropTunnelPopErrorMatch[1], 64)
		if err == nil {
			metrics.DatapathDropTunnelPopError = v
		}
	}

	// Datapath drop recirc error
	datapathDropRecircErrorRegexp := regexp.MustCompile(`(?m)^[ \t]*datapath_drop_recirc_error.*total:\s*(\d+)`)
	datapathDropRecircErrorMatch := datapathDropRecircErrorRegexp.FindStringSubmatch(coverageStats)
	metrics.DatapathDropRecircError = -1
	if len(datapathDropRecircErrorMatch) > 1 {
		v, err := strconv.ParseFloat(datapathDropRecircErrorMatch[1], 64)
		if err == nil {
			metrics.DatapathDropRecircError = v
		}
	}

	// Datapath drop invalid port
	datapathDropInvalidPortRegexp := regexp.MustCompile(`(?m)^[ \t]*datapath_drop_invalid_port.*total:\s*(\d+)`)
	datapathDropInvalidPortMatch := datapathDropInvalidPortRegexp.FindStringSubmatch(coverageStats)
	metrics.DatapathDropInvalidPort = -1
	if len(datapathDropInvalidPortMatch) > 1 {
		v, err := strconv.ParseFloat(datapathDropInvalidPortMatch[1], 64)
		if err == nil {
			metrics.DatapathDropInvalidPort = v
		}
	}

	// Datapath drop invalid tunnel port
	datapathDropInvalidTnlPortRegexp := regexp.MustCompile(`(?m)^[ \t]*datapath_drop_invalid_tnl_port.*total:\s*(\d+)`)
	datapathDropInvalidTnlPortMatch := datapathDropInvalidTnlPortRegexp.FindStringSubmatch(coverageStats)
	metrics.DatapathDropInvalidTnlPort = -1
	if len(datapathDropInvalidTnlPortMatch) > 1 {
		v, err := strconv.ParseFloat(datapathDropInvalidTnlPortMatch[1], 64)
		if err == nil {
			metrics.DatapathDropInvalidTnlPort = v
		}
	}

	// Datapath drop sample error
	datapathDropSampleErrorRegexp := regexp.MustCompile(`(?m)^[ \t]*datapath_drop_sample_error.*total:\s*(\d+)`)
	datapathDropSampleErrorMatch := datapathDropSampleErrorRegexp.FindStringSubmatch(coverageStats)
	metrics.DatapathDropSampleError = -1
	if len(datapathDropSampleErrorMatch) > 1 {
		v, err := strconv.ParseFloat(datapathDropSampleErrorMatch[1], 64)
		if err == nil {
			metrics.DatapathDropSampleError = v
		}
	}

	// Datapath drop NSH decap error
	datapathDropNshDecapErrorRegexp := regexp.MustCompile(`(?m)^[ \t]*datapath_drop_nsh_decap_error.*total:\s*(\d+)`)
	datapathDropNshDecapErrorMatch := datapathDropNshDecapErrorRegexp.FindStringSubmatch(coverageStats)
	metrics.DatapathDropNshDecapError = -1
	if len(datapathDropNshDecapErrorMatch) > 1 {
		v, err := strconv.ParseFloat(datapathDropNshDecapErrorMatch[1], 64)
		if err == nil {
			metrics.DatapathDropNshDecapError = v
		}
	}

	// Drop action of pipeline
	dropActionOfPipelineRegexp := regexp.MustCompile(`(?m)^[ \t]*drop_action_of_pipeline.*total:\s*(\d+)`)
	dropActionOfPipelineMatch := dropActionOfPipelineRegexp.FindStringSubmatch(coverageStats)
	metrics.DropActionOfPipeline = -1
	if len(dropActionOfPipelineMatch) > 1 {
		v, err := strconv.ParseFloat(dropActionOfPipelineMatch[1], 64)
		if err == nil {
			metrics.DropActionOfPipeline = v
		}
	}

	// Drop action bridge not found
	dropActionBridgeNotFoundRegexp := regexp.MustCompile(`(?m)^[ \t]*drop_action_bridge_not_found.*total:\s*(\d+)`)
	dropActionBridgeNotFoundMatch := dropActionBridgeNotFoundRegexp.FindStringSubmatch(coverageStats)
	metrics.DropActionBridgeNotFound = -1
	if len(dropActionBridgeNotFoundMatch) > 1 {
		v, err := strconv.ParseFloat(dropActionBridgeNotFoundMatch[1], 64)
		if err == nil {
			metrics.DropActionBridgeNotFound = v
		}
	}

	// Drop action recursion too deep
	dropActionRecursionTooDeepRegexp := regexp.MustCompile(`(?m)^[ \t]*drop_action_recursion_too_deep.*total:\s*(\d+)`)
	dropActionRecursionTooDeepMatch := dropActionRecursionTooDeepRegexp.FindStringSubmatch(coverageStats)
	metrics.DropActionRecursionTooDeep = -1
	if len(dropActionRecursionTooDeepMatch) > 1 {
		v, err := strconv.ParseFloat(dropActionRecursionTooDeepMatch[1], 64)
		if err == nil {
			metrics.DropActionRecursionTooDeep = v
		}
	}

	// Drop action too many resubmit
	dropActionTooManyResubmitRegexp := regexp.MustCompile(`(?m)^[ \t]*drop_action_too_many_resubmit.*total:\s*(\d+)`)
	dropActionTooManyResubmitMatch := dropActionTooManyResubmitRegexp.FindStringSubmatch(coverageStats)
	metrics.DropActionTooManyResubmit = -1
	if len(dropActionTooManyResubmitMatch) > 1 {
		v, err := strconv.ParseFloat(dropActionTooManyResubmitMatch[1], 64)
		if err == nil {
			metrics.DropActionTooManyResubmit = v
		}
	}

	// Drop action stack too deep
	dropActionStackTooDeepRegexp := regexp.MustCompile(`(?m)^[ \t]*drop_action_stack_too_deep.*total:\s*(\d+)`)
	dropActionStackTooDeepMatch := dropActionStackTooDeepRegexp.FindStringSubmatch(coverageStats)
	metrics.DropActionStackTooDeep = -1
	if len(dropActionStackTooDeepMatch) > 1 {
		v, err := strconv.ParseFloat(dropActionStackTooDeepMatch[1], 64)
		if err == nil {
			metrics.DropActionStackTooDeep = v
		}
	}

	// Drop action no recirculation context
	dropActionNoRecirculationContextRegexp := regexp.MustCompile(`(?m)^[ \t]*drop_action_no_recirculation_context.*total:\s*(\d+)`)
	dropActionNoRecirculationContextMatch := dropActionNoRecirculationContextRegexp.FindStringSubmatch(coverageStats)
	metrics.DropActionNoRecirculationContext = -1
	if len(dropActionNoRecirculationContextMatch) > 1 {
		v, err := strconv.ParseFloat(dropActionNoRecirculationContextMatch[1], 64)
		if err == nil {
			metrics.DropActionNoRecirculationContext = v
		}
	}

	// Drop action recirculation conflict
	dropActionRecirculationConflictRegexp := regexp.MustCompile(`(?m)^[ \t]*drop_action_recirculation_conflict.*total:\s*(\d+)`)
	dropActionRecirculationConflictMatch := dropActionRecirculationConflictRegexp.FindStringSubmatch(coverageStats)
	metrics.DropActionRecirculationConflict = -1
	if len(dropActionRecirculationConflictMatch) > 1 {
		v, err := strconv.ParseFloat(dropActionRecirculationConflictMatch[1], 64)
		if err == nil {
			metrics.DropActionRecirculationConflict = v
		}
	}

	// Drop action too many MPLS labels
	dropActionTooManyMplsLabelsRegexp := regexp.MustCompile(`(?m)^[ \t]*drop_action_too_many_mpls_labels.*total:\s*(\d+)`)
	dropActionTooManyMplsLabelsMatch := dropActionTooManyMplsLabelsRegexp.FindStringSubmatch(coverageStats)
	metrics.DropActionTooManyMplsLabels = -1
	if len(dropActionTooManyMplsLabelsMatch) > 1 {
		v, err := strconv.ParseFloat(dropActionTooManyMplsLabelsMatch[1], 64)
		if err == nil {
			metrics.DropActionTooManyMplsLabels = v
		}
	}

	// Drop action invalid tunnel metadata
	dropActionInvalidTunnelMetadataRegexp := regexp.MustCompile(`(?m)^[ \t]*drop_action_invalid_tunnel_metadata.*total:\s*(\d+)`)
	dropActionInvalidTunnelMetadataMatch := dropActionInvalidTunnelMetadataRegexp.FindStringSubmatch(coverageStats)
	metrics.DropActionInvalidTunnelMetadata = -1
	if len(dropActionInvalidTunnelMetadataMatch) > 1 {
		v, err := strconv.ParseFloat(dropActionInvalidTunnelMetadataMatch[1], 64)
		if err == nil {
			metrics.DropActionInvalidTunnelMetadata = v
		}
	}

	// Drop action unsupported packet type
	dropActionUnsupportedPacketTypeRegexp := regexp.MustCompile(`(?m)^[ \t]*drop_action_unsupported_packet_type.*total:\s*(\d+)`)
	dropActionUnsupportedPacketTypeMatch := dropActionUnsupportedPacketTypeRegexp.FindStringSubmatch(coverageStats)
	metrics.DropActionUnsupportedPacketType = -1
	if len(dropActionUnsupportedPacketTypeMatch) > 1 {
		v, err := strconv.ParseFloat(dropActionUnsupportedPacketTypeMatch[1], 64)
		if err == nil {
			metrics.DropActionUnsupportedPacketType = v
		}
	}

	// Drop action congestion
	dropActionCongestionRegexp := regexp.MustCompile(`(?m)^[ \t]*drop_action_congestion.*total:\s*(\d+)`)
	dropActionCongestionMatch := dropActionCongestionRegexp.FindStringSubmatch(coverageStats)
	metrics.DropActionCongestion = -1
	if len(dropActionCongestionMatch) > 1 {
		v, err := strconv.ParseFloat(dropActionCongestionMatch[1], 64)
		if err == nil {
			metrics.DropActionCongestion = v
		}
	}

	// Drop action forwarding disabled
	dropActionForwardingDisabledRegexp := regexp.MustCompile(`(?m)^[ \t]*drop_action_forwarding_disabled.*total:\s*(\d+)`)
	dropActionForwardingDisabledMatch := dropActionForwardingDisabledRegexp.FindStringSubmatch(coverageStats)
	metrics.DropActionForwardingDisabled = -1
	if len(dropActionForwardingDisabledMatch) > 1 {
		v, err := strconv.ParseFloat(dropActionForwardingDisabledMatch[1], 64)
		if err == nil {
			metrics.DropActionForwardingDisabled = v
		}
	}

	// Drop reasons new
	// Netdev VXLAN TSO drops
	netdevVxlanTsoDropsRegexp := regexp.MustCompile(`(?m)^[ \t]*netdev_vxlan_tso_drops.*total:\s*(\d+)`)
	netdevVxlanTsoDropsMatch := netdevVxlanTsoDropsRegexp.FindStringSubmatch(coverageStats)
	metrics.NetdevVxlanTsoDrops = -1
	if len(netdevVxlanTsoDropsMatch) > 1 {
		v, err := strconv.ParseFloat(netdevVxlanTsoDropsMatch[1], 64)
		if err == nil {
			metrics.NetdevVxlanTsoDrops = v
		}
	}

	// Netdev Geneve TSO drops
	netdevGeneveTsoDropsRegexp := regexp.MustCompile(`(?m)^[ \t]*netdev_geneve_tso_drops.*total:\s*(\d+)`)
	netdevGeneveTsoDropsMatch := netdevGeneveTsoDropsRegexp.FindStringSubmatch(coverageStats)
	metrics.NetdevGeneveTsoDrops = -1
	if len(netdevGeneveTsoDropsMatch) > 1 {
		v, err := strconv.ParseFloat(netdevGeneveTsoDropsMatch[1], 64)
		if err == nil {
			metrics.NetdevGeneveTsoDrops = v
		}
	}

	// Netdev push header drops
	netdevPushHeaderDropsRegexp := regexp.MustCompile(`(?m)^[ \t]*netdev_push_header_drops.*total:\s*(\d+)`)
	netdevPushHeaderDropsMatch := netdevPushHeaderDropsRegexp.FindStringSubmatch(coverageStats)
	metrics.NetdevPushHeaderDrops = -1
	if len(netdevPushHeaderDropsMatch) > 1 {
		v, err := strconv.ParseFloat(netdevPushHeaderDropsMatch[1], 64)
		if err == nil {
			metrics.NetdevPushHeaderDrops = v
		}
	}

	// Netdev soft seg drops
	netdevSoftSegDropsRegexp := regexp.MustCompile(`(?m)^[ \t]*netdev_soft_seg_drops.*total:\s*(\d+)`)
	netdevSoftSegDropsMatch := netdevSoftSegDropsRegexp.FindStringSubmatch(coverageStats)
	metrics.NetdevSoftSegDrops = -1
	if len(netdevSoftSegDropsMatch) > 1 {
		v, err := strconv.ParseFloat(netdevSoftSegDropsMatch[1], 64)
		if err == nil {
			metrics.NetdevSoftSegDrops = v
		}
	}

	// Datapath drop tunnel TSO recirc
	datapathDropTunnelTsoRecircRegexp := regexp.MustCompile(`(?m)^[ \t]*datapath_drop_tunnel_tso_recirc.*total:\s*(\d+)`)
	datapathDropTunnelTsoRecircMatch := datapathDropTunnelTsoRecircRegexp.FindStringSubmatch(coverageStats)
	metrics.DatapathDropTunnelTsoRecirc = -1
	if len(datapathDropTunnelTsoRecircMatch) > 1 {
		v, err := strconv.ParseFloat(datapathDropTunnelTsoRecircMatch[1], 64)
		if err == nil {
			metrics.DatapathDropTunnelTsoRecirc = v
		}
	}

	// Datapath drop invalid bond
	datapathDropInvalidBondRegexp := regexp.MustCompile(`(?m)^[ \t]*datapath_drop_invalid_bond.*total:\s*(\d+)`)
	datapathDropInvalidBondMatch := datapathDropInvalidBondRegexp.FindStringSubmatch(coverageStats)
	metrics.DatapathDropInvalidBond = -1
	if len(datapathDropInvalidBondMatch) > 1 {
		v, err := strconv.ParseFloat(datapathDropInvalidBondMatch[1], 64)
		if err == nil {
			metrics.DatapathDropInvalidBond = v
		}
	}

	// Datapath drop HW miss recover
	datapathDropHwMissRecoverRegexp := regexp.MustCompile(`(?m)^[ \t]*datapath_drop_hw_miss_recover.*total:\s*(\d+)`)
	datapathDropHwMissRecoverMatch := datapathDropHwMissRecoverRegexp.FindStringSubmatch(coverageStats)
	metrics.DatapathDropHwMissRecover = -1
	if len(datapathDropHwMissRecoverMatch) > 1 {
		v, err := strconv.ParseFloat(datapathDropHwMissRecoverMatch[1], 64)
		if err == nil {
			metrics.DatapathDropHwMissRecover = v
		}
	}
}

func parsePMDStats(metrics *OvsMetric, pmdStats string) {

	missWithSuccessUpcallRegexp := regexp.MustCompile(`(?m)^[ \t]*miss\s+with\s+success\s+upcall:\s*(\d+)`)
	missWithSuccessUpcallMatch := missWithSuccessUpcallRegexp.FindStringSubmatch(pmdStats)
	metrics.MissWithSuccessUpcall = -1
	if len(missWithSuccessUpcallMatch) > 1 {
		v, err := strconv.ParseFloat(missWithSuccessUpcallMatch[1], 64)
		if err == nil {
			metrics.MissWithSuccessUpcall = v
		}
	}

	missWithFailedUpcallRegexp := regexp.MustCompile(`(?m)^[ \t]*miss\s+with\s+failed\s+upcall:\s*(\d+)`)
	missWithFailedUpcallMatch := missWithFailedUpcallRegexp.FindStringSubmatch(pmdStats)
	metrics.MissWithFailedUpcall = -1
	if len(missWithFailedUpcallMatch) > 1 {
		v, err := strconv.ParseFloat(missWithFailedUpcallMatch[1], 64)
		if err == nil {
			metrics.MissWithFailedUpcall = v
		}
	}

	processingCyclesRegexp := regexp.MustCompile(`(?m)^[ \t]*processing cycles:.*\((\d{1,3}(?:\.\d+)?)%\)`)
	processingCyclesMatch := processingCyclesRegexp.FindStringSubmatch(pmdStats)
	metrics.ProcessingCycles = -1
	if len(processingCyclesMatch) > 1 {
		v, err := strconv.ParseFloat(processingCyclesMatch[1], 64)
		if err == nil {
			metrics.ProcessingCycles = v
		}
	}

	idleCyclesRegexp := regexp.MustCompile(`(?m)^[ \t]*idle cycles:.*\((\d{1,3}(?:\.\d+)?)%\)`)
	idleCyclesMatch := idleCyclesRegexp.FindStringSubmatch(pmdStats)
	metrics.IdleCycles = -1
	if len(idleCyclesMatch) > 1 {
		v, err := strconv.ParseFloat(idleCyclesMatch[1], 64)
		if err == nil {
			metrics.IdleCycles = v
		}
	}

	avgSubtableLookupsMegaflowRegexp := regexp.MustCompile(`(?m)^[ \t]*avg\.\s+subtable\s+lookups\s+per\s+megaflow\s+hit:[ \t]*(\d+(\.\d+)?)`)
	avgSubtableLookupsMegaflowMatch := avgSubtableLookupsMegaflowRegexp.FindStringSubmatch(pmdStats)
	metrics.AvgSubtableLookupsMegaflow = -1
	if len(avgSubtableLookupsMegaflowMatch) > 1 {
		v, err := strconv.ParseFloat(avgSubtableLookupsMegaflowMatch[1], 64)
		if err == nil {
			metrics.AvgSubtableLookupsMegaflow = v
		}
	}
}

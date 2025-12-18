package main

import (
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

type OvsMetric struct {
	// PMD stats
	MissWithSuccessUpcall      float64
	MissWithFailedUpcall       float64
	AvgSubtableLookupsMegaflow float64
	ProcessingCycles           float64
	IdleCycles                 float64
	// Offload stats
	OffloadEnqueued            float64
	OffloadInserted            float64
	OffloadCtUniDirConnections float64
	OffloadCtBiDirConnections  float64
	OffloadCumAvgLatencyUs     float64
	OffloadCumLatencyStddevUs  float64
	OffloadCumLatencyMaxUs     float64
	OffloadCumLatencyMinUs     float64
	OffloadExpAvgLatencyUs     float64
	OffloadExpLatencyStddevUs  float64
	// Memory/show stats
	MemoryHandlers            float64
	MemoryIdlCellsOpenVSwitch float64
	MemoryOfconns             float64
	MemoryPorts               float64
	MemoryRevalidators        float64
	MemoryRules               float64
	MemoryUdpifKeys           float64
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
	// Upcall Flow Limit
	UpcallFlowLimitGrew    float64
	UpcallFlowLimitHit     float64
	UpcallFlowLimitKill    float64
	UpcallFlowLimitReduced float64
	UpcallFlowLimitScaled  float64
	// DOCA Pipe Group
	DocaUniqueItemTemplates float64
	// Parsed metrics from ovs-appctl metrics/show
	ParsedMetrics []*dto.MetricFamily
}

// parseTextFormat parses Prometheus TEXT exposition into []*MetricFamily.
func parseTextFormat(r io.Reader) ([]*dto.MetricFamily, error) {
	var parser expfmt.TextParser
	mfMap, err := parser.TextToMetricFamilies(r)
	if err != nil {
		return nil, err
	}
	out := make([]*dto.MetricFamily, 0, len(mfMap))
	for _, mf := range mfMap {
		out = append(out, mf)
	}
	return out, nil
}

func getOvsMetric() (*OvsMetric, int) {
	var ovsMetric OvsMetric
	successCount := 0

	cmd := exec.Command("/usr/bin/ovs-appctl", "dpif-netdev/pmd-stats-show")
	pmdStatsOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
	} else {
		parsePMDStats(&ovsMetric, string(pmdStatsOutput))
		successCount++
	}

	cmd = exec.Command("/usr/bin/ovs-appctl", "dpctl/offload-stats-show")
	offloadStatsOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
	} else {
		parseOffloadStats(&ovsMetric, string(offloadStatsOutput))
		successCount++
	}

	cmd = exec.Command("/usr/bin/ovs-appctl", "memory/show")
	memoryShowOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
	} else {
		parseMemoryShow(&ovsMetric, string(memoryShowOutput))
	}

	cmd = exec.Command("/usr/bin/ovs-appctl", "coverage/show")
	coverageOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running command: %v\n", err)
	} else {
		parseCoverageDropReasons(&ovsMetric, string(coverageOutput))
		parseCoverageDoca(&ovsMetric, string(coverageOutput))
		successCount++
	}

	// Parse metrics from ovs-appctl metrics/show
	cmd = exec.Command("/usr/bin/ovs-appctl", "metrics/show")
	metricsOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running ovs-appctl metrics/show: %v\n", err)
	} else {
		parsedMetrics, err := parseTextFormat(strings.NewReader(string(metricsOutput)))
		if err != nil {
			fmt.Printf("Error parsing metrics output: %v\n", err)
		} else {
			ovsMetric.ParsedMetrics = parsedMetrics
			successCount++
		}
	}

	// Parse DOCA pipe group unique item templates
	cmd = exec.Command("/bin/sh", "-c", "ovs-appctl doca-pipe-group/dump | grep -v 'empty_match' | grep -oE 'match.*act' | sort | uniq | wc -l")
	docaPipeGroupOutput, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running doca-pipe-group/dump command: %v\n", err)
		ovsMetric.DocaUniqueItemTemplates = -1
	} else {
		parseDocaUniqueTemplates(&ovsMetric, string(docaPipeGroupOutput))
	}

	return &ovsMetric, successCount
}

func parseDocaUniqueTemplates(metrics *OvsMetric, output string) {
	metrics.DocaUniqueItemTemplates = -1

	trimmed := strings.TrimSpace(output)
	if trimmed == "" {
		return
	}

	count, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		fmt.Printf("Error parsing unique item templates count: %v\n", err)
		return
	}

	metrics.DocaUniqueItemTemplates = count
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

	// Upcall Flow Limit
	// upcall_flow_limit_grew
	upcallFlowLimitGrewRegexp := regexp.MustCompile(`(?m)^[ \t]*upcall_flow_limit_grew.*total:\s*(\d+)`)
	upcallFlowLimitGrewMatch := upcallFlowLimitGrewRegexp.FindStringSubmatch(coverageStats)
	metrics.UpcallFlowLimitGrew = -1
	if len(upcallFlowLimitGrewMatch) > 1 {
		v, err := strconv.ParseFloat(upcallFlowLimitGrewMatch[1], 64)
		if err == nil {
			metrics.UpcallFlowLimitGrew = v
		}
	}

	// upcall_flow_limit_hit
	upcallFlowLimitHitRegexp := regexp.MustCompile(`(?m)^[ \t]*upcall_flow_limit_hit.*total:\s*(\d+)`)
	upcallFlowLimitHitMatch := upcallFlowLimitHitRegexp.FindStringSubmatch(coverageStats)
	metrics.UpcallFlowLimitHit = -1
	if len(upcallFlowLimitHitMatch) > 1 {
		v, err := strconv.ParseFloat(upcallFlowLimitHitMatch[1], 64)
		if err == nil {
			metrics.UpcallFlowLimitHit = v
		}
	}

	// upcall_flow_limit_kill
	upcallFlowLimitKillRegexp := regexp.MustCompile(`(?m)^[ \t]*upcall_flow_limit_kill.*total:\s*(\d+)`)
	upcallFlowLimitKillMatch := upcallFlowLimitKillRegexp.FindStringSubmatch(coverageStats)
	metrics.UpcallFlowLimitKill = -1
	if len(upcallFlowLimitKillMatch) > 1 {
		v, err := strconv.ParseFloat(upcallFlowLimitKillMatch[1], 64)
		if err == nil {
			metrics.UpcallFlowLimitKill = v
		}
	}

	// upcall_flow_limit_reduced
	upcallFlowLimitReducedRegexp := regexp.MustCompile(`(?m)^[ \t]*upcall_flow_limit_reduced.*total:\s*(\d+)`)
	upcallFlowLimitReducedMatch := upcallFlowLimitReducedRegexp.FindStringSubmatch(coverageStats)
	metrics.UpcallFlowLimitReduced = -1
	if len(upcallFlowLimitReducedMatch) > 1 {
		v, err := strconv.ParseFloat(upcallFlowLimitReducedMatch[1], 64)
		if err == nil {
			metrics.UpcallFlowLimitReduced = v
		}
	}

	// upcall_flow_limit_scaled
	upcallFlowLimitScaledRegexp := regexp.MustCompile(`(?m)^[ \t]*upcall_flow_limit_scaled.*total:\s*(\d+)`)
	upcallFlowLimitScaledMatch := upcallFlowLimitScaledRegexp.FindStringSubmatch(coverageStats)
	metrics.UpcallFlowLimitScaled = -1
	if len(upcallFlowLimitScaledMatch) > 1 {
		v, err := strconv.ParseFloat(upcallFlowLimitScaledMatch[1], 64)
		if err == nil {
			metrics.UpcallFlowLimitScaled = v
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

func parseOffloadStats(metrics *OvsMetric, offloadStats string) {
	// Initialize as invalid
	metrics.OffloadEnqueued = -1
	metrics.OffloadInserted = -1
	metrics.OffloadCtUniDirConnections = -1
	metrics.OffloadCtBiDirConnections = -1
	metrics.OffloadCumAvgLatencyUs = -1
	metrics.OffloadCumLatencyStddevUs = -1
	metrics.OffloadCumLatencyMaxUs = -1
	metrics.OffloadCumLatencyMinUs = -1
	metrics.OffloadExpAvgLatencyUs = -1
	metrics.OffloadExpLatencyStddevUs = -1

	enqueuedRegexp := regexp.MustCompile(`(?m)^[ \t]*Total[ \t]+Enqueued offloads:[ \t]*([0-9]+)`)
	insertedRegexp := regexp.MustCompile(`(?m)^[ \t]*Total[ \t]+Inserted offloads:[ \t]*([0-9]+)`)
	ctUniDirRegexp := regexp.MustCompile(`(?m)^[ \t]*Total[ \t]+CT uni-dir Connections:[ \t]*([0-9]+)`)
	ctBiDirRegexp := regexp.MustCompile(`(?m)^[ \t]*Total[ \t]+CT bi-dir Connections:[ \t]*([0-9]+)`)
	cumAvgLatRegexp := regexp.MustCompile(`(?m)^[ \t]*Total[ \t]+Cumulative Average latency \(us\):[ \t]*([0-9]+)`)
	cumStddevRegexp := regexp.MustCompile(`(?m)^[ \t]*Total[ \t]+Cumulative Latency stddev \(us\):[ \t]*([0-9]+)`)
	cumMaxRegexp := regexp.MustCompile(`(?m)^[ \t]*Total[ \t]+Cumulative Latency max \(us\):[ \t]*([0-9]+)`)
	cumMinRegexp := regexp.MustCompile(`(?m)^[ \t]*Total[ \t]+Cumulative Latency min \(us\):[ \t]*([0-9]+)`)
	expAvgLatRegexp := regexp.MustCompile(`(?m)^[ \t]*Total[ \t]+Exponential Average latency \(us\):[ \t]*([0-9]+)`)
	expStddevRegexp := regexp.MustCompile(`(?m)^[ \t]*Total[ \t]+Exponential Latency stddev \(us\):[ \t]*([0-9]+)`)

	if m := enqueuedRegexp.FindStringSubmatch(offloadStats); len(m) > 1 {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			metrics.OffloadEnqueued = v
		}
	}
	if m := insertedRegexp.FindStringSubmatch(offloadStats); len(m) > 1 {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			metrics.OffloadInserted = v
		}
	}
	if m := ctUniDirRegexp.FindStringSubmatch(offloadStats); len(m) > 1 {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			metrics.OffloadCtUniDirConnections = v
		}
	}
	if m := ctBiDirRegexp.FindStringSubmatch(offloadStats); len(m) > 1 {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			metrics.OffloadCtBiDirConnections = v
		}
	}
	if m := cumAvgLatRegexp.FindStringSubmatch(offloadStats); len(m) > 1 {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			metrics.OffloadCumAvgLatencyUs = v
		}
	}
	if m := cumStddevRegexp.FindStringSubmatch(offloadStats); len(m) > 1 {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			metrics.OffloadCumLatencyStddevUs = v
		}
	}
	if m := cumMaxRegexp.FindStringSubmatch(offloadStats); len(m) > 1 {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			metrics.OffloadCumLatencyMaxUs = v
		}
	}
	if m := cumMinRegexp.FindStringSubmatch(offloadStats); len(m) > 1 {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			metrics.OffloadCumLatencyMinUs = v
		}
	}
	if m := expAvgLatRegexp.FindStringSubmatch(offloadStats); len(m) > 1 {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			metrics.OffloadExpAvgLatencyUs = v
		}
	}
	if m := expStddevRegexp.FindStringSubmatch(offloadStats); len(m) > 1 {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			metrics.OffloadExpLatencyStddevUs = v
		}
	}
}

func parseMemoryShow(metrics *OvsMetric, memoryShow string) {
	metrics.MemoryHandlers = -1
	metrics.MemoryIdlCellsOpenVSwitch = -1
	metrics.MemoryOfconns = -1
	metrics.MemoryPorts = -1
	metrics.MemoryRevalidators = -1
	metrics.MemoryRules = -1
	metrics.MemoryUdpifKeys = -1

	handlersRegexp := regexp.MustCompile(`(?m)handlers:\s*(\d+)`)
	idlCellsRegexp := regexp.MustCompile(`(?m)idl-cells-Open_vSwitch:\s*(\d+)`)
	ofconnsRegexp := regexp.MustCompile(`(?m)ofconns:\s*(\d+)`)
	portsRegexp := regexp.MustCompile(`(?m)ports:\s*(\d+)`)
	revalidatorsRegexp := regexp.MustCompile(`(?m)revalidators:\s*(\d+)`)
	rulesRegexp := regexp.MustCompile(`(?m)rules:\s*(\d+)`)
	udpifKeysRegexp := regexp.MustCompile(`(?m)udpif keys:\s*(\d+)`)

	if m := handlersRegexp.FindStringSubmatch(memoryShow); len(m) > 1 {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			metrics.MemoryHandlers = v
		}
	}
	if m := idlCellsRegexp.FindStringSubmatch(memoryShow); len(m) > 1 {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			metrics.MemoryIdlCellsOpenVSwitch = v
		}
	}
	if m := ofconnsRegexp.FindStringSubmatch(memoryShow); len(m) > 1 {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			metrics.MemoryOfconns = v
		}
	}
	if m := portsRegexp.FindStringSubmatch(memoryShow); len(m) > 1 {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			metrics.MemoryPorts = v
		}
	}
	if m := revalidatorsRegexp.FindStringSubmatch(memoryShow); len(m) > 1 {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			metrics.MemoryRevalidators = v
		}
	}
	if m := rulesRegexp.FindStringSubmatch(memoryShow); len(m) > 1 {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			metrics.MemoryRules = v
		}
	}
	if m := udpifKeysRegexp.FindStringSubmatch(memoryShow); len(m) > 1 {
		if v, err := strconv.ParseFloat(m[1], 64); err == nil {
			metrics.MemoryUdpifKeys = v
		}
	}
}

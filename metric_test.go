package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
)

func Test_parseCoverageDoca(t *testing.T) {
	tests := []struct {
		name   string
		output string
		metric OvsMetric
	}{
		{
			name: "output1",
			output: `ovs_doca_no_mark  0.0/sec     0.000/sec        0.0000/sec   total: 5
ovs_doca_invalid_classify_port  0.0/sec     0.000/sec        0.0000/sec   total: 8
doca_queue_empty  0.0/sec     0.000/sec        0.0000/sec   total: 12
doca_queue_none_processed  0.0/sec     0.000/sec        0.0000/sec   total: 15
doca_resize_block  0.0/sec     0.000/sec        0.0000/sec   total: 20
doca_pipe_resize  0.0/sec     0.000/sec        0.0000/sec   total: 25
doca_pipe_resize_over_10_ms  0.0/sec     0.000/sec        0.0000/sec   total: 30`,
			metric: OvsMetric{
				// DOCA
				OvsDocaNoMark:              5,
				OvsDocaInvalidClassifyPort: 8,
				DocaQueueEmpty:             12,
				DocaQueueNoneProcessed:     15,
				DocaResizeBlock:            20,
				DocaPipeResize:             25,
				DocaPipeResizeOver10Ms:     30,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ovsMetric OvsMetric
			parseCoverageDoca(&ovsMetric, tt.output)

			diff := cmp.Diff(ovsMetric, tt.metric)
			// If there's a difference, `cmp.Diff` will return a string representation of the diff
			if diff != "" {
				t.Errorf("Structs are different:\n%s", diff)
			}
		})
	}
}

func Test_parseCoverageDropReasons(t *testing.T) {

	tests := []struct {
		name   string
		output string
		metric OvsMetric
	}{
		{
			name: "output1",
			output: `
datapath_drop_upcall_error   0.0/sec     0.000/sec        0.0000/sec   total: 5
datapath_drop_lock_error   0.0/sec     0.000/sec        0.0000/sec   total: 6
datapath_drop_rx_invalid_packet   0.0/sec     0.000/sec        0.0000/sec   total: 7
datapath_drop_meter   0.0/sec     0.000/sec        0.0000/sec   total: 8
datapath_drop_userspace_action_error   0.0/sec     0.000/sec        0.0000/sec   total: 9
datapath_drop_tunnel_push_error   0.0/sec     0.000/sec        0.0000/sec   total: 10
datapath_drop_tunnel_pop_error   0.0/sec     0.000/sec        0.0000/sec   total: 11
datapath_drop_recirc_error   0.0/sec     0.000/sec        0.0000/sec   total: 12
datapath_drop_invalid_port   0.0/sec     0.000/sec        0.0000/sec   total: 13
datapath_drop_invalid_tnl_port   0.0/sec     0.000/sec        0.0000/sec   total: 14
datapath_drop_sample_error   0.0/sec     0.000/sec        0.0000/sec   total: 15
datapath_drop_nsh_decap_error   0.0/sec     0.000/sec        0.0000/sec   total: 16
drop_action_of_pipeline   0.0/sec     0.000/sec        0.0000/sec   total: 17
drop_action_bridge_not_found   0.0/sec     0.000/sec        0.0000/sec   total: 18
drop_action_recursion_too_deep   0.0/sec     0.000/sec        0.0000/sec   total: 19
drop_action_too_many_resubmit   0.0/sec     0.000/sec        0.0000/sec   total: 20
drop_action_stack_too_deep   0.0/sec     0.000/sec        0.0000/sec   total: 21
drop_action_no_recirculation_context   0.0/sec     0.000/sec        0.0000/sec   total: 22
drop_action_recirculation_conflict   0.0/sec     0.000/sec        0.0000/sec   total: 23
drop_action_too_many_mpls_labels   0.0/sec     0.000/sec        0.0000/sec   total: 24
drop_action_invalid_tunnel_metadata   0.0/sec     0.000/sec        0.0000/sec   total: 25
drop_action_unsupported_packet_type   0.0/sec     0.000/sec        0.0000/sec   total: 26
drop_action_congestion   0.0/sec     0.000/sec        0.0000/sec   total: 27
drop_action_forwarding_disabled   0.0/sec     0.000/sec        0.0000/sec   total: 28
netdev_vxlan_tso_drops   0.0/sec     0.000/sec        0.0000/sec   total: 29
netdev_geneve_tso_drops   0.0/sec     0.000/sec        0.0000/sec   total: 30
netdev_push_header_drops   0.0/sec     0.000/sec        0.0000/sec   total: 31
netdev_soft_seg_drops   0.0/sec     0.000/sec        0.0000/sec   total: 32
datapath_drop_tunnel_tso_recirc   0.0/sec     0.000/sec        0.0000/sec   total: 33
datapath_drop_invalid_bond   0.0/sec     0.000/sec        0.0000/sec   total: 34
datapath_drop_hw_miss_recover   0.0/sec     0.000/sec        0.0000/sec   total: 35
datapath_drop_invalid_mark   0.0/sec     0.000/sec        0.0000/sec   total: 36
upcall_flow_limit_grew       0.0/sec     0.000/sec        0.0000/sec   total: 37
upcall_flow_limit_hit        0.0/sec     0.000/sec        0.0000/sec   total: 38
upcall_flow_limit_kill       0.0/sec     0.000/sec        0.0000/sec   total: 39
upcall_flow_limit_reduced    0.0/sec     0.000/sec        0.0000/sec   total: 40
upcall_flow_limit_scaled     0.0/sec     0.000/sec        0.0000/sec   total: 41`,
			metric: OvsMetric{
				// Drop reasons
				UpcallDrops:                      5,
				UpcallDropsLockError:             6,
				RxDropsInvalidPacket:             7,
				DatapathDropMeter:                8,
				DatapathDropUserspaceActionError: 9,
				DatapathDropTunnelPushError:      10,
				DatapathDropTunnelPopError:       11,
				DatapathDropRecircError:          12,
				DatapathDropInvalidPort:          13,
				DatapathDropInvalidTnlPort:       14,
				DatapathDropSampleError:          15,
				DatapathDropNshDecapError:        16,
				DropActionOfPipeline:             17,
				DropActionBridgeNotFound:         18,
				DropActionRecursionTooDeep:       19,
				DropActionTooManyResubmit:        20,
				DropActionStackTooDeep:           21,
				DropActionNoRecirculationContext: 22,
				DropActionRecirculationConflict:  23,
				DropActionTooManyMplsLabels:      24,
				DropActionInvalidTunnelMetadata:  25,
				DropActionUnsupportedPacketType:  26,
				DropActionCongestion:             27,
				DropActionForwardingDisabled:     28,
				// Drop reasons new
				NetdevVxlanTsoDrops:         29,
				NetdevGeneveTsoDrops:        30,
				NetdevPushHeaderDrops:       31,
				NetdevSoftSegDrops:          32,
				DatapathDropTunnelTsoRecirc: 33,
				DatapathDropInvalidBond:     34,
				DatapathDropHwMissRecover:   35,
				DatapathDropInvalidMark:     36,
				// Upcall Flow Limit
				UpcallFlowLimitGrew:    37,
				UpcallFlowLimitHit:     38,
				UpcallFlowLimitKill:    39,
				UpcallFlowLimitReduced: 40,
				UpcallFlowLimitScaled:  41,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ovsMetric OvsMetric
			parseCoverageDropReasons(&ovsMetric, tt.output)

			diff := cmp.Diff(ovsMetric, tt.metric)
			if diff != "" {
				t.Errorf("Structs are different:\n%s", diff)
			}
		})
	}
}

func Test_metricParsePMDStats(t *testing.T) {

	tests := []struct {
		name   string
		output string
		metric OvsMetric
	}{
		{
			name: "output1",
			output: `
pmd thread numa_id 0 core_id 11:
packets received: 89813835
packet recirculations: 25377014
avg. datapath passes per packet: 1.28
phwol hits: 4596
mfex opt hits: 0
simple match hits: 22
emc hits: 3392099
smc hits: 0
megaflow hits: 78498765
avg. subtable lookups per megaflow hit: 5.38
miss with success upcall: 33284747
miss with failed upcall: 10620
avg. packets per output batch: 1.05
idle cycles: 731072761249336 (99.80%)
processing cycles: 1492654477083 (0.20%)
avg cycles per packet: 8156487.42 (732565415726419/89813835)
avg processing cycles per packet: 16619.43 (1492654477083/89813835)
main thread:
packets received: 4
packet recirculations: 0
avg. datapath passes per packet: 1.00
phwol hits: 0
mfex opt hits: 0
simple match hits: 2
emc hits: 0
smc hits: 0
megaflow hits: 0
avg. subtable lookups per megaflow hit: 0.00
miss with success upcall: 2
miss with failed upcall: 0
avg. packets per output batch: 0.00`,
			metric: OvsMetric{
				MissWithFailedUpcall:       10620,
				IdleCycles:                 99.80,
				ProcessingCycles:           0.20,
				MissWithSuccessUpcall:      33284747,
				AvgSubtableLookupsMegaflow: 5.38,
			},
		},
		{
			name: "output2",
			output: `
pmd thread numa_id 0 core_id 11:
packets received: 7828
packet recirculations: 0
avg. datapath passes per packet: 1.00
phwol hits: 0
mfex opt hits: 0
simple match hits: 6
emc hits: 6662
smc hits: 0
megaflow hits: 81
avg. subtable lookups per megaflow hit: 1.20
miss with success upcall: 1047
miss with failed upcall: 0
avg. packets per output batch: 1.00
idle cycles: 33408495423437 (100.00%)
processing cycles: 468243547 (0.00%)
avg cycles per packet: 4267879875.70 (33408963666984/7828)
avg processing cycles per packet: 59816.50 (468243547/7828)
main thread:
packets received: 3378
packet recirculations: 0
avg. datapath passes per packet: 1.00
phwol hits: 0
mfex opt hits: 0
simple match hits: 3373
emc hits: 0
smc hits: 0
megaflow hits: 2
avg. subtable lookups per megaflow hit: 1.50
miss with success upcall: 3
miss with failed upcall: 0
avg. packets per output batch: 0.00`,
			metric: OvsMetric{
				// PMD stats
				MissWithFailedUpcall:       0,
				IdleCycles:                 100,
				ProcessingCycles:           0,
				MissWithSuccessUpcall:      1047,
				AvgSubtableLookupsMegaflow: 1.2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ovsMetric OvsMetric
			parsePMDStats(&ovsMetric, tt.output)

			diff := cmp.Diff(ovsMetric, tt.metric)

			// If there's a difference, `cmp.Diff` will return a string representation of the diff
			if diff != "" {
				t.Errorf("Structs are different:\n%s", diff)
			}
		})
	}
}

// Minimal test to assert scrape fails (HTTP 500) when all commands fail
// by stubbing fetchOvsMetrics to return successCount=0 and then invoking
// the Prometheus handler.
func Test_ScrapeFailWhenAllCommandsFail(t *testing.T) {
	// Why copy to "old" and restore it later?
	// We replace the global fetch function to simulate total failure for THIS test only.
	// Saving the original (old) and restoring it prevents leaking the stub into other tests
	// (order-dependent bugs, parallel test flakiness, and masking real integrations).
	old := fetchOvsMetrics
	fetchOvsMetrics = func() (*OvsMetric, int) { return &OvsMetric{}, 0 }
	t.Cleanup(func() { fetchOvsMetrics = old })

	reg := prometheus.NewRegistry()
	reg.MustRegister(newOvsDPCollector())

	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	h := promhttp.HandlerFor(reg, promhttp.HandlerOpts{ErrorHandling: promhttp.HTTPErrorOnError})
	h.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 when all commands fail, got %d", w.Code)
	}
}

func Test_parseOffloadStats(t *testing.T) {
	input := `HW Offload stats:
     Total                 Enqueued offloads:       0
     Total                 Inserted offloads:    8283
     Total            CT uni-dir Connections:       0
     Total             CT bi-dir Connections:       0
     Total   Cumulative Average latency (us):   14882
     Total    Cumulative Latency stddev (us):   14689
     Total       Cumulative Latency max (us):  692044
     Total       Cumulative Latency min (us):       3
     Total  Exponential Average latency (us):   14381
     Total   Exponential Latency stddev (us):   11337`

	var got OvsMetric
	parseOffloadStats(&got, input)

	want := OvsMetric{
		OffloadEnqueued:            0,
		OffloadInserted:            8283,
		OffloadCtUniDirConnections: 0,
		OffloadCtBiDirConnections:  0,
		OffloadCumAvgLatencyUs:     14882,
		OffloadCumLatencyStddevUs:  14689,
		OffloadCumLatencyMaxUs:     692044,
		OffloadCumLatencyMinUs:     3,
		OffloadExpAvgLatencyUs:     14381,
		OffloadExpLatencyStddevUs:  11337,
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatalf("parseOffloadStats mismatch (-got +want):\n%s", diff)
	}
}

func Test_parseMemoryShow(t *testing.T) {
	input := `handlers:11 idl-cells-Open_vSwitch:2742 ofconns:2 ports:47 revalidators:5 rules:20658 udpif keys:29`

	var got OvsMetric
	parseMemoryShow(&got, input)

	want := OvsMetric{
		MemoryHandlers:            11,
		MemoryIdlCellsOpenVSwitch: 2742,
		MemoryOfconns:             2,
		MemoryPorts:               47,
		MemoryRevalidators:        5,
		MemoryRules:               20658,
		MemoryUdpifKeys:           29,
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Fatalf("parseMemoryShow mismatch (-got +want):\n%s", diff)
	}
}

func Test_parseTextFormat_ovs_metrics_show(t *testing.T) {
	// Test parsing of ovs-appctl metrics/show output using sample from example-output
	input := `# HELP ovs_vswitchd_bridge A metric with a constant value '1' labeled by bridge name and type present on the instance.
# TYPE ovs_vswitchd_bridge gauge
ovs_vswitchd_bridge{name="br-hbn",type="netdev"} 1
ovs_vswitchd_bridge{name="br-sfc",type="netdev"} 1
# HELP ovs_vswitchd_bridge_n_bridges Number of bridges present in the instance.
# TYPE ovs_vswitchd_bridge_n_bridges gauge
ovs_vswitchd_bridge_n_bridges 2
# HELP ovs_vswitchd_bridge_n_flows Number of flows present on the bridge.
# TYPE ovs_vswitchd_bridge_n_flows gauge
ovs_vswitchd_bridge_n_flows{name="br-hbn",type="netdev"} 10916
ovs_vswitchd_bridge_n_flows{name="br-sfc",type="netdev"} 49
# HELP ovs_vswitchd_bridge_n_ports Number of ports present on the bridge.
# TYPE ovs_vswitchd_bridge_n_ports gauge
ovs_vswitchd_bridge_n_ports{name="br-hbn",type="netdev"} 28
ovs_vswitchd_bridge_n_ports{name="br-sfc",type="netdev"} 21
# HELP ovs_vswitchd_conntrack_connection_limit Maximum number of connections allowed.
# TYPE ovs_vswitchd_conntrack_connection_limit gauge
ovs_vswitchd_conntrack_connection_limit{datapath="ovs-netdev"} 3000000
# HELP ovs_vswitchd_conntrack_n_connections Number of tracked connections.
# TYPE ovs_vswitchd_conntrack_n_connections gauge
ovs_vswitchd_conntrack_n_connections{datapath="ovs-netdev"} 0
# HELP ovs_vswitchd_conntrack_n_dccp Number of tracked DCCP connections.
# TYPE ovs_vswitchd_conntrack_n_dccp gauge
ovs_vswitchd_conntrack_n_dccp{datapath="ovs-netdev"} 0
# HELP ovs_vswitchd_conntrack_n_icmp Number of tracked ICMP connections.
# TYPE ovs_vswitchd_conntrack_n_icmp gauge
ovs_vswitchd_conntrack_n_icmp{datapath="ovs-netdev"} 0
# HELP ovs_vswitchd_conntrack_n_icmp6 Number of tracked ICMPv6 connections.
# TYPE ovs_vswitchd_conntrack_n_icmp6 gauge
ovs_vswitchd_conntrack_n_icmp6{datapath="ovs-netdev"} 0
# HELP ovs_vswitchd_conntrack_n_igmp Number of tracked IGMP connections.
# TYPE ovs_vswitchd_conntrack_n_igmp gauge
ovs_vswitchd_conntrack_n_igmp{datapath="ovs-netdev"} 0
# HELP ovs_vswitchd_conntrack_n_other Number of tracked connections of undefined type.
# TYPE ovs_vswitchd_conntrack_n_other gauge
ovs_vswitchd_conntrack_n_other{datapath="ovs-netdev"} 0
# HELP ovs_vswitchd_conntrack_n_sctp Number of tracked SCTP connections.
# TYPE ovs_vswitchd_conntrack_n_sctp gauge
ovs_vswitchd_conntrack_n_sctp{datapath="ovs-netdev"} 0
# HELP ovs_vswitchd_conntrack_n_tcp Number of tracked TCP connections.
# TYPE ovs_vswitchd_conntrack_n_tcp gauge
ovs_vswitchd_conntrack_n_tcp{datapath="ovs-netdev"} 0
# HELP ovs_vswitchd_conntrack_n_udp Number of tracked UDP connections.
# TYPE ovs_vswitchd_conntrack_n_udp gauge
ovs_vswitchd_conntrack_n_udp{datapath="ovs-netdev"} 0
# HELP ovs_vswitchd_conntrack_n_udplite Number of tracked UDPLite connections.
# TYPE ovs_vswitchd_conntrack_n_udplite gauge
ovs_vswitchd_conntrack_n_udplite{datapath="ovs-netdev"} 0
# HELP ovs_vswitchd_conntrack_tcp_seq_chk The TCP sequence checking mode: disabled(0) or enabled(1).
# TYPE ovs_vswitchd_conntrack_tcp_seq_chk gauge
ovs_vswitchd_conntrack_tcp_seq_chk{datapath="ovs-netdev"} 1
# HELP ovs_vswitchd_datapath_bytes_total Number of bytes processed in total on this datapath.
# TYPE ovs_vswitchd_datapath_bytes_total counter
ovs_vswitchd_datapath_bytes_total{datapath="netdev@ovs-netdev"} 1.844674407370955e+19
# HELP ovs_vswitchd_datapath_cache_hit_total Number of mega flow mask cache hits for flow table matches.
# TYPE ovs_vswitchd_datapath_cache_hit_total counter
ovs_vswitchd_datapath_cache_hit_total{datapath="netdev@ovs-netdev"} 0
# HELP ovs_vswitchd_datapath_hit_total Number of flow table matches.
# TYPE ovs_vswitchd_datapath_hit_total counter
ovs_vswitchd_datapath_hit_total{datapath="netdev@ovs-netdev"} 226700582
# HELP ovs_vswitchd_datapath_hw_offload_n_ct_bidir Number of bi-directional connections offloaded in hardware.
# TYPE ovs_vswitchd_datapath_hw_offload_n_ct_bidir gauge
ovs_vswitchd_datapath_hw_offload_n_ct_bidir{datapath="ovs-netdev"} 0
# HELP ovs_vswitchd_datapath_hw_offload_n_ct_unidir Number of uni-directional connections offloaded in hardware.
# TYPE ovs_vswitchd_datapath_hw_offload_n_ct_unidir gauge
ovs_vswitchd_datapath_hw_offload_n_ct_unidir{datapath="ovs-netdev"} 0
# HELP ovs_vswitchd_datapath_hw_offload_n_enqueued Number of hardware offload requests waiting to be processed.
# TYPE ovs_vswitchd_datapath_hw_offload_n_enqueued gauge
ovs_vswitchd_datapath_hw_offload_n_enqueued{datapath="ovs-netdev"} 0
# HELP ovs_vswitchd_datapath_hw_offload_n_inserted Number of hardware offload rules currently inserted.
# TYPE ovs_vswitchd_datapath_hw_offload_n_inserted gauge
ovs_vswitchd_datapath_hw_offload_n_inserted{datapath="ovs-netdev"} 203
# HELP ovs_vswitchd_datapath_lost_total Number of misses not sent to userspace.
# TYPE ovs_vswitchd_datapath_lost_total counter
ovs_vswitchd_datapath_lost_total{datapath="netdev@ovs-netdev"} 1113733
# HELP ovs_vswitchd_datapath_mask_hit_total Number of mega flow masks visited for flow table matches.
# TYPE ovs_vswitchd_datapath_mask_hit_total counter
ovs_vswitchd_datapath_mask_hit_total{datapath="netdev@ovs-netdev"} 0
# HELP ovs_vswitchd_datapath_missed_total Number of flow table misses.
# TYPE ovs_vswitchd_datapath_missed_total counter
ovs_vswitchd_datapath_missed_total{datapath="netdev@ovs-netdev"} 248464
# HELP ovs_vswitchd_datapath_n_flows Number of flows present.
# TYPE ovs_vswitchd_datapath_n_flows gauge
ovs_vswitchd_datapath_n_flows{datapath="netdev@ovs-netdev"} 77
# HELP ovs_vswitchd_datapath_n_handlers Number of upcall handler threads.
# TYPE ovs_vswitchd_datapath_n_handlers gauge
ovs_vswitchd_datapath_n_handlers{name="netdev@ovs-netdev"} 5
# HELP ovs_vswitchd_datapath_n_masks Number of mega flow masks.
# TYPE ovs_vswitchd_datapath_n_masks gauge
ovs_vswitchd_datapath_n_masks{datapath="netdev@ovs-netdev"} 0
# HELP ovs_vswitchd_datapath_n_revalidators Number of revalidator threads.
# TYPE ovs_vswitchd_datapath_n_revalidators gauge
ovs_vswitchd_datapath_n_revalidators{name="netdev@ovs-netdev"} 3
# HELP ovs_vswitchd_datapath_offloaded_bytes_total Number of bytes processed in hardware on this datapath.
# TYPE ovs_vswitchd_datapath_offloaded_bytes_total counter
ovs_vswitchd_datapath_offloaded_bytes_total{datapath="netdev@ovs-netdev"} 1.844674407370955e+19
# HELP ovs_vswitchd_datapath_offloaded_packets_total Number of packets processed in hardware on this datapath.
# TYPE ovs_vswitchd_datapath_offloaded_packets_total counter
ovs_vswitchd_datapath_offloaded_packets_total{datapath="netdev@ovs-netdev"} 1.844674407370955e+19
# HELP ovs_vswitchd_datapath_packets_total Number of packets processed in total on this datapath.
# TYPE ovs_vswitchd_datapath_packets_total counter
ovs_vswitchd_datapath_packets_total{datapath="netdev@ovs-netdev"} 1.844674407370955e+19
# HELP ovs_vswitchd_datapath_tx_bytes_total Number of bytes emitted in total from this datapath.
# TYPE ovs_vswitchd_datapath_tx_bytes_total counter
ovs_vswitchd_datapath_tx_bytes_total{datapath="netdev@ovs-netdev"} 1.844674407370955e+19
# HELP ovs_vswitchd_datapath_tx_offloaded_bytes_total Total number of bytes emitted from this datapath and fully processed in hardware.
# TYPE ovs_vswitchd_datapath_tx_offloaded_bytes_total counter
ovs_vswitchd_datapath_tx_offloaded_bytes_total{datapath="netdev@ovs-netdev"} 1.844674407370955e+19
# HELP ovs_vswitchd_datapath_tx_offloaded_packets_total Total number of packets emitted from this datapath and fully processed in hardware.
# TYPE ovs_vswitchd_datapath_tx_offloaded_packets_total counter
ovs_vswitchd_datapath_tx_offloaded_packets_total{datapath="netdev@ovs-netdev"} 1.844674407370955e+19
# HELP ovs_vswitchd_datapath_tx_packets_total Number of packets emitted in total from this datapath.
# TYPE ovs_vswitchd_datapath_tx_packets_total counter
ovs_vswitchd_datapath_tx_packets_total{datapath="netdev@ovs-netdev"} 1.844674407370955e+19
# HELP ovs_vswitchd_handler_n_threads Number of upcall handler threads in total.
# TYPE ovs_vswitchd_handler_n_threads gauge
ovs_vswitchd_handler_n_threads 5
# HELP ovs_vswitchd_interface_admin_state The administrative state of the interface: down(0) or up(1).
# TYPE ovs_vswitchd_interface_admin_state gauge
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 1
ovs_vswitchd_interface_admin_state{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 1
ovs_vswitchd_interface_admin_state{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 1
ovs_vswitchd_interface_admin_state{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 1
ovs_vswitchd_interface_admin_state{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 1
ovs_vswitchd_interface_admin_state{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 1
ovs_vswitchd_interface_admin_state{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 1
ovs_vswitchd_interface_admin_state{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 1
ovs_vswitchd_interface_admin_state{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 1
ovs_vswitchd_interface_admin_state{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 1
ovs_vswitchd_interface_admin_state{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 1
ovs_vswitchd_interface_admin_state{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 1
ovs_vswitchd_interface_admin_state{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 1
ovs_vswitchd_interface_admin_state{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 1
ovs_vswitchd_interface_admin_state{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 1
ovs_vswitchd_interface_admin_state{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 1
ovs_vswitchd_interface_admin_state{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 1
ovs_vswitchd_interface_admin_state{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 1
ovs_vswitchd_interface_admin_state{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 1
ovs_vswitchd_interface_admin_state{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 1
ovs_vswitchd_interface_admin_state{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 1
ovs_vswitchd_interface_admin_state{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 1
ovs_vswitchd_interface_admin_state{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 1
# HELP ovs_vswitchd_interface_collisions_total The number of collisions during packet transmission.
# TYPE ovs_vswitchd_interface_collisions_total counter
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 0
ovs_vswitchd_interface_collisions_total{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 0
# HELP ovs_vswitchd_interface_duplex The duplex mode of the interface: half(0) or full(1).
# TYPE ovs_vswitchd_interface_duplex gauge
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 1
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 1
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 0
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 1
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 1
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 0
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 1
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 1
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 0
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 1
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 1
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 1
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 0
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 0
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 1
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 0
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 1
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 1
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 0
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 1
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 1
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 1
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 0
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 0
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 1
ovs_vswitchd_interface_duplex{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_duplex{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 0
ovs_vswitchd_interface_duplex{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 1
ovs_vswitchd_interface_duplex{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 1
ovs_vswitchd_interface_duplex{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 0
ovs_vswitchd_interface_duplex{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 1
ovs_vswitchd_interface_duplex{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 0
ovs_vswitchd_interface_duplex{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 0
ovs_vswitchd_interface_duplex{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_duplex{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 1
ovs_vswitchd_interface_duplex{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 0
ovs_vswitchd_interface_duplex{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 1
ovs_vswitchd_interface_duplex{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_duplex{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 1
ovs_vswitchd_interface_duplex{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 0
ovs_vswitchd_interface_duplex{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 0
ovs_vswitchd_interface_duplex{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 1
ovs_vswitchd_interface_duplex{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 1
ovs_vswitchd_interface_duplex{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 1
ovs_vswitchd_interface_duplex{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 0
ovs_vswitchd_interface_duplex{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 1
ovs_vswitchd_interface_duplex{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 0
# HELP ovs_vswitchd_interface_ifindex The ifindex of the interface.
# TYPE ovs_vswitchd_interface_ifindex gauge
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 12996772
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 10770506
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 0
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 1991636
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 14308061
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 0
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 4891815
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 12815565
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 0
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 11435656
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 4873550
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 8167755
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 0
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 0
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 8684821
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 0
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 14811628
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 10087104
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 0
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 22
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 14518608
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 3581584
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 0
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 0
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 4218687
ovs_vswitchd_interface_ifindex{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_ifindex{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 0
ovs_vswitchd_interface_ifindex{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 6693098
ovs_vswitchd_interface_ifindex{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 16593497
ovs_vswitchd_interface_ifindex{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 0
ovs_vswitchd_interface_ifindex{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 3040122
ovs_vswitchd_interface_ifindex{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 0
ovs_vswitchd_interface_ifindex{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 0
ovs_vswitchd_interface_ifindex{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_ifindex{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 5927161
ovs_vswitchd_interface_ifindex{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 0
ovs_vswitchd_interface_ifindex{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 1211306
ovs_vswitchd_interface_ifindex{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_ifindex{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 13896755
ovs_vswitchd_interface_ifindex{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 0
ovs_vswitchd_interface_ifindex{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 0
ovs_vswitchd_interface_ifindex{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 15376229
ovs_vswitchd_interface_ifindex{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 23
ovs_vswitchd_interface_ifindex{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 9089622
ovs_vswitchd_interface_ifindex{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 0
ovs_vswitchd_interface_ifindex{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 9825410
ovs_vswitchd_interface_ifindex{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 0
# HELP ovs_vswitchd_interface_info A metric with a constant value '1' labeled with the driver name, version and firmware version of the interface.
# TYPE ovs_vswitchd_interface_info gauge
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-hbn",type="dpdk",port="p1_if_r"} 1
ovs_vswitchd_interface_info{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 1
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-hbn",type="dpdk",port="pf0vf5_if_r"} 1
ovs_vswitchd_interface_info{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 1
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-hbn",type="dpdk",port="pf0tss0_if_r"} 1
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-hbn",type="dpdk",port="pf0hpf_if_r"} 1
ovs_vswitchd_interface_info{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 1
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-hbn",type="dpdk",port="pf0vf3_if_r"} 1
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-hbn",type="dpdk",port="pf0pub0_if_r"} 1
ovs_vswitchd_interface_info{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 1
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-hbn",type="dpdk",port="pf0vf7_if_r"} 1
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-hbn",type="dpdk",port="p1"} 1
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-hbn",type="dpdk",port="pf0vf4_if_r"} 1
ovs_vswitchd_interface_info{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 1
ovs_vswitchd_interface_info{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 1
ovs_vswitchd_interface_info{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 1
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-hbn",type="dpdk",port="pf0vf6_if_r"} 1
ovs_vswitchd_interface_info{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 1
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-hbn",type="dpdk",port="pf0vf0_if_r"} 1
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-hbn",type="dpdk",port="p0_if_r"} 1
ovs_vswitchd_interface_info{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 1
ovs_vswitchd_interface_info{driver_name="tun",driver_version="1.6",bridge="br-hbn",type="tap",port="br-hbn"} 1
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-hbn",type="dpdk",port="p0"} 1
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-hbn",type="dpdk",port="pf0vf2_if_r"} 1
ovs_vswitchd_interface_info{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 1
ovs_vswitchd_interface_info{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 1
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-hbn",type="dpdk",port="pf0vf1_if_r"} 1
ovs_vswitchd_interface_info{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 1
ovs_vswitchd_interface_info{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 1
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-sfc",type="dpdk",port="pf0hpf"} 1
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-sfc",type="dpdk",port="pf0vf3"} 1
ovs_vswitchd_interface_info{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 1
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-sfc",type="dpdk",port="pf0vf7"} 1
ovs_vswitchd_interface_info{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 1
ovs_vswitchd_interface_info{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 1
ovs_vswitchd_interface_info{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 1
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-sfc",type="dpdk",port="pf0vf1"} 1
ovs_vswitchd_interface_info{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 1
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-sfc",type="dpdk",port="pf0vf4"} 1
ovs_vswitchd_interface_info{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 1
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-sfc",type="dpdk",port="pf0vf5"} 1
ovs_vswitchd_interface_info{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 1
ovs_vswitchd_interface_info{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 1
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-sfc",type="dpdk",port="pf0vf2"} 1
ovs_vswitchd_interface_info{driver_name="tun",driver_version="1.6",bridge="br-sfc",type="tap",port="br-sfc"} 1
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-sfc",type="dpdk",port="pf0vf0"} 1
ovs_vswitchd_interface_info{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 1
ovs_vswitchd_interface_info{driver_name="mlx5_pci",driver_version="MLNX_DPDK 22.11.2407.1.0",bridge="br-sfc",type="dpdk",port="pf0vf6"} 1
ovs_vswitchd_interface_info{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 1
# HELP ovs_vswitchd_interface_ingress_policy_bit_burst Maximum receive burst size in kb.
# TYPE ovs_vswitchd_interface_ingress_policy_bit_burst gauge
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 0
ovs_vswitchd_interface_ingress_policy_bit_burst{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 0
# HELP ovs_vswitchd_interface_ingress_policy_bit_rate Maximum receive rate in kbps on the interface. Disabled if set to 0.
# TYPE ovs_vswitchd_interface_ingress_policy_bit_rate gauge
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 0
ovs_vswitchd_interface_ingress_policy_bit_rate{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 0
# HELP ovs_vswitchd_interface_ingress_policy_pkt_burst Maximum receive burst size in number of packets.
# TYPE ovs_vswitchd_interface_ingress_policy_pkt_burst gauge
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 0
ovs_vswitchd_interface_ingress_policy_pkt_burst{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 0
# HELP ovs_vswitchd_interface_ingress_policy_pkt_rate Maximum receive rate in pps on the interface. Disabled if set to 0.
# TYPE ovs_vswitchd_interface_ingress_policy_pkt_rate gauge
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 0
ovs_vswitchd_interface_ingress_policy_pkt_rate{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 0
# HELP ovs_vswitchd_interface_link_resets_total The number of time the interface link changed.
# TYPE ovs_vswitchd_interface_link_resets_total counter
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 3
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 3
ovs_vswitchd_interface_link_resets_total{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 0
ovs_vswitchd_interface_link_resets_total{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 0
# HELP ovs_vswitchd_interface_link_speed The current speed of the interface link in Mbps.
# TYPE ovs_vswitchd_interface_link_speed gauge
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 0
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 0
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 0
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 0
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 0
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 0
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 0
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 10
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 0
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 0
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_link_speed{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 0
ovs_vswitchd_interface_link_speed{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 0
ovs_vswitchd_interface_link_speed{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 0
ovs_vswitchd_interface_link_speed{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 0
ovs_vswitchd_interface_link_speed{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_link_speed{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 0
ovs_vswitchd_interface_link_speed{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_link_speed{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 0
ovs_vswitchd_interface_link_speed{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 0
ovs_vswitchd_interface_link_speed{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 10
ovs_vswitchd_interface_link_speed{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 0
ovs_vswitchd_interface_link_speed{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 100000
ovs_vswitchd_interface_link_speed{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 0
# HELP ovs_vswitchd_interface_link_state The state of the interface link: down(0) or up(1).
# TYPE ovs_vswitchd_interface_link_state gauge
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 1
ovs_vswitchd_interface_link_state{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 1
ovs_vswitchd_interface_link_state{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 1
ovs_vswitchd_interface_link_state{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 1
ovs_vswitchd_interface_link_state{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 1
ovs_vswitchd_interface_link_state{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 1
ovs_vswitchd_interface_link_state{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 1
ovs_vswitchd_interface_link_state{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 1
ovs_vswitchd_interface_link_state{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 1
ovs_vswitchd_interface_link_state{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 1
ovs_vswitchd_interface_link_state{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 1
ovs_vswitchd_interface_link_state{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 1
ovs_vswitchd_interface_link_state{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 1
ovs_vswitchd_interface_link_state{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 1
ovs_vswitchd_interface_link_state{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 1
ovs_vswitchd_interface_link_state{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 1
ovs_vswitchd_interface_link_state{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 1
ovs_vswitchd_interface_link_state{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 1
ovs_vswitchd_interface_link_state{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 1
ovs_vswitchd_interface_link_state{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 1
ovs_vswitchd_interface_link_state{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 1
ovs_vswitchd_interface_link_state{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 1
ovs_vswitchd_interface_link_state{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 1
# HELP ovs_vswitchd_interface_mtu The MTU of the interface.
# TYPE ovs_vswitchd_interface_mtu gauge
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 9216
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 9216
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 0
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 9216
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 9216
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 0
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 9216
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 9216
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 0
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 9216
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 9216
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 9216
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 0
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 0
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 9216
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 0
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 9216
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 9216
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 0
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 9216
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 9216
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 9216
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 0
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 0
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 9216
ovs_vswitchd_interface_mtu{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_mtu{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 0
ovs_vswitchd_interface_mtu{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 9216
ovs_vswitchd_interface_mtu{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 9216
ovs_vswitchd_interface_mtu{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 0
ovs_vswitchd_interface_mtu{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 9216
ovs_vswitchd_interface_mtu{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 0
ovs_vswitchd_interface_mtu{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 0
ovs_vswitchd_interface_mtu{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_mtu{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 9216
ovs_vswitchd_interface_mtu{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 0
ovs_vswitchd_interface_mtu{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 9216
ovs_vswitchd_interface_mtu{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_mtu{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 9216
ovs_vswitchd_interface_mtu{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 0
ovs_vswitchd_interface_mtu{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 0
ovs_vswitchd_interface_mtu{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 9216
ovs_vswitchd_interface_mtu{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 9216
ovs_vswitchd_interface_mtu{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 9216
ovs_vswitchd_interface_mtu{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 0
ovs_vswitchd_interface_mtu{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 9216
ovs_vswitchd_interface_mtu{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 0
# HELP ovs_vswitchd_interface_multicast_total The number of multicast packets received by the interface.
# TYPE ovs_vswitchd_interface_multicast_total counter
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 0
ovs_vswitchd_interface_multicast_total{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 0
# HELP ovs_vswitchd_interface_of_port The OpenFlow port ID associated with the interface.
# TYPE ovs_vswitchd_interface_of_port gauge
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 4
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 1
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 18
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 9
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 24
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 6
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 17
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 14
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 26
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 19
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 22
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 5
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 16
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 27
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 21
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 11
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 20
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 13
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 8
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 2
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 7
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 65534
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 3
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 12
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 15
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 23
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 10
ovs_vswitchd_interface_of_port{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 25
ovs_vswitchd_interface_of_port{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 4
ovs_vswitchd_interface_of_port{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 1
ovs_vswitchd_interface_of_port{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 9
ovs_vswitchd_interface_of_port{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 18
ovs_vswitchd_interface_of_port{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 17
ovs_vswitchd_interface_of_port{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 6
ovs_vswitchd_interface_of_port{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 14
ovs_vswitchd_interface_of_port{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 19
ovs_vswitchd_interface_of_port{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 5
ovs_vswitchd_interface_of_port{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 16
ovs_vswitchd_interface_of_port{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 11
ovs_vswitchd_interface_of_port{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 20
ovs_vswitchd_interface_of_port{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 13
ovs_vswitchd_interface_of_port{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 2
ovs_vswitchd_interface_of_port{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 8
ovs_vswitchd_interface_of_port{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 7
ovs_vswitchd_interface_of_port{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 65534
ovs_vswitchd_interface_of_port{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 3
ovs_vswitchd_interface_of_port{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 12
ovs_vswitchd_interface_of_port{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 15
ovs_vswitchd_interface_of_port{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 10
# HELP ovs_vswitchd_interface_rx_bytes_total The number of bytes received.
# TYPE ovs_vswitchd_interface_rx_bytes_total counter
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 511555
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 144
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 170588972
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 0
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 267285
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 0
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 120168
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 60398200
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 0
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 144
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 2230941
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 144
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 58906544
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 0
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 170270829
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 144
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 170332362
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 116384
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 1721647307
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 46951137790101
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 2768
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 69204595
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 115012
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 169956174
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 0
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 120288
ovs_vswitchd_interface_rx_bytes_total{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_rx_bytes_total{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 441760579
ovs_vswitchd_interface_rx_bytes_total{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 103898391537
ovs_vswitchd_interface_rx_bytes_total{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 3040326
ovs_vswitchd_interface_rx_bytes_total{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 180
ovs_vswitchd_interface_rx_bytes_total{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 339948
ovs_vswitchd_interface_rx_bytes_total{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 441542360
ovs_vswitchd_interface_rx_bytes_total{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 180
ovs_vswitchd_interface_rx_bytes_total{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_rx_bytes_total{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 3057090
ovs_vswitchd_interface_rx_bytes_total{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 180
ovs_vswitchd_interface_rx_bytes_total{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 336722
ovs_vswitchd_interface_rx_bytes_total{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 123935004
ovs_vswitchd_interface_rx_bytes_total{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 337922
ovs_vswitchd_interface_rx_bytes_total{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 275037873694
ovs_vswitchd_interface_rx_bytes_total{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 441479962
ovs_vswitchd_interface_rx_bytes_total{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 3066296
ovs_vswitchd_interface_rx_bytes_total{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 112418970
ovs_vswitchd_interface_rx_bytes_total{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 3066478
ovs_vswitchd_interface_rx_bytes_total{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 180
ovs_vswitchd_interface_rx_bytes_total{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 336764
ovs_vswitchd_interface_rx_bytes_total{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 441228848
# HELP ovs_vswitchd_interface_rx_crc_errors_total The number of packets with CRC errors received by the interface.
# TYPE ovs_vswitchd_interface_rx_crc_errors_total counter
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 0
ovs_vswitchd_interface_rx_crc_errors_total{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 0
# HELP ovs_vswitchd_interface_rx_dropped_total Number of packets received but not processed, e.g. due to lack of resources or unsupported protocol. For hardware interface this counter should not include packets dropped by the device due to buffer exhaustion which are counted separately in rx_missed_errors.
# TYPE ovs_vswitchd_interface_rx_dropped_total counter
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 12
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 0
ovs_vswitchd_interface_rx_dropped_total{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 0
# HELP ovs_vswitchd_interface_rx_errors_total Total number of bad packets received on this interface. This counter includes all rx_length_errors, rx_crc_errors, rx_frame_errors and other errors not otherwise counted.
# TYPE ovs_vswitchd_interface_rx_errors_total counter
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 0
ovs_vswitchd_interface_rx_errors_total{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 0
# HELP ovs_vswitchd_interface_rx_fifo_errors_total Receiver FIFO error counter. This statistics was used interchangeably with rx_over_errors but is not recommended for use in drivers for high speed interfaces. This statistics is used on software devices, e.g. to count software packets queue overflow or sequencing errors.
# TYPE ovs_vswitchd_interface_rx_fifo_errors_total counter
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 0
ovs_vswitchd_interface_rx_fifo_errors_total{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 0
# HELP ovs_vswitchd_interface_rx_frame_errors_total The number of received packets with frame alignment errors on the interface.
# TYPE ovs_vswitchd_interface_rx_frame_errors_total counter
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 0
ovs_vswitchd_interface_rx_frame_errors_total{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 0
# HELP ovs_vswitchd_interface_rx_length_errors_total The number of packets dropped due to invalid length.
# TYPE ovs_vswitchd_interface_rx_length_errors_total counter
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 0
ovs_vswitchd_interface_rx_length_errors_total{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 0
# HELP ovs_vswitchd_interface_rx_missed_errors_total The number of packets missed by the host due to lack of buffer space. This usually indicates that the host interface is slower than the hardware interface. This statistics corresponds to hardware events and is not used on software devices.
# TYPE ovs_vswitchd_interface_rx_missed_errors_total counter
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 0
ovs_vswitchd_interface_rx_missed_errors_total{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 0
# HELP ovs_vswitchd_interface_rx_over_errors_total Receiver FIFO overflow event counter. This statistics was used interchangeably with rx_fifo_errors. This statistics corresponds to hardware events and is not commonly used on software devices.
# TYPE ovs_vswitchd_interface_rx_over_errors_total counter
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 0
ovs_vswitchd_interface_rx_over_errors_total{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 0
# HELP ovs_vswitchd_interface_rx_packets_total The number of packets received.
# TYPE ovs_vswitchd_interface_rx_packets_total counter
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 789
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 1
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 2331651
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 0
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 1426
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 0
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 898
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 529785
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 0
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 1
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 13115
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 1
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 796528
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 0
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 2327937
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 1
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 2328572
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 873
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 4976570
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 192506426878
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 40
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 77645
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 866
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 2324167
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 0
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 907
ovs_vswitchd_interface_rx_packets_total{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_rx_packets_total{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 5338993
ovs_vswitchd_interface_rx_packets_total{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 350083717
ovs_vswitchd_interface_rx_packets_total{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 20886
ovs_vswitchd_interface_rx_packets_total{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 2
ovs_vswitchd_interface_rx_packets_total{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 2395
ovs_vswitchd_interface_rx_packets_total{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 5336286
ovs_vswitchd_interface_rx_packets_total{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 2
ovs_vswitchd_interface_rx_packets_total{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_rx_packets_total{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 21004
ovs_vswitchd_interface_rx_packets_total{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 2
ovs_vswitchd_interface_rx_packets_total{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 2388
ovs_vswitchd_interface_rx_packets_total{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 1005106
ovs_vswitchd_interface_rx_packets_total{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 2385
ovs_vswitchd_interface_rx_packets_total{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 48222351
ovs_vswitchd_interface_rx_packets_total{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 5335329
ovs_vswitchd_interface_rx_packets_total{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 21077
ovs_vswitchd_interface_rx_packets_total{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 495246
ovs_vswitchd_interface_rx_packets_total{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 21076
ovs_vswitchd_interface_rx_packets_total{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 2
ovs_vswitchd_interface_rx_packets_total{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 2382
ovs_vswitchd_interface_rx_packets_total{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 5332247
# HELP ovs_vswitchd_interface_tx_bytes_total The number of bytes transmitted.
# TYPE ovs_vswitchd_interface_tx_bytes_total counter
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 124607
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 0
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 441760579
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 0
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 53494247922
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 180
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 1899826
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 31787022
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 180
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 0
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 103581798
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 0
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 123935004
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 180
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 441542360
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 0
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 441479962
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 1913998
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 1914431
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 275037873694
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 0
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 1856100416
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 1914956
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 441228848
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 180
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 1912613
ovs_vswitchd_interface_tx_bytes_total{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_tx_bytes_total{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 170588972
ovs_vswitchd_interface_tx_bytes_total{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 60040366
ovs_vswitchd_interface_tx_bytes_total{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 28660
ovs_vswitchd_interface_tx_bytes_total{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 0
ovs_vswitchd_interface_tx_bytes_total{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 0
ovs_vswitchd_interface_tx_bytes_total{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 170270829
ovs_vswitchd_interface_tx_bytes_total{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 0
ovs_vswitchd_interface_tx_bytes_total{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_tx_bytes_total{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 30568
ovs_vswitchd_interface_tx_bytes_total{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 0
ovs_vswitchd_interface_tx_bytes_total{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 0
ovs_vswitchd_interface_tx_bytes_total{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 58906544
ovs_vswitchd_interface_tx_bytes_total{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 0
ovs_vswitchd_interface_tx_bytes_total{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 46951137790101
ovs_vswitchd_interface_tx_bytes_total{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 170332362
ovs_vswitchd_interface_tx_bytes_total{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 29500
ovs_vswitchd_interface_tx_bytes_total{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 0
ovs_vswitchd_interface_tx_bytes_total{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 29916
ovs_vswitchd_interface_tx_bytes_total{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 0
ovs_vswitchd_interface_tx_bytes_total{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 0
ovs_vswitchd_interface_tx_bytes_total{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 169956174
# HELP ovs_vswitchd_interface_tx_dropped_total The number of packets dropped on their way to transmission, e.g. due to lack of resources.
# TYPE ovs_vswitchd_interface_tx_dropped_total counter
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 0
ovs_vswitchd_interface_tx_dropped_total{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 0
# HELP ovs_vswitchd_interface_tx_errors_total Total number of transmit issues on this interface.
# TYPE ovs_vswitchd_interface_tx_errors_total counter
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 0
ovs_vswitchd_interface_tx_errors_total{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 0
# HELP ovs_vswitchd_interface_tx_packets_total The number of packets transmitted.
# TYPE ovs_vswitchd_interface_tx_packets_total counter
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="p1_if_r",type="dpdk",port="p1_if_r"} 841
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="vxlan0",type="vxlan",port="vxlan0"} 0
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="pf0vf5_if_r",type="dpdk",port="pf0vf5_if_r"} 0
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="p-pf0vf0_if_r-hbn",type="patch",port="p-pf0vf0_if_r-h"} 5338993
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="pf0tss0_if_r",type="dpdk",port="pf0tss0_if_r"} 0
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="pf0hpf_if_r",type="dpdk",port="pf0hpf_if_r"} 221081998
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="p-pf0vf4_if_r-hbn",type="patch",port="p-pf0vf4_if_r-h"} 2
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="pf0vf3_if_r",type="dpdk",port="pf0vf3_if_r"} 20722
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="pf0pub0_if_r",type="dpdk",port="pf0pub0_if_r"} 529743
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="p-pf0vf5_if_r-hbn",type="patch",port="p-pf0vf5_if_r-h"} 2
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="pf0vf7_if_r",type="dpdk",port="pf0vf7_if_r"} 0
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="p1",type="dpdk",port="p1"} 11828
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="pf0vf4_if_r",type="dpdk",port="pf0vf4_if_r"} 0
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="p-pf0pub0_if_r-hbn",type="patch",port="p-pf0pub0_if_r-"} 1005106
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="p-pf0vf6_if_r-hbn",type="patch",port="p-pf0vf6_if_r-h"} 2
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="p-pf0vf1_if_r-hbn",type="patch",port="p-pf0vf1_if_r-h"} 5336286
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="pf0vf6_if_r",type="dpdk",port="pf0vf6_if_r"} 0
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="p-pf0vf2_if_r-hbn",type="patch",port="p-pf0vf2_if_r-h"} 5335329
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="pf0vf0_if_r",type="dpdk",port="pf0vf0_if_r"} 20885
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="p0_if_r",type="dpdk",port="p0_if_r"} 2040
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="p-pf0hpf_if_r-hbn",type="patch",port="p-pf0hpf_if_r-h"} 48222351
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="br-hbn",type="tap",port="br-hbn"} 0
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="p0",type="dpdk",port="p0"} 5737819
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="pf0vf2_if_r",type="dpdk",port="pf0vf2_if_r"} 20904
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="p-pf0vf3_if_r-hbn",type="patch",port="p-pf0vf3_if_r-h"} 5332247
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="p-pf0vf7_if_r-hbn",type="patch",port="p-pf0vf7_if_r-h"} 2
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="pf0vf1_if_r",type="dpdk",port="pf0vf1_if_r"} 20872
ovs_vswitchd_interface_tx_packets_total{bridge="br-hbn",name="p-pf0tss0_if_r-hbn",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_tx_packets_total{bridge="br-sfc",name="p-pf0vf0_if_r-sfc",type="patch",port="p-pf0vf0_if_r-s"} 2331651
ovs_vswitchd_interface_tx_packets_total{bridge="br-sfc",name="pf0hpf",type="dpdk",port="pf0hpf"} 87071
ovs_vswitchd_interface_tx_packets_total{bridge="br-sfc",name="pf0vf3",type="dpdk",port="pf0vf3"} 373
ovs_vswitchd_interface_tx_packets_total{bridge="br-sfc",name="p-pf0vf7_if_r-sfc",type="patch",port="p-pf0vf7_if_r-s"} 0
ovs_vswitchd_interface_tx_packets_total{bridge="br-sfc",name="pf0vf7",type="dpdk",port="pf0vf7"} 0
ovs_vswitchd_interface_tx_packets_total{bridge="br-sfc",name="p-pf0vf1_if_r-sfc",type="patch",port="p-pf0vf1_if_r-s"} 2327937
ovs_vswitchd_interface_tx_packets_total{bridge="br-sfc",name="p-pf0vf5_if_r-sfc",type="patch",port="p-pf0vf5_if_r-s"} 0
ovs_vswitchd_interface_tx_packets_total{bridge="br-sfc",name="p-pf0tss0_if_r-sfc",type="patch",port="p-pf0tss0_if_r-"} 0
ovs_vswitchd_interface_tx_packets_total{bridge="br-sfc",name="pf0vf1",type="dpdk",port="pf0vf1"} 400
ovs_vswitchd_interface_tx_packets_total{bridge="br-sfc",name="p-pf0vf6_if_r-sfc",type="patch",port="p-pf0vf6_if_r-s"} 0
ovs_vswitchd_interface_tx_packets_total{bridge="br-sfc",name="pf0vf4",type="dpdk",port="pf0vf4"} 0
ovs_vswitchd_interface_tx_packets_total{bridge="br-sfc",name="p-pf0pub0_if_r-sfc",type="patch",port="p-pf0pub0_if_r-"} 796528
ovs_vswitchd_interface_tx_packets_total{bridge="br-sfc",name="pf0vf5",type="dpdk",port="pf0vf5"} 0
ovs_vswitchd_interface_tx_packets_total{bridge="br-sfc",name="p-pf0hpf_if_r-sfc",type="patch",port="p-pf0hpf_if_r-s"} 192506426878
ovs_vswitchd_interface_tx_packets_total{bridge="br-sfc",name="p-pf0vf2_if_r-sfc",type="patch",port="p-pf0vf2_if_r-s"} 2328572
ovs_vswitchd_interface_tx_packets_total{bridge="br-sfc",name="pf0vf2",type="dpdk",port="pf0vf2"} 388
ovs_vswitchd_interface_tx_packets_total{bridge="br-sfc",name="br-sfc",type="tap",port="br-sfc"} 0
ovs_vswitchd_interface_tx_packets_total{bridge="br-sfc",name="pf0vf0",type="dpdk",port="pf0vf0"} 387
ovs_vswitchd_interface_tx_packets_total{bridge="br-sfc",name="p-pf0vf4_if_r-sfc",type="patch",port="p-pf0vf4_if_r-s"} 0
ovs_vswitchd_interface_tx_packets_total{bridge="br-sfc",name="pf0vf6",type="dpdk",port="pf0vf6"} 0
ovs_vswitchd_interface_tx_packets_total{bridge="br-sfc",name="p-pf0vf3_if_r-sfc",type="patch",port="p-pf0vf3_if_r-s"} 2324167
# HELP ovs_vswitchd_memory_data The process sum of data and stack size in bytes.
# TYPE ovs_vswitchd_memory_data gauge
ovs_vswitchd_memory_data 765923328
# HELP ovs_vswitchd_memory_frag_factor The fragmentation factor of the process dynamic memory, defined as (rss/in_use).
# TYPE ovs_vswitchd_memory_frag_factor gauge
ovs_vswitchd_memory_frag_factor 3.399481843180951
# HELP ovs_vswitchd_memory_in_use The amount of memory currently allocated in bytes.
# TYPE ovs_vswitchd_memory_in_use gauge
ovs_vswitchd_memory_in_use 233282272
# HELP ovs_vswitchd_memory_rss The process resident set size in bytes.
# TYPE ovs_vswitchd_memory_rss gauge
ovs_vswitchd_memory_rss 793038848
# HELP ovs_vswitchd_memory_vmsize The process virtual memory size in bytes.
# TYPE ovs_vswitchd_memory_vmsize gauge
ovs_vswitchd_memory_vmsize 210071859200
# HELP ovs_vswitchd_metrics_histogram_read_errors_total Number of histogram reads that could not resolve without inconsistencies.
# TYPE ovs_vswitchd_metrics_histogram_read_errors_total counter
ovs_vswitchd_metrics_histogram_read_errors_total 0
# HELP ovs_vswitchd_poll_threads_busy_cycles Percent of useful CPU cycles.
# TYPE ovs_vswitchd_poll_threads_busy_cycles gauge
ovs_vswitchd_poll_threads_busy_cycles{core="7",numa="0",datapath="ovs-netdev"} 0.183127123014794
# HELP ovs_vswitchd_poll_threads_busy_cycles_per_packet Average number of active CPU cycles per packet.
# TYPE ovs_vswitchd_poll_threads_busy_cycles_per_packet gauge
ovs_vswitchd_poll_threads_busy_cycles_per_packet{core="7",numa="0",datapath="ovs-netdev"} 12913.56641880694
# HELP ovs_vswitchd_poll_threads_cycles_per_packet Average number of CPU cycles per packet.
# TYPE ovs_vswitchd_poll_threads_cycles_per_packet gauge
ovs_vswitchd_poll_threads_cycles_per_packet{core="7",numa="0",datapath="ovs-netdev"} 7051695.131891363
# HELP ovs_vswitchd_poll_threads_hit_total Number of flow table matches.
# TYPE ovs_vswitchd_poll_threads_hit_total counter
ovs_vswitchd_poll_threads_hit_total{core="7",numa="0",datapath="ovs-netdev"} 226432991
# HELP ovs_vswitchd_poll_threads_idle_cycles Percent of idle CPU cycles.
# TYPE ovs_vswitchd_poll_threads_idle_cycles gauge
ovs_vswitchd_poll_threads_idle_cycles{core="7",numa="0",datapath="ovs-netdev"} 99.8168728769852
# HELP ovs_vswitchd_poll_threads_lookups_per_hit Average number of lookups per flow table hit.
# TYPE ovs_vswitchd_poll_threads_lookups_per_hit gauge
ovs_vswitchd_poll_threads_lookups_per_hit{core="7",numa="0",datapath="ovs-netdev"} 5.190361654787031
# HELP ovs_vswitchd_poll_threads_lost_total Number of flow table misses and upcall failed.
# TYPE ovs_vswitchd_poll_threads_lost_total counter
ovs_vswitchd_poll_threads_lost_total{core="7",numa="0",datapath="ovs-netdev"} 1113700
# HELP ovs_vswitchd_poll_threads_missed_total Number of flow table misses and upcall succeeded.
# TYPE ovs_vswitchd_poll_threads_missed_total counter
ovs_vswitchd_poll_threads_missed_total{core="7",numa="0",datapath="ovs-netdev"} 248389
# HELP ovs_vswitchd_poll_threads_packets_per_batch Average number of packets per batch.
# TYPE ovs_vswitchd_poll_threads_packets_per_batch gauge
ovs_vswitchd_poll_threads_packets_per_batch{core="7",numa="0",datapath="ovs-netdev"} 5.537417459392985
# HELP ovs_vswitchd_poll_threads_packets_total Number of received packets.
# TYPE ovs_vswitchd_poll_threads_packets_total counter
ovs_vswitchd_poll_threads_packets_total{core="7",numa="0",datapath="ovs-netdev"} 227774665
# HELP ovs_vswitchd_poll_threads_passes_per_packet Average number of datapath passes per packet.
# TYPE ovs_vswitchd_poll_threads_passes_per_packet gauge
ovs_vswitchd_poll_threads_passes_per_packet{core="7",numa="0",datapath="ovs-netdev"} 1.000089632444416
# HELP ovs_vswitchd_poll_threads_recirc_per_packet Average number of recirculations per packet.
# TYPE ovs_vswitchd_poll_threads_recirc_per_packet gauge
ovs_vswitchd_poll_threads_recirc_per_packet{core="7",numa="0",datapath="ovs-netdev"} 8.963244441606357e-05
# HELP ovs_vswitchd_poll_threads_recirculations_total Number of executed packet recirculations.
# TYPE ovs_vswitchd_poll_threads_recirculations_total counter
ovs_vswitchd_poll_threads_recirculations_total{core="7",numa="0",datapath="ovs-netdev"} 20416
# HELP ovs_vswitchd_revalidator_n_threads Number of revalidator threads in total.
# TYPE ovs_vswitchd_revalidator_n_threads gauge
ovs_vswitchd_revalidator_n_threads 3
# HELP ovs_vswitchd_scrape_duration_seconds Time elapsed to process this request in seconds.
# TYPE ovs_vswitchd_scrape_duration_seconds gauge
ovs_vswitchd_scrape_duration_seconds 0.34
`

	got, err := parseTextFormat(strings.NewReader(input))
	if err != nil {
		t.Fatalf("parseTextFormat failed: %v", err)
	}

	// Should parse multiple metric families (the exact count may vary)
	if len(got) < 6 {
		t.Fatalf("expected at least 6 metric families, got %d", len(got))
	}

	// Create map for easy lookup
	mfMap := make(map[string]*dto.MetricFamily)
	for _, mf := range got {
		mfMap[*mf.Name] = mf
	}

	// Check a few metrics in the map for sanity

	// Test bridge metric with labels
	bridgeMF := mfMap["ovs_vswitchd_bridge"]
	if bridgeMF == nil {
		t.Fatal("ovs_vswitchd_bridge metric family not found")
	}
	if *bridgeMF.Type != dto.MetricType_GAUGE {
		t.Errorf("expected GAUGE type, got %v", *bridgeMF.Type)
	}
	if len(bridgeMF.Metric) != 2 {
		t.Errorf("expected 2 bridge metrics, got %d", len(bridgeMF.Metric))
	}

	// Test bridge count metric without labels
	bridgeCountMF := mfMap["ovs_vswitchd_bridge_n_bridges"]
	if bridgeCountMF == nil {
		t.Fatal("ovs_vswitchd_bridge_n_bridges metric family not found")
	}
	if len(bridgeCountMF.Metric) != 1 {
		t.Errorf("expected 1 bridge count metric, got %d", len(bridgeCountMF.Metric))
	}
	if *bridgeCountMF.Metric[0].Gauge.Value != 2.0 {
		t.Errorf("expected bridge count 2.0, got %f", *bridgeCountMF.Metric[0].Gauge.Value)
	}

	// Test scrape duration metric
	scrapeMF := mfMap["ovs_vswitchd_scrape_duration_seconds"]
	if scrapeMF == nil {
		t.Fatal("ovs_vswitchd_scrape_duration_seconds metric family not found")
	}
	if *scrapeMF.Metric[0].Gauge.Value != 0.34 {
		t.Errorf("expected scrape duration 0.34, got %f", *scrapeMF.Metric[0].Gauge.Value)
	}
}

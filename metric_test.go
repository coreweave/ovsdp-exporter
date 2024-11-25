package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
datapath_drop_hw_miss_recover   0.0/sec     0.000/sec        0.0000/sec   total: 35`,
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

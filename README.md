### ovsdp-exporter
OVS datapath metric exporter. It executes a few `ovs-appctl` commands (plus `ovs-vsctl` and `ovs-ofctl` for per-bridge meter stats) and exposes selected values as Prometheus metrics.

### Commands executed and exported metrics

#### `ovs-appctl dpif-netdev/pmd-stats-show`
PMD thread and datapath performance stats.

- `ovsdp_miss_with_success_upcall`: Cache misses with successful upcalls.
- `ovsdp_miss_with_failed_upcall`: Cache misses with failed upcalls.
- `ovsdp_processing_cycles`: CPU cycles spent actively processing packets (percent).
- `ovsdp_idle_cycles`: CPU cycles idle waiting for packets (percent).
- `ovsdp_avg_subtable_lookups_megaflow`: Average subtable lookups per megaflow hit.

#### `ovs-appctl dpctl/offload-stats-show`
Hardware offload statistics and latency.

- `ovsdp_offload_enqueued`: Enqueued offloads total.
- `ovsdp_offload_inserted`: Inserted offloads total.
- `ovsdp_offload_ct_unidir_connections`: CT uni-dir connections offloaded.
- `ovsdp_offload_ct_bidir_connections`: CT bi-dir connections offloaded.
- `ovsdp_offload_cum_avg_latency_us`: Cumulative average latency (microseconds).
- `ovsdp_offload_cum_latency_stddev_us`: Cumulative latency standard deviation (microseconds).
- `ovsdp_offload_cum_latency_max_us`: Cumulative latency maximum observed (microseconds).
- `ovsdp_offload_cum_latency_min_us`: Cumulative latency minimum observed (microseconds).
- `ovsdp_offload_exp_avg_latency_us`: Exponential moving average latency (microseconds).
- `ovsdp_offload_exp_latency_stddev_us`: Exponential moving latency standard deviation (microseconds).

#### `ovs-appctl coverage/show`
Drop reasons, DOCA counters, and upcall flow-limit behavior.

- Drop reasons (datapath and actions):
  - `ovsdp_datapath_drop_upcall_error`: Drop due to error in the Upcall process.
  - `ovsdp_datapath_drop_lock_error`: Drop due to Upcall lock contention.
  - `ovsdp_datapath_drop_rx_invalid_packet`: Drop invalid packet (shorter than Ethernet header indicates).
  - `ovsdp_datapath_drop_meter`: Drop in the OpenFlow Meter Table.
  - `ovsdp_datapath_drop_userspace_action_error`: Drop due to generic action execution error.
  - `ovsdp_datapath_drop_tunnel_push_error`: Drop due to tunnel push (encap) error.
  - `ovsdp_datapath_drop_tunnel_pop_error`: Drop due to tunnel pop (decap) error.
  - `ovsdp_datapath_drop_recirc_error`: Drop due to recirculation error.
  - `ovsdp_datapath_drop_invalid_port`: Drop due to invalid port.
  - `ovsdp_datapath_drop_invalid_tnl_port`: Drop due to invalid tunnel port on pop.
  - `ovsdp_datapath_drop_sample_error`: Drop due to sampling error.
  - `ovsdp_datapath_drop_nsh_decap_error`: Drop due to invalid NSH decapsulation.
  - `ovsdp_drop_action_of_pipeline`: Drop due to pipeline/action parsing errors.
  - `ovsdp_drop_action_bridge_not_found`: Drop due to bridge not found at translation time.
  - `ovsdp_drop_action_recursion_too_deep`: Drop due to excessive translation recursion.
  - `ovsdp_drop_action_too_many_resubmit`: Drop due to too many resubmits.
  - `ovsdp_drop_action_stack_too_deep`: Drop due to excessive stack usage (>64kB).
  - `ovsdp_drop_action_no_recirculation_context`: Drop due to missing recirculation context.
  - `ovsdp_drop_action_recirculation_conflict`: Drop due to recirculation conflict.
  - `ovsdp_drop_action_too_many_mpls_labels`: Drop due to too many MPLS labels to pop.
  - `ovsdp_drop_action_invalid_tunnel_metadata`: Drop due to invalid GENEVE tunnel metadata.
  - `ovsdp_drop_action_unsupported_packet_type`: Drop due to unsupported packet type.
  - `ovsdp_drop_action_congestion`: Drop due to ECN congestion mismatch.
  - `ovsdp_drop_action_forwarding_disabled`: Drop when port forwarding is disabled.

- Additional datapath counters:
  - `ovsdp_netdev_vxlan_tso_drops`: Drops due to VXLAN TSO issues.
  - `ovsdp_netdev_geneve_tso_drops`: Drops due to Geneve TSO issues.
  - `ovsdp_netdev_push_header_drops`: Drops due to push header errors.
  - `ovsdp_netdev_soft_seg_drops`: Drops due to software segmentation issues.
  - `ovsdp_datapath_drop_tunnel_tso_recirc`: Drops due to tunnel TSO recirculation errors.
  - `ovsdp_datapath_drop_invalid_bond`: Drops due to invalid bond configuration.
  - `ovsdp_datapath_drop_hw_miss_recover`: Drops due to hardware miss recovery failure.

- DOCA:
  - `ovsdp_ovs_doca_no_mark`: Packets dropped due to missing mark in OVS-DOCA.
  - `ovsdp_ovs_doca_invalid_classify_port`: Packets dropped due to invalid classify port in OVS-DOCA.
  - `ovsdp_doca_queue_empty`: Times an offload completion queue was found empty.
  - `ovsdp_doca_queue_none_processed`: Times a queue had pending entries but none processed.
  - `ovsdp_doca_resize_block`: Queue processing blocked during pipeline resizing with no entries processed.
  - `ovsdp_doca_pipe_resize`: Times a pipe resize operation began.
  - `ovsdp_doca_pipe_resize_over_10_ms`: Times a pipe resize took longer than 10 ms.

- Upcall Flow Limit behavior:
  - `ovsdp_upcall_flow_limit_grew`: Flow limit increased due to fast processing.
  - `ovsdp_upcall_flow_limit_hit`: Flow limit was hit during upcall processing.
  - `ovsdp_upcall_flow_limit_kill`: Flows killed due to exceeding flow limit.
  - `ovsdp_upcall_flow_limit_reduced`: Flow limit reduced due to high processing time.
  - `ovsdp_upcall_flow_limit_scaled`: Flow limit scaled down due to very long processing time.

#### `ovs-appctl metrics/show`
Prometheus-formatted metrics directly from OVS. All metrics returned by this command are exposed with their original names and labels, supporting both gauge and counter types. This includes a wide range of `ovs_vswitchd_*` prefixed metrics covering:

- **Bridge metrics**: Bridge configuration and flow counts
  - `ovs_vswitchd_bridge`: A metric with a constant value '1' labeled by bridge name and type present on the instance
  - `ovs_vswitchd_bridge_n_bridges`: Number of bridges present in the instance
  - `ovs_vswitchd_bridge_n_flows`: Number of flows present on the bridge
  - `ovs_vswitchd_bridge_n_ports`: Number of ports present on the bridge

- **Connection tracking**: Conntrack statistics by connection type
  - `ovs_vswitchd_conntrack_connection_limit`: Maximum number of connections allowed
  - `ovs_vswitchd_conntrack_n_connections`: Number of tracked connections
  - `ovs_vswitchd_conntrack_n_dccp`: Number of tracked DCCP connections
  - `ovs_vswitchd_conntrack_n_icmp`: Number of tracked ICMP connections
  - `ovs_vswitchd_conntrack_n_icmp6`: Number of tracked ICMPv6 connections
  - `ovs_vswitchd_conntrack_n_igmp`: Number of tracked IGMP connections
  - `ovs_vswitchd_conntrack_n_other`: Number of tracked connections of undefined type
  - `ovs_vswitchd_conntrack_n_sctp`: Number of tracked SCTP connections
  - `ovs_vswitchd_conntrack_n_tcp`: Number of tracked TCP connections
  - `ovs_vswitchd_conntrack_n_udp`: Number of tracked UDP connections
  - `ovs_vswitchd_conntrack_n_udplite`: Number of tracked UDPLite connections
  - `ovs_vswitchd_conntrack_tcp_seq_chk`: The TCP sequence checking mode: disabled(0) or enabled(1)

- **Datapath statistics**: Flow table and packet processing metrics
  - `ovs_vswitchd_datapath_bytes_total`: Number of bytes processed in total on this datapath
  - `ovs_vswitchd_datapath_packets_total`: Number of packets processed in total on this datapath
  - `ovs_vswitchd_datapath_offloaded_bytes_total`: Number of bytes processed in hardware on this datapath
  - `ovs_vswitchd_datapath_offloaded_packets_total`: Number of packets processed in hardware on this datapath
  - `ovs_vswitchd_datapath_tx_bytes_total`: Number of bytes emitted in total from this datapath
  - `ovs_vswitchd_datapath_tx_packets_total`: Number of packets emitted in total from this datapath
  - `ovs_vswitchd_datapath_tx_offloaded_bytes_total`: Total number of bytes emitted from this datapath and fully processed in hardware
  - `ovs_vswitchd_datapath_tx_offloaded_packets_total`: Total number of packets emitted from this datapath and fully processed in hardware
  - `ovs_vswitchd_datapath_hit_total`: Number of flow table matches
  - `ovs_vswitchd_datapath_missed_total`: Number of flow table misses
  - `ovs_vswitchd_datapath_lost_total`: Number of misses not sent to userspace
  - `ovs_vswitchd_datapath_cache_hit_total`: Number of mega flow mask cache hits for flow table matches
  - `ovs_vswitchd_datapath_mask_hit_total`: Number of mega flow masks visited for flow table matches
  - `ovs_vswitchd_datapath_n_flows`: Number of flows present
  - `ovs_vswitchd_datapath_n_masks`: Number of mega flow masks
  - `ovs_vswitchd_datapath_n_handlers`: Number of upcall handler threads
  - `ovs_vswitchd_datapath_n_revalidators`: Number of revalidator threads
  - `ovs_vswitchd_datapath_hw_offload_n_ct_bidir`: Number of bi-directional connections offloaded in hardware
  - `ovs_vswitchd_datapath_hw_offload_n_ct_unidir`: Number of uni-directional connections offloaded in hardware
  - `ovs_vswitchd_datapath_hw_offload_n_enqueued`: Number of hardware offload requests waiting to be processed
  - `ovs_vswitchd_datapath_hw_offload_n_inserted`: Number of hardware offload rules currently inserted

- **Interface statistics**: Per-interface metrics with detailed RX/TX counters
  - `ovs_vswitchd_interface_admin_state`: The administrative state of the interface: down(0) or up(1)
  - `ovs_vswitchd_interface_link_state`: The state of the interface link: down(0) or up(1)
  - `ovs_vswitchd_interface_link_speed`: The current speed of the interface link in Mbps
  - `ovs_vswitchd_interface_duplex`: The duplex mode of the interface: half(0) or full(1)
  - `ovs_vswitchd_interface_mtu`: The MTU of the interface
  - `ovs_vswitchd_interface_link_resets_total`: The number of time the interface link changed
  - `ovs_vswitchd_interface_rx_bytes_total`: The number of bytes received
  - `ovs_vswitchd_interface_rx_packets_total`: The number of packets received
  - `ovs_vswitchd_interface_tx_bytes_total`: The number of bytes transmitted
  - `ovs_vswitchd_interface_tx_packets_total`: The number of packets transmitted
  - `ovs_vswitchd_interface_rx_dropped_total`: Number of packets received but not processed, e.g. due to lack of resources or unsupported protocol. For hardware interface this counter should not include packets dropped by the device due to buffer exhaustion which are counted separately in rx_missed_errors
  - `ovs_vswitchd_interface_tx_dropped_total`: The number of packets dropped on their way to transmission, e.g. due to lack of resources
  - `ovs_vswitchd_interface_rx_errors_total`: Total number of bad packets received on this interface. This counter includes all rx_length_errors, rx_crc_errors, rx_frame_errors and other errors not otherwise counted
  - `ovs_vswitchd_interface_tx_errors_total`: Total number of transmit issues on this interface
  - `ovs_vswitchd_interface_rx_crc_errors_total`: The number of packets with CRC errors received by the interface
  - `ovs_vswitchd_interface_rx_frame_errors_total`: The number of received packets with frame alignment errors on the interface
  - `ovs_vswitchd_interface_rx_fifo_errors_total`: Receiver FIFO error counter. This statistics was used interchangeably with rx_over_errors but is not recommended for use in drivers for high speed interfaces. This statistics is used on software devices, e.g. to count software packets queue overflow or sequencing errors
  - `ovs_vswitchd_interface_rx_length_errors_total`: The number of packets dropped due to invalid length
  - `ovs_vswitchd_interface_rx_missed_errors_total`: The number of packets missed by the host due to lack of buffer space. This usually indicates that the host interface is slower than the hardware interface. This statistics corresponds to hardware events and is not used on software devices
  - `ovs_vswitchd_interface_rx_over_errors_total`: Receiver FIFO overflow event counter. This statistics was used interchangeably with rx_fifo_errors. This statistics corresponds to hardware events and is not commonly used on software devices
  - `ovs_vswitchd_interface_collisions_total`: The number of collisions during packet transmission
  - `ovs_vswitchd_interface_multicast_total`: The number of multicast packets received by the interface
  - `ovs_vswitchd_interface_ingress_policy_bit_rate`: Maximum receive rate in kbps on the interface. Disabled if set to 0
  - `ovs_vswitchd_interface_ingress_policy_bit_burst`: Maximum receive burst size in kb
  - `ovs_vswitchd_interface_ingress_policy_pkt_rate`: Maximum receive rate in pps on the interface. Disabled if set to 0
  - `ovs_vswitchd_interface_ingress_policy_pkt_burst`: Maximum receive burst size in number of packets
  - `ovs_vswitchd_interface_info`: A metric with a constant value '1' labeled with the driver name, version and firmware version of the interface
  - `ovs_vswitchd_interface_ifindex`: The ifindex of the interface
  - `ovs_vswitchd_interface_of_port`: The OpenFlow port ID associated with the interface

- **Poll thread (PMD) metrics**: Performance metrics for datapath poll threads
  - `ovs_vswitchd_poll_threads_n`: Number of polling threads
  - `ovs_vswitchd_poll_threads_packets_total`: Number of received packets
  - `ovs_vswitchd_poll_threads_recirculations_total`: Number of executed packet recirculations
  - `ovs_vswitchd_poll_threads_hit_total`: Number of flow table matches
  - `ovs_vswitchd_poll_threads_missed_total`: Number of flow table misses and upcall succeeded
  - `ovs_vswitchd_poll_threads_lost_total`: Number of flow table misses and upcall failed
  - `ovs_vswitchd_poll_threads_busy_cycles`: Percent of useful CPU cycles
  - `ovs_vswitchd_poll_threads_idle_cycles`: Percent of idle CPU cycles
  - `ovs_vswitchd_poll_threads_cycles_per_packet`: Average number of CPU cycles per packet
  - `ovs_vswitchd_poll_threads_busy_cycles_per_packet`: Average number of active CPU cycles per packet
  - `ovs_vswitchd_poll_threads_passes_per_packet`: Average number of datapath passes per packet
  - `ovs_vswitchd_poll_threads_recirc_per_packet`: Average number of recirculations per packet
  - `ovs_vswitchd_poll_threads_packets_per_batch`: Average number of packets per batch
  - `ovs_vswitchd_poll_threads_lookups_per_hit`: Average number of lookups per flow table hit

- **Memory metrics**: OVS process memory usage
  - `ovs_vswitchd_memory_in_use`: The amount of memory currently allocated in bytes
  - `ovs_vswitchd_memory_rss`: The process resident set size in bytes
  - `ovs_vswitchd_memory_vmsize`: The process virtual memory size in bytes
  - `ovs_vswitchd_memory_data`: The process sum of data and stack size in bytes
  - `ovs_vswitchd_memory_frag_factor`: The fragmentation factor of the process dynamic memory, defined as (rss/in_use)

- **Thread counts**:
  - `ovs_vswitchd_handler_n_threads`: Number of upcall handler threads in total
  - `ovs_vswitchd_revalidator_n_threads`: Number of revalidator threads in total

- **Scrape metadata**:
  - `ovs_vswitchd_scrape_duration_seconds`: Time elapsed to process this request in seconds
  - `ovs_vswitchd_metrics_histogram_read_errors_total`: Number of histogram reads that could not resolve without inconsistencies

#### `ovs-ofctl -O OpenFlow13 meter-stats <bridge>`
Per-meter statistics from the OpenFlow (1.3+) meter table. Bridges are discovered
each scrape with `ovs-vsctl list-br`, and `meter-stats` is run against each one
(there is no instance-wide variant). All metrics carry a `bridge` and `meter`
label; band metrics additionally carry a `band` label.

- `ovsdp_meter_flow_count`: Number of flows referencing this meter (gauge).
- `ovsdp_meter_packet_in_total`: Packets passed through this meter (counter).
- `ovsdp_meter_byte_in_total`: Bytes passed through this meter (counter).
- `ovsdp_meter_duration_seconds`: Time in seconds this meter has existed (gauge).
- `ovsdp_meter_band_packet_total`: Packets acted on by this meter band, e.g. dropped or remarked (counter).
- `ovsdp_meter_band_byte_total`: Bytes acted on by this meter band, e.g. dropped or remarked (counter).

#### `ovs-appctl memory/show`
High-level memory and thread/connection counts.

- `ovsdp_memory_handlers`: Number of OVS handler threads handling OpenFlow connections and upcalls.
- `ovsdp_memory_idl_cells_open_vswitch`: OVSDB cells in use for `Open_vSwitch` table (transaction/monitor memory).
- `ovsdp_memory_ofconns`: Active OpenFlow controller connections.
- `ovsdp_memory_ports`: Configured datapath ports (physical, virtual, and internal).
- `ovsdp_memory_revalidators`: Revalidator threads that periodically revalidate userspace datapath flows.
- `ovsdp_memory_rules`: Installed OpenFlow rules (software and hardware offloaded).
- `ovsdp_memory_udpif_keys`: Unique userspace datapath (udpif) flow keys handled in software.

### Running

- Build: `go build .`
- Run exporter: `./ovsdp-exporter -metrics.host :9000 -metrics.pathname /metrics`
- Scrape: visit `http://<host>:9000/metrics`

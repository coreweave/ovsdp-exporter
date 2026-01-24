### ovsdp-exporter
OVS datapath metric exporter. It executes a few `ovs-appctl` commands and exposes selected values as Prometheus metrics.

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

#### `ovs-appctl doca-pipe-group/dump`
DOCA pipe group unique item templates.

- `ovsdp_doca_unique_item_templates`: Number of unique item templates created from doca-pipe-group/dump.

#### `ovs-appctl metrics/show`
Prometheus-formatted metrics directly from OVS. All metrics returned by this command are exposed with their original names and labels, supporting both gauge and counter types. This includes a wide range of `ovs_vswitchd_*` prefixed metrics covering:

- **Bridge metrics**: Bridge configuration and flow counts
  - `ovs_vswitchd_bridge`: Bridge presence indicator (labeled by name and type)
  - `ovs_vswitchd_bridge_n_bridges`: Number of bridges
  - `ovs_vswitchd_bridge_n_flows`: Number of flows per bridge
  - `ovs_vswitchd_bridge_n_ports`: Number of ports per bridge

- **Connection tracking**: Conntrack statistics by connection type
  - `ovs_vswitchd_conntrack_connection_limit`: Maximum connections allowed
  - `ovs_vswitchd_conntrack_n_connections`: Total tracked connections
  - `ovs_vswitchd_conntrack_n_dccp`, `ovs_vswitchd_conntrack_n_icmp`, `ovs_vswitchd_conntrack_n_icmp6`, `ovs_vswitchd_conntrack_n_igmp`, `ovs_vswitchd_conntrack_n_other`, `ovs_vswitchd_conntrack_n_sctp`, `ovs_vswitchd_conntrack_n_tcp`, `ovs_vswitchd_conntrack_n_udp`, `ovs_vswitchd_conntrack_n_udplite`: Per-protocol connection counts
  - `ovs_vswitchd_conntrack_tcp_seq_chk`: TCP sequence checking mode

- **Datapath statistics**: Flow table and packet processing metrics
  - `ovs_vswitchd_datapath_bytes_total`, `ovs_vswitchd_datapath_packets_total`: Total bytes/packets processed
  - `ovs_vswitchd_datapath_offloaded_bytes_total`, `ovs_vswitchd_datapath_offloaded_packets_total`: Hardware-offloaded traffic
  - `ovs_vswitchd_datapath_tx_bytes_total`, `ovs_vswitchd_datapath_tx_packets_total`: Transmitted traffic
  - `ovs_vswitchd_datapath_tx_offloaded_bytes_total`, `ovs_vswitchd_datapath_tx_offloaded_packets_total`: Hardware-offloaded TX traffic
  - `ovs_vswitchd_datapath_hit_total`, `ovs_vswitchd_datapath_missed_total`, `ovs_vswitchd_datapath_lost_total`: Flow table lookup results
  - `ovs_vswitchd_datapath_cache_hit_total`, `ovs_vswitchd_datapath_mask_hit_total`: Megaflow mask cache statistics
  - `ovs_vswitchd_datapath_n_flows`, `ovs_vswitchd_datapath_n_masks`: Flow and mask counts
  - `ovs_vswitchd_datapath_n_handlers`, `ovs_vswitchd_datapath_n_revalidators`: Thread counts
  - `ovs_vswitchd_datapath_hw_offload_n_ct_bidir`, `ovs_vswitchd_datapath_hw_offload_n_ct_unidir`: Hardware-offloaded connection counts
  - `ovs_vswitchd_datapath_hw_offload_n_enqueued`, `ovs_vswitchd_datapath_hw_offload_n_inserted`: Hardware offload queue statistics

- **Interface statistics**: Per-interface metrics with detailed RX/TX counters
  - `ovs_vswitchd_interface_admin_state`, `ovs_vswitchd_interface_link_state`: Interface state
  - `ovs_vswitchd_interface_link_speed`, `ovs_vswitchd_interface_duplex`, `ovs_vswitchd_interface_mtu`: Link characteristics
  - `ovs_vswitchd_interface_link_resets_total`: Link reset count
  - `ovs_vswitchd_interface_rx_bytes_total`, `ovs_vswitchd_interface_rx_packets_total`: Received traffic
  - `ovs_vswitchd_interface_tx_bytes_total`, `ovs_vswitchd_interface_tx_packets_total`: Transmitted traffic
  - `ovs_vswitchd_interface_rx_dropped_total`, `ovs_vswitchd_interface_tx_dropped_total`: Dropped packets
  - `ovs_vswitchd_interface_rx_errors_total`, `ovs_vswitchd_interface_tx_errors_total`: Error counts
  - `ovs_vswitchd_interface_rx_crc_errors_total`, `ovs_vswitchd_interface_rx_frame_errors_total`, `ovs_vswitchd_interface_rx_fifo_errors_total`, `ovs_vswitchd_interface_rx_length_errors_total`, `ovs_vswitchd_interface_rx_missed_errors_total`, `ovs_vswitchd_interface_rx_over_errors_total`: Detailed RX error types
  - `ovs_vswitchd_interface_collisions_total`, `ovs_vswitchd_interface_multicast_total`: Additional interface statistics
  - `ovs_vswitchd_interface_ingress_policy_bit_rate`, `ovs_vswitchd_interface_ingress_policy_bit_burst`, `ovs_vswitchd_interface_ingress_policy_pkt_rate`, `ovs_vswitchd_interface_ingress_policy_pkt_burst`: Ingress policing parameters
  - `ovs_vswitchd_interface_info`, `ovs_vswitchd_interface_ifindex`, `ovs_vswitchd_interface_of_port`: Interface metadata

- **Poll thread (PMD) metrics**: Performance metrics for datapath poll threads
  - `ovs_vswitchd_poll_threads_packets_total`, `ovs_vswitchd_poll_threads_recirculations_total`: Packet processing counts
  - `ovs_vswitchd_poll_threads_hit_total`, `ovs_vswitchd_poll_threads_missed_total`, `ovs_vswitchd_poll_threads_lost_total`: Flow lookup results per thread
  - `ovs_vswitchd_poll_threads_busy_cycles`, `ovs_vswitchd_poll_threads_idle_cycles`: CPU cycle utilization
  - `ovs_vswitchd_poll_threads_cycles_per_packet`, `ovs_vswitchd_poll_threads_busy_cycles_per_packet`: Processing efficiency
  - `ovs_vswitchd_poll_threads_passes_per_packet`, `ovs_vswitchd_poll_threads_recirc_per_packet`: Pipeline complexity metrics
  - `ovs_vswitchd_poll_threads_packets_per_batch`: Batching efficiency
  - `ovs_vswitchd_poll_threads_lookups_per_hit`: Megaflow lookup efficiency

- **Memory metrics**: OVS process memory usage
  - `ovs_vswitchd_memory_in_use`, `ovs_vswitchd_memory_rss`, `ovs_vswitchd_memory_vmsize`, `ovs_vswitchd_memory_data`: Memory consumption statistics
  - `ovs_vswitchd_memory_frag_factor`: Memory fragmentation factor

- **Thread counts**:
  - `ovs_vswitchd_handler_n_threads`: Total upcall handler threads
  - `ovs_vswitchd_revalidator_n_threads`: Total revalidator threads

- **Scrape metadata**:
  - `ovs_vswitchd_scrape_duration_seconds`: Time taken to collect metrics
  - `ovs_vswitchd_metrics_histogram_read_errors_total`: Errors reading histogram metrics

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

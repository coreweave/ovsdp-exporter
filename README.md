### ovsdp-exporter
OVS datapath metric exporter. It executes a few `ovs-appctl` commands and exposes selected values as Prometheus metrics.

This exporter provides two categories of OVS metrics:
- **ovsdp_***
- **ovs_vswitchd_***

## OVS Metrics Categories

### OVSDP Metrics (ovsdp_*)
Commands executed and parsed to extract specific datapath performance metrics:

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

#### `ovs-appctl memory/show`
High-level memory and thread/connection counts.

- `ovsdp_memory_handlers`: Number of OVS handler threads handling OpenFlow connections and upcalls.
- `ovsdp_memory_idl_cells_open_vswitch`: OVSDB cells in use for `Open_vSwitch` table (transaction/monitor memory).
- `ovsdp_memory_ofconns`: Active OpenFlow controller connections.
- `ovsdp_memory_ports`: Configured datapath ports (physical, virtual, and internal).
- `ovsdp_memory_revalidators`: Revalidator threads that periodically revalidate userspace datapath flows.
- `ovsdp_memory_rules`: Installed OpenFlow rules (software and hardware offloaded).
- `ovsdp_memory_udpif_keys`: Unique userspace datapath (udpif) flow keys handled in software.

### OVS Vswitchd Metrics (ovs_vswitchd_*)
Comprehensive ovs-vswitchd metrics from `ovs-appctl metrics/show` in native Prometheus format:

- **Bridge metrics**: Number of bridges, flows, and ports per bridge
- **Connection tracking**: TCP, UDP, ICMP connection counts and limits  
- **Datapath statistics**: Packets, bytes, hits, misses, offload statistics
- **Interface metrics**: Admin state, link state, and statistics per interface
- **Hardware offload**: Rules inserted, connection tracking, and performance metrics
- **And more!**

### Running

- Build: `go build .`
- Run exporter: `./ovsdp-exporter -metrics.host :9000 -metrics.pathname /metrics`
- Scrape: visit `http://<host>:9000/metrics`

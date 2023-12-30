# devourer

![logo](https://github.com/m-mizutani/devourer/assets/605953/65b12b2d-3d79-4ba0-b312-de171210210b)

## What is this?

`devourer` is a tool to monitor network traffic and log network flows  to BigQuery.

## Usage

### Prerequisites

- **BigQuery dataset**: You need to create a BigQuery dataset to store network flows. See [here](https://cloud.google.com/bigquery/docs/datasets) for more details.
- **Service Account with BigQuery write permission**: You need to create a service account with BigQuery write permission. See [here](https://cloud.google.com/bigquery/docs/reference/libraries) for more details. You need to grant `roles/bigquery.dataEditor` role to the service account.
- **Service Account Key**: You need to create a service account key for the service account. See [here](https://cloud.google.com/iam/docs/creating-managing-service-account-keys) for more details.

### Installation

#### Binary

```bash
$ go install github.com/m-mizutani/devourer@latest
```

#### Docker image

```bash
$ docker pull ghcr.io/m-mizutani/devourer:latest
```

### Run

```bash
$ devourer capture -i <interface> \
    --bq-project-id <project-id> \
    --bq-dataset <dataset-name> \
    --bq-sa-key-file <service-account-key-file>
```

Or you can set environment variables instead of command line options.

```bash
$ export DEVOURER_BQ_PROJECT_ID=<project-id>
$ export DEVOURER_BQ_DATASET=<dataset-name>
$ export DEVOURER_BQ_SA_KEY_FILE=<service-account-key-file>
$ devourer capture -i <interface>
```

## How it works

`devourer` captures network packets and extract network flows from them. A network flow is a sequence of packets that have the same 5-tuple (source IP, destination IP, source port, destination port, protocol). `devourer` aggregates network flows and store them to BigQuery.

`devourer` does not monitor status of transport layer such as TCP, UDP, and ICMP. It only monitors network layer and application layer. So, `devourer` cannot detect TCP connection status such as SYN, SYN-ACK, and FIN. `devourer` determines closing flow by timeout. (default timeout is 120 seconds). After the timeout, `devourer` inserts the flow to BigQuery with `flow_logs` status.

## Network flow schema

`devourer` stores network flows to BigQuery with the following schema as `flow_logs` table. The schema will be created automatically when you run `devourer` for the first time.

| Column name | Type | Description |
| --- | --- | --- |
| id | string | Unique ID of the network flow. |
| protocol | string | Protocol of the network flow. |
| src_addr | string | Source IP address of the network flow. |
| dst_addr | string | Destination IP address of the network flow. |
| src_port | int | Source port of the network flow. |
| dst_port | int | Destination port of the network flow. |
| first_seen_at | timestamp | Timestamp when the network flow was first seen. |
| last_seen_at | timestamp | Timestamp when the network flow was last seen. |
| src_bytes | int | Number of bytes sent from the source to the destination. |
| dst_bytes | int | Number of bytes sent from the destination to the source. |
| src_packets | int | Number of packets sent from the source to the destination. |
| dst_packets | int | Number of packets sent from the destination to the source. |
| status | string | Status of the network flow. |

## License

Apache License 2.0

# kafka-topic-analyser

## Features
- Request-Response pattern tracing
  - Response Time / Latency across a pair of topics
  - Correlation by
    - Key
    - Header
    - Combination of Headers
- Schema Stats for a topic
- Time distribution 
- Consumer Stats / Topic
  - Response times
  - Latency
- Producer Stats / Topic
  - Throughput

## Spec
- Store data in duckdb
- Load from topic to duckdb
- Analyse and store in duckdb
- Query from duckdb
- Single binary/docker-image distribution
- Single dependency: duckdb

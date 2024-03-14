# Ethereum Optimism Load Generation Tool

## Description

A load testing tool designed to simulate and measure the performance of Ethereum Optimism networks under various transaction loads.

## Usage

To generate a load, the following command will dump transactions into `load.json` which configured by `predefined/simpletransfer.toml`

```bash
go run . export \
    --http.rpc-url http://0.0.0.0:9545 \
    --engine.rpc-url http://0.0.0.0:8551 \
    --engine-jwt-secret <JWT token without 0x-prefixed> \
    --output-path load.json \
    --config-path predefined/simpletransfer.toml \
```


To simulate the load, the following command will send block payloads into the chain, every block payload contains `201` transactions


```bash
go run . import \
    --engine.rpc-url http://0.0.0.0:8551 \
    --engine-jwt-secret <JWT token without 0x-prefixed> \
    --input-path load.json \
    --block-txs-count 201
```

## Configuration

See [`predefined/simpletransfer.toml`](./predefined/simpletransfer.toml) for configuration example.

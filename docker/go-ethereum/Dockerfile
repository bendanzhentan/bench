FROM ethereum/client-go:latest

COPY genesis.json /tmp
RUN geth init /tmp/genesis.json

RUN rm -f ~/.ethereum/geth/nodekey

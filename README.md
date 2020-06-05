ckb udt cli
===========

## Build

```bash
export GOPROXY=https://goproxy.io
go mod download
go build .
```

## Usage

### Issue

```bash
./ckb-udt-cli issue -c config.yaml -k YOUR_PRIVATE_KEY -a AMOUNT
```

### Create anyone can pay cell

```bash
./ckb-udt-cli create-cell -c config.yaml -k YOUR_PRIVATE_KEY -u UUID
```

### Transfer

```bash
./ckb-udt-cli transfer -c config.yaml -k YOUR_PRIVATE_KEY -u UUID -t RECIPIENT_ADDRESS -a AMOUNT
```

### Balance

```bash
./ckb-udt-cli balance -c config.yaml -u UUID -a ADDRESS
```

## Example data

https://explorer.nervos.org/aggron/sudt/0xe3be4fb98ec914886c6525abac97e1f8769c59492636a1d35955e9163ef46efa

uuid: `0x94bbc8327e16d195de87815c391e7b9131e80419c51a405a0b21227c6ee05129`


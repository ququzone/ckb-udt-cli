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
./ckb-udt-cli issue -c config.yaml  -k YOUR_PRIVATE_KEY -a AMOUNT
```

### Transfer

```bash
./ckb-udt-cli transfer -c config.yaml -k YOUR_PRIVATE_KEY -u UUID -t RECIPIENT_ADDRESS -a AMOUNT
```


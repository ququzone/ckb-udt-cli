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
./ckb-udt-cli issue -k YOUR_PRIVATE_KEY -a AMOUNT
```

### Transfer

```bash
./ckb-udt-cli transfer -k YOUR_PRIVATE_KEY -t RECIPIENT_ADDRESS -a AMOUNT
```


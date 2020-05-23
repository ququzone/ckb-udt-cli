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

uuid: `0xcd64ecc1fa2570073cbe9b2dfda7974288b564f323b4cd07e9d84fef22d62661`


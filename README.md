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

### transfer

```bash
./ckb-udt-cli transfer -c config.yaml -k YOUR_PRIVATE_KEY -u UUID -t RECIPIENT_ADDRESS -a AMOUNT
```

### balance

```bash
./ckb-udt-cli balance -c config.yaml -u UUID -a ADDRESS
```

## Example data

uuid: `0x6a242b57227484e904b4e08ba96f19a623c367dcbd18675ec6f2a71a0ff4ec26`


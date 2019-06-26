# Nebula CLI

## Usage
```bash
./nebula -h
```

## Build

```bash
make client
make build
```

## Configure
```bash
mkdir -p ${HOME}/.config/nebula
echo 'apiHostAddr: http://api.staging.nebula.insights.puppet.net' > ${HOME}/.config/nebula/config.yaml
```
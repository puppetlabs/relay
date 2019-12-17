# Nebula CLI

## Usage
```bash
./bin/nebula-cli -h
```

## Build

```bash
make
```

Note that you will need these two tools:

- `npm i -g api-spec-converter`
- `GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger`

## Configure
```bash
mkdir -p ${HOME}/.config/nebula
echo 'apiHostAddr: http://api.nebula.puppet.com' > ${HOME}/.config/nebula/config.yaml
```

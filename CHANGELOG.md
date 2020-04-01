# [3.3.0](https://github.com/puppetlabs/relay/compare/v3.2.0...v3.3.0) (2020-02-12)


### New

* Add ability to delete workflow secrets ([ba2cced0126ac897225f0d7851f0e5ac8850fa10](https://github.com/puppetlabs/relay/commit/ba2cced0126ac897225f0d7851f0e5ac8850fa10))

# [3.2.0](https://github.com/puppetlabs/relay/compare/v3.1.1...v3.2.0) (2020-02-12)


### Update

* New API version, support listing workflow secrets, event sources, and canceling workflow run ([49db3ed37269dd408556b569256ad01b906fd80d](https://github.com/puppetlabs/relay/commit/49db3ed37269dd408556b569256ad01b906fd80d))

### Upgrade

* Bump handlebars from 4.1.2 to 4.5.3 ([0117b0db39a30c82bbb40e9570330a45f51f43d9](https://github.com/puppetlabs/relay/commit/0117b0db39a30c82bbb40e9570330a45f51f43d9))
* Bump lodash from 4.17.11 to 4.17.15 ([6da9d5a5b96d78c4957d5c988976461da6bee110](https://github.com/puppetlabs/relay/commit/6da9d5a5b96d78c4957d5c988976461da6bee110))
* Bump npm from 6.9.2 to 6.13.4 ([957205b7b58bddc5c11004e0bd3e24e7cf9f8446](https://github.com/puppetlabs/relay/commit/957205b7b58bddc5c11004e0bd3e24e7cf9f8446))
* Fix dependencies to avoid checksum errors ([baab2cf511beaf4be57cf530b3b2c13919f3a5d0](https://github.com/puppetlabs/relay/commit/baab2cf511beaf4be57cf530b3b2c13919f3a5d0))

## [3.1.1](https://github.com/puppetlabs/relay/compare/v3.1.0...v3.1.1) (2019-10-24)


### Fix

* Correctly encode parameter and secret values with Base64 if needed ([efa9cddedb985cc05357035078c67ec3ab1e2360](https://github.com/puppetlabs/relay/commit/efa9cddedb985cc05357035078c67ec3ab1e2360))

### Upgrade

* Bump go-swagger version to resolve flattening issues ([5c2fad08cff4de95cfc25f76e84be2ee9c1ae454](https://github.com/puppetlabs/relay/commit/5c2fad08cff4de95cfc25f76e84be2ee9c1ae454))

# [3.1.0](https://github.com/puppetlabs/relay/compare/v3.0.0...v3.1.0) (2019-10-06)


### New

* Support parameters ([df3e87f42b749c500896f49b346f4f717fa984a3](https://github.com/puppetlabs/relay/commit/df3e87f42b749c500896f49b346f4f717fa984a3))

### Upgrade

* Bump api-spec-converter back to upstream; update API client ([c13bd2d5aa79165be4e0469edfce7acc749d41c4](https://github.com/puppetlabs/relay/commit/c13bd2d5aa79165be4e0469edfce7acc749d41c4))

# [3.0.0](https://github.com/puppetlabs/relay/compare/v2.1.0...v3.0.0) (2019-08-22)


### Breaking

* Update to latest API endpoints ([5f313dd5120d72f04289af85089a9486dbb7a256](https://github.com/puppetlabs/relay/commit/5f313dd5120d72f04289af85089a9486dbb7a256))

### Update

* Added a config path flag to the global flags list ([120703ca175e9ddfca59e43d59528195d94c67f6](https://github.com/puppetlabs/relay/commit/120703ca175e9ddfca59e43d59528195d94c67f6))

# [2.1.0](https://github.com/puppetlabs/relay/compare/v2.0.0...v2.1.0) (2019-08-13)


### New

* nebula-cli workflow {logs,status} subcommands ([9fd905300119007d3d95e2859d5c1ab4e172ac4d](https://github.com/puppetlabs/relay/commit/9fd905300119007d3d95e2859d5c1ab4e172ac4d))

# [2.0.0](https://github.com/puppetlabs/relay/compare/v1.2.0...v2.0.0) (2019-08-01)


### Breaking

* Renaming output binary to nebula from nebula-cli ([2c9804441024700f3ce343b008f8b25577c71c33](https://github.com/puppetlabs/relay/commit/2c9804441024700f3ce343b008f8b25577c71c33))

# [1.2.0](https://github.com/puppetlabs/relay/compare/v1.1.0...v1.2.0) (2019-07-24)


### Update

* Add timeout to `workflow run` command ([640dcee1d512510edde0e03946a05b56274398bf](https://github.com/puppetlabs/relay/commit/640dcee1d512510edde0e03946a05b56274398bf))
* Add timeout to `workflow run` command ([0a7622b3c613baf5c3d5bfb86c360ace00d7c4b1](https://github.com/puppetlabs/relay/commit/0a7622b3c613baf5c3d5bfb86c360ace00d7c4b1))

# [1.1.0](https://github.com/puppetlabs/relay/compare/v1.0.1...v1.1.0) (2019-07-24)


### Chore

* Removing unused workflow file loading. (#4) ([257e0aa05a76c5954d88b120a761b142b104a18f](https://github.com/puppetlabs/relay/commit/257e0aa05a76c5954d88b120a761b142b104a18f)), closes [#4](https://github.com/puppetlabs/relay/issues/4)

### New

* Add ability to set secrets on a workflow ([35fdadc51bd7ffa781b2b6f168840b10ef0ef8fe](https://github.com/puppetlabs/relay/commit/35fdadc51bd7ffa781b2b6f168840b10ef0ef8fe))

## [1.0.1](https://github.com/puppetlabs/relay/compare/v1.0.0...v1.0.1) (2019-07-03)


### Chore

* Release 1.0.0 ([1a64ea78fce057f19af7dff677434a84fbb6e4ae](https://github.com/puppetlabs/relay/commit/1a64ea78fce057f19af7dff677434a84fbb6e4ae))

### Fix

* Use the correct version when building binaries ([ee7af2e798471a2bd227ae36daad33ba361b4b34](https://github.com/puppetlabs/relay/commit/ee7af2e798471a2bd227ae36daad33ba361b4b34))

# 1.0.0 (2019-07-03)


### Chore

* Release 1.0.0 ([cb10fcef57b918d0f36af8f654c1bdda13b1711b](https://github.com/puppetlabs/relay/commit/cb10fcef57b918d0f36af8f654c1bdda13b1711b))

### Update

* Add GitHub releases ([0cb1424bc79272f98d97dfd362ed9a30128c210a](https://github.com/puppetlabs/relay/commit/0cb1424bc79272f98d97dfd362ed9a30128c210a))

# 1.0.0 (2019-07-03)


### Update

* Add GitHub releases ([0cb1424bc79272f98d97dfd362ed9a30128c210a](https://github.com/puppetlabs/relay/commit/0cb1424bc79272f98d97dfd362ed9a30128c210a))

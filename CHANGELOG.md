### v2.0.0-beta3 -> v2.0.0-beta4

#### Documentation

 - [`32f47b8`](https://github.com/deis/workflow-manager/commit/32f47b89c255d05cbcd084f068b09dec94f7ad7b) badge: added code-beat badge
 - [`23327fa`](https://github.com/deis/workflow-manager/commit/23327faf3659c7a127c300c1a139f8a1f08a26b5) CHANGELOG.md: update for v2.0.0-beta3

### v2.0.0-beta2 -> v2.0.0-beta3

#### Features

 - [`6a7acfe`](https://github.com/deis/workflow-manager/commit/6a7acfe9d13bb7c289d24e3f01b0029c10df0cec) types: add generic data property to types.Version

#### Fixes

 - [`a21b95f`](https://github.com/deis/workflow-manager/commit/a21b95fc764d444a19a375bb1fe99a5c98a67432) rootfs: copy only the built binary in the Dockerfile

#### Maintenance

 - [`61a12e2`](https://github.com/deis/workflow-manager/commit/61a12e2200c7082724fce80213547f04cf50d29e) .travis.yml: Stop deploying images from Travis
 - [`af030e9`](https://github.com/deis/workflow-manager/commit/af030e9bc7b7379ed15228b7ebbe8431674070c6) changelog: update the changelog for the beta2

### v2.0.0-beta1 -> v2.0.0-beta2

#### Features

 - [`9ba6db7`](https://github.com/deis/workflow-manager/commit/9ba6db795f8a490d760f11e48ef7070ca832fe22) types: added "train" field to Version
 - [`b948fbd`](https://github.com/deis/workflow-manager/commit/b948fbd63c93f4d80d928c966bc8679fdcef8cad) _scripts: add CHANGELOG.md and generator script
 - [`a731083`](https://github.com/deis/workflow-manager/commit/a7310831d52fcd9a3a5f937d0e27943488a29d13) README.md: add quay.io container badge

### v2.0.0-beta1

#### Features

 - [`941cb10`](https://github.com/deis/workflow-manager/commit/941cb102db4af4eb9cc3270af5220ab617c08a14) Makefile: enable immutable (git-based

#### Fixes

 - [`61a9d8c`](https://github.com/deis/workflow-manager/commit/61a9d8c11367e21b7f1f9f76708a517cb9828c61) _scripts/deploy.sh: tag dockerhub images with quay.io prefix
 - [`46c1e3c`](https://github.com/deis/workflow-manager/commit/46c1e3cc240c012b9ded3fb6434998135b24858e) ci: CI pipeline now works for local+travis
 - [`29f0a63`](https://github.com/deis/workflow-manager/commit/29f0a63b727b1fa10c811fd794254b577333ae81) README.md: make go report card badge go green
 - [`bd4c6ff`](https://github.com/deis/workflow-manager/commit/bd4c6ffe9150d24b5d29c7a311e46faba2068704) travis: use common config and update go deps

#### Maintenance

 - [`d04aeb9`](https://github.com/deis/workflow-manager/commit/d04aeb9bab83da21b91a4f93132e96fc664e0820) Dockerfile: suspect travis needs explicit reference to bin/boot

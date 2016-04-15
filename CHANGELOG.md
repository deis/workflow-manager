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

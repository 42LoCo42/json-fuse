# json-fuse
A little FUSE experiment. Values are read-only for now, maybe I'll add write support later if I need it.
Should be pretty fast:

- [rancher-data.json](https://github.com/rancher/rancher/releases/download/v2.6.8/rancher-data.json): 50 ms
- [large-file.json](https://github.com/json-iterator/test-data/raw/master/large-file.json): 1.1 s

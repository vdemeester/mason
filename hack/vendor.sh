#!/usr/bin/env bash
set -e

cd "$(dirname "$BASH_SOURCE")/.."
rm -rf vendor/
source 'hack/.vendor-helpers.sh'

clone git github.com/vdemeester/libmason ed03930c0cd794cffd8fc5d9837f14944b32ee4f
clone git github.com/spf13/cobra 4c05eb1145f16d0e6bb4a3e1b6d769f4713cb41f
clone git github.com/spf13/pflag 8f6a28b0916586e7f22fe931ae2fcfc380b1c0e6
clone git github.com/Sirupsen/logrus v0.10.0
clone git github.com/docker/distribution 5bbf65499960b184fe8e0f045397375e1a6722b8
clone git github.com/docker/docker 534753663161334baba06f13b8efa4cad22b5bc5
clone git github.com/docker/go-units f2d77a61e3c169b43402a0a1e84f06daf29b8190
clone git github.com/docker/go-connections 990a1a1a70b0da4c4cb70e117971a4f0babfbf1a
clone git github.com/opencontainers/runc cc29e3dded8e27ba8f65738f40d251c885030a28
clone git github.com/docker/engine-api 62043eb79d581a32ea849645277023c550732e52
clone git golang.org/x/net 47990a1ba55743e6ef1affd3a14e5bac8553615d https://github.com/golang/net.git

clean && mv vendor/src/* vendor

#!/usr/bin/env bash
set -e

cd "$(dirname "$BASH_SOURCE")/.."
rm -rf vendor/
source 'hack/.vendor-helpers.sh'

clone git github.com/vdemeester/libmason 82c2daef2c5f253c1b5a98c1a737d2a68e10c4e2
clone git github.com/spf13/cobra 4c05eb1145f16d0e6bb4a3e1b6d769f4713cb41f
clone git github.com/spf13/pflag 8f6a28b0916586e7f22fe931ae2fcfc380b1c0e6
clone git github.com/Sirupsen/logrus v0.9.0
clone git github.com/docker/distribution d06d6d3b093302c02a93153ac7b06ebc0ffd1793
clone git github.com/docker/docker v1.11.0
clone git github.com/docker/go-units 651fc226e7441360384da338d0fd37f2440ffbe3
clone git github.com/docker/go-connections v0.2.0
clone git github.com/opencontainers/runc 7b6c4c418d5090f4f11eee949fdf49afd15838c9
clone git github.com/docker/engine-api c9b5a47a3a6fea97d0a2242307ad25a3331ce03b
clone git golang.org/x/net 47990a1ba55743e6ef1affd3a14e5bac8553615d https://github.com/golang/net.git

clean && mv vendor/src/* vendor

# 🏓 go-build  [![GoDoc](https://godoc.org/github.com/bobheadxi/go-build?status.svg)](https://godoc.org/github.com/bobheadxi/go-build) [![Build Status](https://travis-ci.com/bobheadxi/go-build.svg?branch=master)](https://travis-ci.com/bobheadxi/go-build) [![codecov](https://codecov.io/gh/bobheadxi/go-build/branch/master/graph/badge.svg)](https://codecov.io/gh/bobheadxi/go-build)

This is an isolated version of the old `build` package I wrote for [Inertia](https://github.com/ubclaunchpad/inertia), a command-line application that enables easy, self-hosted continuous deployment. It provides a Golang API for executing Dockerfile, docker-compose, and Herokuish builds and deployments.

This package is still WIP and in the process of being cleaned up, but the bulk of the code has been in active use in Inertiad for a while - see package [`inertiad/build`](https://github.com/ubclaunchpad/inertia/tree/master/daemon/inertiad/build).

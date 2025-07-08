<!-- markdownlint-disable -->

<div align="center">

<img src="https://rcdn.sunist.cn/avatar/dusk.png" style="width: 200px" alt="dusk-scheduler"/>

<br>

<img src="https://img.shields.io/github/go-mod/go-version/alioth-center/dusk-scheduler" alt="Go Version"/>
<img src="https://goreportcard.com/badge/github.com/alioth-center/dusk-scheduler" alt="Go Report Card"/>
<img src="https://img.shields.io/github/license/alioth-center/dusk-scheduler" alt="License"/>
<img src="https://wakatime.com/badge/github/alioth-center/dusk-scheduler.svg" alt="WakaTime"/>

# **D**usk **U**niversal **S**cene **K**it - Scheduler

</div>

<!-- markdownlint-restore -->


## Summary

`dusk-system` is a distributed graphics rendering system implemented in Golang, designed to convert HTML content into static images (PNG/JPG) via a simple HTTP API. The system comprises two core components:

- [`dusk-scheduler`](https://github.com/alioth-center/dusk-scheduler): Manages rendering job creation, task scheduling, and result retrieval
- [`dusk-painter`](https://github.com/alioth-center/dusk-painter): Executes the actual HTML rendering process

This architecture enables efficient parallelization of rendering workloads, providing scalable on-demand image generation for web content, reports, or visualizations.

## Contributors

![Contributors](https://contrib.rocks/image?repo=alioth-center/dusk-scheduler&max=1000)
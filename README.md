# URL Shortener: Golang vs. Python Benchmark

A cloud-native URL shortener benchmark designed to isolate and measure exactly how different concurrency paradigms impact system performance under heavy socket saturation.

## 🎯 Benchmark Objective

The primary aim of this project is to quantify the performance difference between **asynchronous single-threaded event models** (Python/FastAPI) and **native multi-threaded M:N scheduling models** (Go/Gin).

## System Architecture & Constraints

To ensure a perfectly fair comparison, both application engines were isolated using Docker Linux control groups (`cgroups`) to enforce strict resource caps.

* **Host Processor:** Intel Core i3-1005G1 (2 Cores, 4 Threads)
* **Application CPU Limit:** `1.5 Cores` max (shared kernel allocation via Docker)
* **Database CPU Limit:** `0.5 Cores` max (dedicated to Redis)
* **Traffic Generator:** `wrk` (2 Threads, 100 Concurrent HTTP Sockets, 30-Second Saturation)

**Note** : Cores here refer to logical processors (threads).

## The Showdown: Benchmark Results

The following metrics were captured under raw socket saturation on native host loops using identical connection load thresholds.

| Performance Metric | Python (FastAPI + Uvicorn) | Go (Gin Router + Go 1.26) | Delta / Impact |
| :--- | :--- | :--- | :--- |
| **Throughput (Requests/sec)** | `3,968.71` RPS | **`37,067.88` RPS** | **Go is ~9.34x Faster** |
| **Total Requests Volume** | `119,085` total | **`1,112,453` total** | **+1 Million requests served** |
| **Median Latency (50%)** | `25.00 ms` | **`2.52 ms`** | **10x lower baseline latency** |
| **Tail Latency (P99)** | `32.64 ms` | **`10.15 ms`** | **3.2x lower spike latency** |
| **Data Transfer Rate** | `4.32 MB/sec` | **`4.77 MB/sec`** | **Optimal network efficiency** |

---

## 🧠 Architectural Post-Mortem

### Why Go Achieved a 9.3x Throughput:

1. **M:N Scheduling vs. Single-Threaded Shuffling:** Python uses a single-threaded event loop (`asyncio`) that shuffles tasks sequentially on a single core's time-slice. Go utilizes an internal runtime scheduler that multiplexes thousands of lightweight goroutines (~2KB baseline) across multiple cores, unlocking true hardware parallelism within our `1.5 CPU` limit.

2. **Compiled Machine Code vs. Interpreted Bytecode:** FastAPI processes requests through dynamic python interpretation layers and serialization checks. Go compiles down directly to a native binary, executing raw machine instructions with near-zero runtime abstraction or interpreter overhead.

### Design Choices: Why Gin?

While Go's standard `net/http` library is incredibly fast, it lacks robust wildcard routing parameters natively. To capture dynamic short codes (`/:short_code`) efficiently without manual string parsing boilerplate, **Gin** was selected. It uses an optimized **Radix Tree (Prefix Tree)** algorithm for zero-allocation routing lookups, offering the blistering performance of custom frameworks while retaining 100% compatibility with Go's standard library ecosystem and memory safety guarantees.

---

## How to Reproduce Locally

### 1. Spin up the Infrastructure

To spin up a specific engine along with its isolated Redis backplane, execute:

```bash
# To test the Python Pipeline
docker compose up -d redis-service python-api

# To test the Go Pipeline
docker compose up -d redis-service go-api
```

### 2. Install and Execute the Load Generator

Install `wrk` on your host system to drive the traffic engine.

> **Windows Users:** Run these steps and benchmark commands inside a WSL 2 terminal instance. Since `wrk` relies on native Linux network event notification primitives, it must be executed within a Linux-subsystem environment.

```bash
# Ubuntu / Debian / WSL2
sudo apt update && sudo apt install wrk -y

# Arch Linux
sudo pacman -S wrk

# macOS (Homebrew)
brew install wrk
```

Once installed, run the high-velocity benchmark against the targeted network gateway:

```bash
# Test Python Pipeline (Port 8000)
wrk -t2 -c100 -d30s --latency http://localhost:8000/docs

# Test Go Pipeline (Port 8081)
wrk -t2 -c100 -d30s --latency http://localhost:8081/
```

// Copyright 2025 Omkar Prabhu
#include <benchmark/benchmark.h>
#include <spdlog/spdlog.h>

#include <string>
#include <unordered_map>

#include "worker/components.h"

using components::ComponentConfig;

static void BM_ComponentsParseComponentConfig(
    benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  std::string config = R"([
    {"name":"component1"},
    {"name":"component2","source":"source2"},
    {"name":"component3","source":"source3","params":"params3"}
])";
  for (auto _ : state) {
    std::unordered_map<std::string, ComponentConfig> result =
        components::ParseComponentConfig(config);
    benchmark::DoNotOptimize(result);
  }
}

BENCHMARK(BM_ComponentsParseComponentConfig)->ThreadPerCpu();
BENCHMARK_MAIN();

// Copyright 2025 Omkar Prabhu
#include <benchmark/benchmark.h>
#include <gmock/gmock.h>
#include <gtest/gtest.h>
#include <spdlog/spdlog.h>

#include <memory>
#include <string>
#include <unordered_map>
#include <utility>

#include "places_test.cc" // NOLINT
#include "worker/places.h"

static void BM_PlacesNoOp(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  for (auto _ : state) {
    Places places;
    std::unordered_map<std::string, std::string> output =
        places.ReverseGeocode("", "", std::nullopt, std::nullopt);
    benchmark::DoNotOptimize(output);
  }
}

static void BM_OpenStreetMapNoOp(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  for (auto _ : state) {
    Places places;
    std::unordered_map<std::string, std::string> output =
        places.ReverseGeocode("", "", std::nullopt, std::nullopt);
    benchmark::DoNotOptimize(output);
  }
}

static void BM_OpenStreetMapErrorResponse(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  for (auto _ : state) {
    auto mock_client = std::make_shared<MockHttpClient>();
    OpenStreetMap osm(mock_client);
    cpr::Response mock_response;
    mock_response.error.message = "mock error";
    mock_response.status_code = 200;
    mock_response.text = "{}";
    EXPECT_CALL(*mock_client, Get(::testing::_, ::testing::_))
        .WillOnce(::testing::Return(mock_response));
    std::unordered_map<std::string, std::string> output =
        osm.ReverseGeocode("", "", 1.23, 4.56);
    benchmark::DoNotOptimize(output);
  }
}

static void BM_OpenStreetMapIncorrectStatusCode(
    benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  for (auto _ : state) {
    auto mock_client = std::make_shared<MockHttpClient>();
    OpenStreetMap osm(mock_client);
    cpr::Response mock_response;
    mock_response.status_code = 500;
    mock_response.text = "{}";
    EXPECT_CALL(*mock_client, Get(::testing::_, ::testing::_))
        .WillOnce(::testing::Return(mock_response));
    std::unordered_map<std::string, std::string> output =
        osm.ReverseGeocode("", "", 0.0, 0.0);
    benchmark::DoNotOptimize(output);
  }
}

static void BM_OpenStreetMapSuccess(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  for (auto _ : state) {
    auto mock_client = std::make_shared<MockHttpClient>();
    OpenStreetMap osm(mock_client);
    cpr::Response mock_response;
    mock_response.status_code = 200;
    mock_response.text =
        R"({"address":{"suburb":"K/W Ward","city":"Mumbai","state":"Maharashtra","postcode":"402205","country":"India"}})";
    EXPECT_CALL(*mock_client, Get(::testing::_, ::testing::_))
        .WillOnce(::testing::Return(mock_response));
    std::unordered_map<std::string, std::string> output =
        osm.ReverseGeocode("", "", 1.23, 4.56);
    benchmark::DoNotOptimize(output);
  }
}

BENCHMARK(BM_PlacesNoOp)->ThreadPerCpu();
BENCHMARK(BM_OpenStreetMapNoOp)->ThreadPerCpu();
BENCHMARK(BM_OpenStreetMapErrorResponse)->ThreadPerCpu();
BENCHMARK(BM_OpenStreetMapIncorrectStatusCode)->ThreadPerCpu();
BENCHMARK(BM_OpenStreetMapSuccess)->ThreadPerCpu();
BENCHMARK_MAIN();

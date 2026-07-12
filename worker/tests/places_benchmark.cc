// Copyright 2025 Omkar Prabhu
#include <benchmark/benchmark.h>
#include <cpr/cpr.h>
#include <gmock/gmock.h>
#include <gtest/gtest.h>
#include <spdlog/spdlog.h>

#include <iostream>
#include <memory>
#include <optional>
#include <string>
#include <unordered_map>
#include <utility>

#include "places_test.cc" // NOLINT
#include "protos/api.pb.h"
#include "worker/components.h"
#include "worker/places.h"

using components::ComponentConfig;
using components::places::HttpClientInterface;
using components::places::OpenStreetMap;
using components::places::Places;
using services::api::APIClient;
using ::testing::_;
using ::testing::Invoke;
using ::testing::NiceMock;
using ::testing::Return;

static void BM_PlacesEmptyInput(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  Places places;
  for (auto _ : state) {
    std::unordered_map<std::string, std::string> output =
        places.ReverseGeocode("", "", "", "", "");
    benchmark::DoNotOptimize(output);
  }
}

static void BM_OpenStreetMapEmptyInput(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  Places places;
  for (auto _ : state) {
    std::unordered_map<std::string, std::string> output =
        places.ReverseGeocode("", "", "", "", "");
    benchmark::DoNotOptimize(output);
  }
}

static void BM_OpenStreetMapErrorResponse(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  auto mock_api_stub = std::make_unique<NiceMock<MockAPIStub>>();
  MockAPIStub* mock_stub = mock_api_stub.get();
  std::shared_ptr<APIClient> mock_api_client =
      std::make_shared<APIClient>(std::move(mock_api_stub));
  EXPECT_CALL(*mock_stub, SaveMediaItemPlace(_, _, _))
      .WillRepeatedly(
          Invoke([&](grpc::ClientContext*, const MediaItemPlaceRequest&,
                     google::protobuf::Empty*) { return grpc::Status::OK; }));
  std::shared_ptr<MockHttpClient> mock_http_client =
      std::make_shared<MockHttpClient>();
  OpenStreetMap osm(mock_http_client, mock_api_client);
  cpr::Response mock_response;
  mock_response.error.message = "mock error";
  mock_response.status_code = 200;
  mock_response.text = "{}";
  EXPECT_CALL(*mock_http_client, Get(::testing::_, ::testing::_))
      .WillRepeatedly(::testing::Return(mock_response));
  for (auto _ : state) {
    std::unordered_map<std::string, std::string> output =
        osm.ReverseGeocode("", "", "", "1.23", "4.56");
    benchmark::DoNotOptimize(output);
  }
}

static void BM_OpenStreetMapIncorrectStatusCode(
    benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  auto mock_api_stub = std::make_unique<NiceMock<MockAPIStub>>();
  MockAPIStub* mock_stub = mock_api_stub.get();
  std::shared_ptr<APIClient> mock_api_client =
      std::make_shared<APIClient>(std::move(mock_api_stub));
  EXPECT_CALL(*mock_stub, SaveMediaItemPlace(_, _, _))
      .WillRepeatedly(
          Invoke([&](grpc::ClientContext*, const MediaItemPlaceRequest&,
                     google::protobuf::Empty*) { return grpc::Status::OK; }));
  std::shared_ptr<MockHttpClient> mock_http_client =
      std::make_shared<MockHttpClient>();
  OpenStreetMap osm(mock_http_client, mock_api_client);
  cpr::Response mock_response;
  mock_response.status_code = 500;
  mock_response.text = "{}";
  EXPECT_CALL(*mock_http_client, Get(::testing::_, ::testing::_))
      .WillRepeatedly(::testing::Return(mock_response));
  for (auto _ : state) {
    std::unordered_map<std::string, std::string> output =
        osm.ReverseGeocode("", "", "", "0.0", "0.0");
    benchmark::DoNotOptimize(output);
  }
}

static void BM_OpenStreetMapSuccess(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  auto mock_api_stub = std::make_unique<NiceMock<MockAPIStub>>();
  MockAPIStub* mock_stub = mock_api_stub.get();
  std::shared_ptr<APIClient> mock_api_client =
      std::make_shared<APIClient>(std::move(mock_api_stub));
  EXPECT_CALL(*mock_stub, SaveMediaItemPlace(_, _, _))
      .WillRepeatedly(
          Invoke([&](grpc::ClientContext*, const MediaItemPlaceRequest&,
                     google::protobuf::Empty*) { return grpc::Status::OK; }));
  std::shared_ptr<MockHttpClient> mock_http_client =
      std::make_shared<MockHttpClient>();
  OpenStreetMap osm(mock_http_client, mock_api_client);
  cpr::Response mock_response;
  mock_response.status_code = 200;
  mock_response.text =
      R"({"address":{"suburb":"K/W Ward","city":"Mumbai","state":"Maharashtra","postcode":"402205","country":"India"}})";
  EXPECT_CALL(*mock_http_client, Get(::testing::_, ::testing::_))
      .WillRepeatedly(::testing::Return(mock_response));
  for (auto _ : state) {
    std::unordered_map<std::string, std::string> output =
        osm.ReverseGeocode("", "", "", "1.23", "4.56");
    benchmark::DoNotOptimize(output);
  }
}

BENCHMARK(BM_PlacesEmptyInput)->ThreadPerCpu();
BENCHMARK(BM_OpenStreetMapEmptyInput)->ThreadPerCpu();
BENCHMARK(BM_OpenStreetMapErrorResponse)->ThreadPerCpu();
BENCHMARK(BM_OpenStreetMapIncorrectStatusCode)->ThreadPerCpu();
BENCHMARK(BM_OpenStreetMapSuccess)->ThreadPerCpu();

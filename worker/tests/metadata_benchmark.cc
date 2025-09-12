// Copyright 2025 Omkar Prabhu
#include <benchmark/benchmark.h>
#include <gmock/gmock.h>
#include <gtest/gtest.h>
#include <spdlog/spdlog.h>

#include <memory>
#include <string>
#include <unordered_map>
#include <utility>

#include "metadata_test.cc" // NOLINT
#include "worker/metadata.h"

using components::metadata::ExifToolClientInterface;
using components::metadata::Metadata;

static void BM_MetadataInit(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  ComponentConfig config = ComponentConfig("", "");
  for (auto _ : state) {
    auto metadata = components::metadata::Init(config);
    benchmark::DoNotOptimize(metadata);
  }
}

static void BM_MetadataEmptyData(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  for (auto _ : state) {
    auto mock_client = std::make_shared<MockExifToolClient>();
    Metadata metadata(mock_client);
    std::unordered_map<std::string, std::string> mock_data;
    EXPECT_CALL(*mock_client, Extract(::testing::_))
        .WillOnce(::testing::Return(mock_data));
    std::unordered_map<std::string, std::string> output =
        metadata.Extract("", "", "");
    benchmark::DoNotOptimize(output);
  }
}

static void BM_MetadataSuccess(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  for (auto _ : state) {
    auto mock_client = std::make_shared<MockExifToolClient>();
    Metadata metadata(mock_client);
    std::unordered_map<std::string, std::string> mock_data = {
        {"VideoFrameRate", "30"},
        {"GPSLatitude", "19 deg 13' 11.99\" N"},
        {"GPSLongitude", "73 deg 6' 19.19\" E"},
        {"ExifImageWidth", "4032"},
        {"ImageWidth", "226"},
        {"ExifImageHeight", "3024"},
        {"ImageHeight", "4032"},
        {"Make", "Apple"},
        {"Model", "iPhone 15 Pro"},
        {"FocalLength", "4.2 mm"},
        {"FNumber", "1.6"},
        {"ISO", "640"},
        {"MIMEType", "image/heic"},
        {"ExposureTime", "1/25"},
        {"LivePhotoVideoIndex", "1112547328"},
        {"DateCreated", "2022:04:03 12:56:11"}};
    EXPECT_CALL(*mock_client, Extract(::testing::_))
        .WillOnce(::testing::Return(mock_data));
    std::unordered_map<std::string, std::string> output =
        metadata.Extract("", "", "");
    benchmark::DoNotOptimize(output);
  }
}

BENCHMARK(BM_MetadataInit)->ThreadPerCpu();
BENCHMARK(BM_MetadataEmptyData)->ThreadPerCpu();
BENCHMARK(BM_MetadataSuccess)->ThreadPerCpu();

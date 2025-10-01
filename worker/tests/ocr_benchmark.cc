// Copyright 2025 Omkar Prabhu
#include <benchmark/benchmark.h>
#include <gmock/gmock.h>
#include <gtest/gtest.h>
#include <spdlog/spdlog.h>

#include <iostream>
#include <memory>
#include <optional>
#include <string>
#include <unordered_map>
#include <utility>
#include <vector>

#include "ocr_test.cc" // NOLINT
#include "protos/api.pb.h"
#include "protos/api_mock.grpc.pb.h"
#include "worker/components.h"
#include "worker/ocr.h"

using components::ComponentConfig;
using components::ocr::ModelInferenceInterface;
using components::ocr::OCR;
using components::ocr::PaddlePaddle;
using components::ocr::PaddlePaddleModel;
using services::api::APIClient;
using ::testing::_;
using ::testing::Invoke;
using ::testing::NiceMock;
using ::testing::Return;

static void BM_OCRInit(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  ComponentConfig config = ComponentConfig("paddlepaddle", "params");
  for (auto _ : state) {
    auto ocr = components::ocr::Init(config, nullptr);
    benchmark::DoNotOptimize(ocr);
  }
}

static void BM_OCREmptyInput(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  for (auto _ : state) {
    components::ocr::OCR ocr;
    std::unordered_map<std::string, std::string> output =
        ocr.Extract("", "", "", "");
    benchmark::DoNotOptimize(output);
  }
}

static void BM_OCRPaddlePaddleEmptyInput(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  for (auto _ : state) {
    std::shared_ptr<MockModelInference> mock_model =
        std::make_shared<MockModelInference>();
    PaddlePaddle paddlepaddle(mock_model, nullptr);
    std::unordered_map<std::string, std::string> output =
        paddlepaddle.Extract("", "", "", "");
    benchmark::DoNotOptimize(output);
  }
}

static void BM_OCRPaddlePaddleError(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  for (auto _ : state) {
    std::shared_ptr<MockModelInference> mock_model =
        std::make_shared<MockModelInference>();
    EXPECT_CALL(*mock_model, Run(::testing::_))
        .WillOnce(::testing::Throw(std::runtime_error("some error")));
    PaddlePaddle paddlepaddle(mock_model, nullptr);
    std::unordered_map<std::string, std::string> output =
        paddlepaddle.Extract("", "", "", "file_path");
    benchmark::DoNotOptimize(output);
  }
}

static void BM_OCRPaddlePaddleSuccess(benchmark::State& state) { // NOLINT
  spdlog::set_level(spdlog::level::off);
  for (auto _ : state) {
    std::shared_ptr<MockModelInference> mock_model =
        std::make_shared<MockModelInference>();
    std::vector<std::pair<std::string, float>> mock_response = {
        {"thats what she said", 99.86},
        {"boy have you lost your mind", 98.12},
        {"spiderface", 75.43}};
    EXPECT_CALL(*mock_model, Run(::testing::_))
        .WillOnce(::testing::Return(mock_response));
    PaddlePaddle paddlepaddle(mock_model, nullptr);
    std::unordered_map<std::string, std::string> output =
        paddlepaddle.Extract("", "", "", "file_path");
    benchmark::DoNotOptimize(output);
  }
}

BENCHMARK(BM_OCRInit)->ThreadPerCpu();
BENCHMARK(BM_OCREmptyInput)->ThreadPerCpu();
BENCHMARK(BM_OCRPaddlePaddleEmptyInput)->ThreadPerCpu();
BENCHMARK(BM_OCRPaddlePaddleError)->ThreadPerCpu();
BENCHMARK(BM_OCRPaddlePaddleSuccess)->ThreadPerCpu();

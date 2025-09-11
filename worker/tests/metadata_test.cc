// Copyright 2025 Omkar Prabhu
#include "worker/metadata.h"

#include <gmock/gmock.h>
#include <gtest/gtest.h>
#include <spdlog/spdlog.h>

#include <iostream>
#include <memory>
#include <optional>
#include <string>
#include <unordered_map>
#include <utility>

#include "worker/components.h"

using components::ComponentConfig;
using components::metadata::ExifToolClientInterface;
using components::metadata::Metadata;

class MockExifToolClient : public ExifToolClientInterface {
 public:
  MOCK_METHOD((std::unordered_map<std::string, std::string>), Extract,
              (const std::string& file_path), (override));
};

void assertMetadataResult(std::unordered_map<std::string, std::string> expected,
                          std::unordered_map<std::string, std::string> actual) {
  EXPECT_EQ(expected.size(), actual.size());
  for (const auto& [key, value] : expected) {
    EXPECT_EQ(value, actual[key]);
  }
}

TEST(MetadataTest, Init) {
  spdlog::set_level(spdlog::level::off);
  auto metadata = components::metadata::Init(ComponentConfig("", ""));
  ASSERT_TRUE(metadata != nullptr);
}

TEST(MetadataTest, EmptyData) {
  spdlog::set_level(spdlog::level::off);
  auto mock_client = std::make_shared<MockExifToolClient>();
  Metadata metadata(mock_client);
  std::unordered_map<std::string, std::string> mock_data;
  EXPECT_CALL(*mock_client, Extract(::testing::_))
      .WillOnce(::testing::Return(mock_data));
  std::unordered_map<std::string, std::string> result =
      metadata.Extract("", "", "");
  assertMetadataResult({}, result);
}

TEST(MetadataTest, Success) {
  spdlog::set_level(spdlog::level::off);
  auto mock_client = std::make_shared<MockExifToolClient>();
  Metadata metadata(mock_client);
  std::unordered_map<std::string, std::string> mock_data;
  mock_data["key"] = "value";
  EXPECT_CALL(*mock_client, Extract(::testing::_))
      .WillOnce(::testing::Return(mock_data));
  std::unordered_map<std::string, std::string> result =
      metadata.Extract("", "", "");
  assertMetadataResult(
      {
          {"key", "value"},
      },
      result);
}

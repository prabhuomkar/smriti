// Copyright 2025 Omkar Prabhu
#include "worker/places.h"

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

using components::places::HttpClientInterface;
using components::places::OpenStreetMap;
using components::places::Places;

class MockHttpClient : public HttpClientInterface {
 public:
  MOCK_METHOD(cpr::Response, Get,
              (const cpr::Url& url, const cpr::Header& headers), (override));
};

void assertPlacesResult(std::unordered_map<std::string, std::string> expected,
                        std::unordered_map<std::string, std::string> actual) {
  EXPECT_EQ(expected.size(), actual.size());
  for (const auto& [key, value] : expected) {
    EXPECT_EQ(value, actual[key]);
  }
}

TEST(PlacesTest, EmptyInput) {
  spdlog::set_level(spdlog::level::off);
  Places places;
  std::unordered_map<std::string, std::string> result =
      places.ReverseGeocode("", "", std::nullopt, std::nullopt);
  assertPlacesResult({}, result);
}

TEST(OpenStreetMapTest, EmptyInput) {
  spdlog::set_level(spdlog::level::off);
  auto mock_client = std::make_shared<MockHttpClient>();
  OpenStreetMap osm(mock_client);
  std::unordered_map<std::string, std::string> result =
      osm.ReverseGeocode("", "", std::nullopt, std::nullopt);
  assertPlacesResult({}, result);
}

TEST(OpenStreetMapTest, ErrorResponse) {
  spdlog::set_level(spdlog::level::off);
  auto mock_client = std::make_shared<MockHttpClient>();
  OpenStreetMap osm(mock_client);
  cpr::Response mock_response;
  mock_response.error.message = "mock error";
  mock_response.status_code = 200;
  mock_response.text = "{}";
  EXPECT_CALL(*mock_client, Get(::testing::_, ::testing::_))
      .WillOnce(::testing::Return(mock_response));
  std::unordered_map<std::string, std::string> result =
      osm.ReverseGeocode("", "", 1.23, 4.56);
  assertPlacesResult({}, result);
}

TEST(OpenStreetMapTest, IncorrectStatusCode) {
  spdlog::set_level(spdlog::level::off);
  auto mock_client = std::make_shared<MockHttpClient>();
  OpenStreetMap osm(mock_client);
  cpr::Response mock_response;
  mock_response.status_code = 500;
  mock_response.text = "{}";
  EXPECT_CALL(*mock_client, Get(::testing::_, ::testing::_))
      .WillOnce(::testing::Return(mock_response));
  std::unordered_map<std::string, std::string> result =
      osm.ReverseGeocode("", "", 0.0, 0.0);
  assertPlacesResult({}, result);
}

TEST(OpenStreetMapTest, Success) {
  spdlog::set_level(spdlog::level::off);
  auto mock_client = std::make_shared<MockHttpClient>();
  OpenStreetMap osm(mock_client);
  cpr::Response mock_response;
  mock_response.status_code = 200;
  mock_response.text =
      R"({"address":{"suburb":"K/W Ward","city":"Mumbai","state":"Maharashtra","postcode":"402205","country":"India"}})";
  EXPECT_CALL(*mock_client, Get(::testing::_, ::testing::_))
      .WillOnce(::testing::Return(mock_response));
  std::unordered_map<std::string, std::string> result =
      osm.ReverseGeocode("", "", 1.23, 4.56);
  assertPlacesResult(
      {
          {"postcode", "402205"},
          {"country", "India"},
          {"locality", "Mumbai"},
          {"area", "K/W Ward"},
      },
      result);
}

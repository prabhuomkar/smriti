// Copyright 2025 Omkar Prabhu
#include "worker/components.h"

#include <gtest/gtest.h>
#include <spdlog/spdlog.h>

#include <string>
#include <unordered_map>

using components::ComponentConfig;

TEST(ComponentsTest, ParseComponentConfigSuccess) {
  spdlog::set_level(spdlog::level::off);
  std::string config = R"([
        {"name":"component1"},
        {"name":"component2","source":"source2"},
        {"name":"component3","source":"source3","params":"params3"}
])";
  std::unordered_map<std::string, ComponentConfig> result =
      components::ParseComponentConfig(config);
  ASSERT_EQ(result.size(), 3);
  EXPECT_TRUE(result.at("component1").source.empty());
  EXPECT_TRUE(result.at("component1").params.empty());
  EXPECT_EQ(result.at("component2").source, "source2");
  EXPECT_TRUE(result.at("component2").params.empty());
  EXPECT_EQ(result.at("component3").source, "source3");
  EXPECT_EQ(result.at("component3").params, "params3");
}

TEST(ComponentsTest, ParseComponentConfigFailure) {
  spdlog::set_level(spdlog::level::off);
  std::string config = R"()";
  std::unordered_map<std::string, ComponentConfig> result =
      components::ParseComponentConfig(config);
  ASSERT_EQ(result.size(), 0);
  config = R"([])";
  result = components::ParseComponentConfig(config);
  ASSERT_EQ(result.size(), 0);
}

// Copyright 2025 Omkar Prabhu
#include "worker/components.h"

#include <simdjson.h>
#include <spdlog/spdlog.h>

#include <string>
#include <unordered_map>

namespace components {

ComponentConfig::ComponentConfig(const std::string& source,
                                 const std::string& params)
    : source(source), params(params) {}

std::unordered_map<std::string, ComponentConfig> ParseComponentConfig(
    const std::string& config) {
  if (config.empty()) {
    spdlog::error("empty component config");
    return {};
  }
  std::unordered_map<std::string, ComponentConfig> configs;
  simdjson::ondemand::parser parser;
  simdjson::padded_string padded_config = simdjson::padded_string(config);
  simdjson::ondemand::document doc = parser.iterate(padded_config);
  auto items = doc.get_array();
  if (items.error() != simdjson::SUCCESS) {
    spdlog::error("error parsing component config: {}",
                  simdjson::error_message(items.error()));
    return {};
  }
  for (auto item : items) {
    std::string name = std::string(item["name"].get_string().value());
    std::string source = "";
    if (item["source"].error() == simdjson::SUCCESS) {
      source = std::string(item["source"].get_string().value());
    }
    std::string params = "";
    if (item["params"].error() == simdjson::SUCCESS) {
      params = std::string(item["params"].get_string().value());
    }
    configs.insert({name, ComponentConfig(source, params)});
  }
  return configs;
}

} // namespace components

// Copyright 2025 Omkar Prabhu
#pragma once

#include <string>
#include <unordered_map>

namespace components {

class Component {
 public:
  Component(const std::string& name, const std::string& source,
            const std::string& params)
      : name(name), source(source), params(params) {}
  std::string name;
  std::string source;
  std::string params;
};

class ComponentConfig {
 public:
  ComponentConfig(const std::string& source, const std::string& params);
  std::string source;
  std::string params;
};

std::unordered_map<std::string, ComponentConfig> ParseComponentConfig(
    const std::string& config);

} // namespace components

// Copyright 2025 Omkar Prabhu
#pragma once

#include <cpr/cpr.h>

#include <memory>
#include <optional>
#include <string>
#include <unordered_map>
#include <utility>

#include "worker/components.h"

namespace components {

namespace places {

class HttpClientInterface {
 public:
  virtual ~HttpClientInterface() = default;
  virtual cpr::Response Get(const cpr::Url& url,
                            const cpr::Header& headers) = 0;
};

class HttpClient : public HttpClientInterface {
 public:
  cpr::Response Get(const cpr::Url& url, const cpr::Header& headers) override {
    return cpr::Get(url, headers);
  }
};

class Places {
 public:
  Places() {}
  virtual std::unordered_map<std::string, std::string> ReverseGeocode(
      const std::string& id, const std::string& mediaitem_id,
      const std::string& latitude, const std::string& longitude);
};

class OpenStreetMap : public Places {
 public:
  explicit OpenStreetMap(
      std::shared_ptr<HttpClientInterface> http_client,
      std::string url =
          "https://nominatim.openstreetmap.org/"
          "reverse.php?zoom=18&format=jsonv2&lat={lat}&lon={lon}",
      int timeout = 60)
      : http_client_(http_client), url_(std::move(url)), timeout_(timeout) {}
  std::unordered_map<std::string, std::string> ReverseGeocode(
      const std::string& id, const std::string& mediaitem_id,
      const std::string& latitude, const std::string& longitude) override;

 private:
  std::shared_ptr<HttpClientInterface> http_client_;
  std::string url_;
  int timeout_;
};

std::string Format(std::string input,
                   const std::unordered_map<std::string, std::string>& values);

std::shared_ptr<Places> Init(const ComponentConfig& config);
} // namespace places

} // namespace components

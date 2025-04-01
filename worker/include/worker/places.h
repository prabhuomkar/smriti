// Copyright 2025 Omkar Prabhu
#pragma once

#include <cpr/cpr.h>

#include <memory>
#include <optional>
#include <string>
#include <unordered_map>
#include <utility>

class HttpClientInterface {
 public:
  virtual ~HttpClientInterface() = default;
  virtual cpr::Response Get(const cpr::Url& url,
                            const cpr::Header& headers) = 0;
};

class Places {
 public:
  Places() {}
  virtual std::unordered_map<std::string, std::string> ReverseGeocode(
      const std::string& user_id, const std::string& mediaitem_id,
      std::optional<double> latitude, std::optional<double> longitude);
};

class OpenStreetMap : public Places {
 public:
  OpenStreetMap(std::shared_ptr<HttpClientInterface> http_client,
                std::string url =
                    "https://nominatim.openstreetmap.org/"
                    "reverse.php?zoom=18&format=jsonv2&lat={lat}&lon={lon}",
                int timeout = 60)
      : http_client_(http_client), url_(std::move(url)), timeout_(timeout) {}
  std::unordered_map<std::string, std::string> ReverseGeocode(
      const std::string& user_id, const std::string& mediaitem_id,
      std::optional<double> latitude, std::optional<double> longitude) override;

 private:
  std::shared_ptr<HttpClientInterface> http_client_;
  std::string url_;
  int timeout_;
  std::string Format(
      std::string input,
      const std::unordered_map<std::string, std::string>& values);
};

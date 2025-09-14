// Copyright 2025 Omkar Prabhu
#include "worker/places.h"

#include <cpr/cpr.h>
#include <simdjson.h>
#include <spdlog/spdlog.h>

#include <memory>
#include <optional>
#include <string>
#include <unordered_map>
#include <utility>

#include "worker/api_client.h"
#include "worker/components.h"

using services::api::APIClient;

namespace components {

namespace places {

std::unordered_map<std::string, std::string> Places::ReverseGeocode(
    const std::string& id, const std::string& user_id,
    const std::string& mediaitem_id, const std::string& latitude,
    const std::string& longitude) {
  return {};
}

std::unordered_map<std::string, std::string> OpenStreetMap::ReverseGeocode(
    const std::string& id, const std::string& user_id,
    const std::string& mediaitem_id, const std::string& latitude,
    const std::string& longitude) {
  if (latitude == "" && longitude == "") {
    return {};
  }

  cpr::Response r = http_client_->Get(
      cpr::Url{Format(url_, {{"lat", latitude}, {"lon", longitude}})},
      cpr::Header{{"User-Agent", "smriti-worker"},
                  {"Accept-Language", "en-GB,en-US"}});
  if (r.error.message != "") {
    spdlog::error("error in openstreetmap response: {}", r.error.message);
    return {};
  }
  if (r.status_code != 200) {
    spdlog::error("error in openstreetmap response: {} {}", r.status_code,
                  r.text);
    return {};
  }

  spdlog::debug("openstreetmap response: {} {}", r.status_code, r.text);

  simdjson::ondemand::parser parser;
  simdjson::padded_string padded_body(r.text);
  simdjson::ondemand::document doc = parser.iterate(padded_body);

  std::unordered_map<std::string, std::string> result;

  if (doc["address"].error() == simdjson::SUCCESS) {
    simdjson::ondemand::object address = doc["address"].get_object();

    if (address["postcode"].error() == simdjson::SUCCESS) {
      result["postcode"] =
          std::string(address["postcode"].get_string().value());
    }
    if (address["country"].error() == simdjson::SUCCESS) {
      result["country"] = std::string(address["country"].get_string().value());
    }

    std::string locality = "";
    if (address["village"].error() == simdjson::SUCCESS) {
      locality = std::string(address["village"].get_string().value());
    }
    if (locality == "" && address["town"].error() == simdjson::SUCCESS) {
      locality = std::string(address["town"].get_string().value());
    }
    if (address["city"].error() == simdjson::SUCCESS) {
      locality = std::string(address["city"].get_string().value());
    }
    if (locality == "" &&
        address["municipality"].error() == simdjson::SUCCESS) {
      locality = std::string(address["municipality"].get_string().value());
    }
    if (locality == "" && address["county"].error() == simdjson::SUCCESS) {
      locality = std::string(address["county"].get_string().value());
    }
    result["locality"] = locality;

    std::string area;
    if (address["suburb"].error() == simdjson::SUCCESS) {
      area = std::string(address["suburb"].get_string().value());
    }
    if (area == "" && address["quarter"].error() == simdjson::SUCCESS) {
      area = std::string(address["quarter"].get_string().value());
    }
    if (area == "" && address["neighbourhood"].error() == simdjson::SUCCESS) {
      area = std::string(address["neighbourhood"].get_string().value());
    }
    if (area == "" && address["road"].error() == simdjson::SUCCESS) {
      area = std::string(address["road"].get_string().value());
    }
    result["area"] = area;
  } else {
    spdlog::error("error in openstreetmap response: {}",
                  simdjson::error_message(doc["address"].error()));
    return {};
  }

  MediaItemPlaceRequest request;
  request.set_userid(user_id);
  request.set_mediaitemid(mediaitem_id);
  request.set_postcode(result["postcode"]);
  request.set_country(result["country"]);
  request.set_locality(result["locality"]);
  request.set_area(result["area"]);
  bool ok = api_client_->SaveMediaItemPlace(request);

  return result;
}

std::string Format(std::string input,
                   const std::unordered_map<std::string, std::string>& values) {
  for (const auto& [key, value] : values) {
    std::string placeholder = "{" + key + "}";
    size_t pos = input.find(placeholder);
    if (pos != std::string::npos) {
      input.replace(pos, placeholder.length(), value);
    }
  }
  return input;
}

std::shared_ptr<Places> Init(const ComponentConfig& config,
                             std::shared_ptr<APIClient> api_client) {
  if (config.source == "openstreetmap") {
    return std::make_shared<OpenStreetMap>(std::make_shared<HttpClient>(),
                                           api_client);
  }
  return nullptr;
}

} // namespace places

} // namespace components

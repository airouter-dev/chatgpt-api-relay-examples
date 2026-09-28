#include "image_client.hpp"

#include <curl/curl.h>

#include <algorithm>
#include <cctype>
#include <cstdlib>
#include <fstream>
#include <initializer_list>
#include <iterator>
#include <mutex>
#include <optional>
#include <sstream>
#include <stdexcept>
#include <string_view>
#include <utility>

namespace image_demo {
namespace {

constexpr std::size_t kMaxResponseBytes = 100U * 1024U * 1024U;

struct HttpResponse {
    long status = 0;
    std::string body;
    std::string content_type;
};

void ensure_curl_global() {
    static std::once_flag once;
    std::call_once(once, [] {
        const CURLcode result = curl_global_init(CURL_GLOBAL_DEFAULT);
        if (result != CURLE_OK) {
            throw std::runtime_error("curl_global_init failed");
        }
    });
}

std::size_t write_body(char* data, std::size_t size, std::size_t count, void* user_data) {
    auto* body = static_cast<std::string*>(user_data);
    const std::size_t bytes = size * count;
    if (bytes > kMaxResponseBytes || body->size() > kMaxResponseBytes - bytes) {
        return 0; // causes CURLE_WRITE_ERROR rather than allocating unbounded data
    }
    body->append(data, bytes);
    return bytes;
}

std::size_t write_headers(char* data, std::size_t size, std::size_t count, void* user_data) {
    auto* content_type = static_cast<std::string*>(user_data);
    const std::string line(data, size * count);
    constexpr std::string_view prefix = "content-type:";
    std::string lower;
    lower.reserve(line.size());
    for (const char character : line) {
        lower.push_back(static_cast<char>(std::tolower(static_cast<unsigned char>(character))));
    }
    if (lower.rfind(prefix, 0) == 0) {
        const auto value_start = line.find(':');
        if (value_start != std::string::npos) {
            *content_type = line.substr(value_start + 1);
            while (!content_type->empty() && std::isspace(static_cast<unsigned char>(content_type->front()))) {
                content_type->erase(content_type->begin());
            }
            while (!content_type->empty() && std::isspace(static_cast<unsigned char>(content_type->back()))) {
                content_type->pop_back();
            }
        }
    }
    return size * count;
}

HttpResponse perform_request(const std::string& url,
                             const Config& config,
                             struct curl_slist* headers,
                             const std::function<void(CURL*)>& configure) {
    ensure_curl_global();
    CURL* curl = curl_easy_init();
    if (curl == nullptr) {
        throw std::runtime_error("curl_easy_init failed");
    }
    HttpResponse response;
    char error_buffer[CURL_ERROR_SIZE] = {};
    curl_easy_setopt(curl, CURLOPT_ERRORBUFFER, error_buffer);
    curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
    curl_easy_setopt(curl, CURLOPT_FOLLOWLOCATION, 1L);
    curl_easy_setopt(curl, CURLOPT_CONNECTTIMEOUT, std::min(config.timeout_seconds, 15L));
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, config.timeout_seconds);
    curl_easy_setopt(curl, CURLOPT_NOSIGNAL, 1L);
    curl_easy_setopt(curl, CURLOPT_USERAGENT, "image-generation-multilang-cpp/1.0");
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_body);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response.body);
    curl_easy_setopt(curl, CURLOPT_HEADERFUNCTION, write_headers);
    curl_easy_setopt(curl, CURLOPT_HEADERDATA, &response.content_type);
    if (headers != nullptr) {
        curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
    }
    configure(curl);
    const CURLcode result = curl_easy_perform(curl);
    if (result != CURLE_OK) {
        const std::string details = error_buffer[0] != '\0' ? error_buffer : curl_easy_strerror(result);
        curl_easy_cleanup(curl);
        throw std::runtime_error(std::string("HTTP request failed: ") + details);
    }
    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &response.status);
    curl_easy_cleanup(curl);
    return response;
}

std::string first_environment(std::initializer_list<const char*> names,
                              const std::string& fallback = {}) {
    for (const char* name : names) {
        if (const char* value = std::getenv(name); value != nullptr && *value != '\0') {
            return value;
        }
    }
    return fallback;
}

std::string trim(std::string value) {
    const auto not_space = [](unsigned char character) { return !std::isspace(character); };
    value.erase(value.begin(), std::find_if(value.begin(), value.end(), not_space));
    value.erase(std::find_if(value.rbegin(), value.rend(), not_space).base(), value.end());
    return value;
}

std::optional<std::string> json_string(const std::string& json, const std::string& key) {
    const std::string needle = "\"" + key + "\"";
    std::size_t key_position = json.find(needle);
    while (key_position != std::string::npos) {
        std::size_t position = json.find(':', key_position + needle.size());
        if (position == std::string::npos) {
            return std::nullopt;
        }
        ++position;
        while (position < json.size() && std::isspace(static_cast<unsigned char>(json[position]))) {
            ++position;
        }
        if (position < json.size() && json[position] == '"') {
            ++position;
            std::string result;
            bool escaped = false;
            for (; position < json.size(); ++position) {
                const char character = json[position];
                if (escaped) {
                    switch (character) {
                    case '"': result.push_back('"'); break;
                    case '\\': result.push_back('\\'); break;
                    case '/': result.push_back('/'); break;
                    case 'b': result.push_back('\b'); break;
                    case 'f': result.push_back('\f'); break;
                    case 'n': result.push_back('\n'); break;
                    case 'r': result.push_back('\r'); break;
                    case 't': result.push_back('\t'); break;
                    default: result.push_back(character); break;
                    }
                    escaped = false;
                } else if (character == '\\') {
                    escaped = true;
                } else if (character == '"') {
                    return result;
                } else {
                    result.push_back(character);
                }
            }
            return std::nullopt;
        }
        key_position = json.find(needle, key_position + needle.size());
    }
    return std::nullopt;
}

std::string response_error(long status, const std::string& body) {
    std::ostringstream message;
    message << "image API returned HTTP " << status;
    if (const auto api_message = json_string(body, "message"); api_message.has_value()) {
        message << ": " << *api_message;
    } else if (!trim(body).empty()) {
        std::string snippet = trim(body);
        if (snippet.size() > 512) {
            snippet.resize(512);
            snippet += "...";
        }
        message << ": " << snippet;
    }
    return message.str();
}

} // namespace

Config Config::from_environment() {
    Config config;
    config.api_key = first_environment({"OPENAI_API_KEY", "AI_ROUTER_API_KEY"});
    config.base_url = first_environment({"OPENAI_BASE_URL", "AI_ROUTER_BASE_URL"}, config.base_url);
    config.model = first_environment({"IMAGE_MODEL", "AI_ROUTER_MODEL"}, config.model);
    config.size = first_environment({"IMAGE_SIZE", "AI_ROUTER_SIZE"}, config.size);
    config.response_format = first_environment({"IMAGE_RESPONSE_FORMAT"}, config.response_format);
    if (const std::string timeout = first_environment({"IMAGE_API_TIMEOUT", "AI_ROUTER_IMAGE_TIMEOUT"}); !timeout.empty()) {
        try {
            const long parsed = std::stol(timeout);
            if (parsed > 0) {
                config.timeout_seconds = parsed;
            }
        } catch (const std::exception&) {
            // Keep the safe default when an environment value is malformed.
        }
    }
    while (!config.base_url.empty() && config.base_url.back() == '/') {
        config.base_url.pop_back();
    }
    return config;
}

ImageClient::ImageClient(Config config) : config_(std::move(config)) {
    while (!config_.base_url.empty() && config_.base_url.back() == '/') {
        config_.base_url.pop_back();
    }
    if (config_.timeout_seconds <= 0) {
        config_.timeout_seconds = 120;
    }
}

std::string ImageClient::endpoint(const std::string& path) const {
    return config_.base_url + "/" + (path.empty() || path.front() != '/' ? path : path.substr(1));
}

std::vector<unsigned char> ImageClient::generate(const std::string& prompt) const {
    if (trim(prompt).empty()) {
        throw std::invalid_argument("prompt must not be empty");
    }
    const std::string body = "{\"model\":\"" + detail::escape_json(config_.model) +
                             "\",\"prompt\":\"" + detail::escape_json(prompt) +
                             "\",\"size\":\"" + detail::escape_json(config_.size) +
                             "\",\"response_format\":\"" + detail::escape_json(config_.response_format) + "\"}";
    struct curl_slist* headers = nullptr;
    headers = curl_slist_append(headers, "Content-Type: application/json");
    headers = curl_slist_append(headers, "Accept: application/json");
    if (config_.api_key.empty()) {
        curl_slist_free_all(headers);
        throw std::invalid_argument("API key is missing; set OPENAI_API_KEY or AI_ROUTER_API_KEY");
    }
    headers = curl_slist_append(headers, ("Authorization: Bearer " + config_.api_key).c_str());
    HttpResponse response;
    try {
        response = perform_request(endpoint("/images/generations"), config_, headers,
            [&body](CURL* curl) {
                curl_easy_setopt(curl, CURLOPT_POST, 1L);
                curl_easy_setopt(curl, CURLOPT_POSTFIELDS, body.c_str());
                curl_easy_setopt(curl, CURLOPT_POSTFIELDSIZE, static_cast<long>(body.size()));
            });
    } catch (...) {
        curl_slist_free_all(headers);
        throw;
    }
    curl_slist_free_all(headers);
    return resolve_image_response(response.status, response.body);
}

std::vector<unsigned char> ImageClient::edit(const std::string& image_path,
                                             const std::string& prompt) const {
    if (trim(image_path).empty()) {
        throw std::invalid_argument("image path must not be empty");
    }
    if (trim(prompt).empty()) {
        throw std::invalid_argument("edit prompt must not be empty");
    }
    if (config_.api_key.empty()) {
        throw std::invalid_argument("API key is missing; set OPENAI_API_KEY or AI_ROUTER_API_KEY");
    }
    std::ifstream image(image_path, std::ios::binary);
    if (!image) {
        throw std::runtime_error("cannot open edit image: " + image_path);
    }

    ensure_curl_global();
    CURL* curl = curl_easy_init();
    if (curl == nullptr) {
        throw std::runtime_error("curl_easy_init failed");
    }
    // A MIME handle owns the file path during the transfer; curl reads it only
    // after perform, so this stream is used for an existence check and closed.
    image.close();
    curl_mime* mime = curl_mime_init(curl);
    if (mime == nullptr) {
        curl_easy_cleanup(curl);
        throw std::runtime_error("curl_mime_init failed");
    }
    auto add_field = [mime](const char* name, const std::string& value) {
        curl_mimepart* part = curl_mime_addpart(mime);
        curl_mime_name(part, name);
        curl_mime_data(part, value.c_str(), CURL_ZERO_TERMINATED);
    };
    add_field("model", config_.model);
    add_field("prompt", prompt);
    add_field("size", config_.size);
    add_field("response_format", config_.response_format);
    curl_mimepart* file_part = curl_mime_addpart(mime);
    curl_mime_name(file_part, "image");
    if (curl_mime_filedata(file_part, image_path.c_str()) != CURLE_OK) {
        curl_mime_free(mime);
        curl_easy_cleanup(curl);
        throw std::runtime_error("cannot attach edit image to multipart request");
    }

    struct curl_slist* headers = nullptr;
    headers = curl_slist_append(headers, "Accept: application/json");
    headers = curl_slist_append(headers, ("Authorization: Bearer " + config_.api_key).c_str());
    HttpResponse response;
    char error_buffer[CURL_ERROR_SIZE] = {};
    curl_easy_setopt(curl, CURLOPT_ERRORBUFFER, error_buffer);
    curl_easy_setopt(curl, CURLOPT_URL, endpoint("/images/edits").c_str());
    curl_easy_setopt(curl, CURLOPT_MIMEPOST, mime);
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
    curl_easy_setopt(curl, CURLOPT_FOLLOWLOCATION, 1L);
    curl_easy_setopt(curl, CURLOPT_CONNECTTIMEOUT, std::min(config_.timeout_seconds, 15L));
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, config_.timeout_seconds);
    curl_easy_setopt(curl, CURLOPT_NOSIGNAL, 1L);
    curl_easy_setopt(curl, CURLOPT_USERAGENT, "image-generation-multilang-cpp/1.0");
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_body);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response.body);
    const CURLcode result = curl_easy_perform(curl);
    if (result != CURLE_OK) {
        const std::string details = error_buffer[0] != '\0' ? error_buffer : curl_easy_strerror(result);
        curl_slist_free_all(headers);
        curl_mime_free(mime);
        curl_easy_cleanup(curl);
        throw std::runtime_error(std::string("edit HTTP request failed: ") + details);
    }
    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &response.status);
    curl_slist_free_all(headers);
    curl_mime_free(mime);
    curl_easy_cleanup(curl);
    return resolve_image_response(response.status, response.body);
}

std::vector<unsigned char> ImageClient::resolve_image_response(long status,
                                                                const std::string& body) const {
    if (status < 200 || status >= 300) {
        throw std::runtime_error(response_error(status, body));
    }
    const auto b64 = json_string(body, "b64_json");
    if (b64.has_value() && !b64->empty()) {
        return detail::decode_base64(*b64);
    }
    const auto url = json_string(body, "url");
    if (url.has_value() && !url->empty()) {
        return download_image(*url);
    }
    throw std::runtime_error("image API data item contains neither b64_json nor url");
}

std::vector<unsigned char> ImageClient::download_image(const std::string& url) const {
    // Deliberately omit Authorization: the URL can point at provider object
    // storage and bearer tokens must not be sent to a different host.
    const HttpResponse response = perform_request(url, config_, nullptr, [](CURL* curl) {
        curl_easy_setopt(curl, CURLOPT_HTTPGET, 1L);
    });
    if (response.status < 200 || response.status >= 300) {
        throw std::runtime_error("image download returned HTTP " + std::to_string(response.status));
    }
    if (response.body.empty()) {
        throw std::runtime_error("image download was empty");
    }
    return std::vector<unsigned char>(response.body.begin(), response.body.end());
}

namespace detail {

std::string escape_json(const std::string& value) {
    std::string escaped;
    escaped.reserve(value.size() + 2);
    for (const unsigned char character : value) {
        switch (character) {
        case '"': escaped += "\\\""; break;
        case '\\': escaped += "\\\\"; break;
        case '\b': escaped += "\\b"; break;
        case '\f': escaped += "\\f"; break;
        case '\n': escaped += "\\n"; break;
        case '\r': escaped += "\\r"; break;
        case '\t': escaped += "\\t"; break;
        default:
            if (character < 0x20) {
                static constexpr char hex[] = "0123456789abcdef";
                escaped += "\\u00";
                escaped.push_back(hex[(character >> 4) & 0x0f]);
                escaped.push_back(hex[character & 0x0f]);
            } else {
                escaped.push_back(static_cast<char>(character));
            }
        }
    }
    return escaped;
}

std::vector<unsigned char> decode_base64(const std::string& value) {
    std::string encoded = value;
    if (encoded.rfind("data:", 0) == 0) {
        const auto comma = encoded.find(',');
        if (comma != std::string::npos) {
            encoded = encoded.substr(comma + 1);
        }
    }
    encoded.erase(std::remove_if(encoded.begin(), encoded.end(), [](unsigned char character) {
        return std::isspace(character) != 0;
    }), encoded.end());
    static constexpr char alphabet[] = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
    std::vector<unsigned char> output;
    int accumulator = 0;
    int bits = -8;
    for (const unsigned char character : encoded) {
        if (character == '=') {
            break;
        }
        const char* found = std::find(std::begin(alphabet), std::end(alphabet) - 1,
                                      static_cast<char>(character));
        if (found == std::end(alphabet) - 1) {
            throw std::invalid_argument("invalid base64 image data");
        }
        accumulator = (accumulator << 6) + static_cast<int>(found - alphabet);
        bits += 6;
        if (bits >= 0) {
            output.push_back(static_cast<unsigned char>((accumulator >> bits) & 0xff));
            bits -= 8;
        }
    }
    if (output.empty()) {
        throw std::invalid_argument("base64 image data decoded to empty bytes");
    }
    return output;
}

} // namespace detail
} // namespace image_demo

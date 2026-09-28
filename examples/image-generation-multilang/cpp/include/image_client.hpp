#pragma once

#include <cstddef>
#include <functional>
#include <string>
#include <vector>

namespace image_demo {

struct Config {
    std::string api_key;
    std::string base_url = "https://api.ai-router.dev/v1";
    std::string model = "gpt-image-2";
    std::string size = "1024x1024";
    std::string response_format = "b64_json";
    long timeout_seconds = 120;

    static Config from_environment();
};

class ImageClient {
public:
    explicit ImageClient(Config config);

    // Returns decoded image bytes from data[0].b64_json or data[0].url.
    std::vector<unsigned char> generate(const std::string& prompt) const;
    std::vector<unsigned char> edit(const std::string& image_path,
                                    const std::string& prompt) const;

private:
    Config config_;

    std::string endpoint(const std::string& path) const;
    std::vector<unsigned char> resolve_image_response(long status,
                                                       const std::string& body) const;
    std::vector<unsigned char> download_image(const std::string& url) const;
};

namespace detail {

std::string escape_json(const std::string& value);
std::vector<unsigned char> decode_base64(const std::string& value);

} // namespace detail
} // namespace image_demo

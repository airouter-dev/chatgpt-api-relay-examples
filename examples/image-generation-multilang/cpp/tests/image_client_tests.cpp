#include "image_client.hpp"

#include <cassert>
#include <string>
#include <vector>

int main() {
    using image_demo::detail::decode_base64;
    using image_demo::detail::escape_json;

    const std::vector<unsigned char> decoded = decode_base64("aGVsbG8=");
    assert(std::string(decoded.begin(), decoded.end()) == "hello");
    const std::vector<unsigned char> data_url = decode_base64("data:image/png;base64,aGVsbG8=");
    assert(std::string(data_url.begin(), data_url.end()) == "hello");
    assert(escape_json("quote\" newline\n") == "quote\\\" newline\\n");
    return 0;
}

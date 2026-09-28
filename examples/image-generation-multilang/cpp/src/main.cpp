#include "image_client.hpp"

#include <fstream>
#include <filesystem>
#include <cstdlib>
#include <iostream>
#include <stdexcept>
#include <string>
#include <vector>

namespace {

struct Options {
    image_demo::Config config = image_demo::Config::from_environment();
    std::string prompt;
    std::string edit_prompt;
    std::string output = "generated.png";
    std::string edited_output = "edited.png";
};

void usage(const char* executable) {
    std::cout << "Usage: " << executable << " --prompt TEXT --edit-prompt TEXT [options]\n"
              << "  --api-key KEY       API key (or OPENAI_API_KEY/AI_ROUTER_API_KEY)\n"
              << "  --base-url URL      API root (default from environment)\n"
              << "  --model MODEL       gpt-image-2 (or another supported model)\n"
              << "  --size SIZE         e.g. 1024x1024\n"
              << "  --response-format F b64_json or url\n"
              << "  --timeout SECONDS   HTTP timeout\n"
              << "  --output PATH       generated image path\n"
              << "  --edited-output PATH edited image path\n";
}

std::string required_value(int& index, int argc, char** argv, const std::string& option) {
    if (index + 1 >= argc) {
        throw std::invalid_argument(option + " requires a value");
    }
    return argv[++index];
}

Options parse_options(int argc, char** argv) {
    Options options;
    for (int index = 1; index < argc; ++index) {
        const std::string option = argv[index];
        if (option == "--help" || option == "-h") {
            usage(argv[0]);
            std::exit(0);
        }
        if (option == "--prompt") options.prompt = required_value(index, argc, argv, option);
        else if (option == "--edit-prompt") options.edit_prompt = required_value(index, argc, argv, option);
        else if (option == "--api-key") options.config.api_key = required_value(index, argc, argv, option);
        else if (option == "--base-url") options.config.base_url = required_value(index, argc, argv, option);
        else if (option == "--model") options.config.model = required_value(index, argc, argv, option);
        else if (option == "--size") options.config.size = required_value(index, argc, argv, option);
        else if (option == "--response-format") options.config.response_format = required_value(index, argc, argv, option);
        else if (option == "--timeout") options.config.timeout_seconds = std::stol(required_value(index, argc, argv, option));
        else if (option == "--output") options.output = required_value(index, argc, argv, option);
        else if (option == "--edited-output") options.edited_output = required_value(index, argc, argv, option);
        else throw std::invalid_argument("unknown option: " + option);
    }
    if (options.prompt.empty() || options.edit_prompt.empty()) {
        throw std::invalid_argument("--prompt and --edit-prompt are required");
    }
    if (options.config.api_key.empty()) {
        throw std::invalid_argument("API key is missing; pass --api-key or set OPENAI_API_KEY");
    }
    return options;
}

void write_image(const std::string& path, const std::vector<unsigned char>& bytes) {
    if (bytes.empty()) throw std::runtime_error("image response was empty");
    const std::filesystem::path output_path(path);
    if (output_path.has_parent_path()) {
        std::filesystem::create_directories(output_path.parent_path());
    }
    std::ofstream output(path, std::ios::binary);
    if (!output) throw std::runtime_error("cannot open output: " + path);
    output.write(reinterpret_cast<const char*>(bytes.data()), static_cast<std::streamsize>(bytes.size()));
    if (!output) throw std::runtime_error("cannot write output: " + path);
}

} // namespace

int main(int argc, char** argv) {
    try {
        const Options options = parse_options(argc, argv);
        const image_demo::ImageClient client(options.config);
        const auto generated = client.generate(options.prompt);
        write_image(options.output, generated);
        const auto edited = client.edit(options.output, options.edit_prompt);
        write_image(options.edited_output, edited);
        std::cout << "Generated " << options.output << "\nEdited    " << options.edited_output << "\n";
        return 0;
    } catch (const std::exception& error) {
        std::cerr << "error: " << error.what() << "\n";
        return 1;
    }
}

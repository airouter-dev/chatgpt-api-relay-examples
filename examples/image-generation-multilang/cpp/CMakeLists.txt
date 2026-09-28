cmake_minimum_required(VERSION 3.20)
project(image_generation_cpp_demo LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

find_package(CURL REQUIRED)

add_library(image_client src/image_client.cpp)
target_include_directories(image_client PUBLIC ${CMAKE_CURRENT_SOURCE_DIR}/include)
target_link_libraries(image_client PUBLIC CURL::libcurl)

if (MSVC)
  target_compile_options(image_client PRIVATE /W4)
else()
  target_compile_options(image_client PRIVATE -Wall -Wextra -Wpedantic)
endif()

add_executable(imagegen src/main.cpp)
target_link_libraries(imagegen PRIVATE image_client)

include(CTest)
if (BUILD_TESTING)
  add_executable(image_client_tests tests/image_client_tests.cpp)
  target_link_libraries(image_client_tests PRIVATE image_client)
  add_test(NAME image_client_tests COMMAND image_client_tests)
endif()

import os
import unittest
from unittest.mock import patch

from image_demo.config import ImageConfig, normalize_base_url


class ConfigTests(unittest.TestCase):
    def test_base_url_gets_v1_once(self):
        self.assertEqual(normalize_base_url("https://example.test/"), "https://example.test/v1")
        self.assertEqual(normalize_base_url("https://example.test/v1/"), "https://example.test/v1")

    def test_env_config_supports_openai_key_fallback(self):
        with patch.dict(os.environ, {"OPENAI_API_KEY": "test-key"}, clear=True):
            config = ImageConfig.from_env()
        self.assertEqual(config.api_key, "test-key")
        self.assertEqual(config.model, "gpt-image-2")


if __name__ == "__main__":
    unittest.main()

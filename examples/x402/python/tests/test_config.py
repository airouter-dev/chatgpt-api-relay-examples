import pytest

from ai_router_x402.client import select_requirements
from ai_router_x402.config import Config


def test_probe_defaults_are_safe():
    cfg = Config()
    cfg.validate()
    assert cfg.mode == "probe"
    assert cfg.model == "gpt-image-2"


def test_pay_requires_key():
    with pytest.raises(ValueError, match="PRIVATE_KEY"):
        Config(mode="pay").validate()


def test_payment_selector_accepts_only_configured_base_usdc():
    selected = select_requirements(2, [{
        "scheme": "exact",
        "network": "eip155:8453",
        "asset": "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
        "payTo": "0xEe027fE58800Bcef81469EDF0C2E7f2ab54965De",
        "amount": "4500",
    }])
    assert selected["amount"] == "4500"
    with pytest.raises(ValueError):
        select_requirements(2, [{**selected, "network": "eip155:1"}])

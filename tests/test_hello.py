"""Smoke test for hello.py."""

import os
import sys

sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

import hello  # noqa: E402


def test_main_prints_message(capsys):
    hello.main()
    assert capsys.readouterr().out == "test new world\n"

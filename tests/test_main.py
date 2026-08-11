"""Unit tests for :mod:`main`."""

import os
import sys

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))

from main import add, main


def test_add_positive():
    assert add(2, 3) == 5


def test_add_zero():
    assert add(0, 0) == 0


def test_add_negative():
    assert add(-2, -3) == -5


def test_add_mixed_signs():
    assert add(-2, 5) == 3


def test_main_output(capsys):
    main()
    captured = capsys.readouterr()
    assert "Hello from testupsun!" in captured.out
    assert "add(2, 3) = 5" in captured.out

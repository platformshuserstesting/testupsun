"""Tests for the palindrome checker."""

import os
import sys

import pytest

sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from palindrome import is_palindrome, main, normalize  # noqa: E402


@pytest.mark.parametrize(
    "text, expected",
    [
        ("racecar", True),
        ("a", True),
        ("", True),
        ("abba", True),
        ("hello", False),
        ("ab", False),
        ("Racecar", True),
        ("A man, a plan, a canal: Panama", True),
        ("No 'x' in Nixon", True),
        ("!!!", True),
        ("Was it a car or a cat I saw?", True),
        ("palindrome", False),
    ],
)
def test_is_palindrome_defaults(text, expected):
    assert is_palindrome(text) is expected


def test_numeric_input_is_coerced():
    assert is_palindrome(121) is True
    assert is_palindrome(123) is False


def test_none_raises_type_error():
    with pytest.raises(TypeError):
        is_palindrome(None)


def test_case_sensitive_comparison():
    assert is_palindrome("Racecar", ignore_case=False) is False
    assert is_palindrome("racecar", ignore_case=False) is True


def test_strict_comparison():
    assert is_palindrome("A man, a plan, a canal: Panama", alphanumeric_only=False) is False
    assert is_palindrome("racecar", alphanumeric_only=False) is True


def test_unicode_case_folding():
    assert is_palindrome("Été") is True
    assert is_palindrome("Été", ignore_case=False) is False


def test_normalize_strips_and_folds():
    assert normalize("A man, a plan!") == "amanaplan"
    assert normalize("A man, a plan!", ignore_case=False) == "Amanaplan"
    assert normalize("A man, a plan!", alphanumeric_only=False) == "a man, a plan!"


def test_cli_palindrome(capsys):
    assert main(["racecar"]) == 0
    assert capsys.readouterr().out == '"racecar" is a palindrome.\n'


def test_cli_not_palindrome(capsys):
    assert main(["hello"]) == 1
    assert capsys.readouterr().out == '"hello" is not a palindrome.\n'


def test_cli_sentence_default(capsys):
    assert main(["A man, a plan, a canal: Panama"]) == 0
    assert "is a palindrome." in capsys.readouterr().out


def test_cli_strict_flag(capsys):
    assert main(["--strict", "A man, a plan, a canal: Panama"]) == 1
    assert "is not a palindrome." in capsys.readouterr().out


def test_cli_case_sensitive_flag(capsys):
    assert main(["--case-sensitive", "Racecar"]) == 1
    assert "is not a palindrome." in capsys.readouterr().out


def test_cli_number(capsys):
    assert main(["121"]) == 0
    assert capsys.readouterr().out == '"121" is a palindrome.\n'


def test_cli_empty_string(capsys):
    assert main([""]) == 0
    assert capsys.readouterr().out == '"" is a palindrome.\n'


def test_cli_interactive_prompt(monkeypatch, capsys):
    monkeypatch.setattr("sys.stdin.isatty", lambda: True)
    monkeypatch.setattr("builtins.input", lambda prompt="": "level")
    assert main([]) == 0
    assert capsys.readouterr().out == '"level" is a palindrome.\n'


def test_cli_piped_stdin(monkeypatch, capsys):
    import io

    monkeypatch.setattr("sys.stdin", io.StringIO("racecar\n"))
    assert main([]) == 0
    assert capsys.readouterr().out == '"racecar" is a palindrome.\n'


def test_cli_usage_error():
    with pytest.raises(SystemExit) as excinfo:
        main(["--unknown-flag", "racecar"])
    assert excinfo.value.code == 2

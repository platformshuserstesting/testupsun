"""Tests for the hello_world module."""

from hello_world import main


def test_main_prints_greeting(capsys):
    main()
    assert capsys.readouterr().out == "Hello, World!\n"


def test_import_prints_nothing(capsys):
    import hello_world  # noqa: F401

    assert capsys.readouterr().out == ""

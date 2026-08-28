#!/usr/bin/env python3
"""Tests for hello.py."""

import contextlib
import io
import unittest

import hello


class MainTestCase(unittest.TestCase):
    def test_main_prints_greeting(self) -> None:
        buffer = io.StringIO()
        with contextlib.redirect_stdout(buffer):
            hello.main()
        self.assertEqual(buffer.getvalue(), "Hello, World!\n")


if __name__ == "__main__":
    unittest.main()

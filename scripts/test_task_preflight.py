"""Exercise eligibility and read-only CLI behavior using temporary fixtures."""

import copy
import json
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest


ROOT = Path(__file__).resolve().parents[1]
SCRIPT = ROOT / "scripts" / "task_preflight.py"
TRACKER = ROOT / "feature_list.json"


class TaskPreflightTests(unittest.TestCase):
    def setUp(self):
        self.original = TRACKER.read_bytes()
        self.tracker = json.loads(self.original)
        self.first, self.second = self.tracker["tasks"][:2]

    def tearDown(self):
        self.assertEqual(TRACKER.read_bytes(), self.original)

    def run_cli(self, task="M3.1", path=None, cwd=None, extra=()):
        command = [sys.executable, str(SCRIPT), "--task", task]
        if path is not None:
            command.extend(["--tracker", str(path)])
        return subprocess.run(command + list(extra), cwd=cwd, capture_output=True,
                              text=True, check=False)

    def fixture_cli(self, task="M3.1"):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "tracker.json"
            original = (json.dumps(self.tracker, indent=3) + "\n").encode("utf-8")
            path.write_bytes(original)
            result = self.run_cli(task, path)
            self.assertEqual(path.read_bytes(), original)
            self.assertEqual(list(Path(directory).iterdir()), [path])
            return result

    def assert_rejected(self, result, *messages):
        self.assertEqual(result.returncode, 1)
        self.assertIn("INELIGIBLE", result.stderr)
        self.assertEqual(result.stdout, "")
        for message in messages:
            self.assertIn(message, result.stderr)

    def verify_first(self):
        self.first.update(status="verified", evidence=dict(
            revision="fixture-revision", commands=["fixture check"],
            result="PASS", reviewed_by="fixture reviewer"))

    def test_current_first_is_eligible_but_requires_approval(self):
        result = self.run_cli()
        self.assertEqual(result.returncode, 0)
        self.assertIn("ELIGIBLE — M3.1 meets task-state prerequisites", result.stdout)
        self.assertIn("HUMAN APPROVAL REQUIRED", result.stdout)
        self.assertIn("does not authorize implementation", result.stdout)
        self.assertNotIn("AUTHORIZED", result.stdout.upper())
        self.assertEqual(result.stderr, "")

    def test_current_second_has_unverified_dependency(self):
        self.assert_rejected(self.run_cli("M3.2"), "M3.2", "depends on unverified M3.1")

    def test_unknown_task(self):
        self.assert_rejected(self.run_cli("M3.999"), "M3.999", "unknown task ID")

    def test_invalid_tracker_contract(self):
        cases = [None, [], {}, copy.deepcopy(self.tracker)]
        cases[-1]["active_task_id"] = "M3.1"
        for data in cases:
            with self.subTest(data=data):
                self.tracker = data
                self.assert_rejected(self.fixture_cli(), "M3.1", "invalid tracker")

    def test_malformed_json_and_invalid_encoding(self):
        for raw in (b'{"tasks": [}', b'\xff'):
            with self.subTest(raw=raw), tempfile.TemporaryDirectory() as directory:
                path = Path(directory) / "invalid.json"
                path.write_bytes(raw)
                self.assert_rejected(self.run_cli(path=path), "M3.1", "cannot load tracker")
                self.assertEqual(path.read_bytes(), raw)

    def test_missing_tracker(self):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "missing.json"
            self.assert_rejected(self.run_cli(path=path), "M3.1", str(path))
            self.assertFalse(path.exists())

    def test_selected_in_progress(self):
        self.first["status"] = "in_progress"
        self.tracker["active_task_id"] = "M3.1"
        self.assert_rejected(self.fixture_cli(), "M3.1", "status is in_progress")

    def test_selected_verified(self):
        self.verify_first()
        self.assert_rejected(self.fixture_cli(), "M3.1", "status is verified")

    def test_selected_blocked(self):
        self.first.update(status="blocked", blocker="Fixture blocker")
        self.assert_rejected(self.fixture_cli(), "M3.1", "status is blocked")

    def test_other_in_progress_prevents_selection(self):
        self.first["status"] = "in_progress"
        self.tracker["active_task_id"] = "M3.1"
        self.assert_rejected(self.fixture_cli("M3.2"), "M3.2", "in_progress: M3.1")

    def test_verified_dependencies_allow_pending_selection(self):
        self.verify_first()
        result = self.fixture_cli("M3.2")
        self.assertEqual(result.returncode, 0)
        self.assertIn("ELIGIBLE — M3.2", result.stdout)
        self.assertIn("HUMAN APPROVAL REQUIRED", result.stdout)

    def test_multiple_in_progress_invalidates_tracker(self):
        self.first["status"] = self.second["status"] = "in_progress"
        self.tracker["active_task_id"] = "M3.1"
        self.assert_rejected(self.fixture_cli(), "invalid tracker", "multiple in_progress tasks")

    def test_cli_default_from_nested_directory(self):
        result = self.run_cli(cwd=ROOT / "docs")
        self.assertEqual(result.returncode, 0)
        self.assertIn("ELIGIBLE — M3.1", result.stdout)

    def test_cli_default_from_unrelated_directory(self):
        with tempfile.TemporaryDirectory() as directory:
            self.assertEqual(self.run_cli(cwd=directory).returncode, 0)

    def test_input_unchanged_after_success_and_failure(self):
        self.assertEqual(self.fixture_cli().returncode, 0)
        self.assertEqual(self.fixture_cli("M3.2").returncode, 1)
        self.tracker["schema_version"] = 2
        self.assertEqual(self.fixture_cli().returncode, 1)

    def test_success_does_not_activate(self):
        self.assertEqual(self.fixture_cli().returncode, 0)
        self.assertTrue(all(task["status"] == "pending" for task in self.tracker["tasks"]))
        self.assertIsNone(self.tracker["active_task_id"])
        self.assertIn("No task activated", self.run_cli().stdout)

    def test_no_approval_cli_flag(self):
        result = self.run_cli(extra=("--approve",))
        self.assertEqual(result.returncode, 2)
        self.assertIn("unrecognized arguments", result.stderr)
        self.assertEqual(result.stdout, "")


if __name__ == "__main__":
    unittest.main()

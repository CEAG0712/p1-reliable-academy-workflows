"""Feature tracker contract tests; the published tracker is only ever read."""

import copy
import json
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest

from validate_features import validate_tracker


ROOT = Path(__file__).resolve().parents[1]
SCRIPT = ROOT / "scripts" / "validate_features.py"
TRACKER = ROOT / "feature_list.json"


class FeatureValidationTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.published = json.loads(TRACKER.read_text(encoding="utf-8"))

    def setUp(self):
        self.tracker = copy.deepcopy(self.published)
        # Contract cases describe explicit states, independent of project progress.
        self.tracker["active_task_id"] = None
        for task in self.tracker["tasks"]:
            task.update(status="pending", blocker=None, evidence=dict(
                revision=None, commands=[], result=None, reviewed_by=None))
        self.first, self.second = self.tracker["tasks"][:2]

    def reject(self, text):
        errors = validate_tracker(self.tracker)
        self.assertTrue(errors)
        self.assertIn(text, "\n".join(errors))

    def verify(self, task):
        task["status"] = "verified"
        task["evidence"] = dict(revision="abc123", commands=["go test ./..."],
                                result="PASS", reviewed_by="human reviewer")

    def run_cli(self, path):
        return subprocess.run([sys.executable, str(SCRIPT), str(path)],
                              capture_output=True, text=True, check=False)

    def test_published_tracker_is_valid(self):
        self.assertEqual(validate_tracker(self.published), [])

    def test_single_in_progress_is_valid(self):
        self.first["status"] = "in_progress"
        self.tracker["active_task_id"] = self.first["id"]
        self.assertEqual(validate_tracker(self.tracker), [])

    def test_two_in_progress_rejected(self):
        self.first["status"] = self.second["status"] = "in_progress"
        self.tracker["active_task_id"] = self.first["id"]
        self.reject(f"multiple in_progress tasks: {self.first['id']}, {self.second['id']}")

    def test_mismatched_active_task_rejected(self):
        self.first["status"] = "in_progress"
        self.tracker["active_task_id"] = self.second["id"]
        self.reject("active_task_id")

    def test_active_task_without_in_progress_rejected(self):
        self.tracker["active_task_id"] = self.first["id"]
        self.reject("must be None")

    def test_missing_dependency_rejected(self):
        self.first["depends_on"] = ["missing-task"]
        self.reject(f"task {self.first['id']}: missing dependency missing-task")

    def test_cycle_rejected(self):
        self.first["depends_on"] = [self.second["id"]]
        self.reject(f"dependency cycle: {self.first['id']} -> {self.second['id']} -> {self.first['id']}")

    def test_duplicate_dependency_rejected(self):
        self.second["depends_on"] *= 2
        self.reject(f"task {self.second['id']}: duplicate dependency {self.first['id']}")

    def test_self_dependency_rejected(self):
        self.first["depends_on"] = [self.first["id"]]
        self.reject(f"task {self.first['id']}: self-dependency")

    def test_start_before_prerequisite_verified_rejected(self):
        self.second["status"] = "in_progress"
        self.tracker["active_task_id"] = self.second["id"]
        self.reject(f"task {self.second['id']}: dependency {self.first['id']} must be verified")

    def test_verified_requires_verified_prerequisite(self):
        self.verify(self.second)
        self.reject(f"task {self.second['id']}: dependency {self.first['id']} must be verified")

    def test_verified_without_evidence_rejected(self):
        self.first["status"] = "verified"
        self.reject(f"task {self.first['id']!r}")

    def test_verified_with_complete_evidence_accepted(self):
        self.verify(self.first)
        self.verify(self.second)
        self.assertEqual(validate_tracker(self.tracker), [])

    def test_each_verified_evidence_field_rejects_empty_and_wrong_types(self):
        for field in ("revision", "commands", "result", "reviewed_by"):
            invalid = (None, "", " \t", 1, {}, [])
            if field == "commands":
                invalid += ([""], [" \t"], [1], ["ok", None])
            for value in invalid:
                with self.subTest(field=field, value=value):
                    self.verify(self.first)
                    self.first["evidence"][field] = value
                    self.reject(f"evidence.{field}")

    def test_blocked_without_reason_rejected(self):
        self.first["status"] = "blocked"
        for value in (None, "", " \t", 1, []):
            with self.subTest(value=value):
                self.first["blocker"] = value
                self.reject("nonempty blocker reason")

    def test_blocked_with_reason_accepted(self):
        self.first.update(status="blocked", blocker="Awaiting prerequisite decision")
        self.assertEqual(validate_tracker(self.tracker), [])

    def test_other_statuses_require_null_blocker(self):
        for status in ("pending", "in_progress", "verified"):
            with self.subTest(status=status):
                self.first.update(status=status, blocker="old reason")
                self.reject("blocker must be null")

    def test_required_metadata_fields(self):
        for field in self.published:
            with self.subTest(field=field):
                self.tracker = copy.deepcopy(self.published)
                del self.tracker[field]
                self.reject(field)

    def test_required_task_and_evidence_fields(self):
        for field in self.published["tasks"][0]:
            with self.subTest(field=field):
                self.tracker = copy.deepcopy(self.published)
                del self.tracker["tasks"][0][field]
                self.reject(field)
        for field in self.published["tasks"][0]["evidence"]:
            with self.subTest(evidence_field=field):
                self.tracker = copy.deepcopy(self.published)
                del self.tracker["tasks"][0]["evidence"][field]
                self.reject(f"evidence.{field}")

    def test_invalid_structure_and_metadata(self):
        for value in (None, [], "tracker", 1):
            with self.subTest(root=value):
                self.assertTrue(validate_tracker(value))
        cases = {"schema_version": (True, 1.0, 2, "1"),
                 "work_in_progress_limit": (True, 1.0, 2, "1"),
                 "project": (None, " ", 1), "milestone": (None, " ", []),
                 "source_of_scope": (None, "", {}),
                 "active_task_id": (" ", [], 1),
                 "status_meanings": (None, [], {}, {"pending": " "}),
                 "rules": (None, "rule", [], [" "]),
                 "tasks": (None, {}, "tasks")}
        for field, values in cases.items():
            for value in values:
                with self.subTest(field=field, value=value):
                    self.tracker = copy.deepcopy(self.published)
                    self.tracker[field] = value
                    self.reject(field)

    def test_invalid_task_field_types(self):
        cases = {"id": (None, " ", [], 1), "title": (None, "", 1),
                 "objective": (None, " ", {}), "status": (None, "done", [], {}),
                 "depends_on": (None, "task", [1], [" "]),
                 "acceptance_criteria": (None, [], [""], "criterion"),
                 "evidence": (None, [], "evidence")}
        for field, values in cases.items():
            for value in values:
                with self.subTest(field=field, value=value):
                    self.tracker = copy.deepcopy(self.published)
                    self.tracker["tasks"][0][field] = value
                    self.reject(field)
        self.tracker = copy.deepcopy(self.published)
        self.tracker["tasks"][0] = None
        self.reject("tasks[0]")

    def test_duplicate_task_ids_rejected(self):
        self.second["id"] = self.first["id"]
        self.reject(f"duplicate task ID: {self.first['id']}")

    def test_future_task_ids_are_supported(self):
        task = copy.deepcopy(self.first)
        task["id"] = "FUTURE-42"
        self.tracker.update(tasks=[task], milestone="Future milestone")
        self.assertEqual(validate_tracker(self.tracker), [])

    def test_malformed_json_cli_rejected(self):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "malformed.json"
            path.write_text('{"tasks": [}', encoding="utf-8")
            result = self.run_cli(path)
            self.assertEqual(result.returncode, 1)
            self.assertIn("FAIL", result.stderr)
            self.assertIn(str(path), result.stderr)
            self.assertIn("line 1", result.stderr)

    def test_cli_does_not_modify_input_on_success_or_failure(self):
        for valid in (True, False):
            with self.subTest(valid=valid), tempfile.TemporaryDirectory() as directory:
                data = copy.deepcopy(self.published)
                if not valid:
                    data["tasks"][0]["depends_on"] = ["missing-task"]
                path = Path(directory) / "tracker.json"
                original = (json.dumps(data, indent=3) + "\n").encode("utf-8")
                path.write_bytes(original)
                result = self.run_cli(path)
                self.assertEqual(result.returncode, 0 if valid else 1)
                self.assertEqual(path.read_bytes(), original)
                if valid:
                    self.assertEqual(result.stdout, "Feature tracker validation: PASS\n")
                else:
                    self.assertIn(data["tasks"][0]["id"], result.stderr)

    def test_validation_does_not_mutate_object(self):
        before = copy.deepcopy(self.tracker)
        validate_tracker(self.tracker)
        self.assertEqual(self.tracker, before)

    def test_missing_file_cli_rejected(self):
        with tempfile.TemporaryDirectory() as directory:
            result = self.run_cli(Path(directory) / "missing.json")
            self.assertEqual(result.returncode, 1)
            self.assertIn("FAIL", result.stderr)


class LiveTrackerPreservationTests(unittest.TestCase):
    """Validate current state and immutable M3 history without freezing M4 progress."""

    @classmethod
    def setUpClass(cls):
        cls.original = json.loads(subprocess.check_output(
            ["git", "show", "5457b77e7835cf951b1c71ab332080d7b81fe735:feature_list.json"],
            cwd=ROOT, text=True))
        cls.raw = TRACKER.read_bytes()
        cls.live = json.loads(cls.raw)

    def test_live_tracker_is_valid_without_mutation(self):
        before = copy.deepcopy(self.live)
        self.assertEqual(validate_tracker(self.live), [])
        self.assertEqual(self.live, before)
        self.assertEqual(TRACKER.read_bytes(), self.raw)

    def test_original_m3_objects_evidence_and_order_are_preserved(self):
        self.assertEqual(len(self.original["tasks"]), 10)
        self.assertEqual(self.live["tasks"][:10], self.original["tasks"])
        self.assertTrue(all(t["status"] == "verified" for t in self.live["tasks"][:10]))


class UnifiedTrackerFixtureTests(unittest.TestCase):
    """Control M4 lifecycle state in memory independently of live task progress."""

    def setUp(self):
        self.tracker = json.loads(TRACKER.read_text(encoding="utf-8"))
        self.tracker["active_task_id"] = None
        for task in self.tracker["tasks"][10:]:
            task.update(status="pending", blocker=None, evidence=dict(
                revision=None, commands=[], result=None, reviewed_by=None))

    def verify_m4_first(self):
        self.tracker["tasks"][10].update(status="verified", evidence=dict(
            revision="fixture-revision", commands=["fixture check"],
            result="PASS", reviewed_by="fixture reviewer"))

    def test_unified_twenty_task_tracker_is_valid(self):
        self.assertEqual(validate_tracker(self.tracker), [])
        self.assertEqual(len(self.tracker["tasks"]), 20)
        self.assertEqual([t["id"] for t in self.tracker["tasks"]],
                         [f"M{stage}.{n}" for stage in (3, 4) for n in range(1, 11)])
        self.assertEqual(self.tracker["schema_version"], 1)
        self.assertEqual(self.tracker["work_in_progress_limit"], 1)
        self.assertIsNone(self.tracker["active_task_id"])

    def test_proposed_m4_chain_pending_evidence_and_blockers(self):
        empty = dict(revision=None, commands=[], result=None, reviewed_by=None)
        for n, task in enumerate(self.tracker["tasks"][10:], start=1):
            with self.subTest(task=task["id"]):
                self.assertEqual(task["status"], "pending")
                self.assertEqual(task["evidence"], empty)
                self.assertIsNone(task["blocker"])
                self.assertEqual(task["depends_on"], ["M3.10" if n == 1 else f"M4.{n - 1}"])

    def test_cross_stage_in_progress_tasks_are_rejected(self):
        self.tracker["tasks"][9]["status"] = "in_progress"
        self.tracker["tasks"][10]["status"] = "in_progress"
        self.tracker["active_task_id"] = "M3.10"
        self.assertIn("multiple in_progress tasks: M3.10, M4.1",
                      "\n".join(validate_tracker(self.tracker)))

    def test_active_id_must_match_sole_m4_in_progress_task(self):
        self.tracker["tasks"][10]["status"] = "in_progress"
        for active in (None, "M3.10", "M4.2"):
            with self.subTest(active=active):
                self.tracker["active_task_id"] = active
                self.assertIn("must be 'M4.1'", "\n".join(validate_tracker(self.tracker)))
        self.tracker["active_task_id"] = "M4.1"
        self.assertEqual(validate_tracker(self.tracker), [])

    def test_schema_v1_allows_populated_pending_evidence(self):
        # Empty M4 evidence is a preparation invariant, not a universal schema rule.
        self.tracker["tasks"][10]["evidence"] = dict(
            revision="fixture", commands=["fixture check"], result="draft",
            reviewed_by="fixture reviewer")
        self.assertEqual(validate_tracker(self.tracker), [])

    def test_pending_m4_nonnull_blocker_is_rejected_by_schema(self):
        self.tracker["tasks"][10]["blocker"] = "fixture blocker"
        self.assertIn("blocker must be null", "\n".join(validate_tracker(self.tracker)))

    def test_verified_m4_first_with_pending_second_is_valid(self):
        self.verify_m4_first()
        self.assertEqual(self.tracker["tasks"][11]["status"], "pending")
        self.assertEqual(validate_tracker(self.tracker), [])

    def test_m4_second_in_progress_requires_matching_active_id(self):
        self.verify_m4_first()
        self.tracker["tasks"][11]["status"] = "in_progress"
        self.tracker["active_task_id"] = "M4.2"
        self.assertEqual(validate_tracker(self.tracker), [])
        for active in (None, "M3.10", "M4.1"):
            with self.subTest(active=active):
                self.tracker["active_task_id"] = active
                self.assertIn("must be 'M4.2'", "\n".join(validate_tracker(self.tracker)))


if __name__ == "__main__":
    unittest.main()

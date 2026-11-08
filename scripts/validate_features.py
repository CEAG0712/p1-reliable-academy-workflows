"""Read-only validation of version 1 feature trackers (standard library only)."""

import argparse
import json
from pathlib import Path
import sys


STATUSES = ("pending", "in_progress", "blocked", "verified")


def nonempty_string(value):
    return isinstance(value, str) and bool(value.strip())


def string_list(value, allow_empty=True):
    return (isinstance(value, list) and (allow_empty or bool(value))
            and all(nonempty_string(item) for item in value))


def validate_metadata(data):
    errors = []
    for field in ("schema_version", "work_in_progress_limit"):
        if type(data.get(field)) is not int or data[field] != 1:
            errors.append(f"{field} must be the integer 1")
    for field in ("project", "milestone", "source_of_scope"):
        if not nonempty_string(data.get(field)):
            errors.append(f"{field} must be a nonempty string")
    if "active_task_id" not in data or not (
        data["active_task_id"] is None or nonempty_string(data["active_task_id"])
    ):
        errors.append("active_task_id must be null or a nonempty task ID")
    meanings = data.get("status_meanings")
    if not isinstance(meanings, dict) or any(
        not nonempty_string(meanings.get(status)) for status in STATUSES
    ):
        errors.append("status_meanings must describe pending, in_progress, blocked, verified")
    if not string_list(data.get("rules"), allow_empty=False):
        errors.append("rules must be a nonempty array of nonempty strings")
    return errors


def validate_evidence(task, label):
    evidence = task.get("evidence")
    if not isinstance(evidence, dict):
        return [f"{label}: evidence must be an object"]
    errors = []
    verified = task.get("status") == "verified"
    for field in ("revision", "result", "reviewed_by"):
        value = evidence.get(field)
        if field not in evidence or not (
            nonempty_string(value) or (value is None and not verified)
        ):
            errors.append(f"{label}: evidence.{field} must be a nonempty string"
                          + ("" if verified else " or null"))
    if not string_list(evidence.get("commands"), allow_empty=not verified):
        errors.append(f"{label}: evidence.commands must be an array of nonempty strings"
                      + (" with at least one command" if verified else ""))
    return errors


def validate_task(task, index):
    if not isinstance(task, dict):
        return [f"tasks[{index}] must be an object"]
    label = f"task {task.get('id')!r} (tasks[{index}])"
    errors = []
    for field in ("id", "title", "objective"):
        if not nonempty_string(task.get(field)):
            errors.append(f"{label}: {field} must be a nonempty string")
    if task.get("status") not in STATUSES:
        errors.append(f"{label}: status must be one of {', '.join(STATUSES)}")
    if not string_list(task.get("depends_on")):
        errors.append(f"{label}: depends_on must be an array of nonempty task IDs")
    if not string_list(task.get("acceptance_criteria"), allow_empty=False):
        errors.append(f"{label}: acceptance_criteria must be a nonempty array of nonempty strings")
    errors.extend(validate_evidence(task, label))
    if task.get("status") == "blocked":
        if not nonempty_string(task.get("blocker")):
            errors.append(f"{label}: blocked task requires a nonempty blocker reason")
    elif "blocker" not in task or task["blocker"] is not None:
        errors.append(f"{label}: blocker must be null unless status is blocked")
    return errors


def validate_active_task(data, tasks):
    active = [task_id for task_id, task in tasks.items() if task["status"] == "in_progress"]
    errors = []
    if len(active) > 1:
        errors.append(f"multiple in_progress tasks: {', '.join(active)}")
    expected = active[0] if len(active) == 1 else None
    if data["active_task_id"] != expected:
        errors.append(f"active_task_id {data['active_task_id']!r} must be {expected!r}; "
                      f"in_progress tasks: {active}")
    return errors


def dependency_cycles(tasks):
    """Iterative DFS avoids a recursion limit for long future task lists."""
    errors, finished = [], set()
    for root in tasks:
        if root in finished:
            continue
        path = [root]
        visiting = {root}
        stack = [iter(tasks[root]["depends_on"])]
        while stack:
            dependency = next(stack[-1], None)
            if dependency is None:
                finished.add(path.pop())
                visiting.difference_update(finished)
                stack.pop()
            elif dependency in visiting:
                cycle = path[path.index(dependency):] + [dependency]
                errors.append("dependency cycle: " + " -> ".join(cycle))
            elif dependency in tasks and dependency not in finished:
                path.append(dependency)
                visiting.add(dependency)
                stack.append(iter(tasks[dependency]["depends_on"]))
    return errors


def validate_dependencies(tasks):
    errors = []
    for task_id, task in tasks.items():
        seen = set()
        for dependency in task["depends_on"]:
            if dependency in seen:
                errors.append(f"task {task_id}: duplicate dependency {dependency}")
            seen.add(dependency)
            if dependency == task_id:
                errors.append(f"task {task_id}: self-dependency")
            if dependency not in tasks:
                errors.append(f"task {task_id}: missing dependency {dependency}")
            elif task["status"] in ("in_progress", "verified") and tasks[dependency]["status"] != "verified":
                errors.append(f"task {task_id}: dependency {dependency} must be verified "
                              f"before task is {task['status']}")
    return errors + dependency_cycles(tasks)


def validate_tracker(data):
    """Return diagnostics without modifying the supplied object."""
    if not isinstance(data, dict):
        return ["tracker must be a JSON object"]
    errors = validate_metadata(data)
    if not isinstance(data.get("tasks"), list):
        return errors + ["tasks must be an array"]
    tasks = {}
    for index, task in enumerate(data["tasks"]):
        errors.extend(validate_task(task, index))
        if isinstance(task, dict) and nonempty_string(task.get("id")):
            if task["id"] in tasks:
                errors.append(f"duplicate task ID: {task['id']}")
            tasks[task["id"]] = task
    # Graph and state checks require well-typed, uniquely identified tasks.
    if errors:
        return errors
    return validate_active_task(data, tasks) + validate_dependencies(tasks)


def main(argv=None):
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("tracker", type=Path, help="path to feature_list.json")
    args = parser.parse_args(argv)
    try:
        data = json.loads(args.tracker.read_text(encoding="utf-8"))
    except (OSError, UnicodeError, ValueError) as error:
        print(f"Feature tracker validation: FAIL\n{args.tracker}: {error}", file=sys.stderr)
        return 1
    errors = validate_tracker(data)
    if errors:
        print("Feature tracker validation: FAIL", file=sys.stderr)
        for error in errors:
            print(f"- {error}", file=sys.stderr)
        return 1
    print("Feature tracker validation: PASS")
    return 0


if __name__ == "__main__":
    sys.exit(main())

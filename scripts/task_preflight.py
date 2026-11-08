"""Read-only task eligibility check; success grants no implementation authority."""

import argparse
import json
from pathlib import Path
import sys

# Keep direct CLI invocation read-only, including the validator import.
sys.dont_write_bytecode = True
from validate_features import validate_tracker


DEFAULT_TRACKER = Path(__file__).resolve().parents[1] / "feature_list.json"


def main(argv=None):
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--task", required=True, help="task ID to check")
    parser.add_argument("--tracker", type=Path, default=DEFAULT_TRACKER,
                        help="tracker path (default: repository feature_list.json)")
    args = parser.parse_args(argv)

    def reject(reason):
        print(f"INELIGIBLE — {args.task}: {reason}", file=sys.stderr)
        return 1

    try:
        data = json.loads(args.tracker.read_text(encoding="utf-8"))
    except (OSError, UnicodeError, ValueError) as error:
        return reject(f"cannot load tracker {args.tracker}: {error}")
    errors = validate_tracker(data)
    if errors:
        return reject("invalid tracker:\n" + "\n".join(f"- {error}" for error in errors))

    tasks = {task["id"]: task for task in data["tasks"]}
    if args.task not in tasks:
        return reject("unknown task ID")
    selected = tasks[args.task]
    if selected["status"] != "pending":
        return reject(f"status is {selected['status']}; must be pending")
    active = [task_id for task_id, task in tasks.items()
              if task["status"] == "in_progress"]
    if active:
        return reject(f"another task is in_progress: {', '.join(active)}; WIP limit is 1")
    unverified = [dependency for dependency in selected["depends_on"]
                  if tasks[dependency]["status"] != "verified"]
    if unverified:
        return reject(f"depends on unverified {', '.join(unverified)}")

    print(f"ELIGIBLE — {args.task} meets task-state prerequisites. "
          "HUMAN APPROVAL REQUIRED. No task activated. "
          "This does not authorize implementation.")
    return 0


if __name__ == "__main__":
    sys.exit(main())

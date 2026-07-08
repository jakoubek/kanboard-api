# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed
- `TaskParams.WithDueDate`/`WithStartDate` and `TaskUpdateParams.SetDueDate`/
  `SetStartDate`/`ClearDueDate`/`ClearStartDate` now send `date_due`/
  `date_started` as a formatted string (`2006-01-02 15:04`) instead of a raw
  Unix timestamp. Kanboard's `createTask`/`updateTask` JSON-RPC methods reject
  a raw timestamp (silently, `updateTask` returns `success: false`) but accept
  this string format. **Breaking:** `CreateTaskRequest.DateDue`/`DateStarted`
  are now `string` (was `int64`); `UpdateTaskRequest.DateDue`/`DateStarted` are
  now `*string` (was `*int64`).
- `GetProjectByName` and `GetProjectByID` now return `ErrProjectNotFound` when
  Kanboard reports "not found" as the literal `false` (not just `null`), instead
  of failing with a JSON unmarshal error.
- `GetAllTaskLinks` now returns the linked task's real ID. Kanboard's actual API
  response carries it in the `task_id` field, not `opposite_task_id` (which
  doesn't exist in the real response) — `TaskLink.OppositeTaskID` was therefore
  always `0`. **Breaking:** `TaskLink.TaskID` and `TaskLink.LinkID` are removed
  (neither field exists in the real response, so they were always `0`);
  `OppositeTaskID` now reads `task_id`.

## [v1.6.0] - 2026-07-04

### Added
- `Task.TimeEstimated` and `Task.TimeSpent` fields (hours), exposing Kanboard's
  `time_estimated`/`time_spent` on all task read methods.
- `StringFloat` type that unmarshals a `float64` from a JSON string or number
  (empty string / `null` decode to `0`).

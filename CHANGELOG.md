# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v1.6.0] - 2026-07-04

### Added
- `Task.TimeEstimated` and `Task.TimeSpent` fields (hours), exposing Kanboard's
  `time_estimated`/`time_spent` on all task read methods.
- `StringFloat` type that unmarshals a `float64` from a JSON string or number
  (empty string / `null` decode to `0`).

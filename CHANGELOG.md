# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]


## [1.0.2] - 2019-07-28
### Added
- Testing for: ErrorHandler, TestHandler, Route, Router, Server, TestServer

### Fixed
- TestHandler's WithResponseHeaders and AddResponseHeaders to actually manipulate the response instead of the request

## [1.0.1] - 2019-07-26
### Added
- CHANGELOG.md
- /server/error_test.go

### Fixed
- BasicErrorHandler's HandleError method, to panic only in case of write error

## [1.0.0] - 2019-07-07
### Added
- Initial commit
- Beta version working
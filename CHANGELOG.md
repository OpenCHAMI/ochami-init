# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.0.1] - 2024-01-10

### Added

- First release reads config file and creates databases
- Retries the postgres connection every second for two minutes to give the posgres instance a chance to come up
- Relies on environment variables to find the postgres instance and connect to it

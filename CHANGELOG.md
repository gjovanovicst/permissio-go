# Changelog

All notable changes to the Permissio.io Go SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial SDK implementation
- Permission checking with `Check()` method
- Auto-scope detection from API key
- Full CRUD operations for Users, Roles, Tenants, and Resources
- Gin middleware integration
- Comprehensive examples

## [0.1.0] - 2024-XX-XX

### Added
- Initial release
- `permit.NewPermit()` client initialization
- `config.NewConfigBuilder()` for configuration
- `enforcement.UserBuilder()` and `enforcement.ResourceBuilder()` for building check requests
- API clients for Users, Roles, Tenants, Resources, and Role Assignments
- Gin middleware for permission enforcement
- Examples for common use cases

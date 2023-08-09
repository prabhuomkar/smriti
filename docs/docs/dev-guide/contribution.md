# Contributing Guide

## Development
Here are some best practices we expect developers to follow:
- Write code that is easy to read and modify.
- Use consistent coding standards for naming and formatting.
- Break up large change requests into smaller, manageable ones.
- Use automated testing to catch bugs early.
- Be responsive to community feedback and engage with users.
- Continuously improve the project with new features and bug fixes.

## Pull Request Checklist
- Followed the style guidelines of this project
- Performed self-review of my code
- Commented the code, particularly in hard-to-understand areas
- Made corresponding changes to the documentation
- New changes generate no new warnings
- Added tests that prove your fix is effective or that the feature works
- New and existing unit tests pass locally with your changes
- Any dependent changes have been merged and published in downstream modules

## Testing

### Unit Testing

#### API
- Run unit tests
```
make test
```
- See code coverage
```
make cover
```
- Adding new unit tests  
  - New files/functions should have test cases in their respective `*_test.go` file in the same package.
  - Tests should have both positive and negative cases with mocking if required.
  - Check `api/internal/handlers/albums_test.go` as reference for adding unit tests.

#### Worker
- Run unit tests
```
make test
```
- See code coverage
```
make cover
```
- Adding new unit tests  
  - New files/functions should have test cases in their respective `test_*.py` file in the `worker/tests` directory.
  - Tests should have both positive and negative cases with mocking if required.
  - Check `worker/tests/components/test_place.py` as reference for adding unit tests.

### Integration Testing

#### Setup
- From the root directory of the project, run all services:
```
make start
```
- From `tests` folder, run following command to setup the test suite:
```
make setup
```
- From `tests` folder, run the command for invoking integration tests for all features:
```
make test-all
```
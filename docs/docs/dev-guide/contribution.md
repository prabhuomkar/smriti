# Contributing Guide

## Development
Here are some best practices we expect developers to follow:
- Write code that is easy to read and modify.
- Use consistent coding standards for naming and formatting.
- Break up large change requests into smaller, manageable ones.
- Use automated testing to catch bugs early.
- Be responsive to community feedback and engage with users.
- Continuously improve the project with new features and bug fixes.

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
  - New files/functions should have a respective `*_test.go` file in the same package.
  - Tests should have both positive and negative cases with mocking if required.
  - Check `api/internal/handlers/albums_test.go` for reference on adding unit tests.

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
  - New files/functions should have a respective `test_*.py` file in the `worker/tests` directory.
  - Tests should have both positive and negative cases with mocking if required.
  - Check `worker/tests/components/test_place.py` for reference on adding unit tests.

### Integration Testing

#### Setup
- Run following command to setup the test suite:
```
make setup
```
- Run the command for invoking integration tests for all features:
```
make test-all
```
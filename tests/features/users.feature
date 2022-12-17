Feature: User Management

    Scenario: validate create user
        Given there are no users
        When create user is requested
        Then user is created
        When get user
        Then user is present
        When get all users
        Then user is present in list

    Scenario: validate update user
        Given there is user
        When update user is requested
        Then user is updated
        When get user
        Then user is present
        When get all users
        Then user is present in list

    Scenario: validate delete user
        Given there is user
        When delete user is requested
        Then user is deleted
        When get user
        Then user is not present
        When get all users
        Then user is not present in list
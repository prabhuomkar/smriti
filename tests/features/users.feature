Feature: User Management

    Scenario: Validate Create User
        Given there are no users
        When create user without auth
        Then auth error is found
        When create user with auth
        Then user is created
        When get user without auth
        Then auth error is found
        When get user with auth
        Then user is present
        When get all users without auth
        Then auth error is found
        When get all users with auth
        Then user is present in list

    Scenario: Validate Update User
        Given a user exists
        When update user without auth
        Then auth error is found
        When update user with auth
        Then user is updated
        When get user without auth
        Then auth error is found
        When get user with auth
        Then user is present
        When get all users without auth
        Then auth error is found
        When get all users with auth
        Then user is present in list

    Scenario: Validate Delete User
        Given a user exists
        When delete user without auth
        Then auth error is found
        When delete user with auth
        Then user is deleted
        When get user without auth
        Then auth error is found
        When get user with auth
        Then user is not present
        When get all users without auth
        Then auth error is found
        When get all users with auth
        Then user is not present in list
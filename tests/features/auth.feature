Feature: Authentication

    Scenario: Validate Login, Refresh and Logout
        Given a user is created
        When user logs in
        Then token is generated
        When user refreshes token
        Then token is refreshed
        When user logs out
        Then token is deleted
        Then a user is deleted
        
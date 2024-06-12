Feature: Authentication

    Scenario: Validate Login, Refresh and Logout
        Given a user default is created
        When user default logs in
        Then token is generated
        When user refreshes token
        Then token is refreshed
        When user logs out
        Then token is deleted
        Then a user default is deleted
        
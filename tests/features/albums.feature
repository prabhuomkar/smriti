Feature: Albums

    Scenario: Background
        Given a user is created

    Scenario: Validate Create Album
        When user logs in
        Then token is generated
        Given there are no albums
        When create album without auth
        Then auth error is found
        When create album with auth
        Then album is created
        When get album without auth
        Then auth error is found
        When get album with auth
        Then album is present
        When get all albums without auth
        Then auth error is found
        When get all albums with auth
        Then album is present in list

    Scenario: Validate Update Album
        When user logs in
        Then token is generated
        Given an album exists
        When update album without auth
        Then auth error is found
        When update album with auth
        Then album is updated
        When get album without auth
        Then auth error is found
        When get album with auth
        Then album is present
        When get all albums without auth
        Then auth error is found
        When get all albums with auth
        Then album is present in list

    Scenario: Validate Delete Album
        When user logs in
        Then token is generated
        Given an album exists
        When delete album without auth
        Then auth error is found
        When delete album with auth
        Then album is deleted
        When get album without auth
        Then auth error is found
        When get album with auth
        Then album is not present
        When get all albums without auth
        Then auth error is found
        When get all albums with auth
        Then album is not present in list
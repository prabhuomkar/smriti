Feature: Places

    Background: Setup User and MediaItem
        Given a user is created if does not exist
        When user logs in
        Then token is generated
        When upload default photo mediaitem with auth if does not exist and wait 10 seconds
        Then mediaitem is uploaded or exists

    Scenario: Validate Places
        Given a mediaitem exists with place
        When get all explored places for mediaitem without auth
        Then auth error is found
        When get all explored places for mediaitem with auth
        Then explored place is present in list
        When get explored place without auth
        Then auth error is found
        When get explored place with auth
        Then explored place is present
        When get all explored places without auth
        Then auth error is found
        When get all explored places with auth
        Then explored place is present in list

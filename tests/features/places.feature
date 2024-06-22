Feature: Places

    Background: Setup User and MediaItem
        Given a user default is created if does not exist
        When user default logs in
        Then token is generated
        When upload default photo mediaitem with auth if does not exist and wait 5 seconds
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
        Then explored place is present with cover mediaitem
        When get all explored places without auth
        Then auth error is found
        When get all explored places with auth
        Then explored place is present in list
        When get all mediaitems for place without auth
        Then auth error is found
        When get all mediaitems for place with auth
        Then mediaitem with place is present in list
        When delete mediaitem with auth
        Then mediaitem is deleted
        When get explored place with auth
        Then explored place is present without cover mediaitem

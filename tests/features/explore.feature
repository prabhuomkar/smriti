Feature: MediaItems

    Background: Setup User and MediaItem
        Given a user is created if does not exist
        When user logs in
        Then token is generated
        When upload photo mediaitem with auth and wait 3 seconds
        Then mediaitem is uploaded

    Scenario: Validate Places
        Given a mediaitem exists with place
        When get all places for mediaitem without auth
        Then auth error is found
        When get all places for mediaitem with auth
        Then place is present in list
        When get place without auth
        Then auth error is found
        When get place with auth
        Then place is present
        When get all places without auth
        Then auth error is found
        When get all places with auth
        Then place is present in list
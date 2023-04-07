Feature: MediaItems

    Background: Setup User
        Given a user is created if does not exist
        When user logs in
        Then token is generated
        When upload mediaitem with auth
        Then mediaitem is uploaded

    @clear
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
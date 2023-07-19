Feature: Search

    Background: Setup User and MediaItem
        Given a user is created if does not exist
        When user logs in
        Then token is generated
        When upload photo mediaitem with auth if does not exist and wait 6 seconds
        Then mediaitem is uploaded or exists

    Scenario: Search MediaItems
        Given a mediaitem exists
        When search for mediaitems without auth
        Then auth error is found
        When search for mediaitems with auth
        Then searched mediaitem is present in list

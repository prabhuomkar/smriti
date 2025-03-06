Feature: Search

    Background: Setup User and MediaItem
        Given a user default is created if does not exist
        When user default logs in
        Then token is generated
        When upload default photo mediaitem with auth if does not exist and wait 5 seconds
        Then mediaitem is uploaded or exists

    Scenario: Search MediaItems
        Given a mediaitem exists
        When search query pizza for mediaitems without auth
        Then auth error is found
        When search query pizza for mediaitems with auth
        Then searched mediaitem is present in list

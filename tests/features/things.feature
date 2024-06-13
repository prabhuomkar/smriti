Feature: Things

    Background: Setup User and MediaItem
        Given a user default is created if does not exist
        When user default logs in
        Then token is generated
        When upload default photo mediaitem with auth if does not exist and wait 10 seconds
        Then mediaitem is uploaded or exists

    Scenario: Validate Things
        Given a mediaitem exists with thing
        When get all explored things for mediaitem without auth
        Then auth error is found
        When get all explored things for mediaitem with auth
        Then explored thing is present in list
        When get explored thing without auth
        Then auth error is found
        When get explored thing with auth
        Then explored thing is present with cover mediaitem
        When get all explored things without auth
        Then auth error is found
        When get all explored things with auth
        Then explored thing is present in list
        When get all mediaitems for thing without auth
        Then auth error is found
        When get all mediaitems for thing with auth
        Then mediaitem with thing is present in list
        When delete mediaitem with auth
        Then mediaitem is deleted
        When get explored thing with auth
        Then explored thing is present without cover mediaitem

Feature: People

    Background: Setup User and MediaItem
        Given a user is created if does not exist
        When user logs in
        Then token is generated
        When upload faces.jpg photo mediaitem with auth if does not exist and wait 10 seconds
        Then mediaitem is uploaded or exists

    Scenario: Validate People
        # Given a mediaitem exists with people
        # When get all explored people for mediaitem without auth
        # Then auth error is found
        # When get all explored people for mediaitem with auth
        # Then explored person is present in list
        # When get explored person without auth
        # Then auth error is found
        # When get explored person with auth
        # Then explored person is present
        # When get all explored people without auth
        # Then auth error is found
        # When get all explored people with auth
        # Then explored person is present in list

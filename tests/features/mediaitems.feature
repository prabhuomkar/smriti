Feature: MediaItems

    Background: Setup User and MediaItem
        Given a user is created if does not exist
        When user logs in
        Then token is generated

    Scenario: Validate Create MediaItem
        Given there are no mediaitems
        When upload mediaitem without auth
        Then auth error is found
        When upload mediaitem with auth
        Then mediaitem is uploaded
        When get mediaitem without auth
        Then auth error is found
        When get mediaitem with auth
        Then mediaitem is present
        When get all mediaitems without auth
        Then auth error is found
        When get all mediaitems with auth
        Then mediaitem is present in list

    Scenario: Validate Update MediaItem
        Given a mediaitem exists
        When update mediaitem without auth
        Then auth error is found
        When update mediaitem with auth
        Then mediaitem is updated
        When get mediaitem without auth
        Then auth error is found
        When get mediaitem with auth
        Then mediaitem is present
        When get all mediaitems without auth
        Then auth error is found
        When get all mediaitems with auth
        Then mediaitem is present in list

    Scenario: Validate Delete MediaItem
        Given a mediaitem exists
        When delete mediaitem without auth
        Then auth error is found
        When delete mediaitem with auth
        Then mediaitem is deleted
        When get mediaitem without auth
        Then auth error is found
        When get mediaitem with auth
        Then mediaitem is present
        When get all mediaitems without auth
        Then auth error is found
        When get all mediaitems with auth
        Then mediaitem is not present in list
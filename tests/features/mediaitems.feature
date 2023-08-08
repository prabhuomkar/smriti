Feature: MediaItems

    Background: Setup User
        Given a user is created if does not exist
        When user logs in
        Then token is generated

    Scenario: Validate Create Photo MediaItem
        Given there are no mediaitems
        When upload default photo mediaitem without auth
        Then auth error is found
        When upload default photo mediaitem with auth
        Then mediaitem is uploaded
        When get mediaitem without auth
        Then auth error is found
        When get mediaitem with auth
        Then mediaitem is present
        When get all mediaitems without auth
        Then auth error is found
        When get all mediaitems with auth
        Then mediaitem is present in list

    Scenario: Validate Update Photo MediaItem
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

    Scenario: Validate Delete Photo MediaItem
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

    Scenario: Validate Create Video MediaItem
        Given there are no mediaitems
        When upload default video mediaitem without auth
        Then auth error is found
        When upload default video mediaitem with auth
        Then mediaitem is uploaded
        When get mediaitem without auth
        Then auth error is found
        When get mediaitem with auth
        Then mediaitem is present
        When get all mediaitems without auth
        Then auth error is found
        When get all mediaitems with auth
        Then mediaitem is present in list

    Scenario: Validate Delete Video MediaItem
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
Feature: Library

    Background: Setup User and MediaItem
        Given a user is created if does not exist
        When user logs in
        Then token is generated
        When upload photo mediaitem with auth and wait 3 seconds
        Then mediaitem is uploaded
        When get mediaitem with auth
        Then mediaitem is present

    Scenario: Validate Favourites
        Given a mediaitem exists
        When mark favourite mediaitem without auth
        Then auth error is found
        When mark favourite mediaitem with auth
        Then mediaitem is marked as favourite
        When get mediaitem with auth
        Then mediaitem is present with marked as favourite
        When get all favourite mediaitems without auth
        Then auth error is found
        When get all favourite mediaitems with auth
        Then mediaitem is present in favourites list
        When unmark favourite mediaitem without auth
        Then auth error is found
        When unmark favourite mediaitem with auth
        Then mediaitem is unmarked as favourite
        When get mediaitem with auth
        Then mediaitem is present with unmarked as favourite
        When get all favourite mediaitems with auth
        Then mediaitem is not present in list

    Scenario: Validate Hidden
        Given a mediaitem exists
        When mark hide mediaitem without auth
        Then auth error is found
        When mark hide mediaitem with auth
        Then mediaitem is marked as hidden
        When get mediaitem with auth
        Then mediaitem is present with marked as hidden
        When get all hidden mediaitems without auth
        Then auth error is found
        When get all hidden mediaitems with auth
        Then mediaitem is present in hidden list
        When unmark hide mediaitem without auth
        Then auth error is found
        When unmark hide mediaitem with auth
        Then mediaitem is unmarked as hidden
        When get mediaitem with auth
        Then mediaitem is present with unmarked as hidden
        When get all hidden mediaitems with auth
        Then mediaitem is not present in list

    Scenario: Validate Trash
        Given a mediaitem exists
        When mark delete mediaitem without auth
        Then auth error is found
        When mark delete mediaitem with auth
        Then mediaitem is marked as deleted
        When get mediaitem with auth
        Then mediaitem is present with marked as deleted
        When get all deleted mediaitems without auth
        Then auth error is found
        When get all deleted mediaitems with auth
        Then mediaitem is present in trash list
        When unmark delete mediaitem without auth
        Then auth error is found
        When unmark delete mediaitem with auth
        Then mediaitem is unmarked as deleted
        When get mediaitem with auth
        Then mediaitem is present with unmarked as deleted
        When get all deleted mediaitems with auth
        Then mediaitem is not present in list
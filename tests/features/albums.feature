Feature: Albums

    Background: Setup User and MediaItem
        Given a user default is created if does not exist
        When user default logs in
        Then token is generated
        When upload default photo mediaitem with auth if does not exist and wait 5 seconds
        Then mediaitem is uploaded or exists

    Scenario: Validate Create Album
        Given there are no albums
        When create album without auth
        Then auth error is found
        When create album with auth
        Then album is created
        When get album without auth
        Then auth error is found
        When get album with auth
        Then album is present
        When get all albums without auth
        Then auth error is found
        When get all albums with auth
        Then album is present in list

    Scenario: Validate Update Album
        Given an album exists
        When update album without auth
        Then auth error is found
        When update album with auth
        Then album is updated
        When get album without auth
        Then auth error is found
        When get album with auth
        Then album is present
        When get all albums without auth
        Then auth error is found
        When get all albums with auth
        Then album is present in list

    Scenario: Validate Add and Remove Album MediaItems
        Given an album exists
        When get album mediaitems without auth
        Then auth error is found
        When add album mediaitems without auth
        Then auth error is found
        When add album mediaitems with auth
        Then album mediaitems are added
        When get album mediaitems with auth
        Then album mediaitems are present
        When mark delete mediaitem with auth
        Then mediaitem is marked as deleted
        When get album mediaitems with auth
        Then album mediaitems are absent
        When unmark delete mediaitem with auth
        Then mediaitem is unmarked as deleted
        When get album mediaitems with auth
        Then album mediaitems are present
        When get album with auth
        Then album is updated after add album mediaitems
        When remove album mediaitems without auth
        Then auth error is found
        When remove album mediaitems with auth
        Then album mediaitems are removed
        When get album mediaitems with auth
        Then album mediaitems are absent
        When get album with auth
        Then album is updated after remove album mediaitems

    Scenario: Validate Delete Album
        Given an album exists
        When delete album without auth
        Then auth error is found
        When delete album with auth
        Then album is deleted
        When get album without auth
        Then auth error is found
        When get album with auth
        Then album is not present
        When get all albums without auth
        Then auth error is found
        When get all albums with auth
        Then album is not present in list
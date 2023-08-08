Feature: Sharing

    Background: Setup User and MediaItem
        Given a user is created if does not exist
        When user logs in
        Then token is generated
        When upload default photo mediaitem with auth if does not exist and wait 10 seconds
        Then mediaitem is uploaded or exists

    Scenario: Validate Create Shared Album
        Given there are no shared albums
        When create shared album without auth
        Then auth error is found
        When create shared album with auth
        Then shared album is created
        When get shared album without auth
        Then shared album is present
        When get shared album with auth
        Then shared album is present
        When get all shared albums without auth
        Then auth error is found
        When get all shared albums with auth
        Then shared album is present in list

    Scenario: Validate Add and Remove Shared Album MediaItems
        Given a shared album exists
        When get shared album mediaitems without auth
        Then shared album mediaitems are absent
        When add shared album mediaitems without auth
        Then auth error is found
        When add shared album mediaitems with auth
        Then shared album mediaitems are added
        When get shared album mediaitems with auth
        Then shared album mediaitems are present
        When mark delete mediaitem with auth
        Then mediaitem is marked as deleted
        When get shared album mediaitems with auth
        Then shared album mediaitems are absent
        When unmark delete mediaitem with auth
        Then mediaitem is unmarked as deleted
        When get shared album mediaitems with auth
        Then shared album mediaitems are present
        When get shared album with auth
        Then shared album is updated after add album mediaitems
        When remove shared album mediaitems without auth
        Then auth error is found
        When remove shared album mediaitems with auth
        Then shared album mediaitems are removed
        When get shared album mediaitems with auth
        Then shared album mediaitems are absent
        When get shared album with auth
        Then shared album is updated after remove album mediaitems

    Scenario: Validate Delete Shared Album
        Given a shared album exists
        When delete shared album without auth
        Then auth error is found
        When delete shared album with auth
        Then shared album is deleted
        When get shared album without auth
        Then shared album is not present
        When get shared album with auth
        Then shared album is not present
        When get all shared albums without auth
        Then auth error is found
        When get all shared albums with auth
        Then shared album is not present in list
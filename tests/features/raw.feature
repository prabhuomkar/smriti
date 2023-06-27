@long
Feature: RAW MediaItems

    Background: Setup User
        Given a user is created if does not exist
        When user logs in
        Then token is generated

    Scenario: Validate Create RAW MediaItem
        Given get list of raw mediaitems to upload
            | camera      |
            | adobe dng converter |
            | apple |
            | arashi vision |
            | autel |
            | blackmagic |
            | dji |
            | samsung |

        When upload raw mediaitems
        Then get raw mediaitems with auth and validate it is present

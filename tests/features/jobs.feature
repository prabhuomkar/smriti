@wip
Feature: Jobs

    Background: Setup User and MediaItem
        When create jobs user with auth
        Then user jobs is created
        When user jobs logs in
        Then token is generated
        When upload jobs.jpeg photo mediaitem with auth
        Then mediaitem is uploaded
        When update jobs user with auth
        Then user jobs is updated
        When user jobs logs in
        Then token is generated

    Scenario: Validate Jobs
        
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
        Given a mediaitem exists
        When get job mediaitem things
        Then job mediaitem related things are absent in list 
        When create jobs for places,classification components without auth
        Then auth error is found
        When create jobs for places,classification components with auth
        Then job is created
        When get jobs without auth and wait 0 seconds
        Then auth error is found
        When get jobs with auth and wait 0 seconds
        Then job is scheduled and present in list
        When get jobs with auth and wait until completed
        Then job is completed and present in list
        When get job without auth
        Then auth error is found
        When get job with auth
        Then job is present
        When get job mediaitem things
        Then job mediaitem related things are present in list
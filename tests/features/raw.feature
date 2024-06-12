@long
Feature: RAW MediaItems

    Background: Setup User
        Given a user default is created if does not exist
        When user default logs in
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
            | epson |
            | eyedeas |
            | fimi |
            | google |
            | hasselblad |
            | hmd global |
            | htc |
            | huawei |
            | kandao |
            | lg |
            | leaf |
            | leica |
            | madv |
            | microsoft |
            | nokia |
            | om digital solutions |
            | oneplus |
            | parrot |
            | phase one |
            | plustek |
            | raspberrypi |
            | realme |
            | ricoh |
            | samsung |
            | xiaoyi |
            | xiaomi |
            | yi technology |
            | yuneec |
            | asus |
            | bq |
            | motorola |
        Then raw mediaitems are ready to upload
        When upload raw mediaitems
        Then get raw mediaitems with auth and validate it is present

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
            | arri |
            | autel |
            | blackmagic |
            | canon |
            | casio |
            | dji |
            | epson |
            | eyedeas |
            | fimi |
            | fujifilm |
            | gitup |
            | gopro |
            | google |
            | hasselblad |
            | hmd global |
            | htc |
            | huawei |
            | hasselblad |
            | kandao |
            | kodak |
            | lg |
            | leaf |
            | leica |
            | light |
            | madv |
            | minolta |
            | mamiya |
            | microsoft |
            | minolta |
            | nikon |
            | nokia |
            | olympus |
            | om digital solutions |
            | olympus |
            | oneplus |
            | parrot |
            | panasonic |
            | paralenz |
            | parrot |
            | pentax |
            | phase one |
            | plustek |
            | polaroid |
            | raspberrypi |
            | realme |
            | ricoh |
            | sjcam |
            | samsung |
            | sigma |
            | sony |
            | xiaoyi |
            | xiaomi |
            | xiro |
            | yi technology |
            | yuneec |
            | asus |
            | bq |
            | moto g8plus |
            | motorola |

        When upload raw mediaitems
        Then get raw mediaitems with auth and validate it is present

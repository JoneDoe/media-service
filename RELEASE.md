[![Quality Gate Status](https://sonar.hubber.loc/api/project_badges/measure?project=media-service&metric=alert_status)](https://sonar.hubber.loc/dashboard?id=media-service)

# Release 1.1.0.0
#### Breaking Changes

* Add **Authorization Guard** on _File Upload_. This ability can be **turned on** by setting-up `FIREWALL_ENDPOINT` **env-variable**. 
After that use `Authorization Bearer <token>` to successfully passing an Authorization. `FIREWALL_ENDPOINT` will be using for verification your **token** 

# Release 1.0.0.0
#### Breaking Changes

* Supports `AWS S3` as file storage instead of local `File System`
* Implements a new endpoint for extended information about file  
    ````
    http://hostname:8080/info/:uuid
    ````
    Output is
    ````
    {
        "status": "ok",
        "data": {
            "fileName": "GeoEye_GeoEye1_50cm_8bit_RGB_DRA15.jpg",
            "uuid": "e4106f5d-6269-434e-9b29-700690aa9ea8",
            "key": "e4106f5d-6269-434e-9b29-700690aa9ea8.jpg",
            "url": "https://some-domain/images/e4106f5d-6269-434e-9b29-700690aa9ea8.jpg"
        }
    }
    ````
* New endpoint for upload files in sub-folders
    ```
    http://hostname:8080/files/upload/:folder
    ```
#### Bux Fixes

* Fix error message when does not exist `GRAYLOG_HOST` env-variable 
* Fix error message when does not exist `GRAYLOG_PORT` env-variable
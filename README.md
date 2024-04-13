# Golang storage server
## API Endpoints
### Upload 
Store file into service.
* URL `http://hostname:8080/files/upload`
* Content type `multipart/form-data`
* Method `POST`
* Data Params `files[]`
* Success Response `application/json`
    * Code: 201
    * Content:
    ````
    {
        "files": [
            {
                "fileName": "photo-1541727687969-ce40493cd847.jpeg",
                "uuid": "e9509777-3811-442c-9c7a-51c9f04f63eb"
            }
        ],
        "status": "ok"
    }
  ````
---
### Upload file(s) in sub-folder
Store file(s) into service in sub-folder named as [:folder].
* URL `http://hostname:8080/files/upload/:folder`
* Content type `multipart/form-data`
* Method `POST`
* URL Params
  
  Required: 
 
        folder=[string]
        
* Data Params `files[]`
* Success Response `application/json`
    * Code: 201
    * Content:
    ````
    {
        "files": [
            {
                "fileName": "photo-1541727687969-ce40493cd847.jpeg",
                "uuid": "e9509777-3811-442c-9c7a-51c9f04f63eb"
            }
        ],
        "status": "ok"
    }
  ````
---
### Delete
Delete stored file
* URL `http://hostname:8080/files/:uuid`
* Method `DELETE`
* URL Params
  
  Required: 
 
        uuid=[string]

* Success Response `application/json`
    * Code: 200
    * Content:
    ````
    {
      "status": "ok",
      "data": "e9509777-3811-442c-9c7a-51c9f04f63eb"
    }
    ````
---
### Get stored file
* URL `http://hostname:8080/files/:uuid`
* Method `GET`
* URL Params
  
  Required: 
 
        uuid=[string]
        
* Success Response: `binary-data`
---
### Get resized Image
* URL `http://hostname:8080/files/:uuid/:profile`
* Method `GET`
* URL Params
  
  Required: 
 
        uuid=[string]
        profile=[string]

* Available options of `profile`
        
        - small [500x500 px]
        - medium [1024x768 px]
        - thumbnail [164x164 px]
        
* Success Response: `binary-data`
---
### Get stored file info
* URL `http://hostname:8080/info/:uuid`
* Method `GET`
* URL Params
  
  Required: 
    
        uuid=[string]
        
* Success Response `application/json`
    * Code: 200
    * Content:
    ````
    {
        "status": "ok",
        "data": {
            "fileName": "GeoEye_GeoEye1_50cm_8bit_RGB_DRA15.jpg",
            "uuid": "e4106f5d-6269-434e-9b29-700690aa9ea8",
            "key": "e4106f5d-6269-434e-9b29-700690aa9ea8.jpg",
            "url": "https://cdn.27zxc.com/sc--media--dev/e4106f5d-6269-434e-9b29-700690aa9ea8.jpg"
        }
    }
    ````
___
### Download by source
* URL `http://hostname:8080/download`
* Method `POST`
* Request `application/json`
    ````
    [
        {
            "source": "https://www.fc-moto.de/WebRoot/FCMotoDB/Shops/10207048/5469/C75C/22BA/A5FF/63C3/4DEB/AE59/5639/Airoh-Aviator-21-Valor-AV2VA32_Bianco_ml.jpg",
            "thumbnail-profile": "thumbnail"
        },
        {
            "source": "https://www.fc-moto.de/WebRoot/FCMotoDB/Shops/10207048/5B8D/23E0/A19A/C919/75AC/AC1E/1404/1A54/Airoh_Aviator_2.2_Check_red_gloss_vorne_links_1_ml.jpg",
            "thumbnail-profile": "small"
        },
        {
             "source": "https://www.wallpapertip.com/wmimgs/40-409559_4k-hdr-gallery-hd-wallpapers-4k-wallpapers-4k.jpg",
             "thumbnail-profile": "medium"
        }
    ]
    ````
* Success Response `application/json`
    * Code: 201
    ````
    {
        "errors": [
            {
                "message": "net/http: request canceled (Client.Timeout exceeded while reading body)",
                "source": "https://www.wallpapertip.com/wmimgs/40-409559_4k-hdr-gallery-hd-wallpapers-4k-wallpapers-4k.jpg"
            }
        ],
        "files": [
            {
                "fileName": "Airoh-Aviator-21-Valor-AV2VA32_Bianco_ml.jpg",
                "uuid": "252adfbd-6322-4ece-a9b5-285efcfd0c38",
                "url": "https://lw02crgw11i.hubber.loc:443/sc--media--dev/default/25/2a/df/252adfbd-6322-4ece-a9b5-285efcfd0c38.jpg"
            },
            {
                "fileName": "Airoh_Aviator_2.2_Check_red_gloss_vorne_links_1_ml.jpg",
                "uuid": "bc080c51-de88-4e13-ac52-488bf55d7fa9",
                "url": "https://lw02crgw11i.hubber.loc:443/sc--media--dev/default/bc/08/0c/bc080c51-de88-4e13-ac52-488bf55d7fa9.jpg"
            }
        ]
    }
    ````
___
### Service health-check
* URL `http://hostname:8080/__healthcheck`
* Method `GET`
* Success Response `application/json`
    * Code: 200
    * Content:
    ````
    {
      "service": "media-service",
      "version": "1.0.0"
    }
    ````
---
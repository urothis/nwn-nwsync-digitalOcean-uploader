# nwn-nwsync-digitalOcean-uploader

This is an application that I use in my CI for my personal nwn projects.

This will ingest your nwsync files and upload them to a DigitalOcean Space.


# HOWTOUSE
1. Latest docker
2. To generate the access keys you need, follow [THIS](https://www.digitalocean.com/community/tutorials/how-to-create-a-digitalocean-space-and-api-key) guide
3. These values are all required, they are assigned in [env.list](https://github.com/urothis/nwn-nwsync-digitalOcean-uploader/env.list)

Key | Value
------------ | -------------
endpoint | ENDPOINT
accessKeyID | ACCESS_KEY
secretAccessKey | SECRET_KEY
bucketName | SPACE_NAME default:"nwn"
moduleName | MODULE_NAME default:"MyModule"   

Note: no spaces in MODULENAME. abbreviations are recommended.

4. Move nwsync files into the [nwsync](https://github.com/urothis/nwn-nwsync-digitalOcean-uploader/tree/master/nwsync) folder.
5. Run the relative .bat (windows) .sh (linux) for your os. 
6. Sit back and watch.

# Nwsync uploader for Digital Ocean Spaces
This is an application that I use in my CI for my personal nwn projects.

This will ingest your nwsync files and upload them to a DigitalOcean Space.

# Requirements
1. Latest Docker version 
2. [Digital Ocean Spaces](https://www.digitalocean.com/products/spaces/) account created
3. Access keys created for the above Digital Ocean Spaces account

# Creating Digital Ocean Spaces 
You can find this on the [Digital Ocean Official Documentation](https://www.digitalocean.com/community/tutorials/how-to-create-a-digitalocean-space-and-api-key#creating-an-access-key)

# How do I use it?
You will need to update the [env.list](https://github.com/urothis/nwn-nwsync-digitalOcean-uploader/blob/master/env.list) file with your Digital Ocean Space information. Table below assumes that your Digital Ocean Space URL is `https://nwsync.nyc3.digitaloceanspaces.com/`

Key         | Value
------------|--------------------------------
Endpoint    | Your space endpoint (i.e: nyc3.digitaloceanspaces.com)
ACCESS_KEY  | Your space access key
SECRET_KEY  | Your space secret key
SPACE_NAME  | Your space name (i.e: nwsync)
MODULE_NAME | This is used to support handling multiple modules in one Space. The name can be anything with no spaces (i.e: demo-module)

Once the file is updated with your Space, account information you will need to move all your nwsync generated files to the [`nwsync`](https://github.com/urothis/nwn-nwsync-digitalOcean-uploader/tree/master/nwsync) folder on this project. 

Once this is done, you will only need to run the script depending on your OS:
- If using Windows, use `uploadNWSync.bat`
- If using Linxu, use `uploadNWSync.sh`

# Cleaning up
After your uploads are complete, you might consider deleting the `README.md` from your Digital Ocean Spaces, since this is not used by Nwsync